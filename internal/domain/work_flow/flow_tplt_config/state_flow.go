package config

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"go-cs/pkg/stream"
)

var (
	StateFlow_Default_RestartResumeOptions = []string{
		"任务需要继续进行",
		"误操作",
	}

	StateFlow_Default_TerminatedReasonOptions = []string{
		"取消任务，综合考虑且已同步现在先不做了",
		"任务重复/合并，与其他在做任务一同推进",
		"取消任务，转其他任务流程",
	}
)

func MustFormStateFlowJson(jsonStr string) (*StateFlow, error) {
	wf := &StateFlow{}
	err := json.Unmarshal([]byte(jsonStr), &wf)
	return wf, err
}

// 状态切换规则
type StateFlowTransitionRule struct {
	//状态扭转的时候，需要填写的表单内容
	ConfirmForm []string `json:"confirmForm"`
	//源状态节点编码 (状态key)
	SourceStateKey string `json:"sourceStateKey"`
	//目标状态编码 (状态key)
	TargetStateKey string `json:"targetStateKey"`
}

// 节点 负责人分配方式配置

// 需求类型状态切换 FlowMode
type StateFlow struct {
	//UUID
	Uuid string `json:"uuid"`
	//状态流程编码
	Key string `json:"key"`
	//状态节点
	StateFlowNodes []*StateFlowNode `json:"stateFlowNodes"`
	//各节点切换规则
	StateFlowTransitionRule []*StateFlowTransitionRule `json:"stateFlowTransitionRule"`

	//终止原因
	TerminatedReasonOptions []string `json:"terminatedReasonOptions"`
	//重启原因
	ResumeReasonOptions []string `json:"resumeReasonOptions"`
	//是否允许其他终止原因
	EnableTerminatedReasonOtherOption bool `json:"enableTerminatedReasonOtherOption"`
	//是否允许其他原因
	EnableResumeReasonOtherOption bool `json:"enableResumeReasonOtherOption"`
}

func NewStateFlow(key string) *StateFlow {
	wf := &StateFlow{
		Key:                               key,
		Uuid:                              uuid.NewString(),
		EnableTerminatedReasonOtherOption: true,
		EnableResumeReasonOtherOption:     true,
		TerminatedReasonOptions:           StateFlow_Default_TerminatedReasonOptions,
		ResumeReasonOptions:               StateFlow_Default_RestartResumeOptions,
	}

	return wf
}

func (s *StateFlow) GetNode(key string) *StateFlowNode {
	return s.GetNodeByKey(key)
}

func (s *StateFlow) GetNodeByKey(key string) *StateFlowNode {
	for _, node := range s.StateFlowNodes {
		if node.SubStateKey == key {
			return node
		}
	}

	return nil
}

func (s *StateFlow) GetNodes() []*StateFlowNode {
	return s.StateFlowNodes
}

func (s *StateFlow) ToString() string {
	r, _ := json.Marshal(s)
	return string(r)
}

func (s *StateFlow) IsValidCode(stateKey string) bool {
	for _, node := range s.StateFlowNodes {
		if stateKey == node.Key {
			return true
		}
	}
	return false
}

func (s *StateFlow) CanPass(code string, nextCode string) bool {
	for _, rule := range s.StateFlowTransitionRule {
		if rule.SourceStateKey == code && rule.TargetStateKey == nextCode {
			return true
		}
	}
	return false
}

func (s *StateFlow) GetInitStateNode() *StateFlowNode {
	for _, node := range s.StateFlowNodes {
		if node.IsInitState {
			return node
		}
	}
	return nil
}

func (s *StateFlow) GetFirstProcessNode() *StateFlowNode {
	for _, node := range s.StateFlowNodes {
		if !node.IsArchivedState {
			return node
		}
	}
	return nil
}

func (s *StateFlow) GetFirstArchivedNode() *StateFlowNode {
	for _, node := range s.StateFlowNodes {
		if node.IsArchivedState {
			return node
		}
	}
	return nil
}

func (s *StateFlow) GetNodeKeysByRoleKey(roleKey string) []string {
	var nodeKey []string
	for _, node := range s.StateFlowNodes {
		if node.Owner != nil {
			for _, v := range node.Owner.OwnerRole {
				if v.Key == roleKey {
					nodeKey = append(nodeKey, node.SubStateKey)
				}
			}
		}
	}
	return nodeKey
}

func (s *StateFlow) GetRoleKeysByNodeKey(nodeKey string) []string {
	var roleKeys []string
	for _, node := range s.StateFlowNodes {
		if node.SubStateKey == nodeKey {
			for _, v := range node.Owner.OwnerRole {
				roleKeys = append(roleKeys, v.Key)
			}
			break
		}
	}
	return roleKeys
}

func (s *StateFlow) GetAllRoleId() []int64 {
	var ids []int64
	for _, node := range s.StateFlowNodes {
		if node.Owner != nil {
			for _, confRole := range node.Owner.OwnerRole {
				ids = append(ids, cast.ToInt64(confRole.Id))
			}
		}
	}

	return stream.Unique(ids)
}

func (s *StateFlow) GetAllStatusId() []int64 {
	var ids []int64
	for _, node := range s.StateFlowNodes {
		ids = append(ids, cast.ToInt64(node.SubStateId))
	}

	return stream.Unique(ids)
}

func (s *StateFlow) GetAllRelatedUserId() []int64 {
	var ids []int64
	for _, node := range s.StateFlowNodes {
		ids = append(ids, node.GetRelatedUserIds()...)
	}

	return stream.Unique(ids)
}

type StateFlowNode struct {
	//UUID
	Uuid string `json:"uuid"`
	//状态节点编码
	Key string `json:"key"`
	//状态节点名称
	Name string `json:"name"`
	//对照任务的状态ID
	SubStateUuid string `json:"subStateUuid"`
	//对照任务的状态ID
	SubStateId string `json:"subStateId"`
	//对照任务的状态ID
	SubStateKey string `json:"subStateKey"`
	//对照任务的状态ID
	SubStateVal string `json:"subStateVal"`
	//是否归档状态
	IsArchivedState bool `json:"isArchivedState"`
	//是否启始状态
	IsInitState bool `json:"isInitState"`
	//可操作的角色权限
	OperationRole []string `json:"operationRole"`
	//关联角色
	Owner *OwnerConf `json:"owner"`
}

func (n *StateFlowNode) GetOwnerRoleId() string {
	if n.Owner != nil && len(n.Owner.OwnerRole) > 0 {
		return n.Owner.OwnerRole[0].Id
	}
	return "0"
}

func (n *StateFlowNode) GetOwnerRoleKey() string {
	if n.Owner != nil && len(n.Owner.OwnerRole) > 0 {
		return n.Owner.OwnerRole[0].Key
	}
	return ""
}

func (n *StateFlowNode) CheckOwnerRule(ownerIds []string) bool {
	if n.Owner != nil && n.Owner.UsageMode == UsageMode_Appointed {
		appointed, isOk := n.Owner.GetUsageModeVal().(*OwnerConf_UsageMode_Appointed)
		if !isOk {
			return false
		}

		for _, v := range appointed.AppointedOwner {
			if v.Type == FillOwnerType_User {
				if stream.ContainsArr(ownerIds, []string{cast.ToString(v.Value)}) {
					return true
				}
			}
		}
		return false
	}

	return true
}
func (n *StateFlowNode) GetRelatedUserIds() []int64 {
	var ids []int64

	if n.Owner != nil && n.Owner.UsageMode == UsageMode_Appointed {
		appointed, isOk := n.Owner.GetUsageModeVal().(*OwnerConf_UsageMode_Appointed)
		if !isOk {
			return nil
		}

		for _, v := range appointed.AppointedOwner {
			if v.Type == FillOwnerType_User {
				ids = append(ids, cast.ToInt64(v.Value))
			}
		}
	}

	if n.Owner != nil && n.Owner.UsageMode == UsageMode_None {
		appointed, isOk := n.Owner.GetUsageModeVal().(*OwnerConf_UsageMode_None)
		if !isOk {
			return nil
		}

		for _, v := range appointed.FillOwner {
			if v.Type == FillOwnerType_User {
				ids = append(ids, cast.ToInt64(v.Value))
			}
		}
	}

	return ids
}

func (n *StateFlowNode) GetRelatedStatusId() int64 {
	return cast.ToInt64(n.SubStateId)
}
