package service

import (
	"context"
	pb "go-cs/api/admin/v1"
	"go-cs/api/comm"
	"go-cs/internal/biz"
	"go-cs/internal/consts"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"slices"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type AdminService struct {
	pb.UnimplementedAdminServer

	uc  *biz.AdminUsecase
	log *log.Helper
}

func NewAdminService(AdminUsecase *biz.AdminUsecase, logger log.Logger) *AdminService {
	moduleName := "AdminService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &AdminService{
		uc:  AdminUsecase,
		log: hlog,
	}
}

func (s *AdminService) SearchUser(ctx context.Context, req *pb.SearchUserRequest) (*pb.SearchUserReply, error) {
	reply := func(err error) (*pb.SearchUserReply, error) {
		return &pb.SearchUserReply{Result: &pb.SearchUserReply_Error{Error: errs.Cast(err)}}, nil
	}

	req.Keyword = strings.TrimSpace(req.Keyword)

	if len(req.Sorts) != 0 {
		for _, by := range req.Sorts {
			if by.Order != "ASC" && by.Order != "DESC" {
				return reply(errs.Param(ctx, "Order"))
			}
		}
	}

	uid := utils.GetLoginUser(ctx).UserId

	data, err := s.uc.SearchUser(ctx, uid, req)
	if err != nil {
		return reply(err)
	}

	return &pb.SearchUserReply{Result: &pb.SearchUserReply_Data{Data: data}}, nil
}

func (s *AdminService) ResetPwd(ctx context.Context, req *pb.ResetPwdRequest) (*pb.ResetPwdReply, error) {
	reply := func(err error) (*pb.ResetPwdReply, error) {
		return &pb.ResetPwdReply{Result: &pb.ResetPwdReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)

	data, err := s.uc.ResetPwd(ctx, loginUser.UserId, req.UserId)
	if err != nil {
		return reply(err)
	}

	return &pb.ResetPwdReply{Result: &pb.ResetPwdReply_Data{Data: data}}, nil
}

func (s *AdminService) SetRole(ctx context.Context, req *pb.SetRoleRequest) (*pb.SetRoleReply, error) {
	reply := func(err error) (*pb.SetRoleReply, error) {
		return &pb.SetRoleReply{Result: &pb.SetRoleReply_Error{Error: errs.Cast(err)}}, nil
	}

	if !slices.Contains(consts.GetAllSystemRoles(), consts.SystemRole(req.Role)) {
		return reply(errs.Param(ctx, "Role"))
	}

	uid := utils.GetLoginUser(ctx).UserId

	err := s.uc.SetRole(ctx, uid, req.UserId, consts.SystemRole(req.Role))
	if err != nil {
		return reply(err)
	}

	return &pb.SetRoleReply{Result: nil}, nil
}

func (s *AdminService) SetNickname(ctx context.Context, req *pb.SetNicknameRequest) (*pb.SetNicknameReply, error) {
	reply := func(err error) (*pb.SetNicknameReply, error) {
		return &pb.SetNicknameReply{Result: &pb.SetNicknameReply_Error{Error: errs.Cast(err)}}, nil
	}

	validate := utils.NewValidator()
	if err := validate.Var(req.Value, "required,utf8Len=2-14,nickname"); err != nil {
		return reply(errs.New(ctx, comm.ErrorCode_USER_WRONG_RULE_NICKNAME))
	}

	loginUser := utils.GetLoginUser(ctx)

	err := s.uc.SetNickname(ctx, loginUser.UserId, req.UserId, req.Value)
	if err != nil {
		return reply(err)
	}

	return &pb.SetNicknameReply{Result: nil}, nil
}

func (s *AdminService) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserReply, error) {
	reply := func(err error) (*pb.AddUserReply, error) {
		return &pb.AddUserReply{Result: &pb.AddUserReply_Error{Error: errs.Cast(err)}}, nil
	}

	validate := utils.NewValidator()

	// 用户名
	if err := validate.Var(req.Username, "required,utf8Len=5-20,username"); err != nil {
		return reply(errs.New(ctx, comm.ErrorCode_USER_WRONG_RULE_USERNAME))
	}

	// 昵称
	if err := validate.Var(req.Nickname, "required,utf8Len=2-14,nickname"); err != nil {
		return reply(errs.New(ctx, comm.ErrorCode_USER_WRONG_RULE_NICKNAME))
	}

	// 密码
	if err := validate.Var(req.Password, "required,utf8Len=5-20,password"); err != nil {
		return reply(errs.New(ctx, comm.ErrorCode_USER_WRONG_RULE_PASSWORD))
	}

	loginUser := utils.GetLoginUser(ctx)
	err := s.uc.AddUser(ctx, loginUser.UserId, req)
	if err != nil {
		return reply(err)
	}

	return &pb.AddUserReply{Result: nil}, nil
}
