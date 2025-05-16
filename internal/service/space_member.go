package service

import (
	"context"
	"go-cs/api/comm"
	pb "go-cs/api/space_member/v1"
	"go-cs/internal/bean/vo"
	"go-cs/internal/biz"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type SpaceMemberService struct {
	pb.UnimplementedSpaceMemberServer

	uc *biz.SpaceMemberUsecase

	log *log.Helper
}

func NewSpaceMemberService(spaceMemberUsecase *biz.SpaceMemberUsecase, logger log.Logger) *SpaceMemberService {
	moduleName := "SpaceMemberService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceMemberService{
		uc:  spaceMemberUsecase,
		log: hlog,
	}
}

func (s *SpaceMemberService) AddSpaceMember(ctx context.Context, req *pb.AddSpaceMemeberRequest) (*pb.AddSpaceMemeberReply, error) {
	var vaildErr error

	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "SpaceId")
		errReply := &pb.AddSpaceMemeberReply{Result: &pb.AddSpaceMemeberReply_Error{Error: errInfo}}
		return errReply, nil
	}

	if vaildErr = validate.Var(req.Users, "required,gt=0,dive"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "Users")
		errReply := &pb.AddSpaceMemeberReply{Result: &pb.AddSpaceMemeberReply_Error{Error: errInfo}}
		return errReply, nil
	}

	for i := 0; i < len(req.Users); i++ {
		//检测内容
		if vaildErr = validate.Var(req.Users[i].UserId, "required,number,gt=0"); vaildErr != nil {
			//参数检查失败
			errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "Users.UserId")
			errReply := &pb.AddSpaceMemeberReply{Result: &pb.AddSpaceMemeberReply_Error{Error: errInfo}}
			return errReply, nil
		}

		//检测内容
		if vaildErr = validate.Var(req.Users[i].RoleId, "required,number,gt=0,oneof=1 2 3"); vaildErr != nil {
			//参数检查失败
			errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "Users.RoleId")
			errReply := &pb.AddSpaceMemeberReply{Result: &pb.AddSpaceMemeberReply_Error{Error: errInfo}}
			return errReply, nil
		}
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.AddSpaceMemeberReply{Result: &pb.AddSpaceMemeberReply_Error{Error: errInfo}}
		return errReply, nil
	}

	in := &vo.AddSpaceMembersVo{}
	in.SpaceId = req.SpaceId
	for i := 0; i < len(req.Users); i++ {
		in.AddMemberUser(req.Users[i].UserId, int64(req.Users[i].RoleId))
	}

	_, err := s.uc.AddMySpaceMembers(ctx, loginUser, in)
	if err != nil {
		//新增成员失败
		errReply := &pb.AddSpaceMemeberReply{Result: &pb.AddSpaceMemeberReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.AddSpaceMemeberReply{Result: &pb.AddSpaceMemeberReply_Data{Data: &pb.AddSpaceMemeberReplyData{}}}
	return okReply, nil
}

func (s *SpaceMemberService) RemoveSpaceMember(ctx context.Context, req *pb.RemoveSpaceMemeberRequest) (*pb.RemoveSpaceMemeberReply, error) {
	var validErr error

	validate := utils.NewValidator()
	if validErr = validate.Var(req.SpaceId, "required,number,gt=0"); validErr != nil {
		//参数检查失败
		errInfo := errs.Param(ctx, "SpaceId")
		errReply := &pb.RemoveSpaceMemeberReply{Result: &pb.RemoveSpaceMemeberReply_Error{Error: errInfo}}
		return errReply, nil
	}

	if validErr = validate.Var(req.UserId, "required,number,gt=0"); validErr != nil {
		//参数检查失败
		errInfo := errs.Param(ctx, "UserId")
		errReply := &pb.RemoveSpaceMemeberReply{Result: &pb.RemoveSpaceMemeberReply_Error{Error: errInfo}}
		return errReply, nil
	}

	if req.TargetUserId < 0 {
		//参数检查失败
		errInfo := errs.Param(ctx, "TargetUserId")
		errReply := &pb.RemoveSpaceMemeberReply{Result: &pb.RemoveSpaceMemeberReply_Error{Error: errInfo}}
		return errReply, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := errs.NotLogin(ctx)
		errReply := &pb.RemoveSpaceMemeberReply{Result: &pb.RemoveSpaceMemeberReply_Error{Error: errInfo}}
		return errReply, nil
	}

	err := s.uc.KickOut(ctx, loginUser.UserId, req.SpaceId, req.UserId, req.TargetUserId)
	if err != nil {
		//新增成员失败
		errReply := &pb.RemoveSpaceMemeberReply{Result: &pb.RemoveSpaceMemeberReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.RemoveSpaceMemeberReply{Result: nil}
	return okReply, nil

}

func (s *SpaceMemberService) SpaceMemberList(ctx context.Context, req *pb.SpaceMemeberListRequest) (*pb.SpaceMemeberListReply, error) {

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.SpaceMemeberListReply{Result: &pb.SpaceMemeberListReply_Error{Error: errInfo}}
		return errReply, nil
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER)
		errReply := &pb.SpaceMemeberListReply{Result: &pb.SpaceMemeberListReply_Error{Error: errInfo}}
		return errReply, nil
	}

	req.UserName = strings.TrimSpace(req.UserName)

	list, err := s.uc.GetMySpaceMemberList(ctx, loginUser.UserId, req.SpaceId, req.UserName)
	if err != nil {
		//参数检查失败
		errReply := &pb.SpaceMemeberListReply{Result: &pb.SpaceMemeberListReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.SpaceMemeberListReply{Result: &pb.SpaceMemeberListReply_Data{Data: &pb.SpaceMemeberListReplyData{
		List: list,
	}}}
	return okReply, nil

}

func (s *SpaceMemberService) SetSpaceMemberRole(ctx context.Context, req *pb.SetSpaceMemberRoleRequest) (*pb.SetSpaceMemberRoleReply, error) {
	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.SetSpaceMemberRoleReply{Result: &pb.SetSpaceMemberRoleReply_Error{Error: errInfo}}
		return errReply, nil
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "SpaceId")
		errReply := &pb.SetSpaceMemberRoleReply{Result: &pb.SetSpaceMemberRoleReply_Error{Error: errInfo}}
		return errReply, nil
	}

	if vaildErr = validate.Var(req.UserId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "UserId")
		errReply := &pb.SetSpaceMemberRoleReply{Result: &pb.SetSpaceMemberRoleReply_Error{Error: errInfo}}
		return errReply, nil
	}

	if vaildErr = validate.Var(req.RoleId, "required,number,gt=0,oneof=1 2 3"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "RoleId")
		errReply := &pb.SetSpaceMemberRoleReply{Result: &pb.SetSpaceMemberRoleReply_Error{Error: errInfo}}
		return errReply, nil
	}

	err := s.uc.SetMySpaceMemberRoleId(ctx, loginUser, req.SpaceId, req.UserId, int(req.RoleId))
	if err != nil {
		//参数检查失败
		errReply := &pb.SetSpaceMemberRoleReply{Result: &pb.SetSpaceMemberRoleReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.SetSpaceMemberRoleReply{Result: &pb.SetSpaceMemberRoleReply_Data{Data: &pb.SetSpaceMemberRoleReplyData{}}}
	return okReply, nil

}

func (s *SpaceMemberService) GetSpaceMemberWorkItemCount(ctx context.Context, req *pb.GetSpaceMemberWorkItemCountRequest) (*pb.GetSpaceMemberWorkItemCountReply, error) {

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.GetSpaceMemberWorkItemCountReply{Result: &pb.GetSpaceMemberWorkItemCountReply_Error{Error: errInfo}}
		return errReply, nil
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "SpaceId")
		errReply := &pb.GetSpaceMemberWorkItemCountReply{Result: &pb.GetSpaceMemberWorkItemCountReply_Error{Error: errInfo}}
		return errReply, nil
	}

	if vaildErr = validate.Var(req.UserId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "UserId")
		errReply := &pb.GetSpaceMemberWorkItemCountReply{Result: &pb.GetSpaceMemberWorkItemCountReply_Error{Error: errInfo}}
		return errReply, nil
	}

	out, err := s.uc.GetSpaceMemberWorkItemCountV2(ctx, loginUser.UserId, req.UserId, req.SpaceId)
	if err != nil {
		//参数检查失败
		errReply := &pb.GetSpaceMemberWorkItemCountReply{Result: &pb.GetSpaceMemberWorkItemCountReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.GetSpaceMemberWorkItemCountReply{Result: &pb.GetSpaceMemberWorkItemCountReply_Data{Data: &pb.GetSpaceMemberWorkItemCountReplyData{
		COUNT: out,
	}}}
	return okReply, nil

}

func (s *SpaceMemberService) AddSpaceManager(ctx context.Context, req *pb.AddSpaceManagerRequest) (*pb.AddSpaceManagerReply, error) {
	var vaildErr error

	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "SpaceId")
		errReply := &pb.AddSpaceManagerReply{Result: &pb.AddSpaceManagerReply_Error{Error: errInfo}}
		return errReply, nil
	}

	if vaildErr = validate.Var(req.UserIds, "required,gt=0,dive"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER_FMT, "UserIds")
		errReply := &pb.AddSpaceManagerReply{Result: &pb.AddSpaceManagerReply_Error{Error: errInfo}}
		return errReply, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.AddSpaceManagerReply{Result: &pb.AddSpaceManagerReply_Error{Error: errInfo}}
		return errReply, nil
	}

	_, err := s.uc.AddMySpaceManager(ctx, loginUser, req.SpaceId, req.UserIds)
	if err != nil {
		//新增成员失败
		errReply := &pb.AddSpaceManagerReply{Result: &pb.AddSpaceManagerReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.AddSpaceManagerReply{Result: &pb.AddSpaceManagerReply_Data{Data: &pb.AddSpaceManagerReplyData{}}}
	return okReply, nil
}

func (s *SpaceMemberService) RemoveSpaceManager(ctx context.Context, req *pb.RemoveSpaceManagerRequest) (*pb.RemoveSpaceManagerReply, error) {
	reply := func(err error) (*pb.RemoveSpaceManagerReply, error) {
		return &pb.RemoveSpaceManagerReply{Result: &pb.RemoveSpaceManagerReply_Error{Error: errs.Cast(err)}}, err
	}

	validate := utils.NewValidator()
	if err := validate.Var(req.SpaceId, "required,number,gt=0"); err != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if err := validate.Var(req.UserIds, "required,gt=0,dive"); err != nil {
		return reply(errs.Param(ctx, "UserIds"))
	}

	_, err := s.uc.RemoveMySpaceManager(ctx, utils.GetLoginUser(ctx), req.SpaceId, req.UserIds)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.RemoveSpaceManagerReply{Result: &pb.RemoveSpaceManagerReply_Data{Data: nil}}
	return okReply, nil

}

func (s *SpaceMemberService) SpaceManagerList(ctx context.Context, req *pb.SpaceManagerListRequest) (*pb.SpaceManagerListReply, error) {

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.SpaceManagerListReply{Result: &pb.SpaceManagerListReply_Error{Error: errInfo}}
		return errReply, nil
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER)
		errReply := &pb.SpaceManagerListReply{Result: &pb.SpaceManagerListReply_Error{Error: errInfo}}
		return errReply, nil
	}

	list, err := s.uc.GetMySpaceManagerList(ctx, loginUser.UserId, req.SpaceId)
	if err != nil {
		//参数检查失败
		errReply := &pb.SpaceManagerListReply{Result: &pb.SpaceManagerListReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.SpaceManagerListReply{Result: &pb.SpaceManagerListReply_Data{Data: &pb.SpaceManagerListReplyData{
		List: list,
	}}}
	return okReply, nil

}

func (s *SpaceMemberService) SpaceMemberById(ctx context.Context, req *pb.SpaceMemberByIdRequest) (*pb.SpaceMemberByIdReply, error) {

	reply := func(err error) (*pb.SpaceMemberByIdReply, error) {
		return &pb.SpaceMemberByIdReply{Result: &pb.SpaceMemberByIdReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if vaildErr = validate.Var(req.UserIds, "required,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "UserIds"))
	}

	result, err := s.uc.QSpaceMemberByIds(ctx, loginUser, req.SpaceId, req.UserIds)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SpaceMemberByIdReply{Result: &pb.SpaceMemberByIdReply_Data{Data: result}}
	return okReply, nil

}
