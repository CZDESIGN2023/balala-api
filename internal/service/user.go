package service

import (
	"context"
	"go-cs/api/comm"
	pb "go-cs/api/user/v1"
	"go-cs/internal/biz"
	"go-cs/internal/conf"
	"go-cs/internal/utils"
	"math"
	"slices"
	"strings"

	"go-cs/internal/utils/errs"

	"github.com/go-kratos/kratos/v2/log"
)

type UserService struct {
	pb.UnimplementedUserServer
	log      *log.Helper
	user     *biz.UserUsecase
	upload   *biz.UploadUsecase
	login    *biz.LoginUsecase
	fileConf *conf.FileConfig
}

func NewUserService(fileConf *conf.FileConfig, logger log.Logger, stu *biz.UserUsecase, upload *biz.UploadUsecase, login *biz.LoginUsecase) *UserService {
	moduleName := "UserService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &UserService{
		user:     stu,
		upload:   upload,
		fileConf: fileConf,
		login:    login,
		log:      hlog,
	}
}

func (s *UserService) SetMyAvatar(ctx context.Context, req *pb.SetMyAvatarRequest) (*pb.SetMyAvatarReply, error) {
	var vaildErr error

	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.Avatar, "required"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER)
		errReply := &pb.SetMyAvatarReply{Result: &pb.SetMyAvatarReply_Error{Error: errInfo}}
		return errReply, nil
	}

	//从ctx中获取用户id
	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.SetMyAvatarReply{Result: &pb.SetMyAvatarReply_Error{Error: errInfo}}
		return errReply, nil
	}

	err := s.user.SetUserAvatar(ctx, loginUser.UserId, req.Avatar)
	if err != nil {
		errReply := &pb.SetMyAvatarReply{Result: &pb.SetMyAvatarReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.SetMyAvatarReply{Result: &pb.SetMyAvatarReply_Data{Data: &pb.SetMyAvatarReplyData{}}}
	return okReply, nil
}

func (s *UserService) RegUser(ctx context.Context, req *pb.RegUserRequest) (*comm.CommonReply, error) {
	reply := func(err error) (*comm.CommonReply, error) {
		return &comm.CommonReply{Result: &comm.CommonReply_Error{Error: errs.Cast(err)}}, nil
	}

	req.UserName = strings.TrimSpace(req.UserName)
	req.NickName = strings.TrimSpace(req.NickName)

	var err error
	//校验参数
	validate := utils.NewValidator()

	switch req.Way {
	case pb.RegUserRequest_Pwd:
		if err = validate.Var(req.UserName, "required,utf8Len=5-20,username"); err != nil {
			return reply(errs.New(ctx, comm.ErrorCode_USER_WRONG_RULE_USERNAME))
		}

		if err = validate.Var(req.NickName, "required,utf8Len=2-14,nickname"); err != nil {
			return reply(errs.New(ctx, comm.ErrorCode_USER_WRONG_RULE_NICKNAME))
		}

		if err = validate.Var(req.Password, "required,utf8Len=5-20,password"); err != nil {
			return reply(errs.New(ctx, comm.ErrorCode_USER_WRONG_RULE_PASSWORD))
		}
	case pb.RegUserRequest_IM, pb.RegUserRequest_QL, pb.RegUserRequest_Halala:
		if req.UserName != "" {
			if err = validate.Var(req.UserName, "required,utf8Len=5-20,username"); err != nil {
				return reply(errs.New(ctx, comm.ErrorCode_USER_WRONG_RULE_USERNAME))
			}

			if err = validate.Var(req.NickName, "required,utf8Len=2-14,nickname"); err != nil {
				return reply(errs.New(ctx, comm.ErrorCode_USER_WRONG_RULE_NICKNAME))
			}
		}

		if req.Password != "" {
			if err = validate.Var(req.Password, "required,utf8Len=5-20,password"); err != nil {
				return reply(errs.New(ctx, comm.ErrorCode_USER_WRONG_RULE_PASSWORD))
			}
		}

		if err = validate.Var(req.PfToken, "required"); err != nil {
			return reply(errs.Business(ctx, "PfToken"))
		}
	}

	user, err := s.user.Register(ctx, req)
	if err != nil {
		return reply(err)
	}

	// 清理人机验证码
	s.login.CleanVerificationCode(ctx)

	local, _ := s.upload.DownloadAvatarImgToLocal(ctx, user.Id, req.Avatar, s.fileConf.LocalPath)
	if local != "" {
		s.user.SetUserAvatar(ctx, user.Id, local)
	}

	okReply := utils.NewCommonOkReply(ctx)
	return okReply, nil
}

func (s *UserService) SearchUserList(ctx context.Context, req *pb.SearchUserListRequest) (*pb.SearchUserListReply, error) {
	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.SearchUserListReply{Result: &pb.SearchUserListReply_Error{Error: errInfo}}
		return errReply, nil
	}

	py := strings.TrimSpace(req.Py)

	data, err := s.user.SearchUserList(ctx, req.SpaceIds, py, req.UserIds)
	if err != nil {
		//不报错，返回一个空的数据
		errReply := &pb.SearchUserListReply{Result: &pb.SearchUserListReply_Data{Data: &pb.SearchUserListReplyData{List: nil}}}
		return errReply, nil
	}

	okReply := &pb.SearchUserListReply{Result: &pb.SearchUserListReply_Data{Data: data}}
	return okReply, nil
}

func (s *UserService) ChangeMyPwd(ctx context.Context, req *pb.ChangeMyPwdRequest) (*pb.ChangeMyPwdReply, error) {
	reply := func(err error) (*pb.ChangeMyPwdReply, error) {
		return &pb.ChangeMyPwdReply{Result: &pb.ChangeMyPwdReply_Error{Error: errs.Cast(err)}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId
	validate := utils.NewValidator()

	switch req.Type {
	case pb.ChangeMyPwdRequest_by_old_pwd:
		if err := validate.Var(req.NewPwd, "required,utf8Len=5-20,password"); err != nil {
			return reply(errs.New(ctx, comm.ErrorCode_USER_WRONG_RULE_PASSWORD))
		}
	case pb.ChangeMyPwdRequest_by_forceUpdate:
		if err := validate.Var(req.NewPwd, "required,utf8Len=5-20,stronger_password"); err != nil {
			return reply(errs.New(ctx, comm.ErrorCode_USER_WRONG_RULE_PASSWORD))
		}

	default:
		return reply(errs.Param(ctx, "Type"))
	}

	err := s.user.ChangeUserPwd(ctx, uid, req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.ChangeMyPwdReply{Result: &pb.ChangeMyPwdReply_Data{Data: &pb.ChangeMyPwdReplyData{}}}
	return okReply, nil
}

func (s *UserService) SetMyNickName(ctx context.Context, req *pb.SetMyNickNameRquest) (*pb.SetMyNickNameReply, error) {

	var vaildErr error
	//校验参数
	validate := utils.NewValidator()

	//检查密码, 二次密码是否符合规则
	if vaildErr = validate.Var(req.NickName, "required,utf8Len=2-14,nickname"); vaildErr != nil {
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_USER_WRONG_RULE_USERNAME)
		errReply := &pb.SetMyNickNameReply{Result: &pb.SetMyNickNameReply_Error{Error: errInfo}}
		return errReply, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.SetMyNickNameReply{Result: &pb.SetMyNickNameReply_Error{Error: errInfo}}
		return errReply, nil
	}

	err := s.user.ChangeUserNickName(ctx, loginUser.UserId, req.NickName)
	if err != nil {
		errReply := &pb.SetMyNickNameReply{Result: &pb.SetMyNickNameReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.SetMyNickNameReply{Result: &pb.SetMyNickNameReply_Data{Data: &pb.SetMyNickNameReplyData{}}}
	return okReply, nil
}

func (s *UserService) GetMyUserInfo(ctx context.Context, req *pb.GetMyUserInfoRequest) (*pb.GetMyUserInfoReply, error) {

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.GetMyUserInfoReply{Result: &pb.GetMyUserInfoReply_Error{Error: errInfo}}
		return errReply, nil
	}

	out, err := s.user.MyInfo(ctx, loginUser.UserId)
	if err != nil {
		errReply := &pb.GetMyUserInfoReply{Result: &pb.GetMyUserInfoReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.GetMyUserInfoReply{Result: &pb.GetMyUserInfoReply_Data{Data: out}}
	return okReply, nil
}

func (s *UserService) GetMySpaceMemberInfo(ctx context.Context, req *pb.GetMySpaceMemberInfoRequest) (*pb.GetMySpaceMemberInfoReply, error) {

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.GetMySpaceMemberInfoReply{Result: &pb.GetMySpaceMemberInfoReply_Error{Error: errInfo}}
		return errReply, nil
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "SpaceId")
		errReply := &pb.GetMySpaceMemberInfoReply{Result: &pb.GetMySpaceMemberInfoReply_Error{Error: errInfo}}
		return errReply, nil
	}

	out, err := s.user.GetMySpaceMemberInfo(ctx, loginUser.UserId, req.SpaceId)
	if err != nil {
		//参数检查失败
		errReply := &pb.GetMySpaceMemberInfoReply{Result: &pb.GetMySpaceMemberInfoReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.GetMySpaceMemberInfoReply{Result: &pb.GetMySpaceMemberInfoReply_Data{Data: out}}
	return okReply, nil
}

func (s *UserService) CheckUserName(ctx context.Context, req *pb.CheckUserNameRequest) (*pb.CheckUserNameReply, error) {

	var userName string = strings.TrimSpace(req.Name)
	var vaildErr error
	//校验参数
	validate := utils.NewValidator()

	//检查基本输入
	//判断用户名, 昵称是否符合规则
	if vaildErr = validate.Var(userName, "required,utf8Len=5-20,username"); vaildErr != nil {
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_USER_WRONG_RULE_USERNAME)
		errReply := &pb.CheckUserNameReply{Result: &pb.CheckUserNameReply_Error{Error: errInfo}}
		return errReply, nil
	}

	isExist, err := s.user.CheckUserName(ctx, userName)
	if err != nil {
		errReply := &pb.CheckUserNameReply{Result: &pb.CheckUserNameReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	if isExist {
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_USER_NAME_IS_EXIST)
		errReply := &pb.CheckUserNameReply{Result: &pb.CheckUserNameReply_Error{Error: errInfo}}
		return errReply, nil
	}

	okReply := &pb.CheckUserNameReply{Result: &pb.CheckUserNameReply_Data{Data: ""}}
	return okReply, nil
}

func (s *UserService) Bind(ctx context.Context, req *pb.BindRequest) (*pb.BindReply, error) {
	reply := func(err error) (*pb.BindReply, error) {
		return &pb.BindReply{Result: &pb.BindReply_Error{Error: errs.Cast(err)}}, nil
	}

	if pb.BindRequest_Type_name[int32(req.Type)] == "" {
		return reply(errs.Param(ctx, "Type"))
	}

	switch req.Type {
	case pb.BindRequest_IM, pb.BindRequest_Ql, pb.BindRequest_Halala:
		if req.Key == "" {
			return reply(errs.Param(ctx, "Key"))
		}
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.user.Bind(ctx, uid, req)
	if err != nil {
		return reply(err)
	}
	return &pb.BindReply{}, nil
}

func (s *UserService) Unbind(ctx context.Context, req *pb.UnbindRequest) (*pb.UnbindReply, error) {
	reply := func(err error) (*pb.UnbindReply, error) {
		return &pb.UnbindReply{Result: &pb.UnbindReply_Error{Error: errs.Cast(err)}}, nil
	}

	if pb.UnbindRequest_Type_name[int32(req.Type)] == "" {
		return reply(errs.Param(ctx, "Type"))
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.user.Unbind(ctx, uid, req)
	if err != nil {
		return reply(err)
	}
	return &pb.UnbindReply{}, nil
}

func (s *UserService) Cancel(ctx context.Context, req *pb.CancelRequest) (*pb.CancelReply, error) {
	reply := func(err error) (*pb.CancelReply, error) {
		return &pb.CancelReply{Result: &pb.CancelReply_Error{Error: errs.Cast(err)}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	if req.UserId == 0 {
		req.UserId = uid
	}

	err := s.user.Cancel(ctx, uid, req.UserId)
	if err != nil {
		return reply(err)
	}

	return &pb.CancelReply{}, nil
}

func (s *UserService) SetSpaceNotify(ctx context.Context, req *pb.SetSpaceNotifyRequest) (*pb.SetSpaceNotifyReply, error) {
	reply := func(err error) (*pb.SetSpaceNotifyReply, error) {
		return &pb.SetSpaceNotifyReply{Result: &pb.SetSpaceNotifyReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if !slices.Contains([]int32{0, 1}, req.Notify) {
		return reply(errs.Param(ctx, "Notify"))
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.user.SetSpaceNotify(ctx, uid, req)
	if err != nil {
		return reply(err)
	}

	return &pb.SetSpaceNotifyReply{}, nil
}

func (s *UserService) SetSpaceOrder(ctx context.Context, req *pb.SetSpaceOrderRequest) (*pb.SetSpaceOrderReply, error) {
	reply := func(err error) (*pb.SetSpaceOrderReply, error) {
		return &pb.SetSpaceOrderReply{Result: &pb.SetSpaceOrderReply_Error{Error: errs.Cast(err)}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	if req.FromIdx < 0 {
		return reply(errs.Param(ctx, "FromIdx"))
	}

	if req.ToIdx < 0 {
		return reply(errs.Param(ctx, "ToIdx"))
	}

	err := s.user.SetSpaceOrder(ctx, uid, req)
	if err != nil {
		return reply(err)
	}

	return &pb.SetSpaceOrderReply{}, nil
}

func (s *UserService) MyPendingWorkItem(ctx context.Context, req *pb.MyPendingWorkItemRequest) (*pb.MyPendingWorkItemReply, error) {
	reply := func(err error) (*pb.MyPendingWorkItemReply, error) {
		return &pb.MyPendingWorkItemReply{Result: &pb.MyPendingWorkItemReply_Error{Error: errs.Cast(err)}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	data, err := s.user.MyPendingWorkItem(ctx, uid)
	if err != nil {
		return reply(err)
	}

	return &pb.MyPendingWorkItemReply{Result: &pb.MyPendingWorkItemReply_Data{Data: data}}, nil
}

func (s *UserService) MyRelatedComment(ctx context.Context, req *pb.MyRelatedCommentRequest) (*pb.MyRelatedCommentReply, error) {
	reply := func(err error) (*pb.MyRelatedCommentReply, error) {
		return &pb.MyRelatedCommentReply{Result: &pb.MyRelatedCommentReply_Error{Error: errs.Cast(err)}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	if req.Size <= 0 || req.Size > 200 {
		req.Size = 20
	}

	if req.Pos <= 0 {
		req.Pos = math.MaxInt
	}

	data, err := s.user.MyRelatedComment(ctx, uid, req)
	if err != nil {
		return reply(err)
	}

	return &pb.MyRelatedCommentReply{Result: &pb.MyRelatedCommentReply_Data{Data: data}}, nil
}

func (s *UserService) MyRelatedCommentByIds(ctx context.Context, req *pb.MyRelatedCommentByIdsRequest) (*pb.MyRelatedCommentByIdsReply, error) {
	reply := func(err error) (*pb.MyRelatedCommentByIdsReply, error) {
		return &pb.MyRelatedCommentByIdsReply{Result: &pb.MyRelatedCommentByIdsReply_Error{Error: errs.Cast(err)}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	if len(req.Ids) == 0 {
		return reply(errs.Param(ctx, "ids"))
	}

	data, err := s.user.MyRelatedCommentByIds(ctx, uid, req.Ids)
	if err != nil {
		return reply(err)
	}

	return &pb.MyRelatedCommentByIdsReply{Result: &pb.MyRelatedCommentByIdsReply_Data{Data: data}}, nil
}

func (s *UserService) NotifyCount(ctx context.Context, req *pb.NotifyCountRequest) (*pb.NotifyCountReply, error) {
	reply := func(err error) (*pb.NotifyCountReply, error) {
		return &pb.NotifyCountReply{Result: &pb.NotifyCountReply_Error{Error: errs.Cast(err)}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	data, err := s.user.NotifyCount(ctx, uid)
	if err != nil {
		return reply(err)
	}

	return &pb.NotifyCountReply{Result: &pb.NotifyCountReply_Data{Data: data}}, nil
}

func (s *UserService) GetTempConfig(ctx context.Context, req *pb.GetTempConfigRequest) (*pb.GetTempConfigReply, error) {
	reply := func(err error) (*pb.GetTempConfigReply, error) {
		return &pb.GetTempConfigReply{Error: errs.Cast(err)}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	data, err := s.user.GetTempConfig(ctx, uid, req.GetKeys())
	if err != nil {
		return reply(err)
	}

	return &pb.GetTempConfigReply{
		Data: data,
	}, nil
}

func (s *UserService) SetTempConfig(ctx context.Context, req *pb.SetTempConfigRequest) (*pb.SetTempConfigReply, error) {
	reply := func(err error) (*pb.SetTempConfigReply, error) {
		return &pb.SetTempConfigReply{Error: errs.Cast(err)}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.user.SetTempConfig(ctx, uid, req.GetConfigs())
	if err != nil {
		return reply(err)
	}

	return &pb.SetTempConfigReply{}, nil
}

func (s *UserService) DelTempConfig(ctx context.Context, req *pb.DelTempConfigRequest) (*pb.DelTempConfigReply, error) {
	reply := func(err error) (*pb.DelTempConfigReply, error) {
		return &pb.DelTempConfigReply{Error: errs.Cast(err)}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.user.DelTempConfig(ctx, uid, req.GetKeys())
	if err != nil {
		return reply(err)
	}

	return &pb.DelTempConfigReply{}, nil
}

func (s *UserService) AllSpaceProfile(ctx context.Context, req *pb.AllSpaceProfileRequest) (*pb.AllSpaceProfileReply, error) {
	reply := func(err error) (*pb.AllSpaceProfileReply, error) {
		return &pb.AllSpaceProfileReply{Error: errs.Cast(err)}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	data, err := s.user.AllSpaceInfo(ctx, uid, req.UserId)
	if err != nil {
		return reply(err)
	}

	return &pb.AllSpaceProfileReply{Data: data}, nil
}

func (s *UserService) SetThirdPlatformNotify(ctx context.Context, req *pb.SetThirdPlatformNotifyRequest) (*pb.SetThirdPlatformNotifyReply, error) {
	reply := func(err error) (*pb.SetThirdPlatformNotifyReply, error) {
		return &pb.SetThirdPlatformNotifyReply{Error: errs.Cast(err)}, nil
	}

	if req.PlatformCode.String() == "" {
		return reply(errs.Param(ctx, "PlatformCode"))
	}

	if !slices.Contains([]int32{0, 1}, req.Notify) {
		return reply(errs.Param(ctx, "Notify"))
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.user.SetThirdPlatformNotify(ctx, uid, req)
	if err != nil {
		return reply(err)
	}

	return &pb.SetThirdPlatformNotifyReply{}, nil
}

func (s *UserService) SetUserConfig(ctx context.Context, req *pb.SetUserConfigRequest) (*pb.SetUserConfigReply, error) {
	reply := func(err error) (*pb.SetUserConfigReply, error) {
		return &pb.SetUserConfigReply{Error: errs.Cast(err)}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.user.SetUserConfig(ctx, uid, req)
	if err != nil {
		return reply(err)
	}

	return &pb.SetUserConfigReply{}, nil
}
