package service

import (
	"context"
	pb "go-cs/api/work_flow/v1"
	"go-cs/internal/biz"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type WorkFlowService struct {
	pb.UnimplementedWorkFlowServer
	log *log.Helper
	uc  *biz.WorkFlowUsecase
}

func NewWorkFlowService(
	logger log.Logger,
	uc *biz.WorkFlowUsecase,
) *WorkFlowService {
	moduleName := "WorkFlowService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &WorkFlowService{
		log: hlog,
		uc:  uc,
	}
}

func (s *WorkFlowService) SpaceWorkFlowList(ctx context.Context, req *pb.SpaceWorkFlowListRequest) (*pb.SpaceWorkFlowListReply, error) {

	reply := func(err error) (*pb.SpaceWorkFlowListReply, error) {
		return &pb.SpaceWorkFlowListReply{Result: &pb.SpaceWorkFlowListReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	loginUser := utils.GetLoginUser(ctx)
	result, err := s.uc.QSpaceWorkFlowList(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SpaceWorkFlowListReply{Result: &pb.SpaceWorkFlowListReply_Data{Data: result}}
	return okReply, nil
}

func (s *WorkFlowService) SpaceWorkFlowById(ctx context.Context, req *pb.SpaceWorkFlowByIdRequest) (*pb.SpaceWorkFlowByIdReply, error) {

	reply := func(err error) (*pb.SpaceWorkFlowByIdReply, error) {
		return &pb.SpaceWorkFlowByIdReply{Result: &pb.SpaceWorkFlowByIdReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if len(req.Ids) == 0 {
		return reply(errs.Param(ctx, "Ids"))
	}

	loginUser := utils.GetLoginUser(ctx)
	result, err := s.uc.QSpaceWorkFlowById(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SpaceWorkFlowByIdReply{Result: &pb.SpaceWorkFlowByIdReply_Data{Data: result}}
	return okReply, nil
}

func (s *WorkFlowService) SpaceWorkFlowPageList(ctx context.Context, req *pb.SpaceWorkFlowPageListRequest) (*pb.SpaceWorkFlowPageListReply, error) {

	reply := func(err error) (*pb.SpaceWorkFlowPageListReply, error) {
		return &pb.SpaceWorkFlowPageListReply{Result: &pb.SpaceWorkFlowPageListReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	loginUser := utils.GetLoginUser(ctx)
	result, err := s.uc.QSpaceWorkFlowPageList(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SpaceWorkFlowPageListReply{Result: &pb.SpaceWorkFlowPageListReply_Data{Data: result}}
	return okReply, nil
}

func (s *WorkFlowService) SetWorkFlowRanking(ctx context.Context, req *pb.SetWorkFlowRankingRequest) (*pb.SetWorkFlowRankingReply, error) {
	reply := func(err error) (*pb.SetWorkFlowRankingReply, error) {
		return &pb.SetWorkFlowRankingReply{Result: &pb.SetWorkFlowRankingReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if len(req.List) == 0 {
		return reply(errs.Param(ctx, "List"))
	}

	rankList := make([]map[string]int64, 0)
	for _, v := range req.List {
		rankList = append(rankList, map[string]int64{
			"id":      v.Id,
			"ranking": v.Ranking,
		})
	}

	loginUser := utils.GetLoginUser(ctx)
	err := s.uc.SetWorkFlowRanking(ctx, loginUser, req.SpaceId, rankList)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SetWorkFlowRankingReply{Result: &pb.SetWorkFlowRankingReply_Data{Data: ""}}
	return okReply, nil
}

func (s *WorkFlowService) SaveWorkFlowTemplateConfig(ctx context.Context, req *pb.SaveWorkFlowTemplateConfigRequest) (*pb.SaveWorkFlowTemplateConfigReply, error) {
	reply := func(err error) (*pb.SaveWorkFlowTemplateConfigReply, error) {
		return &pb.SaveWorkFlowTemplateConfigReply{Result: &pb.SaveWorkFlowTemplateConfigReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.FlowId <= 0 {
		return reply(errs.Param(ctx, "FlowId"))
	}

	loginUser := utils.GetLoginUser(ctx)
	err := s.uc.SaveWorkFlowTemplateConfig(ctx, loginUser, req.SpaceId, req.FlowId, req.Config)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SaveWorkFlowTemplateConfigReply{Result: &pb.SaveWorkFlowTemplateConfigReply_Data{Data: ""}}
	return okReply, nil
}

func (s *WorkFlowService) CreateWorkFlow(ctx context.Context, req *pb.CreateWorkFlowRequest) (*pb.CreateWorkFlowReply, error) {

	reply := func(err error) (*pb.CreateWorkFlowReply, error) {
		return &pb.CreateWorkFlowReply{Result: &pb.CreateWorkFlowReply_Error{Error: errs.Cast(err)}}, nil
	}

	var vaildErr error
	validate := utils.NewValidator()
	if req.Name != "" {
		if vaildErr = validate.Var(strings.TrimSpace(req.Name), "required,utf8Len=2-30"); vaildErr != nil {
			return reply(errs.Business(ctx, "请输入2 ~ 30个字符"))
		}
	}

	if vaildErr = validate.Var(strings.TrimSpace(req.FlowMode), "required"); vaildErr != nil {
		return reply(errs.Business(ctx, "FlowMode"))
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	loginUser := utils.GetLoginUser(ctx)
	flowId, err := s.uc.CreateWorkFlow(ctx, loginUser, req.SpaceId, req.Name, int(req.Status), req.FlowMode)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.CreateWorkFlowReply{Result: &pb.CreateWorkFlowReply_Data{
		Data: &pb.CreateWorkFlowReplyData{
			FlowId: flowId,
		},
	}}
	return okReply, nil
}

func (s *WorkFlowService) GetWorkFlow(ctx context.Context, req *pb.GetWorkFlowRequest) (*pb.GetWorkFlowReply, error) {
	reply := func(err error) (*pb.GetWorkFlowReply, error) {
		return &pb.GetWorkFlowReply{Result: &pb.GetWorkFlowReply_Error{Error: errs.Cast(err)}}, nil
	}

	loginUser := utils.GetLoginUser(ctx)
	data, err := s.uc.QWorkFlow(ctx, loginUser, req.SpaceId, req.FlowId)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.GetWorkFlowReply{Result: &pb.GetWorkFlowReply_Data{
		Data: data,
	}}
	return okReply, nil
}
func (s *WorkFlowService) GetWorkFlowTemplate(ctx context.Context, req *pb.GetWorkFlowTemplateRequest) (*pb.GetWorkFlowTemplateReply, error) {
	reply := func(err error) (*pb.GetWorkFlowTemplateReply, error) {
		return &pb.GetWorkFlowTemplateReply{Result: &pb.GetWorkFlowTemplateReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.TemplateId <= 0 {
		return reply(errs.Param(ctx, "TemplateId"))
	}

	loginUser := utils.GetLoginUser(ctx)
	data, err := s.uc.QWorkFlowTemplate(ctx, loginUser, req.SpaceId, req.TemplateId)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.GetWorkFlowTemplateReply{Result: &pb.GetWorkFlowTemplateReply_Data{
		Data: data,
	}}
	return okReply, nil
}

func (s *WorkFlowService) DelWorkFlow(ctx context.Context, req *pb.DelWorkFlowRequest) (*pb.DelWorkFlowReply, error) {
	reply := func(err error) (*pb.DelWorkFlowReply, error) {
		return &pb.DelWorkFlowReply{Result: &pb.DelWorkFlowReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.FlowId <= 0 {
		return reply(errs.Param(ctx, "FlowId"))
	}

	loginUser := utils.GetLoginUser(ctx)
	err := s.uc.DelWorkFlow(ctx, loginUser, req.SpaceId, req.FlowId)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.DelWorkFlowReply{Result: &pb.DelWorkFlowReply_Data{
		Data: "",
	}}
	return okReply, nil
}

func (s *WorkFlowService) SetWorkFlowStatus(ctx context.Context, req *pb.SetWorkFlowStatusRequest) (*pb.SetWorkFlowStatusReply, error) {
	reply := func(err error) (*pb.SetWorkFlowStatusReply, error) {
		return &pb.SetWorkFlowStatusReply{Result: &pb.SetWorkFlowStatusReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.FlowId <= 0 {
		return reply(errs.Param(ctx, "FlowId"))
	}

	if req.Status < -1 || req.Status > 1 {
		return reply(errs.Param(ctx, "Status"))
	}

	loginUser := utils.GetLoginUser(ctx)
	err := s.uc.SetWorkFlowStatus(ctx, loginUser, req.SpaceId, req.FlowId, int64(req.Status))
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SetWorkFlowStatusReply{Result: &pb.SetWorkFlowStatusReply_Data{
		Data: "",
	}}
	return okReply, nil
}

func (s *WorkFlowService) CopyWorkFlow(ctx context.Context, req *pb.CopyWorkFlowRequest) (*pb.CopyWorkFlowReply, error) {
	reply := func(err error) (*pb.CopyWorkFlowReply, error) {
		return &pb.CopyWorkFlowReply{Result: &pb.CopyWorkFlowReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.FlowId <= 0 {
		return reply(errs.Param(ctx, "FlowId"))
	}

	if req.Status < 0 || req.Status > 1 {
		return reply(errs.Param(ctx, "Status"))
	}

	loginUser := utils.GetLoginUser(ctx)
	newFlowId, err := s.uc.CopyWorkFlow(ctx, loginUser, req.SpaceId, req.FlowId, req.Name, req.Status)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.CopyWorkFlowReply{Result: &pb.CopyWorkFlowReply_Data{
		Data: &pb.CopyWorkFlowReplyData{FlowId: newFlowId},
	}}
	return okReply, nil
}

func (s *WorkFlowService) SetWorkFlowName(ctx context.Context, req *pb.SetWorkFlowNameRequest) (*pb.SetWorkFlowNameReply, error) {
	reply := func(err error) (*pb.SetWorkFlowNameReply, error) {
		return &pb.SetWorkFlowNameReply{Result: &pb.SetWorkFlowNameReply_Error{Error: errs.Cast(err)}}, nil
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(strings.TrimSpace(req.Name), "required,utf8Len=2-30"); vaildErr != nil {
		return reply(errs.Business(ctx, "请输入2 ~ 30个字符"))
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.FlowId <= 0 {
		return reply(errs.Param(ctx, "FlowId"))
	}

	loginUser := utils.GetLoginUser(ctx)
	err := s.uc.SetWorkFlowName(ctx, loginUser, req.SpaceId, req.FlowId, req.Name)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SetWorkFlowNameReply{Result: &pb.SetWorkFlowNameReply_Data{
		Data: "",
	}}
	return okReply, nil
}

func (s *WorkFlowService) GetWorkItemRelationCount(ctx context.Context, req *pb.GetWorkItemRelationCountRequest) (*pb.GetWorkItemRelationCountReply, error) {
	reply := func(err error) (*pb.GetWorkItemRelationCountReply, error) {
		return &pb.GetWorkItemRelationCountReply{Result: &pb.GetWorkItemRelationCountReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.Id <= 0 {
		return reply(errs.Param(ctx, "Id"))
	}

	loginUser := utils.GetLoginUser(ctx)
	total, err := s.uc.QSpaceWorkItemRelationCount(ctx, loginUser, req.SpaceId, req.Id, req.Scene)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.GetWorkItemRelationCountReply{Result: &pb.GetWorkItemRelationCountReply_Data{
		Data: &pb.GetWorkItemRelationCountReplyData{Total: total},
	}}
	return okReply, nil

}

func (s *WorkFlowService) GetOwnerRuleRelationTemplate(ctx context.Context, req *pb.GetOwnerRuleRelationTemplateRequest) (*pb.GetOwnerRuleRelationTemplateReply, error) {
	reply := func(err error) (*pb.GetOwnerRuleRelationTemplateReply, error) {
		return &pb.GetOwnerRuleRelationTemplateReply{Result: &pb.GetOwnerRuleRelationTemplateReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.OwnerUid <= 0 {
		return reply(errs.Param(ctx, "OwnerUid"))
	}

	loginUser := utils.GetLoginUser(ctx)
	resultData, err := s.uc.GetOwnerRuleRelationTemplateRequest(ctx, loginUser, req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.GetOwnerRuleRelationTemplateReply{Result: &pb.GetOwnerRuleRelationTemplateReply_Data{
		Data: resultData,
	}}
	return okReply, nil
}
