package service

import (
	"context"
	pb "go-cs/api/work_item_status/v1"
	uc "go-cs/internal/biz"
	"go-cs/internal/consts"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"slices"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type WorkItemStatusService struct {
	pb.UnimplementedWorkItemStatusServer
	log *log.Helper
	uc  *uc.WorkItemStatusUsecase
}

func NewWorkItemStatusService(
	uc *uc.WorkItemStatusUsecase,
	logger log.Logger,
) *WorkItemStatusService {
	moduleName := "WorkItemStatusService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &WorkItemStatusService{
		log: hlog,
		uc:  uc,
	}
}

func (s *WorkItemStatusService) QSpaceWorkItemStatusList(ctx context.Context, req *pb.SpaceWorkItemStatusListRequest) (*pb.SpaceWorkItemStatusListReply, error) {

	reply := func(err error) (*pb.SpaceWorkItemStatusListReply, error) {
		return &pb.SpaceWorkItemStatusListReply{Result: &pb.SpaceWorkItemStatusListReply_Error{Error: errs.Cast(err)}}, nil
	}

	uid := utils.GetLoginUser(ctx).UserId

	result, err := s.uc.QSpaceWorkItemStatusList(ctx, uid, req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SpaceWorkItemStatusListReply{Result: &pb.SpaceWorkItemStatusListReply_Data{Data: result}}
	return okReply, nil
}

func (s *WorkItemStatusService) QSpaceWorkItemStatusById(ctx context.Context, req *pb.QSpaceWorkItemStatusByIdRequest) (*pb.QSpaceWorkItemStatusByIdReply, error) {

	reply := func(err error) (*pb.QSpaceWorkItemStatusByIdReply, error) {
		return &pb.QSpaceWorkItemStatusByIdReply{Result: &pb.QSpaceWorkItemStatusByIdReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if len(req.Ids) == 0 {
		return reply(errs.Param(ctx, "Ids"))
	}

	result, err := s.uc.QSpaceWorkItemStatusById(ctx, utils.GetLoginUser(ctx), req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.QSpaceWorkItemStatusByIdReply{Result: &pb.QSpaceWorkItemStatusByIdReply_Data{Data: result}}
	return okReply, nil
}

func (s *WorkItemStatusService) SetSpaceWorkItemStatusRanking(ctx context.Context, req *pb.SetSpaceWorkItemStatusRankingRequest) (*pb.SetSpaceWorkItemStatusRankingReply, error) {
	reply := func(err error) (*pb.SetSpaceWorkItemStatusRankingReply, error) {
		return &pb.SetSpaceWorkItemStatusRankingReply{Result: &pb.SetSpaceWorkItemStatusRankingReply_Error{Error: errs.Cast(err)}}, nil
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
	err := s.uc.SetSpaceWorkItemStatusRanking(ctx, loginUser, req.SpaceId, rankList)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SetSpaceWorkItemStatusRankingReply{Result: &pb.SetSpaceWorkItemStatusRankingReply_Data{Data: ""}}
	return okReply, nil
}

func (s *WorkItemStatusService) DelSpaceWorkItemStatus(ctx context.Context, req *pb.DelSpaceWorkItemStatusRequest) (*pb.DelSpaceWorkItemStatusReply, error) {
	reply := func(err error) (*pb.DelSpaceWorkItemStatusReply, error) {
		return &pb.DelSpaceWorkItemStatusReply{Result: &pb.DelSpaceWorkItemStatusReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.Id <= 0 {
		return reply(errs.Param(ctx, "Id"))
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	loginUser := utils.GetLoginUser(ctx)
	err := s.uc.DelSpaceWorkItemStatus(ctx, loginUser, req.SpaceId, req.Id)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.DelSpaceWorkItemStatusReply{Result: &pb.DelSpaceWorkItemStatusReply_Data{Data: ""}}
	return okReply, nil
}

func (s *WorkItemStatusService) SetSpaceWorkItemStatusName(ctx context.Context, req *pb.SetSpaceWorkItemStatusNameRequest) (*pb.SetSpaceWorkItemStatusNameReply, error) {
	reply := func(err error) (*pb.SetSpaceWorkItemStatusNameReply, error) {
		return &pb.SetSpaceWorkItemStatusNameReply{Result: &pb.SetSpaceWorkItemStatusNameReply_Error{Error: errs.Cast(err)}}, nil
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(strings.TrimSpace(req.Name), "required,utf8Len=2-8,common_name"); vaildErr != nil {
		return reply(errs.Business(ctx, "请输入2 ~ 8个字符，支持中英文、数字"))
	}

	if req.Id <= 0 {
		return reply(errs.Param(ctx, "Id"))
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	loginUser := utils.GetLoginUser(ctx)
	err := s.uc.SetSpaceWorkItemStatusName(ctx, loginUser, req.SpaceId, req.Id, req.Name)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SetSpaceWorkItemStatusNameReply{Result: &pb.SetSpaceWorkItemStatusNameReply_Data{Data: ""}}
	return okReply, nil
}

func (s *WorkItemStatusService) CreateSpaceWorkItemStatus(ctx context.Context, req *pb.CreateSpaceWorkItemStatusRequest) (*pb.CreateSpaceWorkItemStatusReply, error) {
	reply := func(err error) (*pb.CreateSpaceWorkItemStatusReply, error) {
		return &pb.CreateSpaceWorkItemStatusReply{Result: &pb.CreateSpaceWorkItemStatusReply_Error{Error: errs.Cast(err)}}, nil
	}

	var vaildErr error
	validate := utils.NewValidator()
	if vaildErr = validate.Var(strings.TrimSpace(req.Name), "required,utf8Len=2-8,common_name"); vaildErr != nil {
		return reply(errs.Business(ctx, "请输入2 ~ 8个字符，支持中英文、数字"))
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if !slices.Contains(consts.FlowScopeList, consts.FlowScope(req.FlowScope)) {
		return reply(errs.Param(ctx, "FlowScope"))
	}

	if req.StatusType != int64(consts.WorkItemStatusType_Process) &&
		req.StatusType != int64(consts.WorkItemStatusType_Archived) {
		return reply(errs.Param(ctx, "StatusType"))
	}

	loginUser := utils.GetLoginUser(ctx)
	err := s.uc.CreateSpaceWorkItemStatus(ctx, loginUser, req.SpaceId, req.Name, consts.FlowScope(req.FlowScope), consts.WorkItemStatusType(req.StatusType))
	if err != nil {
		return reply(err)
	}

	okReply := &pb.CreateSpaceWorkItemStatusReply{Result: &pb.CreateSpaceWorkItemStatusReply_Data{Data: ""}}
	return okReply, nil
}

func (s *WorkItemStatusService) GetWorkItemRelationCount(ctx context.Context, req *pb.GetWorkItemRelationCountRequest) (*pb.GetWorkItemRelationCountReply, error) {
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
	total, err := s.uc.QSpaceWorkItemRelationCount(ctx, loginUser, req.SpaceId, req.Id)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.GetWorkItemRelationCountReply{Result: &pb.GetWorkItemRelationCountReply_Data{
		Data: &pb.GetWorkItemRelationCountReplyData{Total: total},
	}}
	return okReply, nil

}

func (s *WorkItemStatusService) GetTemplateRelationCount(ctx context.Context, req *pb.GetTemplateRelationCountRequest) (*pb.GetTemplateRelationCountReply, error) {
	reply := func(err error) (*pb.GetTemplateRelationCountReply, error) {
		return &pb.GetTemplateRelationCountReply{Result: &pb.GetTemplateRelationCountReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.Id <= 0 {
		return reply(errs.Param(ctx, "Id"))
	}

	total, err := s.uc.QSpaceTemplateRelationCount(ctx, req.SpaceId, req.Id)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.GetTemplateRelationCountReply{Result: &pb.GetTemplateRelationCountReply_Data{
		Data: &pb.GetTemplateRelationCountReplyData{Total: total},
	}}
	return okReply, nil

}
