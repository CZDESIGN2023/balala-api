package service

import (
	"context"
	"slices"
	"time"

	"go-cs/internal/consts"
	domain "go-cs/internal/domain/work_item_status"
	"go-cs/internal/domain/work_item_status/repo"
	"go-cs/internal/pkg/biz_id"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/rand"

	"go-cs/internal/utils"

	"github.com/go-kratos/kratos/v2/log"
)

type WorkItemStatusService struct {
	log       *log.Helper
	repo      repo.WorkItemStatusRepo
	idService *biz_id.BusinessIdService
}

func NewWorkItemStatusService(
	repo repo.WorkItemStatusRepo,
	idService *biz_id.BusinessIdService,
	logger log.Logger,
) *WorkItemStatusService {

	moduleName := "WorkItemStatusService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &WorkItemStatusService{
		log:       hlog,
		idService: idService,
		repo:      repo,
	}
}

func (s *WorkItemStatusService) CreateSpaceDefaultStatusInfo(ctx context.Context, spaceId int64, flowScopes []consts.FlowScope, uid int64) *domain.WorkItemStatusInfo {

	statusList := make(domain.WorkItemStatusItems, 0)

	if slices.Contains(flowScopes, consts.FlowScope_All) {
		//系统默认用
		terminatedStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "已终止", "terminated", "3", 10, consts.WorkItemStatusType_Archived, uid, consts.FlowScope_All, nil)
		statusList = append(statusList, terminatedStatus)
	}

	if slices.Contains(flowScopes, consts.FlowScope_Workflow) {
		//流程模式用
		reviewingStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "评审中", "evaluating", "11", 550, consts.WorkItemStatusType_Process, uid, consts.FlowScope_Workflow, nil)
		planningStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "策划中", "planning", "9", 500, consts.WorkItemStatusType_Process, uid, consts.FlowScope_Workflow, nil)
		designingStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "设计中", "designing", "10", 450, consts.WorkItemStatusType_Process, uid, consts.FlowScope_Workflow, nil)
		checkingStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "验收中", "checking", "5", 400, consts.WorkItemStatusType_Process, uid, consts.FlowScope_Workflow, nil)
		waitConfirmStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "待确认", "wait_confirm", "6", 350, consts.WorkItemStatusType_Process, uid, consts.FlowScope_Workflow, nil)
		testingStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "测试中", "testing", "4", 300, consts.WorkItemStatusType_Process, uid, consts.FlowScope_Workflow, nil)
		progressingStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "进行中", "progressing", "1", 250, consts.WorkItemStatusType_Process, uid, consts.FlowScope_Workflow, nil)
		completedStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "已完成", "completed", "2", 200, consts.WorkItemStatusType_Archived, uid, consts.FlowScope_Workflow, nil)
		closeStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "已关闭", "close", "8", 150, consts.WorkItemStatusType_Archived, uid, consts.FlowScope_Workflow, nil)

		statusList = append(statusList, progressingStatus, completedStatus, testingStatus, checkingStatus, waitConfirmStatus, closeStatus, planningStatus, designingStatus, reviewingStatus)
	}

	if slices.Contains(flowScopes, consts.FlowScope_Stateflow) {
		//状态模式用
		//待确认 Pending -> 修复中,不予处理,转需求
		//修复中 Fixing -> 待验证,转需求
		//待验证 Pending_Verification -> 关闭，转需求
		//关闭 Closed -》 重启
		//重启 Restart -》 修复中， 不予处理
		//转需求 Convert_To_Rrequirement -》 关闭，不予处理
		//不予处理 Do_Not_Process-> 关闭，重启

		stPendingStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "待确认", "st_pending", "st_pending", 1300, consts.WorkItemStatusType_Process, uid, consts.FlowScope_Stateflow, nil)
		stFixingStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "修复中", "st_fixing", "st_fixing", 1250, consts.WorkItemStatusType_Process, uid, consts.FlowScope_Stateflow, nil)
		stDoNotProcessStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "不予处理", "st_do_not_process", "st_do_not_process", 1200, consts.WorkItemStatusType_Process, uid, consts.FlowScope_Stateflow, nil)
		stConvertToStoryStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "转需求", "st_convert_to_story", "st_convert_to_story", 1150, consts.WorkItemStatusType_Process, uid, consts.FlowScope_Stateflow, nil)
		stPendingVerificationStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "待验证", "st_pending_verification", "st_pending_verification", 1100, consts.WorkItemStatusType_Process, uid, consts.FlowScope_Stateflow, nil)
		stProgressingStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "进行中", "st_progressing", "st_progressing", 1050, consts.WorkItemStatusType_Process, uid, consts.FlowScope_Stateflow, nil)
		stRestartStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "重启", "st_restart", "st_restart", 1000, consts.WorkItemStatusType_Process, uid, consts.FlowScope_Stateflow, nil)
		stClosedStatus, _ := s.newSysWorkItemStatusItem(ctx, spaceId, "关闭", "st_closed", "st_closed", 950, consts.WorkItemStatusType_Archived, uid, consts.FlowScope_Stateflow, nil)

		statusList = append(statusList, stPendingStatus, stFixingStatus, stPendingVerificationStatus, stClosedStatus, stRestartStatus, stConvertToStoryStatus, stDoNotProcessStatus, stProgressingStatus)
	}

	//---
	stateInfo := domain.BuildWorkItemStatusInfo(spaceId, statusList)
	return stateInfo
}

type CreateWorkItemStatusItemRequest struct {
	SpaceId    int64
	Name       string
	Key        string
	Val        string
	Uid        int64
	StatusType consts.WorkItemStatusType
	Ranking    int64
	IsSys      int32
	FlowScope  consts.FlowScope
}

func (s *WorkItemStatusService) CreateWorkItemStatusItem(ctx context.Context, req CreateWorkItemStatusItemRequest, oper shared.Oper) (*domain.WorkItemStatusItem, error) {

	itemKey := req.Key
	itemVal := req.Val
	statusType := req.StatusType

	if itemKey == "" {
		itemKey = "status_" + rand.Letters(5)
	}

	if itemVal == "" {
		itemVal = itemKey
	}

	if statusType == 0 {
		statusType = consts.WorkItemStatusType_Process
	}

	return s.newWorkItemStatusItem(ctx, req.SpaceId, req.Name, itemKey, itemVal, req.Ranking, statusType, req.IsSys, req.Uid, req.FlowScope, oper)
}

func (s *WorkItemStatusService) newNormalWorkItemStatusItem(ctx context.Context, spaceId int64, name string, key string, val string, Ranking int64, statusType consts.WorkItemStatusType, uid int64, flowScope consts.FlowScope, oper shared.Oper) (*domain.WorkItemStatusItem, error) {
	return s.newWorkItemStatusItem(ctx, spaceId, name, key, val, Ranking, statusType, 0, uid, flowScope, oper)
}

func (s *WorkItemStatusService) newSysWorkItemStatusItem(ctx context.Context, spaceId int64, name string, key string, val string, Ranking int64, statusType consts.WorkItemStatusType, uid int64, flowScope consts.FlowScope, oper shared.Oper) (*domain.WorkItemStatusItem, error) {
	return s.newWorkItemStatusItem(ctx, spaceId, name, key, val, Ranking, statusType, 1, uid, flowScope, oper)
}

func (s *WorkItemStatusService) newWorkItemStatusItem(ctx context.Context, spaceId int64, name string, key string, val string, Ranking int64, statusType consts.WorkItemStatusType, isSys int32, uid int64, flowScope consts.FlowScope, oper shared.Oper) (*domain.WorkItemStatusItem, error) {

	//isExist, err := s.repo.IsExistByWorkItemStatusName(ctx, spaceId, name, flowScope)
	//if err != nil {
	//	return nil, errs.Internal(ctx, err)
	//}
	//
	//if isExist {
	//	return nil, nil
	//}
	var err error

	if Ranking == 0 {
		Ranking, err = s.repo.GetMaxRanking(ctx, spaceId)
		if err != nil {
			Ranking = time.Now().Unix()
		} else {
			Ranking = Ranking + 100
		}
	}

	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkItemStatus)
	if bizId == nil {
		return nil, errs.Business(ctx, "状态ID分配失败")
	}

	statusItem := domain.NewWorkItemStatusItem(bizId.Id, spaceId, name, key, val, isSys, Ranking, statusType, uid, flowScope, oper)
	return statusItem, nil
}

func (s *WorkItemStatusService) GetWorkItemStatusInfo(ctx context.Context, spaceId int64) (*domain.WorkItemStatusInfo, error) {
	return s.repo.GetWorkItemStatusInfo(ctx, spaceId)
}

func (s *WorkItemStatusService) GetWorkItemStatusItem(ctx context.Context, statusId int64) (*domain.WorkItemStatusItem, error) {
	return s.repo.GetWorkItemStatusItem(ctx, statusId)
}

func (s *WorkItemStatusService) CheckAndUpdateSpaceWorkItemStatusName(ctx context.Context, wItemStatus *domain.WorkItemStatusItem, newName string, oper shared.Oper) error {
	err := wItemStatus.ChangeName(newName, oper)
	if err != nil {
		return err
	}

	return nil
}
