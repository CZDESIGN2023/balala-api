package v1_4_9_0

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-cs/internal/consts"
	"go-cs/internal/domain/space/repo"
	member "go-cs/internal/domain/space_member"
	space_member_repo "go-cs/internal/domain/space_member/repo"
	space_view_repo "go-cs/internal/domain/space_view/repo"
	view_service "go-cs/internal/domain/space_view/service"
	flow_repo "go-cs/internal/domain/work_flow/repo"
	wf_service "go-cs/internal/domain/work_flow/service"
	role_repo "go-cs/internal/domain/work_item_role/repo"
	work_item_role "go-cs/internal/domain/work_item_role/service"
	status_repo "go-cs/internal/domain/work_item_status/repo"
	service2 "go-cs/internal/domain/work_item_status/service"
	repo2 "go-cs/internal/domain/work_item_type/repo"
	"go-cs/internal/domain/work_item_type/service"
	"go-cs/internal/pkg/trans"
	"go-cs/pkg/stream"
	"gorm.io/gorm"
)

type DML struct {
	db               *gorm.DB
	spaceRepo        repo.SpaceRepo
	statusRepo       status_repo.WorkItemStatusRepo
	roleRepo         role_repo.WorkItemRoleRepo
	workItemTypeRepo repo2.WorkItemTypeRepo
	flowRepo         flow_repo.WorkFlowRepo
	viewRepo         space_view_repo.SpaceViewRepo
	memberRepo       space_member_repo.SpaceMemberRepo

	workItemTypeService   *service.WorkItemTypeService
	roleService           *work_item_role.WorkItemRoleService
	statusService         *service2.WorkItemStatusService
	workFlowDomainService *wf_service.WorkFlowService
	viewService           *view_service.SpaceViewService

	rdb *redis.Client

	tm trans.Transaction
}

func NewDML(
	spaceRepo repo.SpaceRepo,
	statusRepo status_repo.WorkItemStatusRepo,
	roleRepo role_repo.WorkItemRoleRepo,
	workItemTypeRepo repo2.WorkItemTypeRepo,
	flowRepo flow_repo.WorkFlowRepo,
	viewRepo space_view_repo.SpaceViewRepo,
	memberRepo space_member_repo.SpaceMemberRepo,

	workItemTypeService *service.WorkItemTypeService,
	roleService *work_item_role.WorkItemRoleService,
	statusService *service2.WorkItemStatusService,
	workFlowService *wf_service.WorkFlowService,
	viewService *view_service.SpaceViewService,

	rdb *redis.Client,

	tm trans.Transaction,
) *DML {
	return &DML{
		spaceRepo:        spaceRepo,
		statusRepo:       statusRepo,
		roleRepo:         roleRepo,
		workItemTypeRepo: workItemTypeRepo,
		flowRepo:         flowRepo,
		viewRepo:         viewRepo,
		memberRepo:       memberRepo,

		workItemTypeService:   workItemTypeService,
		roleService:           roleService,
		statusService:         statusService,
		workFlowDomainService: workFlowService,
		viewService:           viewService,

		rdb: rdb,

		tm: tm,
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
		members, err := d.memberRepo.GetSpaceMemberBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		userIds := stream.Map(members, func(v *member.SpaceMember) int64 {
			return v.UserId
		})
		userViews, err := d.viewService.InitUserGlobalView(ctx, spaceId, userIds)
		if err != nil {
			return err
		}

		for _, v := range userViews {
			if v.Key == consts.SystemViewKey_All {
				key := fmt.Sprintf("balala:user_temp_config:%v:columnsConfig_%v", v.UserId, v.SpaceId)
				v.TableConfig = d.rdb.Get(ctx, key).Val()
			}
		}

		err = d.viewRepo.CreateUserViews(ctx, userViews)
		if err != nil {
			return err
		}
	}

	return nil
}
