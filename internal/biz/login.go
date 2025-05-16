package biz

import (
	"context"
	"github.com/spf13/cast"
	"go-cs/internal/bean"
	"go-cs/internal/biz/command"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	domain "go-cs/internal/domain/user"
	user_repo "go-cs/internal/domain/user/repo"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/third_platform"
	"time"

	"go-cs/api/comm"
	pb "go-cs/api/login/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// 定义 Login 的操作接口
type LoginRepo interface {
	GenerateJwtToken(ctx context.Context, userId int64) (string, error)
	ClearAllJwtToken(ctx context.Context, userId int64) error
	ClearJwtTokenByToken(ctx context.Context, token string) error

	//登录验证码
	GenerateVerificationCode(ctx context.Context, name string) (string, error)
	GetVerificationCode(ctx context.Context, name string) (string, error)
	CleanVerificationCode(ctx context.Context, name string) error
	IncrAndGetCountCode(ctx context.Context, name string) (int64, error)

	SavePfTokenInfo(ctx context.Context, pfToken string, info string) error
	GetPfTokenInfo(ctx context.Context, pfToken string) string
	DelPfTokenInfo(ctx context.Context, pfToken string)
}

type LoginUsecase struct {
	repo                  LoginRepo
	userRepo              user_repo.UserRepo
	adminRepo             AdminRepo
	log                   *log.Helper
	tm                    trans.Transaction
	tpClient              *third_platform.Client
	addUserLoginLogCmd    *command.AddUserLoginLogCmd
	domainMessageProducer *domain_message.DomainMessageProducer
}

// NewLoginUsecase 初始化 LoginUsecase
func NewLoginUsecase(repo LoginRepo, tm trans.Transaction, userRepo user_repo.UserRepo, adminRepo AdminRepo,
	tpClient *third_platform.Client, addUserLoginLogCmd *command.AddUserLoginLogCmd, domainMessageProducer *domain_message.DomainMessageProducer, logger log.Logger) *LoginUsecase {
	moduleName := "LoginUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &LoginUsecase{
		repo:                  repo,
		userRepo:              userRepo,
		adminRepo:             adminRepo,
		log:                   hlog,
		tm:                    tm,
		tpClient:              tpClient,
		addUserLoginLogCmd:    addUserLoginLogCmd,
		domainMessageProducer: domainMessageProducer,
	}
}

func (uc *LoginUsecase) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReplyData, error) {
	var replyData *pb.LoginReplyData
	var user *domain.User
	var err error

	switch req.Way {
	default:
		return nil, errs.Business(ctx, "不支持的登陆方式: "+cast.ToString(req.Way))
	case pb.LoginRequest_Pwd: //密码登陆
		replyData, user, err = uc.LoginByPwd(ctx, req)
	case pb.LoginRequest_IM, pb.LoginRequest_Ql, pb.LoginRequest_Halala:
		replyData, user, err = uc.LoginByIM(ctx, req)
	}

	if err != nil {
		return nil, err
	}

	if replyData != nil {
		return replyData, nil
	}

	//账号状态
	if user.UserStatus == 0 {
		uc.addUserLoginLogCmd.Excute(ctx, user.Id, "登录失败:账号已禁用", 2)
		return nil, errs.New(ctx, comm.ErrorCode_LOGIN_ACCOUNT_DISABLED)
	}

	// 被设置为系统管理员
	needResetPwd, _ := uc.adminRepo.GetAdminNeedResetPwdStatus(ctx, user.Id)
	// 超管用弱密码登陆
	isWeakPwd := user.Role >= consts.SystemRole_Admin && utils.NewValidator().Var(req.Password, "stronger_password") != nil
	if req.Way == pb.LoginRequest_Pwd && (needResetPwd || isWeakPwd) {
		return &pb.LoginReplyData{
			NeedUpdatePwd: true,
			Role:          int64(user.Role),
		}, nil
	}

	//生成token，并入信息
	token, err2 := uc.repo.GenerateJwtToken(ctx, user.Id)

	if err2 != nil {
		return nil, errs.New(ctx, comm.ErrorCode_LOGIN_FAIL)
	}

	uc.CleanVerificationCode(ctx)

	ip := utils.GetIpFrom(ctx)
	user.UpdateLastLogin(ip, time.Now())
	uc.userRepo.SaveUser(ctx, user)

	uc.addUserLoginLogCmd.Excute(ctx, user.Id, "登录成功", 1)

	return &pb.LoginReplyData{
		Token:    "",
		JwtToken: token,
		User: &bean.User{
			Id:           user.Id,
			UserName:     user.UserName,
			Avatar:       user.Avatar,
			UserNickname: user.UserNickname,
			Role:         int64(user.Role),
		},
	}, nil
}

func (uc *LoginUsecase) LoginByIM(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReplyData, *domain.User, error) {
	var user *domain.User
	var err error

	platformCode := comm.ThirdPlatformCode(req.Way)

	pfAccount, err := uc.userRepo.GetThirdPfAccountByPfUserKey(ctx, req.PfToken, int32(platformCode))

	if err != nil {
		client := uc.tpClient.ByPfCode(platformCode)
		if client == nil {
			return nil, nil, errs.New(ctx, comm.ErrorCode_LOGIN_FAIL)
		}

		thirdUserInfo, err2 := client.GetUserInfo(req.PfToken)
		if err2 != nil {
			uc.log.Error(err2)
			return nil, nil, errs.Internal(ctx, err2)
		}

		if errs.IsDbRecordNotFoundErr(err) { //跳转到注册页面
			return &pb.LoginReplyData{
				NeedRegister: true,
				PfUserInfo: utils.ToJSON(map[string]any{
					"username": thirdUserInfo.UserName,
					"nickname": thirdUserInfo.NickName,
					"avatar":   thirdUserInfo.HeadAddr,
				}),
			}, nil, nil
		}
		return nil, nil, errs.Internal(ctx, err)
	}

	user, err = uc.userRepo.GetUserByUserId(ctx, pfAccount.UserId)
	if err != nil {
		return nil, nil, errs.New(ctx, comm.ErrorCode_DB_QUERY_FAIL)
	}

	if req.UserName != "" && req.Password != "" {

	}

	return nil, user, nil

}

func (uc *LoginUsecase) LoginByPwd(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReplyData, *domain.User, error) {
	var user *domain.User
	var err error

	user, err = uc.userRepo.GetUserByUserName(ctx, req.UserName)
	if err != nil {
		uc.GenerateVerificationCode(ctx)
		return nil, nil, errs.Custom(ctx, comm.ErrorCode_LOGIN_ACCOUNT_NOT_EXISTS, "用户名或密码错误")
	}

	//验证密码
	if utils.EncryptUserPassword(req.Password, user.UserSalt) != user.UserPassword {
		uc.GenerateVerificationCode(ctx)
		return nil, nil, errs.New(ctx, comm.ErrorCode_LOGIN_WRONG_ACCOUNT_OR_PASSWORD)
	}

	return nil, user, nil
}

func (uc *LoginUsecase) Logout(ctx context.Context) (*comm.CommonReply, error) {
	//清理缓存
	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser != nil {
		uc.addUserLoginLogCmd.Excute(ctx, loginUser.UserId, "退出登录", 3)
		uc.repo.ClearJwtTokenByToken(ctx, loginUser.JwtTokenId)
	}

	okReply := utils.NewCommonOkReply(ctx)
	return okReply, nil
}

func (uc *LoginUsecase) GenerateVerificationCode(ctx context.Context) (string, error) {
	countCode, err := uc.repo.IncrAndGetCountCode(ctx, "c")
	if err != nil {
		return "", err
	}
	if countCode < 3 {
		return "", nil
	}

	code, err := uc.repo.GenerateVerificationCode(ctx, "login_verification_code")
	return code, err
}

func (uc *LoginUsecase) CleanVerificationCode(ctx context.Context) error {
	err := uc.repo.CleanVerificationCode(ctx, "login_verification_code")
	if err != nil {
		return err
	}

	err = uc.repo.CleanVerificationCode(ctx, "c")
	if err != nil {
		return err
	}
	return nil
}

func (uc *LoginUsecase) GetVerificationCode(ctx context.Context) (string, error) {
	code, err := uc.repo.GetVerificationCode(ctx, "login_verification_code")
	return code, err
}

func (uc *LoginUsecase) GetVerificationCode2(ctx context.Context) (string, error) {
	code, err := uc.repo.GetVerificationCode(ctx, "login_verification_code")
	if err != nil {
		return "", err
	}

	if code != "" {
		countCode, err := uc.repo.IncrAndGetCountCode(ctx, "c")
		if err != nil {
			return "", err
		}
		if countCode < 3 {
			return "", nil
		}
	}

	return code, nil
}
