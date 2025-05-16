package work_item

import (
	domain_message "go-cs/internal/domain/pkg/message"
	shared "go-cs/internal/pkg/domain"
	"go-cs/pkg/stream"
	"slices"
	"strconv"
	"time"
)

type FlowNodeStatus int32

const (
	FlowNodeStatus_Undefined FlowNodeStatus = iota
	FlowNodeStatus_Unstarted
	FlowNodeStatus_InProgress
	FlowNodeStatus_Completed
)

type WorkItemFlowNode struct {
	shared.DomainEntity

	Id             int64 `json:"id,omitempty"`
	FlowId         int64 `json:"flow_id,omitempty"`          //流程Id
	FlowTemplateId int64 `json:"flow_template_id,omitempty"` //流程模版Id
	SpaceId        int64 `json:"space_id,omitempty"`         //空间id
	WorkItemId     int64 `json:"work_item_id,omitempty"`     //任务id

	WorkItemRoleId  int64     `json:"work_item_role_id,omitempty"`  //任务角色id
	WorkItemRoleKey string    `json:"work_item_role_key,omitempty"` //任务角色id
	Directors       Directors `json:"directors,omitempty"`          //节点负责人

	FlowNodeUuid    string         `json:"flow_node_uuid,omitempty"`    //节点uuid
	FlowNodeStatus  FlowNodeStatus `json:"flow_node_status,omitempty"`  //节点状态; 0:未定义 , 1: 未开启 2: 进行中 3:已完成
	FlowNodeCode    string         `json:"flow_node_code,omitempty"`    //节点编码
	FlowNodePassed  int32          `json:"flow_node_passed,omitempty"`  //节点是否通过
	FlowNodeReached int32          `json:"flow_node_reached,omitempty"` //节点是否到达
	StartAt         int64          `json:"start_at,omitempty"`          //开始时间
	FinishAt        int64          `json:"finish_at,omitempty"`         //结束时间
	PlanTime        PlanTime       `json:"plan_time,omitempty"`         //计划时间

	CreatedAt int64 `json:"created_at,omitempty"` //创建时间
	UpdatedAt int64 `json:"updated_at,omitempty"` //更新时间
	DeletedAt int64 `json:"deleted_at,omitempty"` //删除时间
}

func (s *WorkItemFlowNode) ResetStatusForm(form *WorkItemFlowNode) {
	s.FlowNodeStatus = form.FlowNodeStatus
	s.FlowNodePassed = form.FlowNodePassed
	s.FlowNodeReached = form.FlowNodeReached
	s.StartAt = form.StartAt
	s.FinishAt = form.FinishAt
	s.AddDiff(Diff_WorkItemFlowNode_Status)
}

func (s *WorkItemFlowNode) ResetStatus() {
	if s.FlowNodeStatus == FlowNodeStatus_Unstarted {
		return
	}

	s.FlowNodeStatus = FlowNodeStatus_Unstarted
	s.FlowNodePassed = 0
	s.FlowNodeReached = 0
	s.StartAt = 0
	s.FinishAt = 0
	s.AddDiff(Diff_WorkItemFlowNode_Status)
}

func (s *WorkItemFlowNode) ResetProgressStatus() {
	if s.FlowNodeStatus == FlowNodeStatus_InProgress {
		return
	}

	s.FlowNodeStatus = FlowNodeStatus_InProgress
	s.FlowNodePassed = 0
	s.FlowNodeReached = 1
	s.StartAt = time.Now().Unix()
	s.FinishAt = 0
	s.AddDiff(Diff_WorkItemFlowNode_Status)
}

func (s *WorkItemFlowNode) UpdateStateOnNodePass() {

	if s.FlowNodePassed == 1 {
		return
	}

	if s.FlowNodeReached == 0 {
		s.FlowNodeReached = 1
		s.StartAt = time.Now().Unix()
	}

	s.FlowNodePassed = 1
	s.FlowNodeStatus = FlowNodeStatus_Completed
	s.FinishAt = time.Now().Unix()

	s.AddDiff(Diff_WorkItemFlowNode_Status)
}

func (s *WorkItemFlowNode) ChangePlanTime(planTime PlanTime, oper shared.Oper) {

	oldStartAt := s.PlanTime.StartAt
	oldCompleteAt := s.PlanTime.CompleteAt

	s.UpdatePlanTime(planTime)
	s.AddMessage(oper, &domain_message.ModifyWorkItemFlowNode{
		SpaceId:    s.SpaceId,
		WorkItemId: s.WorkItemId,

		FlowNodeId:   s.Id,
		FlowNodeCode: s.FlowNodeCode,
		Updates: []domain_message.FieldUpdate{
			{
				Field:    "planTime",
				OldValue: []any{oldStartAt, oldCompleteAt},
				NewValue: []any{planTime.StartAt, planTime.CompleteAt},
			},
		},
	})
}

func (s *WorkItemFlowNode) UpdatePlanTime(planTime PlanTime) {
	s.PlanTime = planTime
	s.AddDiff(Diff_WorkItemFlowNode_PlanTime)
}

func (s *WorkItemFlowNode) PlanTimeHasSet() bool {
	return s.PlanTime.StartAt != 0 && s.PlanTime.CompleteAt != 0
}

func (s *WorkItemFlowNode) UpdateStateOnNodeReach() {

	if s.FlowNodeReached == 1 {
		return
	}

	s.FlowNodeReached = 1
	s.StartAt = time.Now().Unix()
	s.FlowNodeStatus = FlowNodeStatus_InProgress

	s.AddDiff(Diff_WorkItemFlowNode_Status)
}

func (s *WorkItemFlowNode) UpdateDirectors(directors Directors) {
	s.Directors = stream.Unique(directors)
	s.AddDiff(Diff_WorkItemFlowNode_Directors)
}

func (s *WorkItemFlowNode) SetUpdatedAt(ts int64) {
	s.UpdatedAt = ts
	s.AddDiff(Diff_WorkItemFlowNode_UpdatedAt)
}

func (s *WorkItemFlowNode) AddDirectors(directors Directors) {

	addDirectors := stream.Filter(directors, func(id string) bool {
		return !stream.Contains(s.Directors, id)
	})

	if len(addDirectors) == 0 {
		return
	}

	newDirectors := append(s.Directors, addDirectors...)
	s.UpdateDirectors(newDirectors)

}

func (s *WorkItemFlowNode) ContainsDirector(userId int64) bool {
	return stream.Contains(s.Directors, strconv.Itoa(int(userId)))
}

// 移除负责人
func (s *WorkItemFlowNode) RemoveDirectors(directors Directors) {

	removeDirectors := stream.Filter(directors, func(id string) bool {
		return stream.Contains(s.Directors, id)
	})

	if len(removeDirectors) == 0 {
		return
	}

	newDirectors := stream.Filter(s.Directors, func(id string) bool {
		return !stream.Contains(removeDirectors, id)
	})

	s.UpdateDirectors(newDirectors)

}

func (s *WorkItemFlowNode) IsUnstart() bool {
	return s.FlowNodeStatus == FlowNodeStatus_Unstarted
}

func (s *WorkItemFlowNode) IsInProcess() bool {
	return s.FlowNodeStatus == FlowNodeStatus_InProgress
}

func (s *WorkItemFlowNode) IsCompleted() bool {
	return s.FlowNodeStatus == FlowNodeStatus_Completed
}

func (s *WorkItemFlowNode) IsStartNode() bool {
	return s.FlowNodeCode == "started"
}

func (s *WorkItemFlowNode) IsEndNode() bool {
	return s.FlowNodeCode == "ended"
}

type WorkItemFlowNodes []*WorkItemFlowNode

func (s *WorkItemFlowNodes) GetIds() []int64 {
	var ids []int64
	for _, v := range *s {
		ids = append(ids, v.Id)
	}
	return ids
}

func (s *WorkItemFlowNodes) NodeMap() map[string]*WorkItemFlowNode {
	nodeMap := make(map[string]*WorkItemFlowNode)
	for _, v := range *s {
		nodeMap[v.FlowNodeCode] = v
	}
	return nodeMap
}

func (s *WorkItemFlowNodes) ResetStatus() {
	for _, v := range *s {
		v.ResetStatus()
	}
}

func (s *WorkItemFlowNodes) ResetStatusByNodeCode(nodeCode ...string) {
	for _, v := range *s {
		if slices.Contains(nodeCode, v.FlowNodeCode) {
			v.ResetStatus()
		}
	}
}

func (s *WorkItemFlowNodes) GetProcessingNodes() WorkItemFlowNodes {
	var nodes WorkItemFlowNodes
	for _, node := range *s {
		if node.FlowNodeStatus == FlowNodeStatus_InProgress {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (s *WorkItemFlowNodes) GetAllDirectors() []string {
	var ids []string
	for _, node := range *s {
		ids = append(ids, node.Directors...)

	}
	return stream.Unique(ids)
}

func (s *WorkItemFlowNodes) GetNodeByCode(code string) *WorkItemFlowNode {
	for _, node := range *s {
		if node.FlowNodeCode == code {
			return node
		}
	}
	return nil
}

func (s *WorkItemFlowNodes) GetNodeById(id int64) *WorkItemFlowNode {
	for _, node := range *s {
		if node.Id == id {
			return node
		}
	}
	return nil
}
