package v1_4_7_0

import (
	"context"
	"go-cs/internal/consts"
	"go-cs/internal/domain/space/repo"
	"go-cs/internal/domain/work_flow"
	work_flow_facade "go-cs/internal/domain/work_flow/facade"
	flow_repo "go-cs/internal/domain/work_flow/repo"
	wf_service "go-cs/internal/domain/work_flow/service"
	role_repo "go-cs/internal/domain/work_item_role/repo"
	work_item_role "go-cs/internal/domain/work_item_role/service"
	status_repo "go-cs/internal/domain/work_item_status/repo"
	service2 "go-cs/internal/domain/work_item_status/service"
	repo2 "go-cs/internal/domain/work_item_type/repo"
	"go-cs/internal/domain/work_item_type/service"
	"go-cs/internal/pkg/trans"
	"gorm.io/gorm"
	"time"
)

type DML struct {
	db               *gorm.DB
	spaceRepo        repo.SpaceRepo
	statusRepo       status_repo.WorkItemStatusRepo
	roleRepo         role_repo.WorkItemRoleRepo
	workItemTypeRepo repo2.WorkItemTypeRepo
	flowRepo         flow_repo.WorkFlowRepo

	workItemTypeService   *service.WorkItemTypeService
	roleService           *work_item_role.WorkItemRoleService
	statusService         *service2.WorkItemStatusService
	workFlowDomainService *wf_service.WorkFlowService
	tm                    trans.Transaction
}

func NewDML(
	spaceRepo repo.SpaceRepo,
	statusRepo status_repo.WorkItemStatusRepo,
	roleRepo role_repo.WorkItemRoleRepo,
	workItemTypeRepo repo2.WorkItemTypeRepo,
	flowRepo flow_repo.WorkFlowRepo,

	workItemTypeService *service.WorkItemTypeService,
	roleService *work_item_role.WorkItemRoleService,
	statusService *service2.WorkItemStatusService,
	workFlowService *wf_service.WorkFlowService,
	tm trans.Transaction,
) *DML {
	return &DML{
		spaceRepo:        spaceRepo,
		statusRepo:       statusRepo,
		roleRepo:         roleRepo,
		workItemTypeRepo: workItemTypeRepo,
		flowRepo:         flowRepo,

		workItemTypeService:   workItemTypeService,
		roleService:           roleService,
		statusService:         statusService,
		workFlowDomainService: workFlowService,
		tm:                    tm,
	}
}

func (d *DML) HandleData() error {
	ctx := context.Background()
	// 获取所有spaceId
	spaceIds, err := d.spaceRepo.GetAllSpaceIds()
	if err != nil {
		return err
	}

	for _, spaceId := range spaceIds {
		// 创建状态任务类型
		exist, err := d.workItemTypeRepo.IsExistByKey(ctx, spaceId, "state_task")
		if err != nil {
			return err
		}
		if exist {
			continue
		}

		workItemType, err := d.workItemTypeService.CreateStateFlowTaskWorkType(ctx, spaceId, 0)
		if err != nil {
			return err
		}

		// 创建状态
		statusList := d.statusService.CreateSpaceDefaultStatusInfo(ctx, spaceId, []consts.FlowScope{consts.FlowScope_Stateflow}, 0)

		// 创建角色
		roleList := d.roleService.CreateSpaceDefaultRoles(ctx, work_item_role.CreateSpaceDefaultRolesReq{
			SpaceId:              spaceId,
			StateFlowWitemTypeId: workItemType.Id,
		})

		// 创建issue状态流程
		issueStateFlow := d.workFlowDomainService.NewIssueStateFlow(ctx, &wf_service.GenerateWorkFlowReq{
			SpaceId:            spaceId,
			WorkItemTypeId:     workItemType.Id,
			Ranking:            time.Now().Unix(), //放在禁用第一位
			WorkItemStatusInfo: work_flow_facade.BuildWorkItemStatusInfo(statusList),
			WorkItemRoleInfo:   work_flow_facade.BuildWorkItemRoleInfo(roleList),
		})
		issueStateFlow.WorkFlow.Status = work_flow.WorkFlowStatus_Disable

		err = d.tm.InTx(ctx, func(ctx context.Context) error {
			err := d.workItemTypeRepo.CreateWorkItemType(ctx, workItemType)
			if err != nil {
				return err
			}

			err = d.statusRepo.CreateWorkItemStatusItems(ctx, spaceId, statusList.Items)
			if err != nil {
				return err
			}

			err = d.roleRepo.CreateWorkItemRoles(ctx, roleList)
			if err != nil {
				return err
			}

			//保存工作流信息
			err = d.flowRepo.CreateWorkFlow(ctx, issueStateFlow.WorkFlow)
			if err != nil {
				return err
			}

			//保存工作流模版信息
			err = d.flowRepo.CreateWorkFlowTemplate(ctx, issueStateFlow.WorkFlowTemplate)
			if err != nil {
				return err
			}
			return nil
		})
	}

	return err
}
