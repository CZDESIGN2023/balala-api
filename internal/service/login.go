package service

import (
	"context"
	"go-cs/api/comm"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"strings"

	"github.com/go-kratos/kratos/v2/log"

	pb "go-cs/api/login/v1"
	"go-cs/internal/biz"
	"go-cs/internal/conf"
)

type LoginService struct {
	pb.UnimplementedLoginServer
	uc  *biz.LoginUsecase
	log *log.Helper
}

type Geo struct {
	Country string `json:"country"`
	Area    string `json:"area"`
}

func NewLoginService(loginUsecase *biz.LoginUsecase, c *conf.Data,
	logger log.Logger,
) *LoginService {
	moduleName := "LoginService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	// 初始化ip库, 目前用server3 token登入, 停用功能
	// qqwry2.InitQQwry(c.Qqwry.DatPath)

	return &LoginService{
		uc:  loginUsecase,
		log: hlog,
	}
}

func (s *LoginService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {

	reply := func(err error) (*pb.LoginReply, error) {
		return &pb.LoginReply{Result: &pb.LoginReply_Error{Error: errs.Cast(err)}}, nil
	}

	var validErr error

	req.UserName = strings.TrimSpace(req.UserName)

	//校验参数
	validate := utils.NewValidator()

	switch req.Way {
	case pb.LoginRequest_Pwd:
		if validErr = validate.Var(req.UserName, "required,utf8Len=5-20,username"); validErr != nil {
			return reply(errs.New(ctx, comm.ErrorCode_LOGIN_WRONG_ACCOUNT_OR_PASSWORD))
		}

		if validErr = validate.Var(req.Password, "required,utf8Len=5-20,password"); validErr != nil {
			return reply(errs.New(ctx, comm.ErrorCode_LOGIN_WRONG_ACCOUNT_OR_PASSWORD))
		}

		//判断是否要过机器验证
		verificationCode, err := s.uc.GetVerificationCode(ctx)
		if err != nil {
			return reply(errs.Business(ctx, "人机校验失败"))
		}

		if verificationCode != "" && verificationCode != req.VerificationCode {
			return reply(errs.Business(ctx, "人机校验失败"))
		}
	case pb.LoginRequest_IM, pb.LoginRequest_Ql, pb.LoginRequest_Halala:
		if validErr = validate.Var(req.PfToken, "required"); validErr != nil {
			return reply(errs.Param(ctx, "PfToken"))
		}
	default:
		return reply(errs.Param(ctx, "Way"))
	}

	res, resErr := s.uc.Login(ctx, req)
	if res != nil && res.NeedRegister {
		return &pb.LoginReply{Result: &pb.LoginReply_Data{Data: res}}, nil
	}

	if resErr != nil {
		return reply(resErr)
	}

	return &pb.LoginReply{Result: &pb.LoginReply_Data{Data: res}}, nil
}

func (s *LoginService) Logout(ctx context.Context, req *pb.DoLogoutRequest) (*comm.CommonReply, error) {

	return s.uc.Logout(ctx)

}

func (s *LoginService) GetLoginValidCode(ctx context.Context, req *pb.GetLoginValidCodeRequest) (*pb.GetLoginValidCodeReply, error) {

	code, err := s.uc.GetVerificationCode(ctx)
	if err != nil {
		errInfo := errs.Business(ctx, "获取验证码失败")
		errReply := &pb.GetLoginValidCodeReply{Result: &pb.GetLoginValidCodeReply_Error{Error: errInfo}}
		return errReply, nil
	}

	return &pb.GetLoginValidCodeReply{Result: &pb.GetLoginValidCodeReply_Data{Data: code}}, nil
}
