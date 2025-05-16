package service

import (
	"context"
	pb "go-cs/api/work_item_role/v1"
	uc "go-cs/internal/biz"
	"go-cs/internal/consts"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"slices"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type WorkItemRoleService struct {
	pb.UnimplementedWorkItemRoleServer
	log    *log.Helper
	roleUc *uc.WorkItemRoleUsecase
}

func NewWorkItemRoleService(
	roleUc *uc.WorkItemRoleUsecase,
	logger log.Logger,
) *WorkItemRoleService {
	moduleName := "WorkItemRoleService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &WorkItemRoleService{
		log:    hlog,
		roleUc: roleUc,
	}
}

func (s *WorkItemRoleService) QSpaceRoleList(ctx context.Context, req *pb.SpaceRoleListQueryRequest) (*pb.SpaceRoleListQueryReply, error) {

	reply := func(err error) (*pb.SpaceRoleListQueryReply, error) {
		return &pb.SpaceRoleListQueryReply{Result: &pb.SpaceRoleListQueryReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	uid := utils.GetLoginUser(ctx).UserId

	r, err := s.roleUc.QSpaceWorkItemRoleList(ctx, uid, req)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SpaceRoleListQueryReply{Result: &pb.SpaceRoleListQueryReply_Data{Data: r}}
	return okReply, nil
}

func (s *WorkItemRoleService) SetSpaceWorkItemRoleRanking(ctx context.Context, req *pb.SetSpaceWorkItemRoleRankingRequest) (*pb.SetSpaceWorkItemRoleRankingReply, error) {
	reply := func(err error) (*pb.SetSpaceWorkItemRoleRankingReply, error) {
		return &pb.SetSpaceWorkItemRoleRankingReply{Result: &pb.SetSpaceWorkItemRoleRankingReply_Error{Error: errs.Cast(err)}}, nil
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
	err := s.roleUc.SetSpaceWorkItemRoleRanking(ctx, loginUser, req.SpaceId, rankList)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SetSpaceWorkItemRoleRankingReply{Result: &pb.SetSpaceWorkItemRoleRankingReply_Data{Data: ""}}
	return okReply, nil
}

func (s *WorkItemRoleService) DelSpaceWorkItemRole(ctx context.Context, req *pb.DelSpaceWorkItemRoleRequest) (*pb.DelSpaceWorkItemRoleReply, error) {
	reply := func(err error) (*pb.DelSpaceWorkItemRoleReply, error) {
		return &pb.DelSpaceWorkItemRoleReply{Result: &pb.DelSpaceWorkItemRoleReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.Id <= 0 {
		return reply(errs.Param(ctx, "RoleId"))
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	loginUser := utils.GetLoginUser(ctx)
	err := s.roleUc.DelSpaceWorkItemRole(ctx, loginUser, req.SpaceId, req.Id)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.DelSpaceWorkItemRoleReply{Result: &pb.DelSpaceWorkItemRoleReply_Data{Data: ""}}
	return okReply, nil
}

func (s *WorkItemRoleService) SetSpaceWorkItemRoleName(ctx context.Context, req *pb.SetSpaceWorkItemRoleNameRequest) (*pb.SetSpaceWorkItemRoleNameReply, error) {
	reply := func(err error) (*pb.SetSpaceWorkItemRoleNameReply, error) {
		return &pb.SetSpaceWorkItemRoleNameReply{Result: &pb.SetSpaceWorkItemRoleNameReply_Error{Error: errs.Cast(err)}}, nil
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
	err := s.roleUc.SeSpaceWorkItemRoleName(ctx, loginUser, req.SpaceId, req.Id, req.Name)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.SetSpaceWorkItemRoleNameReply{Result: &pb.SetSpaceWorkItemRoleNameReply_Data{Data: ""}}
	return okReply, nil
}

func (s *WorkItemRoleService) CreateSpaceWorkItemRole(ctx context.Context, req *pb.CreateSpaceWorkItemRoleRequest) (*pb.CreateSpaceWorkItemRoleReply, error) {
	reply := func(err error) (*pb.CreateSpaceWorkItemRoleReply, error) {
		return &pb.CreateSpaceWorkItemRoleReply{Result: &pb.CreateSpaceWorkItemRoleReply_Error{Error: errs.Cast(err)}}, nil
	}

	validate := utils.NewValidator()
	if err := validate.Var(strings.TrimSpace(req.Name), "required,utf8Len=2-8,common_name"); err != nil {
		return reply(errs.Business(ctx, "请输入2 ~ 8个字符，支持中英文、数字"))
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if !slices.Contains(consts.FlowScopeList, consts.FlowScope(req.FlowScope)) {
		return reply(errs.Param(ctx, "FlowScope"))
	}

	loginUser := utils.GetLoginUser(ctx)
	err := s.roleUc.CreateSpaceWorkItemRole(ctx, loginUser, req.SpaceId, req.Name, consts.FlowScope(req.FlowScope))
	if err != nil {
		return reply(err)
	}

	okReply := &pb.CreateSpaceWorkItemRoleReply{Result: &pb.CreateSpaceWorkItemRoleReply_Data{Data: ""}}
	return okReply, nil
}

func (s *WorkItemRoleService) GetWorkItemRelationCount(ctx context.Context, req *pb.GetWorkItemRelationCountRequest) (*pb.GetWorkItemRelationCountReply, error) {
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
	total, err := s.roleUc.QSpaceWorkItemRelationCount(ctx, loginUser, req.SpaceId, req.Id)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.GetWorkItemRelationCountReply{Result: &pb.GetWorkItemRelationCountReply_Data{
		Data: &pb.GetWorkItemRelationCountReplyData{Total: total},
	}}
	return okReply, nil

}

func (s *WorkItemRoleService) GetTemplateRelationCount(ctx context.Context, req *pb.GetTemplateRelationCountRequest) (*pb.GetTemplateRelationCountReply, error) {

	reply := func(err error) (*pb.GetTemplateRelationCountReply, error) {
		return &pb.GetTemplateRelationCountReply{Result: &pb.GetTemplateRelationCountReply_Error{Error: errs.Cast(err)}}, nil
	}

	if req.SpaceId <= 0 {
		return reply(errs.Param(ctx, "SpaceId"))
	}

	if req.Id <= 0 {
		return reply(errs.Param(ctx, "Id"))
	}

	loginUser := utils.GetLoginUser(ctx)
	total, err := s.roleUc.QSpaceTemplateRelationCount(ctx, loginUser, req.SpaceId, req.Id)
	if err != nil {
		return reply(err)
	}

	okReply := &pb.GetTemplateRelationCountReply{Result: &pb.GetTemplateRelationCountReply_Data{
		Data: &pb.GetTemplateRelationCountReplyData{Total: total},
	}}
	return okReply, nil

}
