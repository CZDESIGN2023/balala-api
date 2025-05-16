package biz

import (
	"cmp"
	"context"
	"fmt"
	"go-cs/api/comm"
	pb "go-cs/api/work_flow/v1"
	"go-cs/internal/bean/rsp_convert"
	"go-cs/internal/bean/vo/query"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/consts"
	perm_facade "go-cs/internal/domain/perm/facade"
	perm_service "go-cs/internal/domain/perm/service"
	domain_message "go-cs/internal/domain/pkg/message"
	space_repo "go-cs/internal/domain/space/repo"
	member_repo "go-cs/internal/domain/space_member/repo"
	"go-cs/internal/domain/work_flow/facade"
	tplt_config "go-cs/internal/domain/work_flow/flow_tplt_config"
	"go-cs/internal/domain/work_item_role"
	"go-cs/internal/domain/work_item_status"
	witem_type "go-cs/internal/domain/work_item_type"
	"go-cs/internal/pkg/biz_id"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/rand"
	"go-cs/pkg/stream"
	"slices"
	"sort"
	"time"

	"github.com/google/uuid"

	witem_role_repo "go-cs/internal/domain/work_item_role/repo"
	witem_status_repo "go-cs/internal/domain/work_item_status/repo"

	statics_repo "go-cs/internal/domain/statics/repo"
	wf_domain "go-cs/internal/domain/work_flow"
	wf_repo "go-cs/internal/domain/work_flow/repo"
	wf_service "go-cs/internal/domain/work_flow/service"
	witem_repo "go-cs/internal/domain/work_item/repo"
	witem_type_repo "go-cs/internal/domain/work_item_type/repo"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type WorkFlowUsecase struct {
	repo            wf_repo.WorkFlowRepo
	memberRepo      member_repo.SpaceMemberRepo
	spaceRepo       space_repo.SpaceRepo
	wItemTypeRepo   witem_type_repo.WorkItemTypeRepo
	wItemStatusRepo witem_status_repo.WorkItemStatusRepo
	wItemRoleRepo   witem_role_repo.WorkItemRoleRepo
	staticsRepo     statics_repo.StaticsRepo
	wItemRepo       witem_repo.WorkItemRepo

	wfService   *wf_service.WorkFlowService
	permService *perm_service.PermService
	idService   *biz_id.BusinessIdService

	domainMessageProducer *domain_message.DomainMessageProducer

	log *log.Helper
	tm  trans.Transaction
}

func NewWorkFlowUsecase(
	repo wf_repo.WorkFlowRepo,
	memberRepo member_repo.SpaceMemberRepo,
	spaceRepo space_repo.SpaceRepo,
	wItemTypeRepo witem_type_repo.WorkItemTypeRepo,
	wItemStatusRepo witem_status_repo.WorkItemStatusRepo,
	wItemRoleRepo witem_role_repo.WorkItemRoleRepo,
	staticsRepo statics_repo.StaticsRepo,
	wItemRepo witem_repo.WorkItemRepo,

	wfService *wf_service.WorkFlowService,
	permService *perm_service.PermService,
	idService *biz_id.BusinessIdService,

	domainMessageProducer *domain_message.DomainMessageProducer,

	tm trans.Transaction,
	logger log.Logger,
) *WorkFlowUsecase {
	moduleName := "WorkFlowUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &WorkFlowUsecase{
		repo:            repo,
		memberRepo:      memberRepo,
		spaceRepo:       spaceRepo,
		wItemTypeRepo:   wItemTypeRepo,
		wItemStatusRepo: wItemStatusRepo,
		wItemRoleRepo:   wItemRoleRepo,
		staticsRepo:     staticsRepo,
		wItemRepo:       wItemRepo,

		wfService:   wfService,
		permService: permService,
		idService:   idService,

		domainMessageProducer: domainMessageProducer,

		log: hlog,
		tm:  tm,
	}
}

func (uc *WorkFlowUsecase) SaveWorkFlowTemplateConfig(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, flowId int64, flowTpltConfJSON string) error {

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
		Perm:              consts.PERM_ModifySpaceWorkFlow,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	// 查询工作流信息
	wf, err := uc.repo.GetWorkFlow(ctx, flowId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	checkStatus := func(statusIds []int64) error {
		status, err := uc.wItemStatusRepo.GetWorkItemStatusItemsBySpace(ctx, spaceId)
		if err != nil {
			return errs.Internal(ctx, err)
		}
		spaceStatusIds := stream.Map(status, func(v *work_item_status.WorkItemStatusItem) int64 {
			return v.Id
		})

		if !stream.ContainsArr(spaceStatusIds, statusIds) {
			return errs.Custom(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER, "工作流模板状态配置错误")
		}

		return nil
	}

	checkRoles := func(roleIds []int64) error {
		roles, err := uc.wItemRoleRepo.GetWorkItemRoles(ctx, spaceId)
		if err != nil {
			return errs.Internal(ctx, err)
		}
		spaceRoleIds := stream.Map(roles, func(v *work_item_role.WorkItemRole) int64 {
			return v.Id
		})
		if !stream.ContainsArr(spaceRoleIds, roleIds) {
			return errs.Custom(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER, "工作流模板角色配置错误")
		}

		return nil
	}

	checkUsers := func(userIds []int64) error {
		allIsMember, err := uc.memberRepo.AllIsMember(ctx, spaceId, userIds...)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		if !allIsMember {
			return errs.Custom(ctx, comm.ErrorCode_COMMON_WRONG_PARAMETER, "工作流模板成员配置错误")
		}

		return nil
	}

	// 更新模板
	var newWfTemplate *wf_domain.WorkFlowTemplate
	//核对配置信息
	switch wf.FlowMode {
	default:
		return errs.Business(ctx, fmt.Sprintf("unknown flowMode %v", wf.FlowMode))
	case consts.FlowMode_WorkFlow:
		wfConf, err := tplt_config.MustFormWorkFlowJson(flowTpltConfJSON)
		if err != nil {
			return errs.Business(ctx, "工作流模板配置错误")
		}

		// 检查角色，状态，人员是否在项目中
		roleIds := wfConf.GetAllRoleId()
		statusIds := wfConf.GetAllStatusId()
		userIds := wfConf.GetAllRelatedUserId()

		err = checkStatus(statusIds)
		if err != nil {
			return err
		}
		err = checkRoles(roleIds)
		if err != nil {
			return err
		}
		err = checkUsers(userIds)
		if err != nil {
			return err
		}

		newWfTemplate, err = uc.wfService.UpdateWorkFlowTemplateConf(ctx, wf, wfConf, oper)
		if err != nil {
			return errs.Business(ctx, err.Error())
		}
	case consts.FlowMode_StateFlow:
		wfConf, err := tplt_config.MustFormStateFlowJson(flowTpltConfJSON)
		if err != nil {
			return errs.Business(ctx, "状态流模板配置错误")
		}

		// 检查角色，状态，人员是否在项目中
		roleIds := wfConf.GetAllRoleId()
		statusIds := wfConf.GetAllStatusId()
		userIds := wfConf.GetAllRelatedUserId()

		err = checkStatus(statusIds)
		if err != nil {
			return err
		}
		err = checkRoles(roleIds)
		if err != nil {
			return err
		}
		err = checkUsers(userIds)
		if err != nil {
			return err
		}

		newWfTemplate, err = uc.wfService.UpdateStateFlowTemplateConf(ctx, wf, wfConf, oper)
		if err != nil {
			return errs.Business(ctx, err.Error())
		}
	}

	//持久化
	txErr := uc.tm.InTx(ctx, func(ctx context.Context) error {
		err = uc.repo.SaveWorkFlow(ctx, wf)
		if err != nil {
			return err
		}

		err = uc.repo.CreateWorkFlowTemplate(ctx, newWfTemplate)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return errs.Internal(ctx, txErr)
	}

	err = uc.repo.ClearHistoryTemplate(ctx, wf.Id)
	if err != nil {
		uc.log.Error(err)
	}

	uc.domainMessageProducer.Send(ctx, wf.GetMessages())

	return nil
}

func (uc *WorkFlowUsecase) SetWorkFlowRanking(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, rankingList []map[string]int64) error {

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
		Perm:              consts.PERM_ModifySpaceWorkFlow,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	workFlows := make(wf_domain.WorkFlows, 0)
	for _, v := range rankingList {

		flowId := cast.ToInt64(v["id"])
		newRanking := cast.ToInt64(v["ranking"])

		workFlow, err := uc.repo.GetWorkFlow(ctx, flowId)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		workFlow.ChangeRanking(newRanking, oper)
		workFlows = append(workFlows, workFlow)
	}

	txErr := uc.tm.InTx(ctx, func(ctx context.Context) error {

		for _, v := range workFlows {
			err = uc.repo.SaveWorkFlow(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return errs.Internal(ctx, txErr)
	}

	msg := &domain_message.ChangeWorkFlowOrder{
		DomainMessageBase: shared.DomainMessageBase{
			Oper:     oper,
			OperTime: time.Now(),
		},
		SpaceId: spaceId,
	}

	uc.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})

	return nil
}

func (uc *WorkFlowUsecase) QSpaceWorkFlowPageList(ctx context.Context, oper *utils.LoginUserInfo, req *pb.SpaceWorkFlowPageListRequest) (*pb.SpaceWorkFlowPageListReplyResult, error) {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, req.SpaceId, uid)
	if member == nil || err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.NoPerm(ctx)
		return nil, err
	}

	workItemTypeInfoQueryResult, err := uc.wItemTypeRepo.QWorkItemTypeInfo(ctx, query.WorkItemTypeInfoQuery{
		SpaceId: req.SpaceId,
	})

	if err != nil {
		return nil, err
	}

	queryListResult, err := uc.repo.QTaskWorkFlowViewList(ctx, &query.TaskWorkFlowListQuery{
		SpaceId:         req.SpaceId,
		WorkItemTypeIds: workItemTypeInfoQueryResult.GetMainTaskTypeIds(),
	})

	if err != nil {
		return nil, err
	}

	result := &pb.SpaceWorkFlowPageListReplyResult{
		Total: 0,
		List:  make([]*pb.SpaceWorkFlowPageListReplyResult_Item, 0),
	}

	statusMap, err := uc.wItemStatusRepo.StatusMapBySpaceIds(ctx, []int64{req.SpaceId})
	if err != nil {
		return nil, err
	}

	var list []*pb.SpaceWorkFlowPageListReplyResult_Item
	for _, f := range queryListResult.List {
		var workFlowConf *pb.WorkFlowConf
		var stateFlowConf *pb.StateFlowConf

		switch f.WorkFlow.FlowMode {
		case consts.FlowMode_WorkFlow:
			var nodes []*pb.WorkFlowConf_Node
			var connections []*pb.WorkFlowConf_Connection

			for _, v := range f.WorkFlowTemplate.WorkFlowConf().Nodes {
				nodes = append(nodes, &pb.WorkFlowConf_Node{
					Name: v.Name,
					Key:  v.Key,
				})
			}

			for _, v := range f.WorkFlowTemplate.WorkFlowConf().Connections {
				connections = append(connections, &pb.WorkFlowConf_Connection{
					StartNode: v.StartNode,
					EndNode:   v.EndNode,
				})
			}

			workFlowConf = &pb.WorkFlowConf{
				Nodes:       nodes,
				Connections: connections,
			}
		case consts.FlowMode_StateFlow:
			flowNodes := f.WorkFlowTemplate.StateFlowConf().StateFlowNodes
			flowRules := f.WorkFlowTemplate.StateFlowConf().StateFlowTransitionRule

			nodes := stream.Map(flowNodes, func(node *tplt_config.StateFlowNode) *pb.StateFlowConf_Node {
				nodeName := ""
				statusType := consts.WorkItemStatusType_Process

				if v := statusMap[cast.ToInt64(node.SubStateId)]; v != nil {
					nodeName = v.Name
					statusType = v.StatusType
				}
				return &pb.StateFlowConf_Node{
					Name:            nodeName,
					Key:             node.Key,
					StatusType:      int64(statusType),
					IsInitState:     node.IsInitState,
					IsArchivedState: node.IsArchivedState,
				}
			})

			connections := stream.Map(flowRules, func(rule *tplt_config.StateFlowTransitionRule) *pb.StateFlowConf_Connection {
				return &pb.StateFlowConf_Connection{
					SourceStateKey: rule.SourceStateKey,
					TargetStateKey: rule.TargetStateKey,
				}
			})

			stateFlowConf = &pb.StateFlowConf{
				StateFlowNodes:          nodes,
				StateFlowTransitionRule: connections,
			}
		}

		list = append(list, &pb.SpaceWorkFlowPageListReplyResult_Item{
			Id:            f.WorkFlow.Id,
			Name:          f.WorkFlow.Name,
			Status:        int64(f.WorkFlow.Status),
			Version:       int64(f.WorkFlow.Version),
			TemplateId:    f.WorkFlow.LastTemplateId,
			SpaceId:       f.WorkFlow.SpaceId,
			CreatedAt:     f.WorkFlow.CreatedAt,
			UpdatedAt:     f.WorkFlow.UpdatedAt,
			UserId:        f.WorkFlow.UserId,
			IsSys:         int64(f.WorkFlow.IsSys),
			Ranking:       f.WorkFlow.Ranking,
			FlowMode:      string(f.WorkFlow.FlowMode),
			WorkFlowConf:  workFlowConf,
			StateFlowConf: stateFlowConf,
		})
	}

	//重新排序一下
	sort.Slice(list, func(i, j int) bool {
		return list[j].Ranking < list[i].Ranking
	})

	return &pb.SpaceWorkFlowPageListReplyResult{
		Total: cast.ToInt32(len(result.List)),
		List:  list,
	}, nil
}

func (uc *WorkFlowUsecase) CreateWorkFlow(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, flowName string, flowStatus int, flowMode string) (int64, error) {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.NoPerm(ctx)
		return 0, err
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CreateSpaceWorkFlow,
	})
	if err != nil {
		return 0, errs.NoPerm(ctx)
	}

	workItemTypeInfoQueryResult, err := uc.wItemTypeRepo.QWorkItemTypeInfo(ctx, query.WorkItemTypeInfoQuery{
		SpaceId: spaceId,
	})
	if err != nil {
		return 0, err
	}

	statusInfo, err := uc.wItemStatusRepo.GetWorkItemStatusInfo(ctx, spaceId)
	if err != nil {
		return 0, err
	}

	roleInfo, err := uc.wItemRoleRepo.GetWorkItemRoles(ctx, spaceId)
	if err != nil {
		return 0, err
	}

	workFlowStatus := wf_domain.WorkFlowStatus_Disable
	if flowStatus == 1 {
		workFlowStatus = wf_domain.WorkFlowStatus_Enable
	}

	var workItemTypeId int64
	switch consts.WorkFlowMode(flowMode) {
	case consts.FlowMode_WorkFlow:
		workItemTypeId = workItemTypeInfoQueryResult.GetWorkFlowTaskType().Id
	case consts.FlowMode_StateFlow:
		workItemTypeId = workItemTypeInfoQueryResult.GetStateFlowTaskType().Id
	default:
		return 0, errs.Business(ctx, "流程类型错误")
	}

	result, err := uc.wfService.NewDefaultWorkFlow(ctx, wf_service.GenerateDefaultWorkFlowReq{
		WorkItemStatusInfo: facade.BuildWorkItemStatusInfo(statusInfo),
		WorkItemRoleInfo:   facade.BuildWorkItemRoleInfo(roleInfo),

		SpaceId:        spaceId,
		Uid:            uid,
		FlowName:       flowName,
		WorkItemTypeId: workItemTypeId,
		WorkFlowStatus: workFlowStatus,
		FlowMode:       flowMode,
	}, oper)

	if err != nil {
		return 0, err
	}

	txErr := uc.tm.InTx(ctx, func(ctx context.Context) error {
		err := uc.repo.CreateWorkFlow(ctx, result.WorkFlow)
		if err != nil {
			return err
		}

		err = uc.repo.CreateWorkFlowTemplate(ctx, result.WorkFlowTemplate)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return 0, txErr
	}

	uc.domainMessageProducer.Send(ctx, result.WorkFlow.GetMessages())

	uc.resetRank(ctx, spaceId, result.WorkFlow.Id)

	return result.WorkFlow.Id, nil

}

func (uc *WorkFlowUsecase) QWorkFlow(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workFlowId int64) (*rsp.WorkFlowInfo, error) {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		err := errs.NoPerm(ctx)
		return nil, err
	}

	workFlow, err := uc.repo.GetWorkFlow(ctx, workFlowId)
	if err != nil || !workFlow.IsSameSpace(spaceId) {
		err := errs.NoPerm(ctx)
		return nil, err
	}

	workFlowTemplate, err := uc.repo.GetFlowTemplate(ctx, workFlow.LastTemplateId)
	if err != nil {
		err := errs.NoPerm(ctx)
		return nil, err
	}

	result := &rsp.WorkFlowInfo{
		Id:             workFlow.Id,
		WorkItemTypeId: workFlow.WorkItemTypeId,
		SpaceId:        workFlow.SpaceId,
		Name:           workFlow.Name,
		Ranking:        workFlow.Ranking,
		Version:        workFlow.Version,
		TemplateId:     workFlow.LastTemplateId,
		FlowMode:       string(workFlow.FlowMode),
	}

	switch workFlow.FlowMode {
	case consts.FlowMode_WorkFlow:
		result.TemplateConf = rsp_convert.WorkFlowConfToRsp(workFlowTemplate.WorkFlowConf())
		result.FlowConf = rsp_convert.WorkFlowDefaultConf(tplt_config.GetWorkFlowDefaultConf())
	case consts.FlowMode_StateFlow:
		statusIds := stream.Map(workFlowTemplate.StateFlowConf().StateFlowNodes, func(node *tplt_config.StateFlowNode) int64 {
			return cast.ToInt64(node.SubStateId)
		})

		statusMap, err := uc.wItemStatusRepo.StatusMap(ctx, statusIds)
		if err != nil {
			return nil, errs.Internal(ctx, err)
		}

		// 设置状态名称
		for _, node := range workFlowTemplate.StateFlowConfig.StateFlowNodes {
			var nodeName string
			if status := statusMap[cast.ToInt64(node.SubStateId)]; status != nil {
				nodeName = status.Name
			}
			node.Name = nodeName
		}

		result.StateFlowTemplateConf = rsp_convert.StateFlowConfToRsp(workFlowTemplate.StateFlowConf())
	}

	return result, nil
}

func (uc *WorkFlowUsecase) QWorkFlowTemplate(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workFlowTpltId int64) (*rsp.WorkFlowInfo, error) {

	uid := oper.UserId

	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		err := errs.NoPerm(ctx)
		return nil, err
	}

	workFlowTemplate, err := uc.repo.GetFlowTemplate(ctx, workFlowTpltId)
	if err != nil {
		err := errs.NoPerm(ctx)
		return nil, err
	}

	if workFlowTemplate.SpaceId != spaceId {
		err := errs.NoPerm(ctx)
		return nil, err
	}

	workFlow, err := uc.repo.GetWorkFlow(ctx, workFlowTemplate.WorkFlowId)
	if err != nil || !workFlow.IsSameSpace(spaceId) {
		err := errs.NoPerm(ctx)
		return nil, err
	}

	result := &rsp.WorkFlowInfo{
		Id:             workFlow.Id,
		WorkItemTypeId: workFlow.WorkItemTypeId,
		SpaceId:        workFlow.SpaceId,
		Name:           workFlow.Name,
		Ranking:        workFlow.Ranking,
		Version:        workFlow.Version,
		TemplateId:     workFlow.LastTemplateId,
		FlowMode:       string(workFlow.FlowMode),
	}

	switch workFlowTemplate.FlowMode {
	case consts.FlowMode_WorkFlow:
		result.TemplateConf = rsp_convert.WorkFlowConfToRsp(workFlowTemplate.WorkFlowConf())
		result.FlowConf = rsp_convert.WorkFlowDefaultConf(tplt_config.GetWorkFlowDefaultConf())
	case consts.FlowMode_StateFlow:
		result.StateFlowTemplateConf = rsp_convert.StateFlowConfToRsp(workFlowTemplate.StateFlowConf())
	}

	return result, nil
}

func (uc *WorkFlowUsecase) DelWorkFlow(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workFlowId int64) error {

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
		Perm:              consts.PERM_DeleteSpaceWorkFlow,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	workFlow, err := uc.repo.GetWorkFlow(ctx, workFlowId)
	if err != nil || !workFlow.IsSameSpace(spaceId) {
		err := errs.NoPerm(ctx)
		return err
	}

	if workFlow.IsSyPreset() {
		return errs.Business(ctx, "系统预置流程不允许删除")
	}

	totalNum, err := uc.wItemRepo.CountWorkFlowRelatedSpaceWorkItem(ctx, spaceId, workFlow.Id, nil)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if totalNum > 0 {
		return errs.Business(ctx, "存在已使用当前流程的任务，不允许删除")
	}

	workFlow.OnDelete(oper)

	txErr := uc.tm.InTx(ctx, func(ctx context.Context) error {

		err := uc.repo.DelWorkFlow(ctx, workFlowId)
		if err != nil {
			return err
		}

		err = uc.repo.DelWorkFlowTemplateByFlowId(ctx, workFlowId)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	uc.domainMessageProducer.Send(ctx, workFlow.GetMessages())

	return err
}

func (uc *WorkFlowUsecase) SetWorkFlowStatus(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workFlowId int64, newStatus int64) error {
	//调整排序
	workFlow, err := uc.repo.GetWorkFlow(ctx, workFlowId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.NoPerm(ctx)
		return err
	}

	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_ModifySpaceWorkFlow,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	if !workFlow.IsSameSpace(spaceId) {
		return errs.NoPerm(ctx)
	}

	oldStatus := workFlow.Status
	if newStatus == int64(wf_domain.WorkFlowStatus_Hide) {
		relatedWorkItem, err := uc.wItemRepo.HasWorkItemRelateFlow(ctx, workFlow.SpaceId, workFlow.Id)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		if relatedWorkItem {
			return errs.Business(ctx, "存在关联任务，不允许隐藏")
		}
	}

	err = uc.wfService.ChangeWorkFlowStatus(ctx, workFlow, wf_domain.WorkFlowStatus(newStatus), oper)
	if err != nil {
		return errs.Business(ctx, err.Error())
	}

	//持久化
	err = uc.repo.SaveWorkFlow(ctx, workFlow)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	uc.domainMessageProducer.Send(ctx, workFlow.GetMessages())

	//重置排序
	if oldStatus == wf_domain.WorkFlowStatus_Enable || newStatus == int64(wf_domain.WorkFlowStatus_Enable) {
		uc.resetRank(ctx, spaceId, workFlowId)
	}

	return nil
}

func (uc *WorkFlowUsecase) resetRank(ctx context.Context, spaceId int64, workFlowId int64) error {
	// 全量调整排序值
	list, err := uc.repo.GetWorkFlowBySpaceId(ctx, spaceId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	slices.SortFunc(list, func(a, b *wf_domain.WorkFlow) int {
		return cmp.Compare(b.Ranking, a.Ranking)
	})

	target := stream.FilterOne(list, func(v *wf_domain.WorkFlow) bool {
		return v.Id == workFlowId
	})

	m := stream.GroupBy(list, func(v *wf_domain.WorkFlow) wf_domain.WorkFlowStatus {
		if v.Status == wf_domain.WorkFlowStatus_Hide {
			return wf_domain.WorkFlowStatus_Disable
		}
		return v.Status
	})

	enableList := m[wf_domain.WorkFlowStatus_Enable]
	disableList := m[wf_domain.WorkFlowStatus_Disable]

	moveToFirst := func(list []*wf_domain.WorkFlow, targetId int64) {
		idx := slices.IndexFunc(list, func(v *wf_domain.WorkFlow) bool {
			return v.Id == targetId
		})

		if idx == -1 {
			return
		}

		stream.Move(list, idx, 0)
	}

	switch target.Status {
	case wf_domain.WorkFlowStatus_Enable:
		moveToFirst(enableList, target.Id)
	case wf_domain.WorkFlowStatus_Disable:
		moveToFirst(disableList, target.Id)
	}

	list = append(enableList, disableList...)
	listLen := len(list)

	err = uc.tm.InTx(ctx, func(ctx context.Context) error {
		for i, v := range list {
			v.UpdateRanking(int64(listLen-i) * 100)
			err = uc.repo.SaveWorkFlow(ctx, v)
			if err != nil {
				return errs.Internal(ctx, err)
			}
		}

		return nil
	})

	return err
}

func (uc *WorkFlowUsecase) SetWorkFlowName(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workFlowId int64, newName string) error {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.NoPerm(ctx)
		return err
	}

	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_ModifySpaceWorkFlow,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	//调整排序
	workFlow, err := uc.repo.GetWorkFlow(ctx, workFlowId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if !workFlow.IsSameSpace(spaceId) {
		return errs.NoPerm(ctx)
	}

	err = uc.wfService.ChangeWorkFlowName(ctx, workFlow, newName, oper)
	if err != nil {
		return err
	}

	//持久化
	err = uc.repo.SaveWorkFlow(ctx, workFlow)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	uc.domainMessageProducer.Send(ctx, workFlow.GetMessages())

	return nil
}

func (uc *WorkFlowUsecase) CopyWorkFlow(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workFlowId int64, flowNewName string, flowDefaultStatus int32) (int64, error) {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.NoPerm(ctx)
		return 0, err
	}

	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_ModifySpaceWorkFlow,
	})
	if err != nil {
		return 0, errs.NoPerm(ctx)
	}

	workFlow, err := uc.repo.GetWorkFlow(ctx, workFlowId)
	if err != nil {
		return 0, errs.Internal(ctx, err)
	}

	// 从其他空间复制
	if !workFlow.IsSameSpace(spaceId) {
		srcSpaceId := workFlow.SpaceId

		srcMember, err := uc.memberRepo.GetSpaceMember(ctx, srcSpaceId, uid)
		if srcMember == nil || err != nil {
			err := errs.NoPerm(ctx)
			return 0, err
		}

		err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
			SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(srcMember),
			Perm:              consts.PERM_ModifySpaceWorkFlow,
		})
		if err != nil {
			return 0, errs.NoPerm(ctx)
		}

		flowId, err := uc.CopyWorkFlowFromOtherSpace(ctx, oper, spaceId, workFlowId, flowNewName, flowDefaultStatus)
		if err != nil {
			return 0, err
		}

		uc.resetRank(ctx, spaceId, flowId)

		return flowId, nil
	}

	result, err := uc.wfService.CopyWorkFlow(ctx, workFlow, flowNewName, wf_domain.WorkFlowStatus(flowDefaultStatus), oper)
	if err != nil {
		return 0, errs.Business(ctx, err.Error())
	}

	//持久化
	txErr := uc.tm.InTx(ctx, func(ctx context.Context) error {
		uc.repo.CreateWorkFlow(ctx, result.WorkFlow)
		uc.repo.CreateWorkFlowTemplate(ctx, result.WorkFlowTemplate)
		return nil
	})

	if txErr != nil {
		return 0, errs.Internal(ctx, err)
	}

	uc.domainMessageProducer.Send(ctx, result.WorkFlow.GetMessages())

	uc.resetRank(ctx, spaceId, result.WorkFlow.Id)

	return result.WorkFlow.Id, nil
}

func (uc *WorkFlowUsecase) CopyWorkFlowFromOtherSpace(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workFlowId int64, flowNewName string, flowDefaultStatus int32) (int64, error) {
	uid := oper.UserId

	srcFlow, err := uc.repo.GetWorkFlow(ctx, workFlowId)
	if err != nil {
		return 0, errs.Internal(ctx, err)
	}

	srcSpaceId := srcFlow.SpaceId
	srcSpace, err := uc.spaceRepo.GetSpace(ctx, srcSpaceId)
	if err != nil {
		return 0, errs.Internal(ctx, err)
	}

	srcTemplate, err := uc.repo.GetFlowTemplate(ctx, srcFlow.LastTemplateId)
	if err != nil {
		return 0, errs.Internal(ctx, err)
	}

	typeInfoQueryResult, err := uc.wItemTypeRepo.QWorkItemTypeInfo(ctx, query.WorkItemTypeInfoQuery{SpaceId: spaceId})
	if err != nil {
		return 0, errs.Internal(ctx, err)
	}

	var dstWorkItemType *witem_type.WorkItemType
	switch srcFlow.FlowMode {
	case consts.FlowMode_WorkFlow:
		dstWorkItemType = typeInfoQueryResult.GetWorkFlowTaskType()
	case consts.FlowMode_StateFlow:
		dstWorkItemType = typeInfoQueryResult.GetStateFlowTaskType()
	}

	newFlowId := uc.idService.NewId(ctx, consts.BusinessId_Type_WorkFlow).Id
	newTemplateId := uc.idService.NewId(ctx, consts.BusinessId_Type_WorkFlowTemplate).Id

	newTemplate := wf_domain.NewWorkFlowTemplate(
		newTemplateId,
		spaceId,
		dstWorkItemType.Id,
		newFlowId,
		1,
		srcTemplate.FlowMode,
		srcTemplate.WorkFLowConfig,
		srcTemplate.StateFlowConfig,
		wf_domain.WorkFlowTemplateStatus_Enable,
		uid,
		oper,
	)

	if flowNewName == "" {
		flowNewName = srcFlow.Name
	}

	newFlow := wf_domain.NewWorkFlow(
		newFlowId,
		spaceId,
		dstWorkItemType.Id,
		flowNewName,
		srcFlow.Key,
		srcFlow.Ranking,
		srcFlow.FlowMode,
		newTemplate,
		wf_domain.WorkFlowStatus(flowDefaultStatus),
		0,
		uid,
		oper,
	)

	var templateRoleIds []int64
	var templateStatusIds []int64

	switch srcFlow.FlowMode {
	case consts.FlowMode_WorkFlow:
		conf := newTemplate.WorkFlowConf()
		conf.Uuid = uuid.New().String()
		for _, node := range conf.Nodes {
			for _, v := range node.Owner.OwnerRole {
				templateRoleIds = append(templateRoleIds, cast.ToInt64(v.Id))
			}

			for _, v := range node.OnPass {
				templateStatusIds = append(templateStatusIds, cast.ToInt64(v.TargetSubState.Id))
			}

			for _, v := range node.OnReach {
				templateStatusIds = append(templateStatusIds, cast.ToInt64(v.TargetSubState.Id))
			}
		}
	case consts.FlowMode_StateFlow:
		conf := newTemplate.StateFlowConf()
		for _, node := range conf.StateFlowNodes {
			for _, v := range node.Owner.OwnerRole {
				templateRoleIds = append(templateRoleIds, cast.ToInt64(v.Id))
			}
			templateStatusIds = append(templateStatusIds, cast.ToInt64(node.SubStateId))
		}
	}

	templateRoleIds = stream.Unique(templateRoleIds)
	templateStatusIds = stream.Unique(templateStatusIds)

	mapRole, newRoles, err := uc.mapRole(ctx, uid, srcSpaceId, spaceId, templateRoleIds)
	if err != nil {
		return 0, errs.Internal(ctx, err)
	}

	mapStatus, newStatus, err := uc.mapStatus(ctx, uid, srcSpaceId, spaceId, templateStatusIds)
	if err != nil {
		return 0, errs.Internal(ctx, err)
	}

	switch srcFlow.FlowMode {
	case consts.FlowMode_WorkFlow:
		conf := newTemplate.WorkFlowConf()
		for _, node := range conf.Nodes {
			node.Owner.UsageMode = tplt_config.UsageMode_None
			node.Owner.Value = nil

			for _, v := range node.Owner.OwnerRole {
				if e := mapRole[cast.ToInt64(v.Id)]; e != nil {
					v.Id = cast.ToString(e.Id)
					v.Key = e.Key
				}
			}

			for _, v := range node.OnPass {
				if e := mapStatus[cast.ToInt64(v.TargetSubState.Id)]; e != nil {
					v.TargetSubState.Id = cast.ToString(e.Id)
					v.TargetSubState.Key = e.Key
					v.TargetSubState.Val = e.Val
				}
			}

			for _, v := range node.OnReach {
				if e := mapStatus[cast.ToInt64(v.TargetSubState.Id)]; e != nil {
					v.TargetSubState.Id = cast.ToString(e.Id)
					v.TargetSubState.Key = e.Key
					v.TargetSubState.Val = e.Val
				}
			}
		}
	case consts.FlowMode_StateFlow:
		conf := newTemplate.StateFlowConf()
		for _, node := range conf.StateFlowNodes {
			node.Owner.UsageMode = tplt_config.UsageMode_None
			node.Owner.Value = nil

			for _, v := range node.Owner.OwnerRole {
				if e := mapRole[cast.ToInt64(v.Id)]; e != nil {
					v.Id = cast.ToString(e.Id)
					v.Key = e.Key
				}
			}

			oldKey := node.Key

			status := mapStatus[cast.ToInt64(node.SubStateId)]
			node.Key = status.Key
			node.Name = status.Name
			node.SubStateId = cast.ToString(status.Id)
			node.SubStateKey = status.Key
			node.SubStateVal = status.Val

			// 修改连线
			for _, connection := range conf.StateFlowTransitionRule {
				if oldKey == connection.SourceStateKey {
					connection.SourceStateKey = status.Key
				}
				if oldKey == connection.TargetStateKey {
					connection.TargetStateKey = status.Key
				}
			}
		}
	}

	err = uc.tm.InTx(ctx, func(ctx context.Context) error {
		err := uc.repo.CreateWorkFlowTemplate(ctx, newTemplate)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = uc.repo.CreateWorkFlow(ctx, newFlow)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = uc.wItemRoleRepo.CreateWorkItemRoles(ctx, newRoles)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = uc.wItemStatusRepo.CreateWorkItemStatusItems(ctx, spaceId, newStatus)
		if err != nil {
			return errs.Internal(ctx, err)
		}
		return nil
	})
	if err != nil {
		return 0, errs.Internal(ctx, err)
	}

	uc.domainMessageProducer.Send(ctx, newRoles.GetMessages())
	uc.domainMessageProducer.Send(ctx, newStatus.GetMessages())

	msg := &domain_message.CreateWorkFlow{
		DomainMessageBase: shared.DomainMessageBase{
			Oper:     oper,
			OperTime: time.Now(),
		},
		SpaceId:         spaceId,
		WorkFlowName:    newFlow.Name,
		WorkFlowId:      newFlow.Id,
		SrcSpaceName:    srcSpace.SpaceName,
		SrcWorkFlowName: srcFlow.Name,
		FlowMode:        newFlow.FlowMode,
	}
	uc.domainMessageProducer.Send(ctx, []shared.DomainMessage{msg})

	return newFlowId, nil
}

func (uc *WorkFlowUsecase) mapRole(ctx context.Context, uid, srcSpaceId, dstSpaceId int64, ids []int64) (m map[int64]*work_item_role.WorkItemRole, newRoles work_item_role.WorkItemRoles, err error) {
	srcMap, err := uc.wItemRoleRepo.WorkItemRoleMap(ctx, srcSpaceId)
	if err != nil {
		return nil, nil, errs.Internal(ctx, err)
	}

	dstMap, err := uc.wItemRoleRepo.WorkItemRoleMap(ctx, dstSpaceId)
	if err != nil {
		return nil, nil, errs.Internal(ctx, err)
	}

	dstKeyMap := stream.MapKV(dstMap, func(k int64, v *work_item_role.WorkItemRole) (string, *work_item_role.WorkItemRole) {
		return v.Key, v
	})

	var newItemList []*work_item_role.WorkItemRole

	mapping := map[int64]*work_item_role.WorkItemRole{}

	for _, id := range ids {
		if v := srcMap[id]; v != nil {
			// 是预设 && key相同 && 名称相同
			if v.IsSys == 1 && dstKeyMap[v.Key] != nil && v.Name == dstKeyMap[v.Key].Name {
				mapping[id] = dstKeyMap[v.Key]
			} else {
				// 是自定义

				key := v.Key
				for {
					if _, ok := dstKeyMap[key]; !ok {
						break
					}

					key = "role_" + rand.Letters(5)
				}

				newItem := work_item_role.NewWorkItemRole(
					uc.idService.NewId(ctx, consts.BusinessId_Type_WorkItemRole).Id,
					dstSpaceId,
					0,
					v.Name,
					key,
					time.Now().Unix()+v.Ranking,
					0,
					uid,
					v.FlowScope,
					utils.GetLoginUser(ctx),
				)

				newItemList = append(newItemList, newItem)
				mapping[id] = newItem
			}
		}
	}

	return mapping, newItemList, nil
}

func (uc *WorkFlowUsecase) mapStatus(ctx context.Context, uid, srcSpaceId, dstSpaceId int64, ids []int64) (m map[int64]*work_item_status.WorkItemStatusItem, newItems work_item_status.WorkItemStatusItems, err error) {
	srcMap, err := uc.wItemStatusRepo.WorkItemStatusMap(ctx, srcSpaceId)
	if err != nil {
		return nil, nil, errs.Internal(ctx, err)
	}

	dstMap, err := uc.wItemStatusRepo.WorkItemStatusMap(ctx, dstSpaceId)
	if err != nil {
		return nil, nil, errs.Internal(ctx, err)
	}

	dstKeyMap := stream.MapKV(dstMap, func(k int64, v *work_item_status.WorkItemStatusItem) (string, *work_item_status.WorkItemStatusItem) {
		return v.Key, v
	})

	var newItemList []*work_item_status.WorkItemStatusItem

	mapping := map[int64]*work_item_status.WorkItemStatusItem{}

	for _, id := range ids {
		if v := srcMap[id]; v != nil {
			// 是预设 && key相同 && 名称相同
			if v.IsSys == 1 && dstKeyMap[v.Key] != nil && v.Name == dstKeyMap[v.Key].Name {
				mapping[id] = dstKeyMap[v.Key]
				continue
			} else {
				key := v.Key
				for {
					if _, ok := dstKeyMap[key]; !ok {
						break
					}

					key = "status_" + rand.Letters(5)
				}

				newItem := work_item_status.NewWorkItemStatusItem(
					uc.idService.NewId(ctx, consts.BusinessId_Type_WorkItemStatus).Id,
					dstSpaceId,
					v.Name,
					key,
					key,
					0,
					time.Now().Unix()+v.Ranking,
					v.StatusType,
					uid,
					v.FlowScope,
					utils.GetLoginUser(ctx),
				)

				newItemList = append(newItemList, newItem)
				mapping[id] = newItem
			}
		}
	}

	return mapping, newItemList, nil
}

func (uc *WorkFlowUsecase) QSpaceWorkFlowList(ctx context.Context, oper *utils.LoginUserInfo, req *pb.SpaceWorkFlowListRequest) (*pb.SpaceWorkFlowListReplyReplyData, error) {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, req.SpaceId, uid)
	if member == nil || err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.NoPerm(ctx)
		return nil, err
	}

	typeInfoQueryResult, err := uc.wItemTypeRepo.QWorkItemTypeInfo(ctx, query.WorkItemTypeInfoQuery{SpaceId: req.SpaceId})
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	workFlows, err := uc.repo.QTaskWorkFlowList(ctx, req.SpaceId, typeInfoQueryResult.GetMainTaskTypeIds(), consts.WorkFlowMode(req.FlowMode))
	if err != nil {
		return nil, err
	}

	var list []*pb.SpaceWorkFlowListReplyReplyData_Item
	for _, flow := range workFlows {
		item := &pb.SpaceWorkFlowListReplyReplyData_Item{
			Id:             flow.Id,
			Name:           flow.Name,
			Version:        int64(flow.Version),
			TemplateId:     flow.LastTemplateId,
			SpaceId:        flow.SpaceId,
			CreatedAt:      flow.CreatedAt,
			UpdatedAt:      flow.UpdatedAt,
			WorkItemTypeId: flow.WorkItemTypeId,
			Ranking:        flow.Ranking,
			Key:            flow.Key,
			Status:         int64(flow.Status),
			FlowMode:       string(flow.FlowMode),
		}
		list = append(list, item)
	}

	return &pb.SpaceWorkFlowListReplyReplyData{
		Total: cast.ToInt32(len(list)),
		List:  list,
	}, nil
}

func (uc *WorkFlowUsecase) QSpaceWorkFlowById(ctx context.Context, oper *utils.LoginUserInfo, req *pb.SpaceWorkFlowByIdRequest) (*pb.SpaceWorkFlowByIdReplyData, error) {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, req.SpaceId, uid)
	if member == nil || err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.NoPerm(ctx)
		return nil, err
	}

	workFlows, err := uc.repo.QTaskWorkFlowById(ctx, req.SpaceId, req.Ids)
	if err != nil {
		return nil, err
	}

	var list []*pb.SpaceWorkFlowByIdReplyData_Item
	for _, flow := range workFlows {
		item := &pb.SpaceWorkFlowByIdReplyData_Item{
			Id:             flow.Id,
			Name:           flow.Name,
			Version:        int64(flow.Version),
			TemplateId:     flow.LastTemplateId,
			SpaceId:        flow.SpaceId,
			CreatedAt:      flow.CreatedAt,
			UpdatedAt:      flow.UpdatedAt,
			WorkItemTypeId: flow.WorkItemTypeId,
			Key:            flow.Key,
		}
		list = append(list, item)
	}

	return &pb.SpaceWorkFlowByIdReplyData{
		Total: cast.ToInt32(len(list)),
		List:  make([]*pb.SpaceWorkFlowByIdReplyData_Item, 0),
	}, nil
}

func (uc *WorkFlowUsecase) QSpaceWorkItemRelationCount(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, flowId int64, scene string) (int64, error) {

	// 判断当前用户是否在要查询的项目空间内
	_, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if err != nil {
		return 0, errs.NoPerm(ctx)
	}

	wf, err := uc.repo.GetWorkFlow(ctx, flowId)
	if err != nil {
		return 0, errs.Internal(ctx, err)
	}

	if scene == "status_relation" {
		statusInfo, err := uc.wItemStatusRepo.GetWorkItemStatusInfo(ctx, spaceId)
		if err != nil {
			return 0, errs.Internal(ctx, err)
		}

		flowScope := consts.ConvertFlowModeToFlowScope(wf.FlowMode)

		exStatus := make([]string, 0)
		archivedStatusItems := statusInfo.GetArchivedTypeItems(flowScope)
		if archivedStatusItems != nil {
			exStatus = archivedStatusItems.GetStatusKeys()
		}

		totalNum, err := uc.wItemRepo.CountWorkFlowRelatedSpaceWorkItem(ctx, spaceId, flowId, exStatus)
		return totalNum, err
	}

	totalNum, err := uc.wItemRepo.CountWorkFlowRelatedSpaceWorkItem(ctx, spaceId, flowId, nil)
	return totalNum, err

}

func (uc *WorkFlowUsecase) GetOwnerRuleRelationTemplateRequest(ctx context.Context, oper *utils.LoginUserInfo, req *pb.GetOwnerRuleRelationTemplateRequest) (*pb.GetOwnerRuleRelationTemplateReplyData, error) {

	// 判断当前用户是否在要查询的项目空间内
	_, err := uc.memberRepo.GetSpaceMember(ctx, req.SpaceId, oper.UserId)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	result := &pb.GetOwnerRuleRelationTemplateReplyData{
		List: make([]*pb.GetOwnerRuleRelationTemplateReplyData_Template, 0),
	}

	tpltIds, err := uc.wfService.FindWorkFlowTemplateByAppointedOwnerRule(ctx, req.SpaceId, req.OwnerUid)
	for _, v := range tpltIds {
		tplt, err := uc.repo.GetWorkFlowTemplateFormMemoryCache(ctx, v)
		if err != nil {
			return nil, errs.Internal(ctx, err)
		}

		flow, err := uc.repo.GetWorkFlow(ctx, tplt.WorkFlowId)
		if err != nil {
			return nil, errs.Internal(ctx, err)
		}

		result.List = append(result.List, &pb.GetOwnerRuleRelationTemplateReplyData_Template{
			Id:   flow.Id,
			Name: flow.Name,
		})
	}

	return result, err
}
