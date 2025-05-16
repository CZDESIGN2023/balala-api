package biz

import (
	"context"
	pb "go-cs/api/space_work_item_flow/v1"
	"go-cs/internal/consts"
	perm_facade "go-cs/internal/domain/perm/facade"
	perm_service "go-cs/internal/domain/perm/service"
	domain_message "go-cs/internal/domain/pkg/message"
	space_repo "go-cs/internal/domain/space/repo"
	member_repo "go-cs/internal/domain/space_member/repo"
	user_repo "go-cs/internal/domain/user/repo"
	wf_repo "go-cs/internal/domain/work_flow/repo"
	wf_service "go-cs/internal/domain/work_flow/service"
	witem "go-cs/internal/domain/work_item"
	witem_facade "go-cs/internal/domain/work_item/facade"
	witem_repo "go-cs/internal/domain/work_item/repo"
	witem_service "go-cs/internal/domain/work_item/service"
	witem_status_repo "go-cs/internal/domain/work_item_status/repo"
	witem_status_service "go-cs/internal/domain/work_item_status/service"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/pkg/stream"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type SpaceWorkItemFlowUsecase struct {
	tm  trans.Transaction
	log *log.Helper

	spaceWorkItemRepo  witem_repo.WorkItemRepo
	workItemStatusRepo witem_status_repo.WorkItemStatusRepo
	workFlowRepo       wf_repo.WorkFlowRepo
	userRepo           user_repo.UserRepo
	spaceRepo          space_repo.SpaceRepo
	spaceMemberRepo    member_repo.SpaceMemberRepo

	permService           *perm_service.PermService
	workItemService       *witem_service.WorkItemService
	workFlowService       *wf_service.WorkFlowService
	workItemStatusService *witem_status_service.WorkItemStatusService

	domainMessageProducer *domain_message.DomainMessageProducer
}

func NewSpaceWorkItemFlowUsecase(
	tm trans.Transaction,
	logger log.Logger,

	spaceWorkItemRepo witem_repo.WorkItemRepo,
	workItemStatusRepo witem_status_repo.WorkItemStatusRepo,
	workFlowRepo wf_repo.WorkFlowRepo,
	userRepo user_repo.UserRepo,
	spaceRepo space_repo.SpaceRepo,
	spaceMemberRepo member_repo.SpaceMemberRepo,

	permService *perm_service.PermService,
	workItemService *witem_service.WorkItemService,
	workFlowService *wf_service.WorkFlowService,
	workItemStatusService *witem_status_service.WorkItemStatusService,

	domainMessageProducer *domain_message.DomainMessageProducer,

) *SpaceWorkItemFlowUsecase {
	moduleName := "SpaceWorkItemFlowUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceWorkItemFlowUsecase{
		log: hlog,
		tm:  tm,

		spaceMemberRepo:    spaceMemberRepo,
		spaceRepo:          spaceRepo,
		spaceWorkItemRepo:  spaceWorkItemRepo,
		userRepo:           userRepo,
		workItemStatusRepo: workItemStatusRepo,
		workFlowRepo:       workFlowRepo,

		permService:           permService,
		workItemService:       workItemService,
		workFlowService:       workFlowService,
		workItemStatusService: workItemStatusService,

		domainMessageProducer: domainMessageProducer,
	}
}

func (s *SpaceWorkItemFlowUsecase) SetWorkFlowDirectors(ctx context.Context, oper *utils.LoginUserInfo, workItemId int64, workFlowNodeCode string, addUserIds, removeUserIds []int64) error {

	addUserIds = stream.Unique(addUserIds)
	removeUserIds = stream.Unique(removeUserIds)

	workItem, err := s.spaceWorkItemRepo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{
		Directors:     true,
		Participators: true,
	}, &witem_repo.WithOption{
		FlowNodes: true,
		FlowRoles: true,
	})

	if err != nil {
		return errs.Internal(ctx, err)
	}

	_, err = s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	ok, err := s.spaceMemberRepo.AllIsMember(ctx, workItem.SpaceId, append(addUserIds, oper.UserId)...)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if !ok {
		return errs.NoPerm(ctx)
	}

	err = s.workItemService.SetDirectorsForWorkFlowMainTaskByNodeKey(
		ctx,
		workItem,
		&witem_service.SetDirectorsForTaskByFlowNodeRequest{
			NodeKey:                     workFlowNodeCode,
			AddDirectors:                utils.ToStrArray(addUserIds),
			RemoveDirectors:             utils.ToStrArray(removeUserIds),
			WorkItemStatusServiceFacade: witem_facade.BuildWorkItemStatusServiceFacade(s.workItemStatusService),
		},
		oper,
	)

	if err != nil {
		return err
	}

	err = s.tm.InTx(ctx, func(ctx context.Context) error {

		err = s.spaceWorkItemRepo.SaveWorkItem(ctx, workItem)
		if err != nil {
			return err
		}

		for _, v := range workItem.WorkItemFlowNodes {
			err = s.spaceWorkItemRepo.SaveWorkItemFlowNode(ctx, v)
			if err != nil {
				return err
			}
		}

		for _, v := range workItem.WorkItemFlowRoles {
			err = s.spaceWorkItemRepo.SaveWorkItemFlowRole(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	s.domainMessageProducer.Send(ctx, workItem.GetMessages())

	return nil
}

func (s *SpaceWorkItemFlowUsecase) SetWorkFlowPlanTime(ctx context.Context, oper *utils.LoginUserInfo, workItemId int64, workFlowNodeCode string, startAt, completeAt int64) error {

	uid := oper.UserId

	workItem, err := s.spaceWorkItemRepo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{
		PlanTime:  true,
		Directors: true,
	}, nil)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if startAt != 0 && completeAt != 0 && (startAt < workItem.Doc.PlanStartAt || completeAt > workItem.Doc.PlanCompleteAt) {
		if completeAt > workItem.Doc.PlanCompleteAt {
			return errs.Business(ctx, "不可大于已选总排期")
		}

		if startAt < workItem.Doc.PlanStartAt {
			return errs.Business(ctx, "不可小于已选总排期")
		}
	}

	space, err := s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, workItem.SpaceId, uid)
	if member == nil || err != nil {
		return errs.NoPerm(ctx)
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE_WORK_ITEM,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	stateInfo, err := s.workItemStatusRepo.GetWorkItemStatusInfo(ctx, space.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if stateInfo.HasArchivedItem(workItem.WorkItemStatus.Key) {
		return errs.Business(ctx, "已归档任务不可操作")
	}

	flowNodeInfo, err := s.spaceWorkItemRepo.GetWorkItemFlowNodeByNodeCode(ctx, workItemId, workFlowNodeCode)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	flowNodeInfo.ChangePlanTime(witem.PlanTime{
		StartAt:    startAt,
		CompleteAt: completeAt,
	}, oper)

	//设置节点的排期
	err = s.spaceWorkItemRepo.SaveWorkItemFlowNode(ctx, flowNodeInfo)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	s.domainMessageProducer.Send(ctx, flowNodeInfo.GetMessages())

	return nil
}

func (s *SpaceWorkItemFlowUsecase) UpgradeWorkItemFlow(ctx context.Context, oper *utils.LoginUserInfo, req *pb.UpgradeTaskWorkFlowRequest) error {
	uid := oper.UserId

	workItem, err := s.spaceWorkItemRepo.GetWorkItem(ctx, req.WorkItemId, &witem_repo.WithDocOption{
		Directors:     true,
		Participators: true,
	}, &witem_repo.WithOption{
		FlowRoles: true,
		FlowNodes: true,
	})
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if !workItem.IsSameSpace(req.SpaceId) {
		return errs.NoPerm(ctx)
	}

	if workItem.IsSubTask() {
		return errs.NoPerm(ctx)
	}

	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, workItem.SpaceId, uid)
	if member == nil || err != nil {
		return errs.NoPerm(ctx)
	}

	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_UPGRADE_SPACE_WORK_ITEM_FLOW,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	witemFlow, err := s.workFlowRepo.GetWorkFlow(ctx, workItem.WorkFlowId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	statusInfo, err := s.workItemStatusRepo.GetWorkItemStatusInfo(ctx, req.SpaceId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	upgradeTaskWorkFlowReq := &witem_service.UpgradeTaskFlowRequest{
		UpgradeTemplateId:             witemFlow.LastTemplateId,
		WorkItemStatusFacade:          witem_facade.BuildWorkItemStatusFacade(statusInfo),
		WorkFlowTemplateServiceFacade: witem_facade.BuildWorkFlowTemplateServiceFacade(s.workFlowService),
		Directors:                     make([]*witem_service.UpgradeTaskWorkFlowRequest_Directors, 0),
	}

	for _, v := range req.RoleDirectors {
		upgradeTaskWorkFlowReq.Directors = append(upgradeTaskWorkFlowReq.Directors, &witem_service.UpgradeTaskWorkFlowRequest_Directors{
			RoleId:    cast.ToString(v.RoleId),
			RoleKey:   v.RoleKey,
			Directors: v.Directors,
		})
	}

	var result *witem_service.UpgradeTaskFlowResult
	switch {
	case workItem.IsWorkFlowMainTask():
		result, err = s.workItemService.UpgradeTaskWorkFlow(ctx, workItem, upgradeTaskWorkFlowReq, oper)
	case workItem.IsStateFlowMainTask():
		result, err = s.workItemService.UpgradeTaskStateFlow(ctx, workItem, upgradeTaskWorkFlowReq, oper)
	}
	if err != nil {
		return err
	}

	err = s.tm.InTx(ctx, func(ctx context.Context) error {

		err = s.spaceWorkItemRepo.SaveWorkItem(ctx, workItem)
		if err != nil {
			return err
		}

		removeFlowRoleIds := result.DeleteWorkItemFlowRoles.GetIds()
		if len(removeFlowRoleIds) > 0 {
			_, err = s.spaceWorkItemRepo.DelWorkItemFlowRoleByIds(ctx, removeFlowRoleIds...)
			if err != nil {
				return err
			}
		}

		removeFlowNodeIds := result.DeleteWorkItemFlowNodes.GetIds()
		if len(removeFlowNodeIds) > 0 {
			_, err = s.spaceWorkItemRepo.DelWorkItemFlowNodeByIds(ctx, removeFlowNodeIds...)
			if err != nil {
				return err
			}
		}

		for _, v := range result.NewWorkItemFlowRoles {
			err = s.spaceWorkItemRepo.CreateWorkItemFlowRole(ctx, v)
			if err != nil {
				return err
			}
		}

		for _, v := range result.NewWorkItemFlowNodes {
			err = s.spaceWorkItemRepo.CreateWorkItemFlowNode(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	err = s.workFlowRepo.ClearHistoryTemplate(ctx, workItem.WorkItemFlowId)
	if err != nil {
		s.log.Error(err)
	}

	s.domainMessageProducer.Send(ctx, workItem.GetMessages())

	return nil
}

func (s *SpaceWorkItemFlowUsecase) BatUpgradeWorkItemFlowPrepare(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, flowId int64) ([]int64, error) {
	uid := oper.UserId

	_, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		return nil, errs.NoPerm(ctx)
	}

	// 验证权限
	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_UPGRADE_SPACE_WORK_ITEM_FLOW,
	})
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	flowInfo, err := s.workFlowRepo.GetWorkFlow(ctx, flowId)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	if !flowInfo.IsSameSpace(spaceId) {
		return nil, errs.NoPerm(ctx)
	}

	workItemIds, err := s.spaceWorkItemRepo.GetSpaceWorkItemIdsForUpgradeFlow(ctx, spaceId, flowInfo.Id, int64(flowInfo.Version))
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	return workItemIds, nil
}

func (s *SpaceWorkItemFlowUsecase) BatchUpgradeWorkItemFlow(ctx context.Context, oper *utils.LoginUserInfo, req *pb.BatchUpgradeTaskWorkFlowRequest) (*pb.BatchUpgradeTaskWorkFlowReplyData, error) {

	uid := oper.UserId

	upgradeWorkItemIds := stream.Unique(req.WorkItemIds)
	if len(upgradeWorkItemIds) == 0 || len(upgradeWorkItemIds) > 20 {
		return nil, errs.Business(ctx, "批次升级数量不能超过20")
	}

	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, req.SpaceId, uid)
	if member == nil || err != nil {
		return nil, errs.NoPerm(ctx)
	}

	// 验证权限
	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_UPGRADE_SPACE_WORK_ITEM_FLOW,
	})
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	space, err := s.spaceRepo.GetSpace(ctx, req.SpaceId)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	witemFlow, err := s.workFlowRepo.GetWorkFlow(ctx, req.FlowId)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	if !witemFlow.IsSameSpace(req.SpaceId) {
		return nil, errs.NoPerm(ctx)
	}

	statusInfo, err := s.workItemStatusRepo.GetWorkItemStatusInfo(ctx, space.Id)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	workItems, err := s.spaceWorkItemRepo.GetWorkItemByIds(ctx, upgradeWorkItemIds, &witem_repo.WithDocOption{
		Directors:     true,
		Participators: true,
	}, &witem_repo.WithOption{
		FlowRoles: true,
		FlowNodes: true,
	})
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	workItemMap := stream.ToMap(workItems, func(_ int, v *witem.WorkItem) (int64, *witem.WorkItem) { return v.Id, v })

	var resultData []*pb.BatchUpgradeTaskWorkFlowReplyData_Result
	for _, v := range upgradeWorkItemIds {

		resultItem := &pb.BatchUpgradeTaskWorkFlowReplyData_Result{WorkItemId: v, Code: 0, Message: ""}
		resultData = append(resultData, resultItem)

		workItem := workItemMap[v]

		if !workItem.IsSameSpace(req.SpaceId) {
			resultItem.Message = "任务与空间不匹配"
			continue
		}

		if workItem.WorkFlowId != witemFlow.Id {
			resultItem.Message = "任务与流程不匹配"
			continue
		}

		upgradeTaskWorkFlowReq := &witem_service.UpgradeTaskFlowRequest{
			UpgradeTemplateId:             witemFlow.LastTemplateId,
			WorkItemStatusFacade:          witem_facade.BuildWorkItemStatusFacade(statusInfo),
			WorkFlowTemplateServiceFacade: witem_facade.BuildWorkFlowTemplateServiceFacade(s.workFlowService),
		}

		for _, v := range req.RoleDirectors {
			upgradeTaskWorkFlowReq.Directors = append(upgradeTaskWorkFlowReq.Directors, &witem_service.UpgradeTaskWorkFlowRequest_Directors{
				RoleId:    cast.ToString(v.RoleId),
				RoleKey:   v.RoleKey,
				Directors: v.Directors,
			})
		}

		var result *witem_service.UpgradeTaskFlowResult
		switch {
		case workItem.IsWorkFlowMainTask():
			result, err = s.workItemService.UpgradeTaskWorkFlow(ctx, workItem, upgradeTaskWorkFlowReq, oper)
		case workItem.IsStateFlowMainTask():
			result, err = s.workItemService.UpgradeTaskStateFlow(ctx, workItem, upgradeTaskWorkFlowReq, oper)
		}
		if err != nil {
			resultItem.Message = "升级失败:" + err.Error()
			continue
		}

		err = s.tm.InTx(ctx, func(ctx context.Context) error {

			err = s.spaceWorkItemRepo.SaveWorkItem(ctx, workItem)
			if err != nil {
				return err
			}

			removeFlowRoleIds := result.DeleteWorkItemFlowRoles.GetIds()
			if len(removeFlowRoleIds) > 0 {
				_, err = s.spaceWorkItemRepo.DelWorkItemFlowRoleByIds(ctx, removeFlowRoleIds...)
				if err != nil {
					return err
				}
			}

			removeFlowNodeIds := result.DeleteWorkItemFlowNodes.GetIds()
			if len(removeFlowNodeIds) > 0 {
				_, err = s.spaceWorkItemRepo.DelWorkItemFlowNodeByIds(ctx, removeFlowNodeIds...)
				if err != nil {
					return err
				}
			}

			for _, v := range result.NewWorkItemFlowRoles {
				err = s.spaceWorkItemRepo.CreateWorkItemFlowRole(ctx, v)
				if err != nil {
					return err
				}
			}

			for _, v := range result.NewWorkItemFlowNodes {
				err = s.spaceWorkItemRepo.CreateWorkItemFlowNode(ctx, v)
				if err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			resultItem.Message = "升级失败:持久化异常，" + err.Error()
			continue
		}

		resultItem.Code = 200
		resultItem.Message = "升级成功"

		s.domainMessageProducer.Send(ctx, workItem.GetMessages())
	}

	flowIds := stream.Map(workItems, func(v *witem.WorkItem) int64 {
		return v.WorkFlowId
	})

	for _, id := range stream.Unique(flowIds) {
		err := s.workFlowRepo.ClearHistoryTemplate(ctx, id)
		s.log.Error(err)
	}

	return &pb.BatchUpgradeTaskWorkFlowReplyData{
		Result: resultData,
	}, nil
}
