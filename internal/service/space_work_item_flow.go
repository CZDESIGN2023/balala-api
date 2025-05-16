package service

import (
	"context"
	"go-cs/api/comm"
	pb "go-cs/api/space_work_item_flow/v1"
	"go-cs/internal/biz"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type SpaceWorkItemFlowService struct {
	pb.SpaceWorkItemFlowHTTPServer
	uc  *biz.SpaceWorkItemFlowUsecase
	log *log.Helper
}

func NewSpaceWorkItemFlowService(stu *biz.SpaceWorkItemFlowUsecase, logger log.Logger) *SpaceWorkItemFlowService {
	moduleName := "SpaceWorkItemFlowService"
	_, helper := utils.InitModuleLogger(logger, moduleName)

	return &SpaceWorkItemFlowService{
		uc:  stu,
		log: helper,
	}
}

func (s *SpaceWorkItemFlowService) SetSpaceWorkItemFlowDirector(ctx context.Context, req *pb.SetSpaceWorkItemFlowDirectorRequest) (*pb.SetSpaceWorkItemFlowDirectorReply, error) {

	var reply = func(err *comm.ErrorInfo) (*pb.SetSpaceWorkItemFlowDirectorReply, error) {
		return &pb.SetSpaceWorkItemFlowDirectorReply{Result: &pb.SetSpaceWorkItemFlowDirectorReply_Error{Error: err}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var err error
	validate := utils.NewValidator()

	if err = validate.Var(req.SpaceId, "required,gt=0,number"); err != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if err = validate.Var(req.WorkItemId, "required,gt=0,number"); err != nil {
		return reply(errs.Param(ctx, "WorkFlowId"))
	}
	if req.WorkFlowNodeCode == "" {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if err = validate.Var(req.Director, "required"); err != nil {
		return reply(errs.Param(ctx, "Director"))
	}

	if len(req.Director.Remove) == 0 && len(req.Director.Add) == 0 {
		return reply(errs.Param(ctx, "Director.Remove/Add"))
	}

	err = s.uc.SetWorkFlowDirectors(ctx, loginUser, req.WorkItemId, req.WorkFlowNodeCode, req.Director.Add, req.Director.Remove)
	if err != nil {
		return reply(errs.Cast(err))
	}

	okReply := &pb.SetSpaceWorkItemFlowDirectorReply{Result: &pb.SetSpaceWorkItemFlowDirectorReply_Data{Data: &pb.SetSpaceWorkItemFlowDirectorData{}}}
	return okReply, nil
}

func (s *SpaceWorkItemFlowService) SetSpaceWorkItemFlowPlanTime(ctx context.Context, req *pb.SetSpaceWorkItemFlowPlanTimeRequest) (*pb.SetSpaceWorkItemFlowPlanTimeReply, error) {

	var reply = func(err *comm.ErrorInfo) (*pb.SetSpaceWorkItemFlowPlanTimeReply, error) {
		return &pb.SetSpaceWorkItemFlowPlanTimeReply{Result: &pb.SetSpaceWorkItemFlowPlanTimeReply_Error{Error: err}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var err error
	validate := utils.NewValidator()
	if err = validate.Var(req.SpaceId, "required,gt=0,number"); err != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if err = validate.Var(req.WorkItemId, "required,gt=0,number"); err != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if req.WorkFlowNodeCode == "" {
		return reply(errs.Param(ctx, "WorkFlowNodeCode"))
	}

	var planStartAt, planCompleteAt int64

	if req.PlanTimeAt != nil && req.PlanTimeAt.Start != "" {
		if planStartTime, err := time.ParseInLocation("2006/01/02 15:04:05", req.PlanTimeAt.Start, time.Local); err == nil {
			planStartAt = planStartTime.Unix()
		}

	}

	if req.PlanTimeAt != nil && req.PlanTimeAt.Complete != "" {
		if planCompleteTime, err := time.ParseInLocation("2006/01/02 15:04:05", req.PlanTimeAt.Complete, time.Local); err == nil {
			planCompleteAt = planCompleteTime.Unix()
		}
	}

	err = s.uc.SetWorkFlowPlanTime(ctx, loginUser, req.WorkItemId, req.WorkFlowNodeCode, planStartAt, planCompleteAt)
	if err != nil {
		return reply(errs.Cast(err))
	}

	okReply := &pb.SetSpaceWorkItemFlowPlanTimeReply{Result: &pb.SetSpaceWorkItemFlowPlanTimeReply_Data{Data: ""}}
	return okReply, nil
}

func (s *SpaceWorkItemFlowService) UpgradeWorkItemFlow(ctx context.Context, req *pb.UpgradeTaskWorkFlowRequest) (*pb.UpgradeTaskWorkFlowReply, error) {

	var reply = func(err error) (*pb.UpgradeTaskWorkFlowReply, error) {
		return &pb.UpgradeTaskWorkFlowReply{Result: &pb.UpgradeTaskWorkFlowReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var err error
	validate := utils.NewValidator()
	if err = validate.Var(req.SpaceId, "required,gt=0,number"); err != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if err = validate.Var(req.WorkItemId, "required,gt=0,number"); err != nil {
		return reply(errs.Param(ctx, "WorkItemId"))
	}

	if len(req.RoleDirectors) == 0 {
		return reply(errs.Param(ctx, "RoleDirectors"))
	}

	err = s.uc.UpgradeWorkItemFlow(ctx, loginUser, req)

	result := &pb.UpgradeTaskWorkFlowReplyData{
		Result: make([]*pb.UpgradeTaskWorkFlowReplyData_Result, 0),
	}

	if err != nil {
		result.Result = append(result.Result, &pb.UpgradeTaskWorkFlowReplyData_Result{
			WorkItemId: req.WorkItemId,
			Code:       0,
			Message:    err.Error(),
		})
	} else {
		result.Result = append(result.Result, &pb.UpgradeTaskWorkFlowReplyData_Result{
			WorkItemId: req.WorkItemId,
			Code:       200,
		})
	}

	okReply := &pb.UpgradeTaskWorkFlowReply{Result: &pb.UpgradeTaskWorkFlowReply_Data{Data: result}}
	return okReply, nil
}

func (s *SpaceWorkItemFlowService) BatUpgradeWorkItemFlowPrepare(ctx context.Context, req *pb.BatUpgradeWorkItemFlowPrepareRequest) (*pb.BatUpgradeWorkItemFlowPrepareReply, error) {

	var reply = func(err error) (*pb.BatUpgradeWorkItemFlowPrepareReply, error) {
		return &pb.BatUpgradeWorkItemFlowPrepareReply{Result: &pb.BatUpgradeWorkItemFlowPrepareReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var err error
	validate := utils.NewValidator()
	if err = validate.Var(req.SpaceId, "required,gt=0,number"); err != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if err = validate.Var(req.FlowId, "required,gt=0,number"); err != nil {
		return reply(errs.Param(ctx, "FlowId"))
	}

	ids, err := s.uc.BatUpgradeWorkItemFlowPrepare(ctx, loginUser, req.SpaceId, req.FlowId)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.BatUpgradeWorkItemFlowPrepareReply{Result: &pb.BatUpgradeWorkItemFlowPrepareReply_Data{
		Data: &pb.BatUpgradeWorkItemFlowPrepareReplyData{
			Ids: ids,
		},
	}}
	return okReply, nil
}

func (s *SpaceWorkItemFlowService) BatchUpgradeWorkItemFlow(ctx context.Context, req *pb.BatchUpgradeTaskWorkFlowRequest) (*pb.BatchUpgradeTaskWorkFlowReply, error) {

	var reply = func(err error) (*pb.BatchUpgradeTaskWorkFlowReply, error) {
		return &pb.BatchUpgradeTaskWorkFlowReply{Result: &pb.BatchUpgradeTaskWorkFlowReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser, _ := utils.GetLoginUserInfo(ctx)
	if loginUser == nil {
		return reply(errs.NotLogin(ctx))
	}

	var err error
	validate := utils.NewValidator()
	if err = validate.Var(req.SpaceId, "required,gt=0,number"); err != nil {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if err = validate.Var(req.FlowId, "required,gt=0,number"); err != nil {
		return reply(errs.Param(ctx, "FlowId"))
	}

	if len(req.RoleDirectors) == 0 {
		return reply(errs.Param(ctx, "RoleDirectors"))
	}

	result, err := s.uc.BatchUpgradeWorkItemFlow(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.BatchUpgradeTaskWorkFlowReply{Result: &pb.BatchUpgradeTaskWorkFlowReply_Data{Data: result}}
	return okReply, nil
}
