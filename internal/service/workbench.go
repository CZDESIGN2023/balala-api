package service

import (
	"context"
	"go-cs/api/comm"
	pb "go-cs/api/workbench/v1"
	"go-cs/internal/biz"
	"go-cs/internal/utils"
	"go-cs/internal/utils/date"
	"go-cs/internal/utils/errs"

	"github.com/go-kratos/kratos/v2/log"
)

type WorkbenchService struct {
	pb.WorkbenchHTTPClientImpl

	staticsUc *biz.StaticsUsecase
	log       *log.Helper
}

func NewWorkbenchService(staticsUc *biz.StaticsUsecase, logger log.Logger) *WorkbenchService {
	moduleName := "WorkbenchService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &WorkbenchService{
		staticsUc: staticsUc,
		log:       hlog,
	}
}

func (s *WorkbenchService) GetWorkBenchCount(ctx context.Context, req *pb.GetWorkBenchCountRequest) (*pb.GetWorkBenchCountReply, error) {

	reply := func(err *comm.ErrorInfo) (*pb.GetWorkBenchCountReply, error) {
		return &pb.GetWorkBenchCountReply{Result: &pb.GetWorkBenchCountReply_Error{Error: err}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	data, err := s.staticsUc.GetWorkbenchCount(ctx, loginUser.UserId)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &pb.GetWorkBenchCountReply{Result: &pb.GetWorkBenchCountReply_Data{Data: data}}, nil
}

func (s *WorkbenchService) GetSpaceWorkBenchCount(ctx context.Context, req *pb.GetSpaceWorkBenchCountRequest) (*pb.GetSpaceWorkBenchCountReply, error) {

	reply := func(err *comm.ErrorInfo) (*pb.GetSpaceWorkBenchCountReply, error) {
		return &pb.GetSpaceWorkBenchCountReply{Result: &pb.GetSpaceWorkBenchCountReply_Error{Error: err}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	data, err := s.staticsUc.GetSpaceWorkbenchCount(ctx, loginUser.UserId, req.SpaceId)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &pb.GetSpaceWorkBenchCountReply{Result: &pb.GetSpaceWorkBenchCountReply_Data{Data: data}}, nil
}

func (s *WorkbenchService) GetSpaceWorkBenchCount2(ctx context.Context, req *pb.GetSpaceWorkBenchCountRequest2) (*pb.GetSpaceWorkBenchCountReply2, error) {

	reply := func(err *comm.ErrorInfo) (*pb.GetSpaceWorkBenchCountReply2, error) {
		return &pb.GetSpaceWorkBenchCountReply2{Error: err}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	data, err := s.staticsUc.GetSpaceWorkbenchCount2(ctx, loginUser.UserId, req.SpaceId, req.ConditionGroups)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &pb.GetSpaceWorkBenchCountReply2{Data: data}, nil
}

func (s *WorkbenchService) GetSpaceWorkObjectCountByIds(ctx context.Context, req *pb.GetSpaceWorkObjectCountRequest) (*pb.GetSpaceWorkObjectCountReply, error) {

	reply := func(err *comm.ErrorInfo) (*pb.GetSpaceWorkObjectCountReply, error) {
		return &pb.GetSpaceWorkObjectCountReply{Result: &pb.GetSpaceWorkObjectCountReply_Error{Error: err}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)

	startTime := date.Parse(req.StartTime).Unix()
	endTime := date.Parse(req.EndTime).Unix()
	if req.StartTime == "" {
		startTime = 0
	}
	if req.EndTime == "" {
		endTime = 0
	}

	data, err := s.staticsUc.GetSpaceWorkObjectCountByIds(ctx, loginUser.UserId, req.SpaceId, req.Ids, startTime, endTime)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &pb.GetSpaceWorkObjectCountReply{Result: &pb.GetSpaceWorkObjectCountReply_Data{Data: data}}, nil
}

func (s *WorkbenchService) GetSpaceUserCount(ctx context.Context, req *pb.GetSpaceUserCountRequest) (*pb.GetSpaceUserCountReply, error) {

	reply := func(err *comm.ErrorInfo) (*pb.GetSpaceUserCountReply, error) {
		return &pb.GetSpaceUserCountReply{Result: &pb.GetSpaceUserCountReply_Error{Error: err}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)

	startTime := date.Parse(req.StartTime).Unix()
	endTime := date.Parse(req.EndTime).Unix()
	if req.StartTime == "" {
		startTime = 0
	}
	if req.EndTime == "" {
		endTime = 0
	}

	data, err := s.staticsUc.GetSpaceUserCount(ctx, loginUser.UserId, req.SpaceId, startTime, endTime)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &pb.GetSpaceUserCountReply{Result: &pb.GetSpaceUserCountReply_Data{Data: data}}, nil
}

func (s *WorkbenchService) GetSpaceVersionCount(ctx context.Context, req *pb.GetSpaceVersionCountRequest) (*pb.GetSpaceVersionCountReply, error) {

	reply := func(err *comm.ErrorInfo) (*pb.GetSpaceVersionCountReply, error) {
		return &pb.GetSpaceVersionCountReply{Result: &pb.GetSpaceVersionCountReply_Error{Error: err}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)

	startTime := date.Parse(req.StartTime).Unix()
	endTime := date.Parse(req.EndTime).Unix()
	if req.StartTime == "" {
		startTime = 0
	}
	if req.EndTime == "" {
		endTime = 0
	}

	data, err := s.staticsUc.GetSpaceVersionCount(ctx, loginUser.UserId, req.SpaceId, startTime, endTime)
	if err != nil {
		return reply(errs.Cast(err))
	}

	return &pb.GetSpaceVersionCountReply{Result: &pb.GetSpaceVersionCountReply_Data{Data: data}}, nil
}
