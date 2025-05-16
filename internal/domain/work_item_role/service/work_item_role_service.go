package service

import (
	"context"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/work_item_role"
	"go-cs/internal/domain/work_item_role/repo"
	"go-cs/internal/pkg/biz_id"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/rand"

	"github.com/go-kratos/kratos/v2/log"
)

type WorkItemRoleService struct {
	log       *log.Helper
	idService *biz_id.BusinessIdService
	repo      repo.WorkItemRoleRepo
}

func NewWorkItemRoleService(
	idService *biz_id.BusinessIdService,
	repo repo.WorkItemRoleRepo,
	logger log.Logger,
) *WorkItemRoleService {

	moduleName := "WorkItemRoleService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &WorkItemRoleService{
		repo:      repo,
		log:       hlog,
		idService: idService,
	}
}

func (s *WorkItemRoleService) maxRanking(ctx context.Context, spaceId int64) (int64, error) {

	maxRanking, err := s.repo.GetMaxRanking(ctx, spaceId)
	if err != nil {
		return 0, err
	}

	return maxRanking + 100, nil
}

type CreateSpaceDefaultRolesReq struct {
	SpaceId              int64
	OperUid              int64
	WorkFlowWitemTypeId  int64 //流程模式的工作项类型id
	StateFlowWitemTypeId int64 //状态模式的工作项类型id
}

// 只用流程模式的工作流
func (s *WorkItemRoleService) CreateSpaceDefaultRoles(ctx context.Context, req CreateSpaceDefaultRolesReq) []*domain.WorkItemRole {

	spaceId := req.SpaceId
	uid := req.OperUid
	workItemTypeId := req.WorkFlowWitemTypeId
	stateFlowWitemTypeId := req.StateFlowWitemTypeId

	roles := make(domain.WorkItemRoles, 0)

	if workItemTypeId != 0 {
		// 流程模式预设角色
		evaluator, _ := s.newWorkItemRole(ctx, spaceId, workItemTypeId, "评审组", consts.WorkflowOwnerRole_Evaluator, 700, 1, uid, consts.FlowScope_Workflow, nil)
		producer, _ := s.newWorkItemRole(ctx, spaceId, workItemTypeId, "策划", consts.WorkflowOwnerRole_Producer, 600, 1, uid, consts.FlowScope_Workflow, nil)
		uiDesigner, _ := s.newWorkItemRole(ctx, spaceId, workItemTypeId, "设计", consts.WorkflowOwnerRole_UIDesigner, 500, 1, uid, consts.FlowScope_Workflow, nil)
		acceptor, _ := s.newWorkItemRole(ctx, spaceId, workItemTypeId, "验收", consts.WorkflowOwnerRole_Acceptor, 400, 1, uid, consts.FlowScope_Workflow, nil)
		reviewer, _ := s.newWorkItemRole(ctx, spaceId, workItemTypeId, "审查", consts.WorkflowOwnerRole_Reviewer, 300, 1, uid, consts.FlowScope_Workflow, nil)
		developer, _ := s.newWorkItemRole(ctx, spaceId, workItemTypeId, "开发", consts.WorkflowOwnerRole_Developer, 200, 1, uid, consts.FlowScope_Workflow, nil)
		qa, _ := s.newWorkItemRole(ctx, spaceId, workItemTypeId, "测试", consts.WorkflowOwnerRole_Qa, 100, 1, uid, consts.FlowScope_Workflow, nil)

		roles = append(roles, evaluator, producer, uiDesigner, acceptor, reviewer, developer, qa)
	}

	if stateFlowWitemTypeId != 0 {
		// 状态模式预设角色
		stReviewer, _ := s.newWorkItemRole(ctx, spaceId, stateFlowWitemTypeId, "审核人", consts.StateflowOwnerRole_Reviewer, 1000, 1, uid, consts.FlowScope_Stateflow, nil)
		stOperator, _ := s.newWorkItemRole(ctx, spaceId, stateFlowWitemTypeId, "经办人", consts.StateflowOwnerRole_Operator, 900, 1, uid, consts.FlowScope_Stateflow, nil)
		stReporter, _ := s.newWorkItemRole(ctx, spaceId, stateFlowWitemTypeId, "报告人", consts.StateflowOwnerRole_Reporter, 800, 1, uid, consts.FlowScope_Stateflow, nil)

		roles = append(roles, stOperator, stReporter, stReviewer)
	}

	return roles
}

type CreateSpaceRoleRequest struct {
	SpaceId        int64
	WorkItemTypeId int64
	Name           string
	Key            string
	Ranking        int64
	Uid            int64
	IsSys          int32
	FlowScope      consts.FlowScope
}

func (s *WorkItemRoleService) CreateSpaceRole(ctx context.Context, req CreateSpaceRoleRequest, oper shared.Oper) (*domain.WorkItemRole, error) {
	// 创建空间角色
	itemKey := "role_" + rand.Letters(5)
	if req.Key != "" {
		itemKey = req.Key
	}

	return s.newWorkItemRole(ctx, req.SpaceId, req.WorkItemTypeId, req.Name, itemKey, req.Ranking, req.IsSys, req.Uid, req.FlowScope, oper)
}

func (s *WorkItemRoleService) newWorkItemRole(ctx context.Context, spaceId int64, workItemTypeId int64, name string, key string, ranking int64, isSys int32, uid int64, flowScope consts.FlowScope, oper shared.Oper) (*domain.WorkItemRole, error) {
	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkItemRole)
	if bizId == nil {
		return nil, errs.Business(ctx, "分配Id失败")
	}

	if key == "" {
		key = "role_" + rand.S(5)
	}

	if ranking == 0 {
		maxRanking, err := s.maxRanking(ctx, spaceId)
		if err != nil {
			return nil, errs.Internal(ctx, err)
		}

		ranking = maxRanking
	}

	role := domain.NewWorkItemRole(bizId.Id, spaceId, workItemTypeId, name, key, ranking, isSys, uid, flowScope, oper)
	return role, nil
}

func (s *WorkItemRoleService) CheckAndUpdateSpaceWorkItemRoleName(ctx context.Context, wItemRole *domain.WorkItemRole, newName string, oper shared.Oper) error {
	err := wItemRole.ChangeName(newName, oper)
	if err != nil {
		return err
	}

	return nil
}
