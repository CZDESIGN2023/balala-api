package config

import (
	"encoding/json"
	"go-cs/pkg/stream"

	"github.com/google/uuid"
	"github.com/spf13/cast"
)

type NodePassMode string
type FlowDirection string
type SimulateEventType string
type NodeStartMode string
type UsageMode string
type FillOwnerType string

const (
	WorkflowNodeCode_Started = "started"
	WorkflowNodeCode_Ended   = "ended"

	WorkflowNodeStartMode_PreAllDone NodeStartMode = "pre_node_all_done"

	WorkflowNodePassMode_Auto   NodePassMode = "auto_confirm"
	WorkflowNodePassMode_Single NodePassMode = "single_user_confirm"

	WorkflowNodeFlowDirection_Reach FlowDirection = "reach"
	WorkflowNodeFlowDirection_Pass  FlowDirection = "pass"

	WorkflowSimulateEventType_OnReach       SimulateEventType = "ONREACH"
	WorkflowSimulateEventType_OnPass        SimulateEventType = "ONPASS"
	WorkflowSimulateEventType_OnStatechange SimulateEventType = "ONSTATECHANGE"
)

var (
	Default_CloseReasonOptions = []string{
		"已解决，待同步更新",
		"任务已变更，待同步更新",
		"关闭任务，转其他任务流程",
		"关闭任务，与其他任务重复",
	}

	Default_RebootReasonOptions = []string{
		"任务需要继续进行",
		"误操作",
	}

	Default_RollbackReasonOptions = []string{
		"开发 基础功能 的完成度，不足以完成期望测试",
		"开发 交互操作逻辑 的完成度，不足以完成期望测试",
		"开发 UI/视觉 的完成度，不足以完成期望测试",
	}

	Default_RestartReasonOptions = []string{
		"任务需要继续进行",
		"误操作",
	}

	Default_TerminatedReasonOptions = []string{
		"取消任务，综合考虑且已同步现在先不做了",
		"任务重复/合并，与其他在做任务一同推进",
		"取消任务，转其他任务流程",
	}
)

type WorkFlowNodeStateEvent func(nodeKey string, evType SimulateEventType, changeState *WorkFlowSubState) error

type DefaultReasonOption struct {
	Name    string   `json:"name"`
	Options []string `json:"options"`
}

type WorkFlowDefaultConf struct {
	ReasonOptions []DefaultReasonOption `json:"reasonOptions"`
}

func GetWorkFlowDefaultConf() *WorkFlowDefaultConf {
	return &WorkFlowDefaultConf{
		ReasonOptions: []DefaultReasonOption{
			{Name: "terminated", Options: Default_TerminatedReasonOptions},
			{Name: "reboot", Options: Default_RebootReasonOptions},
			{Name: "resume", Options: Default_RestartReasonOptions},
			{Name: "node.rollback", Options: Default_RollbackReasonOptions},
			{Name: "node.restart", Options: Default_RestartReasonOptions},
			{Name: "node.close", Options: Default_CloseReasonOptions},
		},
	}
}

type WorkFlow struct {

	//客户端不一定要用
	Uuid string `json:"uuid"`
	//编码 客户端不一定要用
	Key         string                `json:"key"`
	Nodes       []*WorkFlowNode       `json:"nodes"`
	Connections []*WorkFlowConnection `json:"connections"`
	//终止原因
	TerminatedReasonOptions []string `json:"terminatedReasonOptions"`
	//是否允许其他终止原因
	EnableTerminatedReasonOtherOption bool `json:"enableTerminatedReasonOtherOption"`
	//重启原因
	RebootReasonOptions []string `json:"rebootReasonOptions"`
	//是否允许其他原因
	EnableRebootReasonOtherOption bool `json:"enableRebootReasonOtherOption"`
	//重启原因
	ResumeReasonOptions []string `json:"resumeReasonOptions"`
	//是否允许其他原因
	EnableResumeReasonOtherOption bool `json:"enableResumeReasonOtherOption"`
	//FormFields
	FormFields WorkFlowFormFields `json:"formFields"`
}

func NewWorkFlow(key string) *WorkFlow {
	wf := &WorkFlow{
		Key:                               key,
		Uuid:                              uuid.NewString(),
		Nodes:                             make([]*WorkFlowNode, 0),
		Connections:                       make([]*WorkFlowConnection, 0),
		EnableTerminatedReasonOtherOption: true,
		EnableResumeReasonOtherOption:     true,
		EnableRebootReasonOtherOption:     true,
		RebootReasonOptions:               make([]string, 0),
		TerminatedReasonOptions:           make([]string, 0),
		FormFields:                        make(WorkFlowFormFields, 0),
	}

	wf.FillDefaultReasonOptions()

	return wf
}

// 按节点顺序的连接, 现在只按一条路径来实现
func (w *WorkFlow) ConnectionNodeLink() []*WorkFlowNode {

	type LinkNode struct {
		Node     *WorkFlowNode
		NextNode *WorkFlowNode
	}

	nodeMap := w.GetNodeMap()

	linkNodeMap := make(map[string]*LinkNode)
	for _, connection := range w.Connections {
		linkNodeMap[connection.StartNode] = &LinkNode{
			NextNode: nodeMap[connection.EndNode],
			Node:     nodeMap[connection.StartNode],
		}
	}

	var nextLinkNodeFn func(linkNode *LinkNode) []*LinkNode
	nextLinkNodeFn = func(linkNode *LinkNode) []*LinkNode {
		linkNodes := make([]*LinkNode, 0)
		if linkNode == nil {
			return linkNodes
		}

		linkNodes = append(linkNodes, linkNode)
		if linkNode.NextNode != nil {
			nextLinkNode := nextLinkNodeFn(linkNodeMap[linkNode.NextNode.Key])
			linkNodes = append(linkNodes, nextLinkNode...)
		}
		return linkNodes
	}

	startLinkNode := linkNodeMap[WorkflowNodeCode_Started]
	linkNodes := nextLinkNodeFn(startLinkNode)
	flowNodes := make([]*WorkFlowNode, 0)
	for _, linkNode := range linkNodes {
		flowNodes = append(flowNodes, linkNode.Node)
	}
	return flowNodes

}

func (w *WorkFlow) FillDefaultReasonOptions() {
	w.RebootReasonOptions = Default_RebootReasonOptions
	w.TerminatedReasonOptions = Default_TerminatedReasonOptions
	w.ResumeReasonOptions = Default_RestartReasonOptions
}

func MustFormWorkFlowJson(jsonStr string) (*WorkFlow, error) {
	wf := &WorkFlow{}
	err := json.Unmarshal([]byte(jsonStr), &wf)
	return wf, err
}

func (w *WorkFlow) ToString() string {
	r, _ := json.Marshal(w)
	return string(r)
}

func (w *WorkFlow) GetNodesByRole(role string) []*WorkFlowNode {
	var ret []*WorkFlowNode
	for _, node := range w.Nodes {
		if node.Owner.CheckOwnerRole(role) {
			ret = append(ret, node)
		}
	}

	return ret
}
func (w *WorkFlow) GetAllRoleId() []int64 {
	var ids []int64
	for _, node := range w.Nodes {
		if !node.IsStartedNode() && !node.IsEndedNode() && node.Owner != nil {
			for _, confRole := range node.Owner.OwnerRole {
				ids = append(ids, cast.ToInt64(confRole.Id))
			}
		}
	}

	return stream.Unique(ids)
}

func (w *WorkFlow) GetAllStatusId() []int64 {
	var ids []int64
	for _, node := range w.Nodes {
		if node.OnPass != nil {
			for _, v := range node.OnPass {
				if v.TargetSubState != nil {
					ids = append(ids, cast.ToInt64(v.TargetSubState.Id))
				}
			}
		}

		if node.OnReach != nil {
			for _, v := range node.OnReach {
				if v.TargetSubState != nil {
					ids = append(ids, cast.ToInt64(v.TargetSubState.Id))
				}
			}
		}
	}

	return stream.Unique(ids)
}

func (w *WorkFlow) GetAllRelatedUserId() []int64 {
	var ids []int64
	for _, node := range w.Nodes {
		ids = append(ids, node.GetRelatedUserIds()...)
	}

	return stream.Unique(ids)
}

type WorkFlowFormFields []*WorkFlowFormField

type WorkFlowFormField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type WorkFlowFormFieldSetting struct {
	Filed []*WorkFlowFormField `json:"filed"`
}

type WorkFlowSubState struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Uuid string `json:"uuid"`
	Val  string `json:"val"`
}

type WorkFlowNode struct {
	Uuid string `json:"id"`
	//名称
	Name string `json:"name"`
	//编码 系统预设 started 开始节点  ended 结束节点
	Key string `json:"key"`
	//节点被激活模式 pre_node_all_done 前置节点均完成
	StartMode NodeStartMode `json:"startMode"`
	// 流程开始的第一个节点 "started"
	BelongStatus string `json:"belongStatus"`
	//节点完成方式 是否需要指定负责人
	NeedDoneOperator bool `json:"needDoneOperator"`
	//节点完成方式可操作的相关角色 _node_owner 节点负责人 _space_manager 空间管理员 _creator 任务创建人 _space_editor 空间编辑人员   ["_node_owner", "_space_manager", "_creator", "_space_editor"]
	DoneOperationRole []string `json:"doneOperationRole"`
	//确认操作按钮名称
	DoneOperationDisplayName string `json:"doneOperationDisplayName"`
	//节点完成方式可操作的原因选项
	DoneOperationReasonOptions []string `json:"doneOperationReasonOptions"`
	//节点完成模式 auto_confirm 自动完成[忽略角色，直接完成] single_user_confirm 单人完成(需配合角色控制) all_user_confirm 多人确认完成(所有节点负责人)
	PassMode NodePassMode `json:"passMode"`
	//到达事件
	OnReach []*WorkFlowNodeEvent `json:"onReach"`
	//完成事件
	OnPass []*WorkFlowNodeEvent `json:"onPass"`

	//是否可关闭
	EnableClose bool `json:"enableClose"`
	//是否支持回滚
	EnableRollback bool `json:"enableRollback"`
	// 是否强制填写计划完成时间
	ForcePlanTime bool `json:"forcePlanTime"`

	//回滚原因选项
	RollbackReasonOptions []string `json:"rollbackReasonOptions"`
	//是否允许其他回滚原因选项
	EnableRollbackReasonOtherOption bool `json:"enableRollbackReasonOtherOption"`

	//关闭原因选项
	CloseReasonOptions []string `json:"closeReasonOptions"`
	//是否允许其他原因选项
	EnableCloseReasonOtherOption bool `json:"enableCloseReasonOtherOption"`

	//重启原因选项
	RestartReasonOptions []string `json:"restartReasonOptions"`
	//是否允许其他原因选项
	EnableRestartReasonOtherOption bool `json:"enableRestartReasonOtherOption"`

	//负责人分配方式
	Owner *OwnerConf `json:"owner"`
}

func NewWorkFlowNode(name string, key string) *WorkFlowNode {

	n := &WorkFlowNode{
		CloseReasonOptions:              make([]string, 0),
		RestartReasonOptions:            make([]string, 0),
		RollbackReasonOptions:           make([]string, 0),
		EnableRestartReasonOtherOption:  true,
		EnableCloseReasonOtherOption:    true,
		EnableRollbackReasonOtherOption: true,
		Uuid:                            uuid.NewString(),
		Key:                             key,
		Name:                            name,
		OnReach:                         make([]*WorkFlowNodeEvent, 0),
		OnPass:                          make([]*WorkFlowNodeEvent, 0),
	}

	n.Owner = &OwnerConf{
		UsageMode:  UsageMode_None,
		ForceOwner: true,
		OwnerRole:  make([]*OwnerConf_Role, 0),
	}
	n.PassMode = WorkflowNodePassMode_Single
	n.NeedDoneOperator = true
	n.DoneOperationRole = []string{"_node_owner", "_space_manager", "_creator", "_space_editor"}

	return n
}

func (n *WorkFlowNode) FillDefaultReasonOptions() {
	n.CloseReasonOptions = Default_CloseReasonOptions
	n.RestartReasonOptions = Default_RestartReasonOptions
	n.RollbackReasonOptions = Default_RollbackReasonOptions
}

type WorkFlowNodeEvent struct {
	//事件类型 changeStoryStage 状态流转
	EventType string `json:"eventType"`
	// 状态流转附件条件
	Condition any `json:"condition"`
	// 状态流转为
	TargetSubState *WorkFlowSubState `json:"targetSubState"`
}

// 节点关系线
type WorkFlowConnection struct {
	//开始节点编码
	StartNode string `json:"startNode"`
	//结束节点编码
	EndNode string `json:"endNode"`
}

// 获取下个要进入的节点
func (w *WorkFlow) GetNextNodeFormConnections(nodeKey string) []*WorkFlowNode {
	var nodes []*WorkFlowNode
	for i := 0; i < len(w.Connections); i++ {
		if w.Connections[i].StartNode == nodeKey {
			node := w.GetNode(w.Connections[i].EndNode)
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (w *WorkFlow) GetStartedNode() *WorkFlowNode {
	for i := 0; i < len(w.Nodes); i++ {
		if w.Nodes[i].IsStartedNode() {
			return w.Nodes[i]
		}
	}
	return nil
}

// 获取节点
func (w *WorkFlow) GetNode(key string) *WorkFlowNode {
	for i := 0; i < len(w.Nodes); i++ {
		node := w.Nodes[i]
		if node.Key == key {
			return node
		}
	}
	return nil
}

func (w *WorkFlow) GetNodeMap() map[string]*WorkFlowNode {
	return stream.ToMap(w.Nodes, func(_ int, v *WorkFlowNode) (string, *WorkFlowNode) {
		return v.Key, v
	})
}

func (w *WorkFlow) GetNodeKeysByRoleKey(roleKey string) []string {
	var nodeKey []string
	for _, node := range w.Nodes {
		if node.Owner != nil {
			for _, v := range node.Owner.OwnerRole {
				if v.Key == roleKey {
					nodeKey = append(nodeKey, node.Key)
				}
			}
		}
	}
	return nodeKey
}

func (w *WorkFlow) GetRoleKeysByNodeKey(nodeKey string) []string {
	var roleKeys []string
	for _, node := range w.Nodes {
		if node.Key == nodeKey {
			for _, v := range node.Owner.OwnerRole {
				roleKeys = append(roleKeys, v.Key)
			}
			break
		}
	}
	return roleKeys
}

// 是否为需要填写负责人的节点
func (n *WorkFlowNode) IsNeedOwnerCharge() bool {
	return n.PassMode != WorkflowNodePassMode_Auto || (n.Owner != nil && n.Owner.ForceOwner)
}

func (n *WorkFlowNode) IsStartedNode() bool {
	return n.Key == WorkflowNodeCode_Started
}

func (n *WorkFlowNode) IsEndedNode() bool {
	return n.Key == WorkflowNodeCode_Ended
}

func (n *WorkFlowNode) GetOwnerRoleId() string {
	if n.Owner != nil && len(n.Owner.OwnerRole) > 0 {
		return n.Owner.OwnerRole[0].Id
	}
	return "0"
}

func (n *WorkFlowNode) GetOwnerRoleKey() string {
	if n.Owner != nil && len(n.Owner.OwnerRole) > 0 {
		return n.Owner.OwnerRole[0].Key
	}
	return ""
}

func (n *WorkFlowNode) IsAutoPass() bool {
	return n.PassMode == WorkflowNodePassMode_Auto
}

func (n *WorkFlowNode) CheckOwnerRule(ownerIds []string) bool {
	if !n.IsNeedOwnerCharge() {
		return true
	}

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

func (n *WorkFlowNode) GetRelatedUserIds() []int64 {
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
	} else if n.Owner != nil && n.Owner.UsageMode == UsageMode_None {
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

func (n *WorkFlowNode) GetRelatedStatusId() int64 {
	for _, v := range n.OnReach {
		if v.TargetSubState != nil {
			return cast.ToInt64(v.TargetSubState.Id)
		}
	}

	return 0
}
