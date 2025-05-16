package biz

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	v1 "go-cs/api/search/v1"
	"go-cs/internal/bean/rsp_convert"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/consts"
	checker "go-cs/internal/domain/perm"
	perm_service "go-cs/internal/domain/perm/service"
	"go-cs/internal/domain/pkg/flow_simulator"
	"go-cs/internal/domain/search/condition_translater"
	"go-cs/internal/domain/search/search_es"
	"go-cs/internal/domain/work_flow"
	config "go-cs/internal/domain/work_flow/flow_tplt_config"
	witem_status "go-cs/internal/domain/work_item_status"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/pkg/stream"
	"google.golang.org/protobuf/types/known/anypb"
	"slices"

	search_repo "go-cs/internal/domain/search/repo"
	space_repo "go-cs/internal/domain/space/repo"
	member "go-cs/internal/domain/space_member"
	member_repo "go-cs/internal/domain/space_member/repo"
	tag_repo "go-cs/internal/domain/space_tag/repo"
	comment_repo "go-cs/internal/domain/space_work_item_comment/repo"
	workObj_repo "go-cs/internal/domain/space_work_object/repo"
	workVersion_repo "go-cs/internal/domain/space_work_version/repo"
	user_repo "go-cs/internal/domain/user/repo"
	wf_repo "go-cs/internal/domain/work_flow/repo"
	witem "go-cs/internal/domain/work_item"
	witem_repo "go-cs/internal/domain/work_item/repo"
	status_repo "go-cs/internal/domain/work_item_status/repo"
	witem_type_repo "go-cs/internal/domain/work_item_type/repo"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type SearchUsecase struct {
	tm  trans.Transaction
	log *log.Helper

	repo                search_repo.SearchRepo
	userRepo            user_repo.UserRepo
	spaceRepo           space_repo.SpaceRepo
	spaceMemberRepo     member_repo.SpaceMemberRepo
	spaceWorkItemRepo   witem_repo.WorkItemRepo
	workItemStatusRepo  status_repo.WorkItemStatusRepo
	spaceWorkObjectRepo workObj_repo.SpaceWorkObjectRepo
	spaceTagRepo        tag_repo.SpaceTagRepo
	commentRepo         comment_repo.SpaceWorkItemCommentRepo
	workVersionRepo     workVersion_repo.SpaceWorkVersionRepo
	workFlowRepo        wf_repo.WorkFlowRepo
	workItemTypeRepo    witem_type_repo.WorkItemTypeRepo

	permService *perm_service.PermService
}

func NewSearchUsecase(
	tm trans.Transaction,
	logger log.Logger,

	repo search_repo.SearchRepo,
	userRepo user_repo.UserRepo,
	spaceRepo space_repo.SpaceRepo,
	spaceMemberRepo member_repo.SpaceMemberRepo,
	spaceWorkItemRepo witem_repo.WorkItemRepo,
	spaceWorkObjectRepo workObj_repo.SpaceWorkObjectRepo,
	spaceTagRepo tag_repo.SpaceTagRepo,
	commentRepo comment_repo.SpaceWorkItemCommentRepo,
	workVersionRepo workVersion_repo.SpaceWorkVersionRepo,
	workFlowRepo wf_repo.WorkFlowRepo,
	workItemTypeRepo witem_type_repo.WorkItemTypeRepo,
	workItemStatusRepo status_repo.WorkItemStatusRepo,

	permService *perm_service.PermService,

) *SearchUsecase {
	moduleName := "SearchUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SearchUsecase{
		tm:                  tm,
		log:                 hlog,
		repo:                repo,
		userRepo:            userRepo,
		spaceRepo:           spaceRepo,
		spaceMemberRepo:     spaceMemberRepo,
		spaceWorkItemRepo:   spaceWorkItemRepo,
		spaceWorkObjectRepo: spaceWorkObjectRepo,
		spaceTagRepo:        spaceTagRepo,
		commentRepo:         commentRepo,
		workVersionRepo:     workVersionRepo,
		workFlowRepo:        workFlowRepo,
		workItemTypeRepo:    workItemTypeRepo,
		workItemStatusRepo:  workItemStatusRepo,
		permService:         permService,
	}
}

func (uc *SearchUsecase) SearchGroupInfo(ctx context.Context, uid int64, req *v1.SearchSpaceWorkItemGroupInfoRequestV2) (*v1.SearchSpaceWorkItemGroupInfoReplyV2Data, error) {
	if req.SpaceId != 0 {
		_, err := uc.spaceMemberRepo.GetSpaceMember(ctx, req.SpaceId, uid)
		if err != nil {
			return nil, errs.NoPerm(ctx)
		}
	}

	if err := search_es.CheckSortField(req.Sorts); err != nil {
		return nil, errs.Param(ctx, err.Error())
	}

	if err := search_es.CheckGroupField(req.Groups); err != nil {
		return nil, errs.Param(ctx, err.Error())
	}

	// 是否是工作台的筛选
	isWorkBench := req.SpaceId == 0

	// 没传spaceId，则获取用户全部的spaceId
	var spaceIds = []int64{req.SpaceId}
	var err error
	if req.SpaceId == 0 {
		spaceIds, err = uc.spaceRepo.GetUserSpaceIds(ctx, uid)
		if err != nil {
			return nil, err
		}
	}

	// 翻译条件
	translateCtx := condition_translater.Ctx{
		Ctx:         ctx,
		StatusRepo:  uc.workItemStatusRepo,
		SpaceIds:    spaceIds,
		IsWorkBench: isWorkBench,
	}
	req.ConditionGroup = translateCtx.TranslateCond(req.ConditionGroup)

	// 构建分组排序
	var groupSorts []Sort
	for _, group := range req.Groups {
		groupSorts = append(groupSorts, Sort{
			Field: group.Field,
			Order: "DESC",
		})
	}

	var paramSorts []Sort
	for _, v := range req.Sorts {
		paramSorts = append(paramSorts, Sort{
			Field: v.Field,
			Order: v.Order,
		})
	}

	// 调整分组排序
	for _, v := range paramSorts {
		idx := slices.IndexFunc(groupSorts, func(s Sort) bool {
			return v.Field == s.Field
		})
		if idx >= 0 {
			groupSorts[idx] = v
		}
	}

	var finalSorts []Sort
	var defaultSorts []Sort
	if isWorkBench {
		defaultSorts = []Sort{
			{Field: "space_id", Order: "DESC"},
			{Field: "work_item_id", Order: "DESC"},
		}
	} else {
		defaultSorts = []Sort{
			{Field: "work_item_id", Order: "DESC"},
		}
	}
	finalSorts = stream.Diff(paramSorts, groupSorts)

	for _, v := range defaultSorts {
		f := func(s Sort) bool {
			return v.Field == s.Field
		}

		if slices.ContainsFunc(groupSorts, f) {
			continue
		}

		if !slices.ContainsFunc(finalSorts, f) {
			finalSorts = append(finalSorts, v)
		}
	}

	//data, err := uc.SearchGroupInfoByDB(ctx, uid, req, finalSorts, groupSorts)
	data, err := uc.SearchGroupInfoByEs(ctx, uid, spaceIds, req, finalSorts, groupSorts)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return &v1.SearchSpaceWorkItemGroupInfoReplyV2Data{
			List:     nil,
			TotalNum: 0,
		}, nil
	}

	return data, nil
}

func (uc *SearchUsecase) SearchMySpaceWorkItemsByIdV2(ctx context.Context, uid int64, workItemIds []int64) (*v1.SearchMySpaceWorkItemsByIdReplyV2Data, error) {

	workItemIds = stream.Filter(workItemIds, func(e int64) bool {
		return e > 0
	})

	//查询基本信息
	list, err := uc.spaceWorkItemRepo.GetWorkItemByIds(ctx, workItemIds, &witem_repo.WithDocOption{
		PlanTime:      true,
		Priority:      true,
		ProcessRate:   true,
		Directors:     true,
		Participators: true,
		Followers:     true,
	}, &witem_repo.WithOption{
		FlowNodes: true,
	})
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	itemMap := stream.ToMap(list, func(i int, v *witem.WorkItem) (int64, *witem.WorkItem) {
		return v.Id, v
	})

	var pids []int64
	for _, v := range list {
		pid := v.Pid
		if pid != 0 && !slices.Contains(workItemIds, pid) { //排除已存在的
			pids = append(pids, v.Pid)
		}
	}

	parentWorkItems, err := uc.spaceWorkItemRepo.GetWorkItemByIds(ctx, pids, nil, nil)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	parentWorkItemMap := stream.ToMap(parentWorkItems, func(i int, v *witem.WorkItem) (int64, *witem.WorkItem) {
		return v.Id, v
	})

	itemMap = stream.MergeMap(parentWorkItemMap, itemMap)

	if len(list) == 0 {
		return nil, nil
	}

	// 获取任务的参与人
	participatorMap := stream.ToMap(list, func(_ int, v *witem.WorkItem) (int64, []int64) {
		return v.Id, utils.ToInt64Array(v.Doc.Participators)
	})

	var (
		allUserIds  []int64
		spaceIds    []int64
		versionIds  []int64
		objectIds   []int64
		flowIds     []int64
		statusIds   []int64
		templateIds []int64
	)
	for _, v := range list {
		allUserIds = append(allUserIds, participatorMap[v.Id]...)
		allUserIds = append(allUserIds, v.UserId, v.Restart.RestartUserId)
		spaceIds = append(spaceIds, v.SpaceId)
		versionIds = append(versionIds, v.VersionId)
		objectIds = append(objectIds, v.WorkObjectId)
		flowIds = append(flowIds, v.WorkItemFlowId)
		statusIds = append(statusIds, v.WorkItemStatus.Id)
		templateIds = append(templateIds, v.WorkFlowTemplateId)
	}

	for _, v := range parentWorkItems {
		statusIds = append(statusIds, v.WorkItemStatus.Id)
	}

	spaceIds = stream.Unique(spaceIds)
	objectIds = stream.Unique(objectIds)
	allUserIds = stream.Unique(allUserIds)
	flowIds = stream.Unique(flowIds)
	statusIds = stream.Unique(statusIds)
	templateIds = stream.Unique(templateIds)

	memberMap, err := uc.spaceMemberRepo.UserSpaceMemberMap(ctx, uid, spaceIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}
	if len(memberMap) != len(spaceIds) {
		return nil, errs.Business(ctx, "不能查询未参与的空间")
	}

	userMap, err := uc.userRepo.UserMap(ctx, allUserIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	spaceMap, err := uc.spaceRepo.SpaceMap(ctx, spaceIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	objectMap, err := uc.spaceWorkObjectRepo.SpaceWorkObjectMapByObjectIds(ctx, objectIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	versionMap, err := uc.workVersionRepo.SpaceWorkVersionMapByVersionIds(ctx, versionIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	flowMap, err := uc.workFlowRepo.WorkFlowMap(ctx, flowIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	statusMap, err := uc.workItemStatusRepo.StatusMapBySpaceIds(ctx, spaceIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	templateMap, err := uc.workFlowRepo.FlowTemplateMap(ctx, templateIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	// 流程节点排序值
	flowNodeRankingMap := stream.MapValue(templateMap, func(v *work_flow.WorkFlowTemplate) map[string]int {
		if v.IsWorkFlow() {
			nodes := v.WorkFlowConf().Nodes
			return stream.ToMap(nodes, func(i int, v *config.WorkFlowNode) (string, int) {
				return v.Key, i
			})
		}

		return nil
	})

	var curFlowNodeMap = map[int64]*witem.WorkItemFlowNode{}
	var curFlowNodeConfMap = map[int64]*anypb.Any{}

	for _, v := range list {
		template, _ := templateMap[v.WorkFlowTemplateId]
		switch {
		case v.IsWorkFlowMainTask():
			processingNodes := v.WorkItemFlowNodes.GetProcessingNodes()
			if len(processingNodes) == 0 {
				continue
			}
			if template != nil && template.WorkFlowConf() != nil {
				graph, _ := flow_simulator.NewWorkFlowGraph(v, template)
				nodeCode := graph.FindRightFirstInProcessNode()

				flowNode := graph.GetNodeInfo(nodeCode)
				if flowNode == nil {
					continue
				}

				nodeConfig := graph.GetNodeConfig(nodeCode)
				if nodeConfig == nil {
					continue
				}

				nodeToRsp := rsp_convert.WorkFlowNodeToRsp(nodeConfig)
				nodeToRsp.StartAt = flowNode.StartAt
				anyVal, _ := anypb.New(nodeToRsp)

				curFlowNodeMap[v.Id] = flowNode
				curFlowNodeConfMap[v.Id] = anyVal
			}
		case v.IsStateFlowMainTask():
			if template != nil && template.StateFlowConf() != nil {
				flowNode := v.WorkItemFlowNodes.GetNodeByCode(v.WorkItemStatus.Key)

				nodeConfig := template.StateFlowConf().GetNode(v.WorkItemStatus.Key)
				if nodeConfig == nil {
					continue
				}

				nodeToRsp := rsp_convert.StateFlowNodeToRsp(nodeConfig)
				nodeToRsp.StartAt = flowNode.StartAt
				anyVal, _ := anypb.New(nodeToRsp)

				curFlowNodeMap[v.Id] = flowNode
				curFlowNodeConfMap[v.Id] = anyVal
			}
		}
	}

	retUserMap := map[int64]*rsp.SimpleUserInfo{}
	for id, v := range userMap {
		retUserMap[id] = &rsp.SimpleUserInfo{
			Id:           v.Id,
			UserId:       v.Id,
			UserName:     v.UserName,
			UserNickname: v.UserNickname,
			Avatar:       v.Avatar,
		}
	}
	roleMap := stream.MapValue(memberMap, func(v *member.SpaceMember) int64 {
		return v.RoleId
	})

	// 过滤当前出用户参与的任务
	curUserParticipate := stream.Filter(workItemIds, func(id int64) bool {
		return slices.Contains(participatorMap[id], uid)
	})
	unreadNumMap, _ := uc.commentRepo.UserUnreadNumMapByWorkItemIds(ctx, uid, curUserParticipate)

	// 进行中的节点
	nodeMap, err := uc.spaceWorkItemRepo.ProcessingNodeMapByWorkItemIds(ctx, workItemIds)
	if err != nil {
		return nil, err
	}

	nodesMap, err := uc.getNodes(uid, itemMap, templateMap, nodeMap, statusMap, flowNodeRankingMap)

	var items []*rsp.SearchWorkItemInfo
	for _, v := range list {
		var parentStatus *rsp.SearchWorkItemInfo_StatusInfo
		if v.Pid != 0 && itemMap[v.Pid] != nil {
			status := statusMap[itemMap[v.Pid].WorkItemStatus.Id]
			parentStatus = &rsp.SearchWorkItemInfo_StatusInfo{
				Id:   status.Id,
				Key:  status.Key,
				Val:  status.Val,
				Name: status.Name,
				Type: int64(status.StatusType),
			}
		}

		var retDirectors []*rsp.SimpleUserInfo
		for _, v := range v.Doc.Directors {
			u := retUserMap[cast.ToInt64(v)]
			if u != nil {
				retDirectors = append(retDirectors, u)
			}
		}

		var retParticipators []*rsp.SimpleUserInfo
		for _, v := range participatorMap[v.Id] {
			u := retUserMap[cast.ToInt64(v)]
			if u != nil {
				retParticipators = append(retParticipators, u)
			}
		}

		// 任务表单细分权限
		workItemRole := ""
		if v.UserId == uid {
			workItemRole = consts.WORK_ITEM_ROLE_CREATOR
		} else {
			if curFlowNodeMap[v.Id] != nil && slices.Contains(curFlowNodeMap[v.Id].Directors, cast.ToString(uid)) {
				workItemRole = consts.WORK_ITEM_ROLE_NODE_OWNER
			}
		}

		//对应父级或子级任务对应的权限
		var permSubFuncRoleJsonV []byte
		if v.Pid == 0 {
			permSubFuncRole := uc.permService.GetPermissionWithScene(roleMap[v.SpaceId], workItemRole, checker.EditWorkItemScene)
			permSubFuncRoleJsonV, _ = json.Marshal(permSubFuncRole)
		} else {
			permSubFuncRole := uc.permService.GetTaskPermission(roleMap[v.SpaceId], workItemRole)
			permSubFuncRoleJsonV, _ = json.Marshal(permSubFuncRole)
		}

		if flowMap[v.WorkItemFlowId] == nil {
			panic(fmt.Sprintf("flow %d not found %d  %d", v.WorkItemFlowId, v.Pid, v.Id))
		}

		// 当前节点排期
		var nodePlanTime *rsp.SearchWorkItemInfo_PlanTime
		if node, ok := curFlowNodeMap[v.Id]; ok {
			nodePlanTime = &rsp.SearchWorkItemInfo_PlanTime{
				StartAt:    node.PlanTime.StartAt,
				CompleteAt: node.PlanTime.CompleteAt,
			}
		}

		item := &rsp.SearchWorkItemInfo{
			Id:                 v.Id,
			SpaceId:            v.SpaceId,
			UserId:             v.UserId,
			WorkObjectId:       v.WorkObjectId,
			WorkObjectName:     objectMap[v.WorkObjectId].WorkObjectName,
			WorkItemGuid:       v.WorkItemGuid,
			WorkItemName:       v.WorkItemName,
			WorkItemTypeKey:    string(v.WorkItemTypeKey),
			FlowId:             v.WorkItemFlowId,
			FlowName:           flowMap[v.WorkItemFlowId].Name,
			FlowMode:           string(flowMap[v.WorkItemFlowId].FlowMode),
			WorkFlowTemplateId: v.WorkFlowTemplateId,
			WorkItemStatus: &rsp.SearchWorkItemInfo_StatusInfo{
				Id:   v.WorkItemStatus.Id,
				Key:  v.WorkItemStatus.Key,
				Val:  v.WorkItemStatus.Val,
				Name: statusMap[v.WorkItemStatus.Id].Name,
				Type: int64(statusMap[v.WorkItemStatus.Id].StatusType),
			},
			PlanTime: &rsp.SearchWorkItemInfo_PlanTime{
				StartAt:    v.Doc.PlanStartAt,
				CompleteAt: v.Doc.PlanCompleteAt,
			},
			ProcessRate:          v.Doc.ProcessRate,
			CreatedAt:            v.CreatedAt,
			UpdatedAt:            v.UpdatedAt,
			Priority:             v.Doc.Priority,
			Pid:                  v.Pid,
			SpaceName:            spaceMap[v.SpaceId].SpaceName,
			LastStatusAt:         v.LastWorkItemStatus.LastAt,
			IsRestart:            v.Restart.IsRestart,
			RestartAt:            v.Restart.RestartAt,
			IconFlags:            v.IconFlags.ToFlags(),
			SpaceRoleId:          roleMap[v.SpaceId],
			RestartUserId:        v.Restart.RestartUserId,
			UnreadCommentNum:     int32(unreadNumMap[v.Id]),
			CommentNum:           v.CommentNum,
			ResumeAt:             v.Resume.ResumeAt,
			ParentStatus:         parentStatus,
			Followed:             v.Doc.Followers.Contains(cast.ToString(uid)),
			VersionId:            v.VersionId,
			VersionName:          versionMap[v.VersionId].VersionName,
			Directors:            retDirectors,
			Creator:              retUserMap[v.UserId],
			RestartUser:          retUserMap[v.Restart.RestartUserId],
			Participators:        retParticipators,
			OperationPermissions: string(permSubFuncRoleJsonV),
			FlowNode:             curFlowNodeConfMap[v.Id],
			NodePlanTime:         nodePlanTime,
			Nodes:                nodesMap[v.Id],
		}

		items = append(items, item)
	}

	return &v1.SearchMySpaceWorkItemsByIdReplyV2Data{
		List:     items,
		TotalNum: int32(len(items)),
	}, nil
}

func (uc *SearchUsecase) getNodes(
	uid int64,
	itemMap map[int64]*witem.WorkItem,
	templateMap map[int64]*work_flow.WorkFlowTemplate,
	nodeMap map[int64][]*witem.WorkItemFlowNode,
	statusMap map[int64]*witem_status.WorkItemStatusItem,
	flowNodeRankingMap map[int64]map[string]int,
) (map[int64][]*rsp.FlowNode, error) {
	// 流程节点名
	flowNodeNameMap := stream.MapValue(templateMap, func(v *work_flow.WorkFlowTemplate) map[string]string {
		switch {
		case v.IsStateFlow():
			nodes := v.StateFlowConf().StateFlowNodes
			return stream.ToMap(nodes, func(i int, v *config.StateFlowNode) (string, string) {
				return v.Key, v.Name
			})
		case v.IsWorkFlow():
			nodes := v.WorkFlowConf().Nodes
			return stream.ToMap(nodes, func(i int, v *config.WorkFlowNode) (string, string) {
				return v.Key, v.Name
			})
		}

		return nil
	})

	// 对进行中对流程任务节点进行排序
	sortByNodeRanking := func(a, b *witem.WorkItemFlowNode) int {
		return cmp.Compare(flowNodeRankingMap[a.FlowTemplateId][a.FlowNodeCode], flowNodeRankingMap[b.FlowTemplateId][b.FlowNodeCode])
	}

	// 如果是节点负责人
	sortByIsDirector := func(a, b *witem.WorkItemFlowNode) int {
		var aVal int
		var bVal int
		if a.ContainsDirector(uid) {
			aVal = -1
		}
		if b.ContainsDirector(uid) {
			bVal = -1
		}
		return cmp.Compare(aVal, bVal)
	}

	finalSortFunc := stream.MergeCmpFunc(sortByIsDirector, sortByNodeRanking)
	nodeMap = stream.MapKV(nodeMap, func(id int64, v []*witem.WorkItemFlowNode) (int64, []*witem.WorkItemFlowNode) {
		workItem := itemMap[id]
		template := templateMap[workItem.WorkFlowTemplateId]

		if workItem.IsWorkFlowMainTask() {
			template.WorkFlowConf()
			slices.SortStableFunc(v, finalSortFunc)
		}

		return id, v
	})

	nodesMap := stream.MapKV(nodeMap, func(id int64, v []*witem.WorkItemFlowNode) (int64, []*rsp.FlowNode) {
		var ret []*rsp.FlowNode
		for _, v := range v {
			ret = append(ret, &rsp.FlowNode{
				FlowNodeName: flowNodeNameMap[v.FlowTemplateId][v.FlowNodeCode],
				FlowNodeCode: v.FlowNodeCode,
				Directors:    v.Directors,
				PlanTime: &rsp.FlowNode_PlanTime{
					StartAt:    v.PlanTime.StartAt,
					CompleteAt: v.PlanTime.CompleteAt,
				},
			})
		}

		return id, ret
	})

	return nodesMap, nil
}
