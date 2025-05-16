package work_item

import (
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils"
	"go-cs/pkg/stream"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"
)

func NewWorkItem(
	id int64,
	spaceId int64,
	userId int64,
	workObjectId int64,
	versionId int64,
	name string,
	status WorkItemStatus,
	planTime PlanTime,
	processRate int32,
	remark string,
	describe string,
	Priority string,
	IconFlag IconFlag,
	tags Tags,
	directors Directors,
	followers []int64,
	workItemFlowNodes WorkItemFlowNodes,
	workItemFlowRoles WorkItemFlowRoles,
	workItemFiles WorkItemFiles,
	oper shared.Oper,
) *WorkItem {
	ins := &WorkItem{
		Id:           id,
		WorkItemGuid: uuid.NewString(),
		SpaceId:      spaceId,
		UserId:       userId,
		WorkObjectId: workObjectId,
		VersionId:    versionId,
		WorkItemName: name,
		Doc: WorkItemDoc{
			PlanStartAt:    planTime.StartAt,
			PlanCompleteAt: planTime.CompleteAt,
			ProcessRate:    processRate,
			Remark:         remark,
			Describe:       describe,
			Priority:       Priority,
			Tags:           tags,
			Directors:      directors,
			Followers:      utils.ToStrArray(followers),
			Participators:  make(Directors, 0),
		},
		Restart:            Restart{},
		Resume:             Resume{},
		LastWorkItemStatus: LastWorkItemStatus{},
		IconFlags:          IconFlag,
		WorkItemStatus:     status,
		WorkItemFlowNodes:  workItemFlowNodes,
		WorkItemFlowRoles:  workItemFlowRoles,
		WorkItemFiles:      workItemFiles,
	}

	if len(ins.WorkItemFiles) > 0 {
		ins.IconFlags.AddFlag(consts.ICON_FLAG_File)
	}

	ins.CreatedAt = time.Now().Unix()

	ins.AddMessage(oper, &domain_message.CreateWorkItem{
		SpaceId:      spaceId,
		WorkItemId:   id,
		WorkItemName: name,
	})

	return ins
}

func NewWorkItemSubTask(
	id int64,
	pid int64,
	spaceId int64,
	userId int64,
	workObjectId int64,
	versionId int64,
	name string,
	status WorkItemStatus,
	planTime PlanTime,
	processRate int32,
	Priority string,
	directors Directors,
	workItemFlowNodes WorkItemFlowNodes,
	oper shared.Oper,
) *WorkItem {
	ins := &WorkItem{
		Id:           id,
		SpaceId:      spaceId,
		Pid:          pid,
		UserId:       userId,
		WorkObjectId: workObjectId,
		VersionId:    versionId,
		WorkItemName: name,
		Doc: WorkItemDoc{
			PlanStartAt:    planTime.StartAt,
			PlanCompleteAt: planTime.CompleteAt,
			ProcessRate:    processRate,
			Priority:       Priority,
			Directors:      directors,
			Participators:  stream.Unique(append(directors, cast.ToString(userId))),
			NodeDirectors:  directors,
			Followers:      make(Directors, 0),
			Tags:           make(Tags, 0),
		},
		Restart:            Restart{},
		Resume:             Resume{},
		LastWorkItemStatus: LastWorkItemStatus{},
		IconFlags:          0,
		WorkItemStatus:     status,
		WorkItemFlowNodes:  workItemFlowNodes,
		WorkItemFlowRoles:  make(WorkItemFlowRoles, 0),
		WorkItemFiles:      make(WorkItemFiles, 0),
	}

	ins.CreatedAt = time.Now().Unix()

	ins.AddMessage(oper, &domain_message.CreateWorkItemSubTask{
		SpaceId:          spaceId,
		ParentWorkItemId: pid,
		WorkItemId:       id,
		WorkItemName:     name,
	})

	return ins
}

func NewWorkItemFlowNode(
	workFlowId int64,
	workFlowTemplateId int64,
	spaceId int64,
	workItemId int64,
	nodeCode string,
	workItemRoleId int64,
	workItemRoleKey string,
	directors Directors,
	createdAt int64,
) *WorkItemFlowNode {
	return &WorkItemFlowNode{
		FlowId:          workFlowId,
		FlowTemplateId:  workFlowTemplateId,
		SpaceId:         spaceId,
		WorkItemId:      workItemId,
		FlowNodeUuid:    uuid.NewString(),
		FlowNodeStatus:  FlowNodeStatus_Unstarted,
		FlowNodeCode:    nodeCode,
		FlowNodePassed:  0,
		FlowNodeReached: 0,
		WorkItemRoleId:  workItemRoleId,
		WorkItemRoleKey: workItemRoleKey,
		Directors:       directors,
		CreatedAt:       createdAt,
	}
}

func NewWorkItemFlowRole(
	spaceId int64,
	workItemId int64,
	workFlowId int64,
	workFlowTemplateId int64,
	workItemRoleId int64,
	workItemRoleKey string,
	directors Directors,
	createdAt int64,
) *WorkItemFlowRole {
	return &WorkItemFlowRole{
		SpaceId:         spaceId,
		WorkItemId:      workItemId,
		FlowId:          workFlowId,
		FlowTemplateId:  workFlowTemplateId,
		WorkItemRoleId:  workItemRoleId,
		WorkItemRoleKey: workItemRoleKey,
		Directors:       directors,
		CreatedAt:       createdAt,
	}
}

func NewWorkItemFile(
	spaceId int64,
	workItemId int64,
	fileInfo FileInfo,
) *WorkItemFile {
	return &WorkItemFile{
		SpaceId:    spaceId,
		WorkItemId: workItemId,
		FileInfo:   fileInfo,
		CreatedAt:  time.Now().Unix(),
		Status:     1,
	}
}
