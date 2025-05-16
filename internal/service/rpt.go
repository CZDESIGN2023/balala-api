package service

import (
	"context"
	pb "go-cs/api/rpt/v1"
	"go-cs/internal/biz"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type RptService struct {
	pb.UnimplementedRptServer

	uc  *biz.RptUsecase
	log *log.Helper
}

func NewRptService(uc *biz.RptUsecase, logger log.Logger) *RptService {
	moduleName := "RptService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &RptService{
		uc:  uc,
		log: hlog,
	}
}

func (s *RptService) SearchRptVersionWitem(ctx context.Context, req *pb.SearchRptVersionWitemRequest) (*pb.SearchRptVersionWitemReply, error) {

	reply := func(err error) (*pb.SearchRptVersionWitemReply, error) {
		return &pb.SearchRptVersionWitemReply{Result: &pb.SearchRptVersionWitemReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(strings.TrimSpace(req.TimeSplitType), "required"); vaildErr != nil {
		return reply(errs.Param(ctx, "TimeSplitType"))
	}

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	replyData, err := s.uc.SearchRptVersionWitem(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	return &pb.SearchRptVersionWitemReply{Result: &pb.SearchRptVersionWitemReply_Data{Data: replyData}}, nil
}

func (s *RptService) SearchRptMemberWitem(ctx context.Context, req *pb.SearchRptMemberWitemRequest) (*pb.SearchRptMemberWitemReply, error) {

	reply := func(err error) (*pb.SearchRptMemberWitemReply, error) {
		return &pb.SearchRptMemberWitemReply{Result: &pb.SearchRptMemberWitemReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(strings.TrimSpace(req.TimeSplitType), "required"); vaildErr != nil {
		return reply(errs.Param(ctx, "TimeSplitType"))
	}

	if vaildErr = validate.Var(req.SpaceId, "required,number,gt=0"); vaildErr != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	replyData, err := s.uc.SearchRptMemberWitem(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	return &pb.SearchRptMemberWitemReply{Result: &pb.SearchRptMemberWitemReply_Data{Data: replyData}}, nil
}

func (s *RptService) DashboardRptSpaceWitem(ctx context.Context, req *pb.DashboardRptSpaceWitemRequest) (*pb.DashboardRptSpaceWitemReply, error) {

	reply := func(err error) (*pb.DashboardRptSpaceWitemReply, error) {
		return &pb.DashboardRptSpaceWitemReply{Result: &pb.DashboardRptSpaceWitemReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(strings.TrimSpace(req.TimeSplitType), "required"); vaildErr != nil {
		return reply(errs.Param(ctx, "TimeSplitType"))
	}

	replyData, err := s.uc.DashboardRptSpaceWitem(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	return &pb.DashboardRptSpaceWitemReply{Result: &pb.DashboardRptSpaceWitemReply_Data{Data: replyData}}, nil
}

func (s *RptService) DashboardRptMemberWitem(ctx context.Context, req *pb.DashboardRptMemberWitemRequest) (*pb.DashboardRptMemberWitemReply, error) {

	reply := func(err error) (*pb.DashboardRptMemberWitemReply, error) {
		return &pb.DashboardRptMemberWitemReply{Result: &pb.DashboardRptMemberWitemReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(strings.TrimSpace(req.TimeSplitType), "required"); vaildErr != nil {
		return reply(errs.Param(ctx, "TimeSplitType"))
	}

	replyData, err := s.uc.DashboardRptMemberWitem(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	return &pb.DashboardRptMemberWitemReply{Result: &pb.DashboardRptMemberWitemReply_Data{Data: replyData}}, nil
}

func (s *RptService) DashboardSpaceList(ctx context.Context, req *pb.DashboardSpaceListRequest) (*pb.DashboardSpaceListReply, error) {

	reply := func(err error) (*pb.DashboardSpaceListReply, error) {
		return &pb.DashboardSpaceListReply{Error: errs.Cast(err)}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	replyData, err := s.uc.DashboardSpaceList(ctx, loginUser.UserId)
	if err != nil {
		return reply(err)
	}

	return &pb.DashboardSpaceListReply{Data: replyData}, nil
}

func (s *RptService) DashboardSpaceIncrWitem(ctx context.Context, req *pb.DashboardSpaceIncrWitemRequest) (*pb.DashboardSpaceIncrWitemReply, error) {

	reply := func(err error) (*pb.DashboardSpaceIncrWitemReply, error) {
		return &pb.DashboardSpaceIncrWitemReply{Error: errs.Cast(err)}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	replyData, err := s.uc.DashboardRptSpaceIncrWitem(ctx, loginUser.UserId, req)
	if err != nil {
		return reply(err)
	}

	return &pb.DashboardSpaceIncrWitemReply{Data: replyData}, nil
}
func (s *RptService) DashboardMemberIncrWitem(ctx context.Context, req *pb.DashboardMemberIncrWitemRequest) (*pb.DashboardMemberIncrWitemReply, error) {

	reply := func(err error) (*pb.DashboardMemberIncrWitemReply, error) {
		return &pb.DashboardMemberIncrWitemReply{Error: errs.Cast(err)}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	replyData, err := s.uc.DashboardRptMemberIncrWitem(ctx, loginUser.UserId, req)
	if err != nil {
		return reply(err)
	}

	return &pb.DashboardMemberIncrWitemReply{Data: replyData}, nil
}
