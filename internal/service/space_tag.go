package service

import (
	"context"
	"go-cs/api/comm"
	pb "go-cs/api/space_tag/v1"
	"go-cs/internal/biz"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/pkg/char"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type SpaceTagService struct {
	pb.UnimplementedSpaceTagServer
	log *log.Helper
	uc  *biz.SpaceTagUsecase
}

func NewSpaceTagService(stu *biz.SpaceTagUsecase, logger log.Logger) *SpaceTagService {
	moduleName := "SpaceTagService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceTagService{
		uc:  stu,
		log: hlog,
	}
}

func (s *SpaceTagService) CreateSpaceTag(ctx context.Context, req *pb.CreateSpaceTagRequest) (*pb.CreateSpaceTagReply, error) {
	var vaildErr error

	req.TagName = char.Filter(req.TagName)

	validate := utils.NewValidator()
	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER)
		errReply := &pb.CreateSpaceTagReply{Result: &pb.CreateSpaceTagReply_Error{Error: errInfo}}
		return errReply, nil
	}

	if vaildErr = validate.Var(strings.TrimSpace(req.TagName), "required,utf8Len=2-12"); vaildErr != nil {
		//参数检查失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER)
		errReply := &pb.CreateSpaceTagReply{Result: &pb.CreateSpaceTagReply_Error{Error: errInfo}}
		return errReply, nil
	}

	//从ctx中获取用户id
	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		//用户信息获取失败
		errInfo := utils.NewErrorInfo(ctx, comm.ErrorCode_LOGIN_USER_NOT_LOGIN)
		errReply := &pb.CreateSpaceTagReply{Result: &pb.CreateSpaceTagReply_Error{Error: errInfo}}
		return errReply, nil
	}

	out, err := s.uc.CreateMySpaceTag(ctx, loginUser, req.SpaceId, strings.TrimSpace(req.TagName))
	if err != nil {
		errReply := &pb.CreateSpaceTagReply{Result: &pb.CreateSpaceTagReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.CreateSpaceTagReply{Result: &pb.CreateSpaceTagReply_Data{Data: out.ToProto()}}
	return okReply, nil
}

func (s *SpaceTagService) ModifySpaceTagName(ctx context.Context, req *pb.ModifySpaceTagNameRequest) (*pb.ModifySpaceTagNameReply, error) {

	var reply = func(err error) (*pb.ModifySpaceTagNameReply, error) {
		return &pb.ModifySpaceTagNameReply{Result: &pb.ModifySpaceTagNameReply_Error{Error: errs.Cast(err)}}, nil
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

	if vaildErr = validate.Var(req.TagId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "TagId"))
	}

	if vaildErr = validate.Var(strings.TrimSpace(req.TagName), "required,utf8Len=2-12"); vaildErr != nil {
		return reply(errs.Param(ctx, "TagName"))
	}

	out, err := s.uc.ModifyMySpaceTagName(ctx, loginUser, req.SpaceId, req.TagId, strings.TrimSpace(req.TagName))
	if err != nil {
		//用户信息获取失败
		errReply := &pb.ModifySpaceTagNameReply{Result: &pb.ModifySpaceTagNameReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.ModifySpaceTagNameReply{Result: &pb.ModifySpaceTagNameReply_Data{Data: out.ToProto()}}
	return okReply, nil
}

func (s *SpaceTagService) DelSpaceTag(ctx context.Context, req *pb.DelSpaceTagRequest) (*pb.DelSpaceTagReply, error) {

	var reply = func(err error) (*pb.DelSpaceTagReply, error) {
		return &pb.DelSpaceTagReply{Result: &pb.DelSpaceTagReply_Error{Error: errs.Cast(err)}}, nil
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

	if vaildErr = validate.Var(req.TagId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "TagId"))
	}

	out, err := s.uc.DelMySpaceTagV2(ctx, loginUser, req.SpaceId, req.TagId)
	if err != nil {
		//用户信息获取失败
		errReply := &pb.DelSpaceTagReply{Result: &pb.DelSpaceTagReply_Error{Error: errs.Cast(err)}}
		return errReply, nil
	}

	okReply := &pb.DelSpaceTagReply{Result: &pb.DelSpaceTagReply_Data{Data: out.ToProto()}}
	return okReply, nil
}

func (s *SpaceTagService) GetSpaceTagListV2(ctx context.Context, req *pb.GetSpaceTagListRequestV2) (*pb.GetSpaceTagListReplyV2, error) {

	reply := func(err *comm.ErrorInfo) (*pb.GetSpaceTagListReplyV2, error) {
		return &pb.GetSpaceTagListReplyV2{Result: &pb.GetSpaceTagListReplyV2_Error{Error: err}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var err error
	validate := utils.NewValidator()
	if err = validate.Var(req.SpaceId, "required,number,gt=0"); err != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	out, err := s.uc.GetMySpaceTagListV2(ctx, loginUser.UserId, req.SpaceId)
	if err != nil {
		return reply(errs.Cast(err))
	}

	okReply := &pb.GetSpaceTagListReplyV2{Result: &pb.GetSpaceTagListReplyV2_Data{Data: &pb.GetSpaceTagListReplyV2Data{List: out}}}
	return okReply, nil
}
