package biz

import (
	"context"
	"go-cs/internal/domain/work_item"
	"go-cs/internal/domain/work_item/repo"
	"go-cs/internal/domain/work_item/service"
	"go-cs/internal/test/dbt"
	"testing"
	"time"

	"go-cs/internal/domain/work_item/facade"
)

func TestSetWorkFlowPlanTime(t *testing.T) {
	err := dbt.UC.SpaceWorkFlowUsecase.SetWorkFlowPlanTime(context.Background(), MockUserInfo(152), 1729, "state_2", time.Now().Unix(), time.Now().Unix())
	if err != nil {
		t.Error(err)
	}
}

func TestUpgradeWorkItemFlow(t *testing.T) {
	ctx := context.Background()

	witem, err := dbt.R.SpaceWorkItemRepo.GetWorkItem(ctx, 3801, &repo.WithDocOption{All: true}, &repo.WithOption{FlowRoles: true, FlowNodes: true})
	if err != nil {
		t.Error(err)
		return
	}

	wStatusInfo, err := dbt.R.WorkStatusRepo.GetWorkItemStatusInfo(ctx, witem.SpaceId)
	if err != nil {
		t.Error(err)
		return
	}

	result, err := dbt.S_Domain.WorkItemService.UpgradeTaskWorkFlow(ctx, witem, &service.UpgradeTaskFlowRequest{
		UpgradeTemplateId: 7192,
		Directors: []*service.UpgradeTaskWorkFlowRequest_Directors{
			{RoleId: "5254", RoleKey: "_productor", Directors: work_item.Directors{"83", "71"}},
			{RoleId: "5252", RoleKey: "_qa", Directors: work_item.Directors{"197", "83"}},
			{RoleId: "5256", RoleKey: "_reviewer", Directors: work_item.Directors{"37"}},
			{RoleId: "5255", RoleKey: "_ui_designner", Directors: work_item.Directors{"37", "197"}},
		},
		WorkItemStatusFacade:          facade.BuildWorkItemStatusFacade(wStatusInfo),
		WorkFlowTemplateServiceFacade: facade.BuildWorkFlowTemplateServiceFacade(dbt.S_Domain.WorkFlowService),
	}, MockUserInfo(1))

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(result)
}
