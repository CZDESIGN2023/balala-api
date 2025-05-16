package work_item

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils"
	"go-cs/pkg/stream"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cast"
)

type WorkItemDoc struct {
	PlanStartAt    int64     `json:"plan_start_at"`
	PlanCompleteAt int64     `json:"plan_complete_at"`
	ProcessRate    int32     `json:"process_rate"`
	Remark         string    `json:"remark"`
	Describe       string    `json:"describe"`
	Priority       string    `json:"priority"`
	Tags           Tags      `json:"tags"`
	Directors      Directors `json:"directors"`
	Followers      Directors `json:"followers"`
	Participators  Directors `json:"participators"`
	NodeDirectors  Directors `json:"node_directors"`
}

type WorkItemReason struct {
	Restart   string `json:"restart,omitempty"`
	Resume    string `json:"resume,omitempty"`
	Rollback  string `json:"rollback,omitempty"`
	Terminate string `json:"terminate,omitempty"`
	Close     string `json:"close,omitempty"`
}

func (r WorkItemReason) ToJSON() string {
	marshal, _ := json.Marshal(r)
	return string(marshal)
}

type WorkItems []*WorkItem

func (w *WorkItems) GetIds() []int64 {
	var ids []int64
	for _, item := range *w {
		ids = append(ids, item.Id)
	}
	return ids
}

func (w *WorkItems) GetMessages() shared.DomainMessages {

	var logs shared.DomainMessages
	for _, item := range *w {
		logs = append(logs, item.GetMessages()...)
	}
	return logs
}

type WorkItem struct {
	shared.AggregateRoot

	Id  int64 ` json:"id"`
	Pid int64 ` json:"pid"` //上级任务id, >0 表示当前为子任务

	SpaceId                 int64                  ` json:"space_id"`                  //空间id
	UserId                  int64                  ` json:"user_id"`                   //创建用户id
	WorkObjectId            int64                  ` json:"work_object_id"`            //模块id
	VersionId               int64                  ` json:"version_id"`                //版本id -> 版本表
	WorkFlowId              int64                  `json:"flow_id"`                    //工作项类型Id
	WorkFlowKey             string                 `json:"flow_key"`                   //工作项类型Key
	WorkFlowTemplateId      int64                  `json:"work_flow_template_id"`      //流程模版Id
	WorkFlowTemplateVersion int64                  `json:"work_flow_template_version"` //流程模版Id
	WorkItemTypeId          int64                  `json:"work_item_type_id"`          //工作项类型Id
	WorkItemTypeKey         consts.WorkItemTypeKey `json:"work_item_type_key"`         //工作项类型Id
	FlowMode                consts.WorkFlowMode    `json:"flow_mode"`                  //流程模式

	WorkItemFlowId  int64  ` json:"work_item_flow_id"` // 子任务时为父级任务flowId
	WorkItemFlowKey string `json:"work_item_flow_key"`

	WorkItemGuid string   ` json:"work_item_guid"` //任务Guid
	WorkItemName string   `json:"work_item_name"`  //任务名称
	Restart      Restart  `json:"restart"`         //是否为重启任务
	IconFlags    IconFlag `json:"icon_flags"`      //图标标记
	CommentNum   int32    `json:"comment_num"`     //评论数
	ChildNum     int32    `json:"child_num"`       //子任务数量
	Resume       Resume   `json:"resume"`          //恢复时间

	WorkItemStatus     WorkItemStatus     `json:"work_item_status"`      //任务状态
	LastWorkItemStatus LastWorkItemStatus `json:"last_work_item_status"` //任务状态
	CountAt            int64              `json:"count_at"`

	CreatedAt int64 ` json:"created_at"` //创建时间
	UpdatedAt int64 ` json:"Updated_at"` //更新时间
	DeletedAt int64 ` json:"deleted_at"` //删除时间

	//------- doc 部分
	Doc WorkItemDoc `json:"doc"`

	// reason
	Reason WorkItemReason `json:"reason"`

	//-------- 流程节点信息部分
	WorkItemFlowNodes WorkItemFlowNodes
	//-------- 流程角色负责人信息
	WorkItemFlowRoles WorkItemFlowRoles
	//-------- 附件
	WorkItemFiles WorkItemFiles
}

//--- 属性变化 ---

func (s *WorkItem) UpdateWorkFlowTemplate(templateId int64, templateVersion int64) {
	if s.WorkFlowTemplateId == templateId {
		return
	}

	s.WorkFlowTemplateId = templateId
	s.WorkFlowTemplateVersion = templateVersion

	s.AddDiff(Diff_WorkFlowTemplate)
}

func (s *WorkItem) UpdateRestart(uid int64, isRestart int32) {
	s.Restart.RestartAt = time.Now().Unix()
	s.Restart.RestartUserId = uid
	s.Restart.IsRestart = isRestart
	s.AddDiff(Diff_Restart)
}

func (s *WorkItem) UpdateStatus(newStatus WorkItemStatus) error {

	if s.WorkItemStatus.Id == newStatus.Id {
		return errors.New("状态未改变")
	}

	oldStatus := s.WorkItemStatus

	s.WorkItemStatus = newStatus
	s.LastWorkItemStatus = LastWorkItemStatus{
		Id:     oldStatus.Id,
		Key:    oldStatus.Key,
		Val:    oldStatus.Val,
		LastAt: time.Now().Unix(),
	}

	s.AddDiff(Diff_Status)
	return nil
}

func (s *WorkItem) ForceUpdateLastStatusTime() {

	s.LastWorkItemStatus.LastAt = time.Now().Unix()
	s.AddDiff(Diff_Status)
}

func (s *WorkItem) UpdateParticipators() {
	var nodeDirectors []string

	for _, v := range s.WorkItemFlowNodes {
		nodeDirectors = append(nodeDirectors, v.Directors...)
	}

	for _, v := range s.WorkItemFlowRoles {
		nodeDirectors = append(nodeDirectors, v.Directors...)
	}

	// 子任务不一定有节点负责人，所有一定要加上当前负责人
	participators := stream.Concat(nodeDirectors, []string(s.Doc.Directors), []string{cast.ToString(s.UserId)})

	s.Doc.NodeDirectors = stream.Unique(nodeDirectors)
	s.Doc.Participators = stream.Unique(participators)
	s.AddDiff(Diff_Participators, Diff_NodeDirectors)
}

func (s *WorkItem) UpdateFollower(directors Directors, oper shared.Oper) {
	oldValue := s.Doc.Followers.ToStrings()
	s.Doc.Followers = directors
	s.AddDiff(Diff_Followers)

	s.AddMessage(oper, &domain_message.ModifyWorkItem{
		SpaceId:      s.SpaceId,
		WorkItemId:   s.Id,
		WorkItemName: s.WorkItemName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "followers",
				OldValue: utils.ToInt64Array(oldValue),
				NewValue: utils.ToInt64Array(directors),
			},
		},
	})
}

func (s *WorkItem) UpdateDirectors(directors Directors) {
	s.Doc.Directors = directors
	s.AddDiff(Diff_Directors)
}

func (s *WorkItem) ReCalcDirectors() {
	processNodes := s.WorkItemFlowNodes.GetProcessingNodes()
	if len(processNodes) > 0 {
		directors := processNodes.GetAllDirectors()
		s.UpdateDirectors(directors)
	}
}

func (s *WorkItem) UpdateName(name string) error {

	if s.WorkItemName == name {
		return errors.New("work item name is same")
	}

	s.WorkItemName = name
	s.AddDiff(Diff_Name)

	return nil
}

func (s *WorkItem) UpdatePriority(priority string) error {

	if s.Doc.Priority == priority {
		return errors.New("priority is same")
	}

	s.Doc.Priority = priority
	s.AddDiff(Diff_Priority)
	return nil
}

func (s *WorkItem) UpdateProcessRate(processRate int32) {
	s.Doc.ProcessRate = processRate
	s.AddDiff(Diff_ProcessRate)
}

func (s *WorkItem) UpdateDescribe(describe string) error {

	if s.Doc.Describe == describe {
		return errors.New("describe is same")
	}

	s.Doc.Describe = describe
	s.AddDiff(Diff_Describe)
	return nil
}

func (s *WorkItem) UpdateRemark(remark string) error {

	if s.Doc.Remark == remark {
		return errors.New("remark is same")
	}

	s.Doc.Remark = remark
	s.AddDiff(Diff_Remark)
	return nil
}

func (s *WorkItem) UpdateVersionId(versionId int64) error {
	if s.VersionId == versionId {
		return errors.New("versionId is same")
	}

	s.VersionId = versionId
	s.AddDiff(Diff_VersionId)
	return nil
}

func (s *WorkItem) UpdateWorkObjectId(objectId int64) error {

	if s.WorkObjectId == objectId {
		return errors.New("workObjectId is same")
	}

	s.WorkObjectId = objectId
	s.AddDiff(Diff_ObjectId)
	return nil
}

func (s *WorkItem) UpdateTags(tags Tags) {
	s.Doc.Tags = tags
	s.AddDiff(Diff_Tags)
}

func (s *WorkItem) UpdatePlanTime(startAt int64, completeAt int64) {

	s.Doc.PlanStartAt = startAt
	s.Doc.PlanCompleteAt = completeAt

	s.AddDiff(Diff_PlanTime)
}

func (s *WorkItem) UpdateRoleDirectors(wItemRoleId int64, directors Directors) {
	if s.WorkItemFlowRoles == nil {
		return
	}

	for _, v := range s.WorkItemFlowRoles {
		if v.WorkItemRoleId == wItemRoleId {
			v.UpdateDirectors(directors)
		}
	}
}

func (s *WorkItem) RemoveAllWorkItemFlowNodes() WorkItemFlowNodes {
	var wItemFlowNodes WorkItemFlowNodes
	for i := 0; i < len(s.WorkItemFlowNodes); i++ {
		wItemFlowNodes = append(wItemFlowNodes, s.WorkItemFlowNodes[i])
	}
	s.WorkItemFlowNodes = make(WorkItemFlowNodes, 0)
	return wItemFlowNodes
}

func (s *WorkItem) RemoveAllWorkItemFlowRoles() WorkItemFlowRoles {
	var wItemFlowRoles WorkItemFlowRoles
	for i := 0; i < len(s.WorkItemFlowRoles); i++ {
		wItemFlowRoles = append(wItemFlowRoles, s.WorkItemFlowRoles[i])
	}
	s.WorkItemFlowRoles = make(WorkItemFlowRoles, 0)
	return wItemFlowRoles
}

func (s *WorkItem) AddWorkItemFlowNode(node ...*WorkItemFlowNode) {
	s.WorkItemFlowNodes = append(s.WorkItemFlowNodes, node...)
}

func (s *WorkItem) AddWorkItemFlowRole(node ...*WorkItemFlowRole) {
	s.WorkItemFlowRoles = append(s.WorkItemFlowRoles, node...)
}

//--- 领域业务相关---

func (s *WorkItem) ChangeVersionId(versionId int64, oper shared.Oper) {

	oldValue := s.VersionId
	err := s.UpdateVersionId(versionId)
	if err != nil {
		return
	}

	s.AddMessage(oper, &domain_message.ModifyWorkItem{
		SpaceId:      s.SpaceId,
		WorkItemId:   s.Id,
		WorkItemName: s.WorkItemName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "versionId",
				OldValue: oldValue,
				NewValue: s.VersionId,
			},
		},
	})

}

func (s *WorkItem) ChangeWorkObjectId(objectId int64, oper shared.Oper) {

	oldValue := s.WorkObjectId

	err := s.UpdateWorkObjectId(objectId)
	if err != nil {
		return
	}

	s.AddMessage(oper, &domain_message.ModifyWorkItem{
		SpaceId:      s.SpaceId,
		WorkItemId:   s.Id,
		WorkItemName: s.WorkItemName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "workObjectId",
				OldValue: oldValue,
				NewValue: objectId,
			},
		},
	})
}

// 更新任务单状态
func (s *WorkItem) ChangeStatus(newStatus WorkItemStatus, reason string, affectByParent bool, oper shared.Oper) {

	oldStatus := s.WorkItemStatus
	err := s.UpdateStatus(newStatus)
	if err != nil {
		return
	}

	// 更新countAt，针对终止，关闭，重启，恢复
	if string(consts.WorkItemStatus_TerminatedKey) == newStatus.Key || string(consts.WorkItemStatus_WorkFlowCloseDefaultKey) == newStatus.Key ||
		string(consts.WorkItemStatus_TerminatedKey) == oldStatus.Key || string(consts.WorkItemStatus_WorkFlowCloseDefaultKey) == oldStatus.Key {
		s.SetCountAt(time.Now())
	}

	s.AddMessage(oper, &domain_message.ChangeWorkItemStatus{
		AffectByParent:       affectByParent,
		SpaceId:              s.SpaceId,
		Pid:                  s.Pid,
		WorkItemId:           s.Id,
		WorkItemName:         s.WorkItemName,
		WorkItemTypeKey:      s.WorkItemTypeKey,
		OldWorkItemStatusKey: oldStatus.Key,
		OldWorkItemStatusVal: oldStatus.Val,
		OldWorkItemStatusId:  oldStatus.Id,
		NewWorkItemStatusKey: newStatus.Key,
		NewWorkItemStatusId:  newStatus.Id,
		NewWorkItemStatusVal: newStatus.Val,
		Reason:               reason,
	})
}

func (s *WorkItem) ChangeStateFlowMainStatus(newStatus WorkItemStatus, reason, remark string, oper shared.Oper) {

	oldStatus := s.WorkItemStatus
	err := s.UpdateStatus(newStatus)
	if err != nil {
		return
	}

	s.AddMessage(oper, &domain_message.ChangeWorkItemStatus{
		SpaceId:              s.SpaceId,
		Pid:                  s.Pid,
		WorkItemId:           s.Id,
		WorkItemName:         s.WorkItemName,
		WorkItemTypeKey:      s.WorkItemTypeKey,
		OldWorkItemStatusKey: oldStatus.Key,
		OldWorkItemStatusVal: oldStatus.Val,
		OldWorkItemStatusId:  oldStatus.Id,
		NewWorkItemStatusKey: newStatus.Key,
		NewWorkItemStatusId:  newStatus.Id,
		NewWorkItemStatusVal: newStatus.Val,
		Reason:               reason,
		Remark:               remark,
	})
}

func (s *WorkItem) UpdateStateFlowStateToFinished() {
	if !s.IsSubTask() {
		return
	}

	for _, v := range s.WorkItemFlowNodes {
		v.UpdateStateOnNodePass()
	}
}

func (s *WorkItem) UpdateStateFlowStateToProgressing() {
	if !s.IsSubTask() {
		return
	}

	for _, v := range s.WorkItemFlowNodes {
		v.UpdateStateOnNodeReach()
	}
}

// 更新计划时间
func (s *WorkItem) ChangePlanTime(startAt int64, completeAt int64, oper shared.Oper) {

	oldStartAt := s.Doc.PlanStartAt
	oldCompleteAt := s.Doc.PlanCompleteAt

	s.UpdatePlanTime(startAt, completeAt)

	s.AddMessage(oper, &domain_message.ModifyWorkItem{
		SpaceId:      s.SpaceId,
		WorkItemId:   s.Id,
		WorkItemName: s.WorkItemName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "planTime",
				OldValue: []any{oldStartAt, oldCompleteAt},
				NewValue: []any{startAt, completeAt},
			},
		},
	})
}

// 增加tag
func (s *WorkItem) ChangeTags(adds, removes Tags, oper shared.Oper) {

	oldTags := s.Doc.Tags

	// 计算出可被移除的负责人
	removes = stream.Filter(removes, func(id string) bool {
		return stream.Contains(s.Doc.Tags, id)
	})

	// 计算出被添加的负责人
	adds = stream.Filter(adds, func(id string) bool {
		return !stream.Contains(removes, id)
	})

	adds = stream.Filter(adds, func(id string) bool {
		return !stream.Contains(s.Doc.Tags, id)
	})

	// 最终被处理的负责人
	tags := s.Doc.Tags
	tags = append(tags, adds...)
	tags = stream.Filter(tags, func(id string) bool {
		return !stream.Contains(removes, id)
	})
	tags = stream.Unique(tags)
	s.UpdateTags(tags)

	s.AddMessage(oper, &domain_message.ChangeWorkItemTag{
		SpaceId:      s.SpaceId,
		WorkItemId:   s.Id,
		WorkItemName: s.WorkItemName,
		AddTags:      adds,
		RemoveTags:   removes,
		OldTags:      oldTags,
		NewTags:      tags,
	})
}

func (s *WorkItem) Follow(directors Directors, oper shared.Oper) {

	addDirectors := stream.Filter(directors, func(id string) bool {
		return !stream.Contains(s.Doc.Followers, id)
	})

	if len(addDirectors) == 0 {
		return
	}

	followers := append(s.Doc.Followers, addDirectors...)
	followers = stream.Unique(followers)

	s.UpdateFollower(followers, oper)
}

func (s *WorkItem) UnFollow(directors Directors, oper shared.Oper) {

	removeDirectors := stream.Filter(directors, func(id string) bool {
		return stream.Contains(s.Doc.Followers, id)
	})

	if len(removeDirectors) == 0 {
		return
	}

	followers := stream.Filter(s.Doc.Followers, func(id string) bool {
		return !stream.Contains(removeDirectors, id)
	})
	followers = stream.Unique(followers)

	s.UpdateFollower(followers, oper)
}

func (s *WorkItem) ChangeName(name string, oper shared.Oper) {

	oldValue := s.WorkItemName
	err := s.UpdateName(name)
	if err != nil {
		return
	}

	s.AddMessage(oper, &domain_message.ModifyWorkItem{
		SpaceId:      s.SpaceId,
		WorkItemId:   s.Id,
		WorkItemName: s.WorkItemName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "workItemName",
				OldValue: oldValue,
				NewValue: name,
			},
		},
	})
}

func (s *WorkItem) ChangePriority(priority string, oper shared.Oper) {

	oldValue := s.Doc.Priority

	err := s.UpdatePriority(priority)
	if err != nil {
		return
	}

	s.AddMessage(oper, &domain_message.ModifyWorkItem{
		SpaceId:      s.SpaceId,
		WorkItemId:   s.Id,
		WorkItemName: s.WorkItemName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "priority",
				OldValue: oldValue,
				NewValue: priority,
			},
		},
	})
}

func (s *WorkItem) ChangeProcessRate(processRate int32, oper shared.Oper) {

	oldValue := s.Doc.ProcessRate

	if oldValue == processRate {
		return
	}

	s.UpdateProcessRate(processRate)

	s.AddMessage(oper, &domain_message.ModifyWorkItem{
		SpaceId:      s.SpaceId,
		WorkItemId:   s.Id,
		WorkItemName: s.WorkItemName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "processRate",
				OldValue: cast.ToString(oldValue),
				NewValue: cast.ToString(processRate),
			},
		},
	})
}

func (s *WorkItem) ChangeDescribe(describe string, oper shared.Oper) {

	oldValue := s.Doc.Describe

	err := s.UpdateDescribe(describe)
	if err != nil {
		return
	}

	s.AddMessage(oper, &domain_message.ModifyWorkItem{
		SpaceId:      s.SpaceId,
		WorkItemId:   s.Id,
		WorkItemName: s.WorkItemName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "describe",
				OldValue: oldValue,
				NewValue: describe,
			},
		},
	})
}

func (s *WorkItem) ChangeRemark(remark string, oper shared.Oper) {

	oldValue := s.Doc.Remark

	err := s.UpdateRemark(remark)
	if err != nil {
		return
	}

	s.AddMessage(oper, &domain_message.ModifyWorkItem{
		SpaceId:      s.SpaceId,
		WorkItemId:   s.Id,
		WorkItemName: s.WorkItemName,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "remark",
				OldValue: oldValue,
				NewValue: remark,
			},
		},
	})
}

func (s *WorkItem) UpdateIconFlag(flagUpdates ...*IconFlagUpdate) {
	s.IconFlags.SetFlag(flagUpdates...)
	s.AddDiff(Diff_IconFlags)
}

func (s *WorkItem) OnDelete(oper shared.Oper) {

	if s.DeletedAt > 0 {
		return
	}

	s.DeletedAt = time.Now().Unix()

	s.AddMessage(oper, &domain_message.DeleteWorkItem{
		ParentWorkItemId: s.Pid,
		SpaceId:          s.SpaceId,
		WorkItemId:       s.Id,
		WorkItemName:     s.WorkItemName,
		PlanStartAt:      s.Doc.PlanStartAt,
		PlanCompleteAt:   s.Doc.PlanCompleteAt,
		ProcessRate:      s.Doc.ProcessRate,
		Directors:        utils.ToInt64Array(s.Doc.Directors),
	})
}

func (s *WorkItem) IsSubTask() bool {
	return s.Pid != 0
}

func (s *WorkItem) IsStateFlowMainTask() bool {
	return s.FlowMode == consts.FlowMode_StateFlow && s.Pid == 0
}

func (s *WorkItem) IsWorkFlowMainTask() bool {
	return s.FlowMode == consts.FlowMode_WorkFlow && s.Pid == 0
}

func (s *WorkItem) HasChild() bool {
	return s.Pid == 0 && s.ChildNum != 0
}

func (s *WorkItem) IsSameSpace(spaceId int64) bool {
	return s.SpaceId == spaceId
}

func (s *WorkItem) IsCreator(userId int64) bool {
	return s.UserId == userId
}

func (s *WorkItem) GetRole(userId int64) string {

	isWorkItemCreator := s.IsCreator(userId)
	isNodeOwner := slices.Contains(s.Doc.Directors, cast.ToString(userId))

	workItemRole := ""
	switch {
	case isWorkItemCreator:
		workItemRole = consts.WORK_ITEM_ROLE_CREATOR
	case isNodeOwner:
		workItemRole = consts.WORK_ITEM_ROLE_NODE_OWNER
	}

	return workItemRole
}

func (s *WorkItem) RemoveWorkFile(fileInfoId []int64) WorkItemFiles {
	removeFiles := stream.Filter(s.WorkItemFiles, func(file *WorkItemFile) bool {
		return stream.Contains(fileInfoId, file.FileInfo.FileInfoId)
	})

	if len(removeFiles) > 0 {
		s.WorkItemFiles = stream.Filter(s.WorkItemFiles, func(file *WorkItemFile) bool {
			return !stream.Contains(removeFiles, file)
		})
	}

	//清理掉flag标记
	s.CalIconFlag()
	return removeFiles
}

func (s *WorkItem) AddWorkFile(files WorkItemFiles) WorkItemFiles {
	var curFileIds []int64
	for _, v := range s.WorkItemFiles {
		curFileIds = append(curFileIds, v.Id)
	}

	addWorkFiles := stream.Filter(files, func(file *WorkItemFile) bool {
		return !stream.Contains(curFileIds, file.Id)
	})

	if len(addWorkFiles) > 0 {
		s.WorkItemFiles = append(s.WorkItemFiles, addWorkFiles...)
		s.UpdateIconFlag(&IconFlagUpdate{
			Flag: consts.ICON_FLAG_File,
			Val:  1,
		})
	}

	return addWorkFiles
}

func (s *WorkItem) HasWorkFile(fileInfoId int64) bool {
	for _, v := range s.WorkItemFiles {
		if v.FileInfo.FileInfoId == fileInfoId {
			return true
		}
	}
	return false
}

func (s *WorkItem) HasWorkItemFlowNodeById(flowNodeId int64) bool {
	for _, v := range s.WorkItemFlowNodes {

		if v.Id == flowNodeId {
			return true
		}
	}
	return false
}

func (s *WorkItem) SetFollowers(userIds []int64) {
	s.Doc.Followers = utils.ToStrArray(userIds)
	s.AddDiff(Diff_Followers)
}

func (s *WorkItem) SetRestartReason(reason string) {
	s.Reason.Restart = reason
	s.AddDiff(Diff_Reason)
}

func (s *WorkItem) SetResumeReason(reason string) {
	s.Reason.Resume = reason
	s.Resume.ResumeAt = time.Now().Unix()
	s.AddDiff(Diff_Reason, Diff_Resume)
}

func (s *WorkItem) SetRollbackReason(reason string) {
	s.Reason.Rollback = reason
	s.AddDiff(Diff_Reason)
}

func (s *WorkItem) SetTerminateReason(reason string) {
	s.Reason.Terminate = reason
	s.AddDiff(Diff_Reason)
}

func (s *WorkItem) SetCloseReason(reason string) {
	s.Reason.Close = reason
	s.AddDiff(Diff_Reason)
}

func (s *WorkItem) SetCountAt(t time.Time) {
	s.CountAt = t.Unix()
	s.AddDiff(Diff_CountAt)
}

func (s *WorkItem) ChangeNodeDirectors(directors Directors, nodeKeys ...string) {
	for _, v := range s.WorkItemFlowNodes {
		if slices.Contains(nodeKeys, v.FlowNodeCode) {
			v.UpdateDirectors(directors)
		}
	}
}

func (s *WorkItem) ChangeRoleDirectors(directors Directors, roleKeys ...string) {
	for _, v := range s.WorkItemFlowRoles {
		if slices.Contains(roleKeys, v.WorkItemRoleKey) {
			v.UpdateDirectors(directors)
		}
	}
}

func (s *WorkItem) IsTerminateStatus() bool {
	return s.WorkItemStatus.Key == string(consts.WorkItemStatus_TerminatedKey)
}

func richTextIconFlags(text string) []uint32 {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(text))
	if err != nil {
		return nil
	}

	var flags []uint32

	videoDoc := doc.Find(".file-video")
	if len(videoDoc.Nodes) > 0 {
		flags = append(flags, ICON_FLAG_File)
		videoDoc.Remove()
	}

	imageDoc := doc.Find(".file-image")
	if len(imageDoc.Nodes) > 0 {
		flags = append(flags, ICON_FLAG_Pic)
		imageDoc.Remove()
	}
	fileDoc := doc.Find(".inserted-file")
	if len(fileDoc.Nodes) > 0 {
		flags = append(flags, ICON_FLAG_File)
		fileDoc.Remove()
	}

	if len(doc.Find("img").Nodes) > 0 {
		flags = append(flags, ICON_FLAG_Pic)
	}

	if len(doc.Find("a").Nodes) > 0 {
		flags = append(flags, ICON_FLAG_Url)
	}

	return flags
}

func (s *WorkItem) CalIconFlag() {
	descFlags := richTextIconFlags(s.Doc.Describe)
	remarkFlags := richTextIconFlags(s.Doc.Remark)

	allFlags := append(descFlags, remarkFlags...)

	if len(s.WorkItemFiles) != 0 {
		allFlags = append(allFlags, ICON_FLAG_File)
	}

	s.IconFlags = IconFlag(consts.MergeIconFlags(stream.Unique(allFlags)))
	s.AddDiff(Diff_IconFlags)
}
