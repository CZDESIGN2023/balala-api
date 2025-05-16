package biz

import (
	"context"
	pb "go-cs/api/work_item_status/v1"
	"go-cs/internal/consts"
	perm_facade "go-cs/internal/domain/perm/facade"
	perm_service "go-cs/internal/domain/perm/service"
	domain_message "go-cs/internal/domain/pkg/message"
	member_repo "go-cs/internal/domain/space_member/repo"
	flow_repo "go-cs/internal/domain/work_flow/repo"
	witem_repo "go-cs/internal/domain/work_item/repo"
	"go-cs/internal/domain/work_item_status"
	"go-cs/internal/domain/work_item_status/repo"
	witem_status_service "go-cs/internal/domain/work_item_status/service"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/pkg/stream"
	"slices"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type WorkItemStatusUsecase struct {
	repo       repo.WorkItemStatusRepo
	memberRepo member_repo.SpaceMemberRepo
	wItemRepo  witem_repo.WorkItemRepo
	flowRepo   flow_repo.WorkFlowRepo

	wItemStatusService *witem_status_service.WorkItemStatusService
	permService        *perm_service.PermService

	domainMessageProducer *domain_message.DomainMessageProducer

	log *log.Helper
	tm  trans.Transaction
}

func NewWorkItemStatusUsecase(
	repo repo.WorkItemStatusRepo,
	memberRepo member_repo.SpaceMemberRepo,
	wItemRepo witem_repo.WorkItemRepo,
	flowRepo flow_repo.WorkFlowRepo,

	wItemStatusService *witem_status_service.WorkItemStatusService,
	permService *perm_service.PermService,

	domainMessageProducer *domain_message.DomainMessageProducer,

	tm trans.Transaction,
	logger log.Logger,
) *WorkItemStatusUsecase {
	moduleName := "WorkItemStatusUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &WorkItemStatusUsecase{
		log:        hlog,
		tm:         tm,
		repo:       repo,
		memberRepo: memberRepo,
		wItemRepo:  wItemRepo,
		flowRepo:   flowRepo,

		wItemStatusService: wItemStatusService,
		permService:        permService,

		domainMessageProducer: domainMessageProducer,
	}
}

func (uc *WorkItemStatusUsecase) SetSpaceWorkItemStatusRanking(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, rankingList []map[string]int64) error {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.NoPerm(ctx)
		return err
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_ModifyWorkFlowStatus,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	statusIds := stream.Map(rankingList, func(v map[string]int64) int64 {
		return v["id"]
	})

	statusMap, err := uc.repo.StatusMap(ctx, statusIds)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	rankingMap := stream.ToMap(rankingList, func(i int, v map[string]int64) (int64, int64) {
		return v["id"], v["ranking"]
	})

	statusList := stream.Values(statusMap)

	//slices.SortFunc(statusList, func(i, j *work_item_status.WorkItemStatusItem) int {
	//	return cmp.Compare(rankingMap[i.Id], rankingMap[j.Id])
	//})
	//
	flowScope := statusList[0].FlowScope
	//var startPos int
	//switch flowScope {
	//case consts.FlowScope_Stateflow:
	//	startPos = 100000
	//case consts.FlowScope_Workflow:
	//	startPos = 1
	//}

	for _, v := range statusList {
		ranking := rankingMap[v.Id]
		v.ChangeRanking(int64(ranking), oper)
	}

	txErr := uc.tm.InTx(ctx, func(ctx context.Context) error {
		for _, v := range statusList {
			err = uc.repo.SaveWorkItemStatusItem(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return errs.Internal(ctx, txErr)
	}

	msg := &domain_message.ChangeStatusOrder{
		DomainMessageBase: shared.DomainMessageBase{
			Oper:     oper,
			OperTime: time.Now(),
		},
		SpaceId:   spaceId,
		FlowScope: flowScope,
	}

	uc.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})

	return nil
}

func (uc *WorkItemStatusUsecase) QSpaceWorkItemStatusList(ctx context.Context, uid int64, req *pb.SpaceWorkItemStatusListRequest) (*pb.SpaceWorkItemStatusListReplyResult, error) {

	// 判断当前用户是否在要查询的项目空间内
	if req.SpaceId != 0 {
		member, err := uc.memberRepo.GetSpaceMember(ctx, req.SpaceId, uid)
		if member == nil || err != nil {
			//不是改空间成员，不允许查看该空间的其它成员列表
			err := errs.NoPerm(ctx)
			return nil, err
		}
	}

	spaceIds := []int64{req.SpaceId}
	if req.SpaceId == 0 {
		spaceIds, _ = uc.memberRepo.GetUserSpaceIdList(ctx, uid)
	}

	list, err := uc.repo.QSpaceWorkItemStatusList(ctx, spaceIds, consts.FlowScope(req.FlowScope))
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	result := &pb.SpaceWorkItemStatusListReplyResult{
		List:  make([]*pb.SpaceWorkItemStatusListReplyResult_Item, 0),
		Total: 0,
	}

	slices.SortFunc(list, func(i, j *work_item_status.WorkItemStatusItem) int {
		return i.Compare(j)
	})

	for _, v := range list {
		result.List = append(result.List, &pb.SpaceWorkItemStatusListReplyResult_Item{
			Id:             v.Id,
			Name:           v.Name,
			Key:            v.Key,
			Val:            v.Val,
			WorkItemTypeId: v.WorkItemTypeId,
			IsSys:          v.IsSys,
			Ranking:        v.Ranking,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
			StatusType:     int32(v.StatusType),
			FlowScope:      string(v.FlowScope),
		})
	}

	result.Total = int32(len(result.List))

	return result, nil
}

func (uc *WorkItemStatusUsecase) QSpaceWorkItemStatusById(ctx context.Context, oper *utils.LoginUserInfo, req *pb.QSpaceWorkItemStatusByIdRequest) (*pb.QSpaceWorkItemStatusByIdReplyData, error) {

	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, req.SpaceId, oper.UserId)
	if member == nil || err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.NoPerm(ctx)
		return nil, err
	}

	list, err := uc.repo.QSpaceWorkItemStatusById(ctx, req.SpaceId, req.Ids)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	result := &pb.QSpaceWorkItemStatusByIdReplyData{
		List:  make([]*pb.QSpaceWorkItemStatusByIdReplyData_ListItem, 0),
		Total: 0,
	}

	for _, v := range list {
		result.List = append(result.List, &pb.QSpaceWorkItemStatusByIdReplyData_ListItem{
			Id:             v.Id,
			Name:           v.Name,
			Key:            v.Key,
			WorkItemTypeId: v.WorkItemTypeId,
			IsSys:          v.IsSys,
			Ranking:        v.Ranking,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
			StatusType:     int32(v.StatusType),
			FlowScope:      string(v.FlowScope),
		})
	}

	result.Total = int32(len(result.List))

	return result, nil
}

func (uc *WorkItemStatusUsecase) SetSpaceWorkItemStatusName(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, statusId int64, newName string) error {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		err := errs.NoPerm(ctx)
		return err
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_ModifyWorkFlowStatus,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	wItemStatus, err := uc.repo.GetWorkItemStatusItem(ctx, statusId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if !wItemStatus.IsSameSpace(spaceId) {
		return errs.NoPerm(ctx)
	}

	if wItemStatus.IsSysFixStatus() {
		return errs.NoPerm(ctx)
	}

	wItemStatus.ChangeName(newName, oper)

	err = uc.repo.SaveWorkItemStatusItem(ctx, wItemStatus)
	if err != nil {
		return err
	}

	uc.domainMessageProducer.Send(ctx, wItemStatus.GetMessages())

	return nil
}

func (uc *WorkItemStatusUsecase) DelSpaceWorkItemStatus(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, statusId int64) error {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		err := errs.NoPerm(ctx)
		return err
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_DeleteWorkFlowStatus,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	wItemStatus, err := uc.repo.GetWorkItemStatusItem(ctx, statusId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if wItemStatus.IsSysDefaultStatus() {
		return errs.Business(ctx, "系统内置状态无法删除")
	}

	if !wItemStatus.IsSameSpace(spaceId) {
		return errs.NoPerm(ctx)
	}

	//检查是不是关联了流程配置
	tpltIds, err := uc.flowRepo.SearchTaskWorkFlowTemplateByNodeStateEvent(ctx, spaceId, cast.ToString(wItemStatus.Id))
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if len(tpltIds) > 0 {
		return errs.Business(ctx, "该任务状态已被应用至任务流程设置中，无法删除")
	}

	err = wItemStatus.OnDelete(oper)
	if err != nil {
		return errs.Business(ctx, err.Error())
	}

	err = uc.repo.DelSpaceWorkItemStatusItem(ctx, wItemStatus.SpaceId, wItemStatus.Id)
	if err != nil {
		return err
	}

	uc.domainMessageProducer.Send(ctx, wItemStatus.GetMessages())

	return nil
}

func (uc *WorkItemStatusUsecase) CreateSpaceWorkItemStatus(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, name string, flowScope consts.FlowScope, statusType consts.WorkItemStatusType) error {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		err := errs.NoPerm(ctx)
		return err
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CreateWorkFlowStatus,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	// 节点流只能创建过程状态
	if flowScope == consts.FlowScope_Workflow {
		statusType = consts.WorkItemStatusType_Process
	}

	newStatusItem, err := uc.wItemStatusService.CreateWorkItemStatusItem(ctx, witem_status_service.CreateWorkItemStatusItemRequest{
		Uid:        uid,
		SpaceId:    spaceId,
		Name:       name,
		FlowScope:  flowScope,
		StatusType: statusType,
	}, oper)
	if err != nil {
		return err
	}

	err = uc.repo.CreateWorkItemStatusItem(ctx, newStatusItem)
	if err != nil {
		return err
	}

	uc.domainMessageProducer.Send(ctx, newStatusItem.GetMessages())

	return nil
}

func (uc *WorkItemStatusUsecase) QSpaceWorkItemRelationCount(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, statusId int64) (int64, error) {

	// 判断当前用户是否在要查询的项目空间内
	_, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if err != nil {
		return 0, errs.NoPerm(ctx)
	}

	totalNum, err := uc.wItemRepo.CountWorkItemStatusRelatedSpaceWorkItem(ctx, spaceId, statusId)
	return totalNum, err
}

func (uc *WorkItemStatusUsecase) QSpaceTemplateRelationCount(ctx context.Context, spaceId int64, statusId int64) (int64, error) {

	//检查是不是关联了流程配置
	tpltIds, err := uc.flowRepo.SearchTaskWorkFlowTemplateByNodeStateEvent(ctx, spaceId, cast.ToString(statusId))
	if err != nil {
		return 0, err
	}

	return int64(len(tpltIds)), err
}
