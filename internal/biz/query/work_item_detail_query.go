package query

import (
	"context"
	"encoding/json"
	"fmt"
	"go-cs/api/comm"
	"go-cs/internal/bean/vo/rsp"
	file_domain "go-cs/internal/domain/file_info"
	file_repo "go-cs/internal/domain/file_info/repo"
	perm "go-cs/internal/domain/perm"
	perm_service "go-cs/internal/domain/perm/service"
	space_repo "go-cs/internal/domain/space/repo"
	space_file_repo "go-cs/internal/domain/space_file_info/repo"
	space_file_service "go-cs/internal/domain/space_file_info/service"
	member_repo "go-cs/internal/domain/space_member/repo"
	member_service "go-cs/internal/domain/space_member/service"
	"go-cs/internal/domain/space_tag"
	tag_repo "go-cs/internal/domain/space_tag/repo"
	comment_repo "go-cs/internal/domain/space_work_item_comment/repo"
	workObj_repo "go-cs/internal/domain/space_work_object/repo"
	user_repo "go-cs/internal/domain/user/repo"
	wf_repo "go-cs/internal/domain/work_flow/repo"
	"go-cs/internal/domain/work_item_role"
	"go-cs/pkg/stream"
	"slices"

	user_domain "go-cs/internal/domain/user"

	witem_domain "go-cs/internal/domain/work_item"
	witem_repo "go-cs/internal/domain/work_item/repo"
	witem_service "go-cs/internal/domain/work_item/service"

	witem_role_repo "go-cs/internal/domain/work_item_role/repo"

	witem_status_repo "go-cs/internal/domain/work_item_status/repo"
	witem_type_repo "go-cs/internal/domain/work_item_type/repo"

	tag_service "go-cs/internal/domain/space_tag/service"
	workVersion_repo "go-cs/internal/domain/space_work_version/repo"
	statics_repo "go-cs/internal/domain/statics/repo"

	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type WorkItemDetailQuery struct {
	repo                witem_repo.WorkItemRepo
	spaceRepo           space_repo.SpaceRepo
	spaceMemberRepo     member_repo.SpaceMemberRepo
	spaceWorkObjectRepo workObj_repo.SpaceWorkObjectRepo
	workItemRoleRepo    witem_role_repo.WorkItemRoleRepo
	workVersionRepo     workVersion_repo.SpaceWorkVersionRepo
	workItemStatusRepo  witem_status_repo.WorkItemStatusRepo
	workItemTypeRepo    witem_type_repo.WorkItemTypeRepo
	workFlowRepo        wf_repo.WorkFlowRepo
	fileInfoRepo        file_repo.FileInfoRepo
	spaceFileInfoRepo   space_file_repo.SpaceFileInfoRepo
	staticsRepo         statics_repo.StaticsRepo
	tagRepo             tag_repo.SpaceTagRepo
	userRepo            user_repo.UserRepo
	commentRepo         comment_repo.SpaceWorkItemCommentRepo

	permService          *perm_service.PermService
	spaceMemberService   *member_service.SpaceMemberService
	spaceTagService      *tag_service.SpaceTagService
	workItemService      *witem_service.WorkItemService
	spaceFileInfoService *space_file_service.SpaceFileInfoService
	log                  *log.Helper
	tm                   trans.Transaction
}

func NewWorkItemDetailQuery(
	repo witem_repo.WorkItemRepo,
	spaceRepo space_repo.SpaceRepo,
	spaceMemberRepo member_repo.SpaceMemberRepo,
	spaceWorkObjectRepo workObj_repo.SpaceWorkObjectRepo,
	workItemRoleRepo witem_role_repo.WorkItemRoleRepo,
	workVersionRepo workVersion_repo.SpaceWorkVersionRepo,
	workItemStatusRepo witem_status_repo.WorkItemStatusRepo,
	workItemTypeRepo witem_type_repo.WorkItemTypeRepo,
	workFlowRepo wf_repo.WorkFlowRepo,
	fileInfoRepo file_repo.FileInfoRepo,
	staticsRepo statics_repo.StaticsRepo,
	tagRepo tag_repo.SpaceTagRepo,
	userRepo user_repo.UserRepo,
	commentRepo comment_repo.SpaceWorkItemCommentRepo,
	spaceFileInfoRepo space_file_repo.SpaceFileInfoRepo,

	permService *perm_service.PermService,
	spaceMemberService *member_service.SpaceMemberService,
	spaceTagService *tag_service.SpaceTagService,
	workItemService *witem_service.WorkItemService,
	spaceFileInfoService *space_file_service.SpaceFileInfoService,

	logger log.Logger,
	tm trans.Transaction,
) *WorkItemDetailQuery {

	moduleName := "biz.WorkItemDetailQuery"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &WorkItemDetailQuery{
		repo:                repo,
		spaceRepo:           spaceRepo,
		spaceMemberRepo:     spaceMemberRepo,
		spaceWorkObjectRepo: spaceWorkObjectRepo,
		workItemRoleRepo:    workItemRoleRepo,
		workVersionRepo:     workVersionRepo,
		staticsRepo:         staticsRepo,
		tagRepo:             tagRepo,
		userRepo:            userRepo,
		commentRepo:         commentRepo,
		spaceFileInfoRepo:   spaceFileInfoRepo,
		fileInfoRepo:        fileInfoRepo,

		workItemStatusRepo: workItemStatusRepo,
		workItemTypeRepo:   workItemTypeRepo,
		workFlowRepo:       workFlowRepo,
		permService:        permService,
		workItemService:    workItemService,
		spaceMemberService: spaceMemberService,
		spaceTagService:    spaceTagService,

		log: hlog,
		tm:  tm,
	}
}

func (s *WorkItemDetailQuery) Execute(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64) (*rsp.SpaceWorkItemDetail, error) {

	//判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if member == nil || err != nil {
		//成员不存在 不允许操作
		errInfo := errs.New(ctx, comm.ErrorCode_PERMISSION_INSUFFICIENT_DATA_PERMISSIONS)
		return nil, errInfo
	}

	_, err = s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		errInfo := errs.Business(ctx, "项目不存在")
		return nil, errInfo
	}

	workItem, err := s.repo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{
		All: true,
	}, &witem_repo.WithOption{
		FlowNodes: true,
		FlowRoles: true,
	})
	if err != nil {
		return nil, errs.RepoErr(ctx, err)
	}

	//任务状态信息
	statusInfo, _ := s.workItemStatusRepo.GetWorkItemStatusInfo(ctx, workItem.SpaceId)

	//r如果是子任务, 获取父级任务
	var parentWorkItem *witem_domain.WorkItem
	if workItem.IsSubTask() {
		parentWorkItem, err = s.repo.GetWorkItem(ctx, workItem.Pid, nil, nil)
		if err != nil {
			return nil, errs.Internal(ctx, err)
		}
	}

	//工作项信息
	workItemTypeInfo, _ := s.workItemTypeRepo.GetWorkItemType(ctx, workItem.WorkItemTypeId)

	//获取模版信息
	wfInfo, err := s.workFlowRepo.GetWorkFlow(ctx, workItem.WorkItemFlowId)
	if err != nil {
		return nil, errs.Business(ctx, "流程模版不存在")
	}

	wfTplt, err := s.workFlowRepo.GetFlowTemplate(ctx, workItem.WorkFlowTemplateId)
	if err != nil {
		return nil, errs.Business(ctx, "流程模版不存在")
	}

	spaceFiles, _ := s.spaceFileInfoRepo.QFileInfoList(ctx, spaceId, workItemId)
	var _allFileInfoIds []int64
	for _, v := range spaceFiles {
		_allFileInfoIds = append(_allFileInfoIds, v.FileInfo.FileInfoId)
	}

	fileInfos, _ := s.fileInfoRepo.GetFileInfoByIds(ctx, _allFileInfoIds)
	fileCoverMap := stream.ToMap(fileInfos, func(_ int, v *file_domain.FileInfo) (int64, string) {
		return v.Id, v.Cover
	})
	fileMetaMap := stream.ToMap(fileInfos, func(_ int, v *file_domain.FileInfo) (int64, string) {
		return v.Id, v.Meta
	})

	//获取子任务信息
	var subTasks []*witem_domain.WorkItem
	if !workItem.IsSubTask() {
		subTasks, _ = s.repo.GetWorkItemByPid(ctx, workItem.Id, &witem_repo.WithDocOption{
			PlanTime:    true,
			ProcessRate: true,
			Priority:    true,
			Directors:   true,
		}, nil)
	}

	var _allDirectors []string
	var _allTags []string
	var _allCreator []int64
	var _allRestartUser []int64
	var _allFollowers []int64
	{
		_allDirectors = workItem.Doc.Directors
		_allTags = workItem.Doc.Tags
		_allCreator = []int64{workItem.UserId}
		if workItem.Restart.RestartUserId != 0 {
			_allRestartUser = []int64{workItem.Restart.RestartUserId}
		}

		for _, v := range workItem.WorkItemFlowNodes { //添加节点负责人
			_allDirectors = append(_allDirectors, v.Directors...)
		}

		for _, v := range workItem.WorkItemFlowRoles { //添加节点负责人
			_allDirectors = append(_allDirectors, v.Directors...)
		}

		for _, task := range subTasks {
			_allDirectors = append(_allDirectors, task.Doc.Directors...)
			//_allTags = append(_allTags, task.Tags...) //子任务不需要返回标签
			_allCreator = append(_allCreator, task.UserId) //子任务不需要返回创建者
		}

		_allFollowers = utils.ToInt64Array(workItem.Doc.Followers)
	}

	allTagIds := utils.ToInt64Array(_allTags)
	allTagIds = stream.Unique(allTagIds)

	allDirectorIds := utils.ToInt64Array(_allDirectors)
	allCreatorIds := _allCreator
	allUserIds := append(allDirectorIds, allCreatorIds...)
	allUserIds = append(allUserIds, _allRestartUser...)
	allUserIds = append(allUserIds, oper.UserId)
	allUserIds = stream.Unique(allUserIds)
	userMap, err := s.userRepo.UserMap(ctx, allUserIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	//获取Tag信息
	tagInfos, _ := s.tagRepo.GetTagByIds(ctx, workItem.Doc.Tags.ToInt64s())
	tagMaps := stream.ToMap(tagInfos, func(i int, t *space_tag.SpaceTag) (int64, *space_tag.SpaceTag) {
		return t.Id, t
	})

	// 计算评论数
	var commentNum int64
	var unreadCommentNum int64
	if workItem.CommentNum > 0 {
		commentNum = int64(workItem.CommentNum)
		unreadCommentNum, _ = s.commentRepo.GetUserUnreadNum(ctx, oper.UserId, workItemId)
	}

	//流程角色信息
	wItemRoles, _ := s.workItemRoleRepo.GetWorkItemRoles(ctx, workItem.SpaceId)
	wItemRolesMap := stream.ToMap(wItemRoles, func(i int, w *work_item_role.WorkItemRole) (string, *work_item_role.WorkItemRole) {
		return w.Key, w
	})

	//---------- 组织返回内容 -------------

	buildSimpleUserInfoRsp := func(u *user_domain.User) *rsp.SimpleUserInfo {
		if u == nil {
			return nil
		}
		return &rsp.SimpleUserInfo{
			Id:           u.Id,
			UserId:       u.Id,
			UserName:     u.UserName,
			UserNickname: u.UserNickname,
			Avatar:       u.Avatar,
		}
	}

	buildTagInfoRsp := func(tag *space_tag.SpaceTag) *rsp.TagInfo {
		if tag == nil {
			return nil
		}

		return &rsp.TagInfo{
			TagName: tag.TagName,
			Id:      tag.Id,
		}
	}

	detailRsp := &rsp.SpaceWorkItemDetail{}
	detailRsp.Id = workItem.Id
	detailRsp.Pid = workItem.Pid
	detailRsp.SpaceId = workItem.SpaceId
	detailRsp.VersionId = workItem.VersionId
	detailRsp.WorkObjectId = workItem.WorkObjectId
	detailRsp.WorkItemGuid = workItem.WorkItemGuid
	detailRsp.WorkItemName = workItem.WorkItemName
	detailRsp.IconFlags = uint32(workItem.IconFlags)
	detailRsp.Priority = workItem.Doc.Priority
	detailRsp.ProcessRate = workItem.Doc.ProcessRate
	detailRsp.Describe = workItem.Doc.Describe
	detailRsp.Remark = workItem.Doc.Remark
	detailRsp.CommentNum = commentNum
	detailRsp.UnreadCommentNum = unreadCommentNum
	detailRsp.Creator = buildSimpleUserInfoRsp(userMap[workItem.UserId])
	detailRsp.RestartAt = workItem.Restart.RestartAt
	detailRsp.RestartBy = buildSimpleUserInfoRsp(userMap[workItem.Restart.RestartUserId])
	detailRsp.IsRestart = workItem.Restart.IsRestart == 1
	detailRsp.ResumeAt = workItem.Resume.ResumeAt

	detailRsp.Followed = slices.Contains(workItem.Doc.Followers, cast.ToString(oper.UserId))
	detailRsp.CreatedAt = workItem.CreatedAt
	detailRsp.UpdatedAt = workItem.UpdatedAt
	detailRsp.PlanTime = &rsp.SpaceWorkItemPlanTime{
		StartAt:    workItem.Doc.PlanStartAt,
		CompleteAt: workItem.Doc.PlanCompleteAt,
	}

	detailRsp.Followers = _allFollowers

	detailRsp.Template = &rsp.SpaceWorkItemDetail_Template{
		TemplateId:          wfTplt.Id,
		TemplateVersion:     int64(wfTplt.Version),
		LastTemplateVersion: int64(wfInfo.Version),
		IsLatestVersion:     wfInfo.IsLastTemplate(wfTplt),
		WorkFlowId:          wfInfo.Id,
		WorkFlowKey:         wfInfo.Key,
		WorkFlowMode:        string(wfInfo.FlowMode),
		WorkFlowName:        wfInfo.Name,
	}

	//工作项
	detailRsp.WorkItemTypeId = workItem.WorkItemTypeId
	detailRsp.WorkItemTypeKey = string(workItem.WorkItemTypeKey)
	if workItemTypeInfo != nil {
		detailRsp.WorkItemTypeName = workItemTypeInfo.Name
	}

	//任务状态
	detailRsp.WorkItemStatus = &rsp.SpaceWorkItemStatus{
		Val: workItem.WorkItemStatus.Val,
		Key: workItem.WorkItemStatus.Key,
		Id:  workItem.WorkItemStatus.Id,
		Last: &rsp.SpaceWorkItemStatus_Last{
			Val:       workItem.LastWorkItemStatus.Val,
			Key:       workItem.LastWorkItemStatus.Key,
			Id:        workItem.LastWorkItemStatus.Id,
			UpdatedAt: workItem.LastWorkItemStatus.LastAt,
		},
	}

	if statusInfo != nil {
		curStatus := statusInfo.GetItemByKey(detailRsp.WorkItemStatus.Key)
		if curStatus != nil {
			detailRsp.WorkItemStatus.Name = curStatus.Name
			detailRsp.WorkItemStatus.Type = int64(curStatus.StatusType)
		}

		lastStatus := statusInfo.GetItemByKey(detailRsp.WorkItemStatus.Last.Key)
		if lastStatus != nil {
			detailRsp.WorkItemStatus.Last.Name = lastStatus.Name
			detailRsp.WorkItemStatus.Last.Type = fmt.Sprintf("%v", lastStatus.StatusType)
		}
	}

	//当前负责人
	detailRsp.Owners = make([]*rsp.SimpleUserInfo, 0)
	for _, v := range workItem.Doc.Directors {
		ownerRsp := buildSimpleUserInfoRsp(userMap[cast.ToInt64(v)])
		if ownerRsp != nil {
			detailRsp.Owners = append(detailRsp.Owners, ownerRsp)
		}
	}

	//TAG信息
	detailRsp.Tags = make([]*rsp.TagInfo, 0)
	for _, v := range workItem.Doc.Tags {
		tagRsp := buildTagInfoRsp(tagMaps[cast.ToInt64(v)])
		if tagRsp != nil {
			detailRsp.Tags = append(detailRsp.Tags, tagRsp)
		}
	}

	//附件信息
	detailRsp.Files = make([]*rsp.FileInfo, 0)
	for _, v := range spaceFiles {
		detailRsp.Files = append(detailRsp.Files, &rsp.FileInfo{
			Id:         v.Id,
			SpaceId:    v.SpaceId,
			FileInfoId: v.FileInfo.FileInfoId,
			FileName:   v.FileInfo.FileName,
			FileSize:   v.FileInfo.FileSize,
			FileUri:    v.FileInfo.FileUri,
			FileCover:  fileCoverMap[v.FileInfo.FileInfoId],
			FileMeta:   fileMetaMap[v.FileInfo.FileInfoId],
			CreatedAt:  v.CreatedAt,
		})
	}

	//父级任务信息
	detailRsp.ParentWorkItem = &rsp.SpaceWorkItemDetail_ParentWorkItem{}
	if parentWorkItem != nil {
		detailRsp.ParentWorkItem.Id = parentWorkItem.Id
		detailRsp.ParentWorkItem.Status = &rsp.SpaceWorkItemDetail_ParentWorkItem_Status{}
		detailRsp.ParentWorkItem.Status.Val = parentWorkItem.WorkItemStatus.Val
		detailRsp.ParentWorkItem.Status.Key = parentWorkItem.WorkItemStatus.Key
		detailRsp.ParentWorkItem.Status.Id = parentWorkItem.WorkItemStatus.Id

		curParentStatus := statusInfo.GetItemByKey(parentWorkItem.WorkItemStatus.Key)
		if curParentStatus != nil {
			detailRsp.ParentWorkItem.Status.Name = curParentStatus.Name
			detailRsp.ParentWorkItem.Status.Type = fmt.Sprintf("%v", curParentStatus.StatusType)
		}
	}

	//节点信息
	detailRsp.Nodes = make([]*rsp.SpaceWorkFlowNode, 0)
	for _, v := range workItem.WorkItemFlowNodes {
		nodeRsp := &rsp.SpaceWorkFlowNode{
			Id:              v.Id,
			FlowNodeUuid:    v.FlowNodeUuid,
			FlowNodeStatus:  int32(v.FlowNodeStatus),
			FlowNodeCode:    v.FlowNodeCode,
			FlowNodePassed:  v.FlowNodePassed,
			FlowNodeReached: v.FlowNodeReached,
			StartAt:         v.StartAt,
			FinishAt:        v.FinishAt,
			Role:            v.WorkItemRoleKey,
			RoleId:          v.WorkItemRoleId,
			CreatedAt:       v.CreatedAt,
			UpdatedAt:       v.UpdatedAt,
			PlanTime: &rsp.SpaceWorkItemPlanTime{
				StartAt:    v.PlanTime.StartAt,
				CompleteAt: v.PlanTime.CompleteAt,
			},
		}

		//节点负责人
		nodeRsp.Owners = make([]*rsp.SimpleUserInfo, 0)
		for _, uid := range v.Directors {
			u := userMap[cast.ToInt64(uid)]
			if u == nil {
				continue
			}
			nodeRsp.Owners = append(nodeRsp.Owners, buildSimpleUserInfoRsp(u))
		}

		detailRsp.Nodes = append(detailRsp.Nodes, nodeRsp)
	}

	// 计算countAt
	detailRsp.CountAt = detailRsp.WorkItemStatus.Last.UpdatedAt
	if workItem.CountAt != 0 {
		detailRsp.CountAt = workItem.CountAt
	}

	//角色负责人
	detailRsp.RoleOwners = make([]*rsp.SpaceWorkItemOwnerRole, 0)
	for _, v := range workItem.WorkItemFlowRoles {
		roleOwnerRsp := &rsp.SpaceWorkItemOwnerRole{
			Owners: make([]*rsp.SimpleUserInfo, 0),
		}

		roleInfo := wItemRolesMap[v.WorkItemRoleKey]
		if roleInfo != nil {
			roleOwnerRsp.RoleId = roleInfo.Id
			roleOwnerRsp.RoleKey = roleInfo.Key
			roleOwnerRsp.RoleName = roleInfo.Name
		} else {
			roleOwnerRsp.RoleId = v.WorkItemRoleId
			roleOwnerRsp.RoleKey = v.WorkItemRoleKey
			roleOwnerRsp.RoleName = v.WorkItemRoleKey
		}

		for _, uid := range v.Directors {
			u := userMap[cast.ToInt64(uid)]
			if u == nil {
				continue
			}
			roleOwnerRsp.Owners = append(roleOwnerRsp.Owners, buildSimpleUserInfoRsp(u))
		}

		detailRsp.RoleOwners = append(detailRsp.RoleOwners, roleOwnerRsp)
	}

	//任务表单权限角色
	var permSubFuncRole map[string]interface{}
	if workItem.Pid == 0 {
		permSubFuncRole = s.permService.GetPermissionWithScene(member.GetRole(), workItem.GetRole(oper.UserId), perm.EditWorkItemScene)
	} else {
		permSubFuncRole = s.permService.GetTaskPermission(member.GetRole(), workItem.GetRole(oper.UserId))
	}

	permSubFuncRoleJsonV, _ := json.Marshal(permSubFuncRole)
	detailRsp.OperationPermissions = string(permSubFuncRoleJsonV)

	//子任务相关
	detailRsp.SubTasks = make([]*rsp.SpaceWorkItemSubTaskListItem, 0)
	for _, v := range subTasks {

		subTaskRsp := &rsp.SpaceWorkItemSubTaskListItem{
			Id:           v.Id,
			Pid:          v.Pid,
			WorkItemName: v.WorkItemName,
			WorkItemGuid: v.WorkItemGuid,
			SpaceId:      v.SpaceId,
			Creator:      buildSimpleUserInfoRsp(userMap[v.UserId]),
			ProcessRate:  v.Doc.ProcessRate,
			Priority:     v.Doc.Priority,
			PlanTime: &rsp.SpaceWorkItemPlanTime{
				StartAt:    v.Doc.PlanStartAt,
				CompleteAt: v.Doc.PlanCompleteAt,
			},
			WorkItemStatus: &rsp.SpaceWorkItemStatus{
				Val:  v.WorkItemStatus.Val,
				Key:  v.WorkItemStatus.Key,
				Id:   v.WorkItemStatus.Id,
				Type: int64(statusInfo.GetItemByKey(v.WorkItemStatus.Key).StatusType),
				Last: &rsp.SpaceWorkItemStatus_Last{
					Val:       v.LastWorkItemStatus.Val,
					Key:       v.LastWorkItemStatus.Key,
					Id:        v.LastWorkItemStatus.Id,
					UpdatedAt: v.LastWorkItemStatus.LastAt,
				},
			},
		}

		curSubTaskStatus := statusInfo.GetItemByKey(v.WorkItemStatus.Key)
		if curSubTaskStatus != nil {
			subTaskRsp.WorkItemStatus.Name = curSubTaskStatus.Name
			subTaskRsp.WorkItemStatus.Type = int64(curSubTaskStatus.StatusType)
		}

		curSubTaskLastStatus := statusInfo.GetItemByKey(v.LastWorkItemStatus.Key)
		if curSubTaskLastStatus != nil {
			subTaskRsp.WorkItemStatus.Last.Name = curSubTaskLastStatus.Name
			subTaskRsp.WorkItemStatus.Last.Type = fmt.Sprintf("%v", curSubTaskLastStatus.StatusType)
		}

		// 子任务负责人
		subTaskRsp.Owners = make([]*rsp.SimpleUserInfo, 0, len(v.Doc.Directors))
		for _, uid := range v.Doc.Directors {
			u := userMap[cast.ToInt64(uid)]
			if u == nil {
				continue
			}
			subTaskRsp.Owners = append(subTaskRsp.Owners, buildSimpleUserInfoRsp(u))
		}

		// 子任务表单权限角色
		taskPermission := s.permService.GetTaskPermission(member.GetRole(), v.GetRole(oper.UserId))
		taskPermissionJSON, _ := json.Marshal(taskPermission)
		subTaskRsp.OperationPermissions = string(taskPermissionJSON)
		detailRsp.SubTasks = append(detailRsp.SubTasks, subTaskRsp)
	}

	return detailRsp, nil
}
