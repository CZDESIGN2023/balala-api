package ods

import "go-cs/internal/dwh/pkg/model"

type OdsWitemFlowNode struct {
	model.OdsModel

	Id                  int64  `protobuf:"zigzag64,1,opt,name=id,proto3" json:"id,omitempty"`
	WorkItemId          int64  `protobuf:"zigzag64,2,opt,name=work_item_id,json=workItemId,proto3" json:"work_item_id,omitempty"` //任务id
	WorkItemRoleId      int64  `protobuf:"zigzag64,25,opt,name=work_item_role_id,json=workItemRoleId,proto3" json:"work_item_role_id,omitempty"`
	WorkItemRoleKey     string `protobuf:"bytes,26,opt,name=work_item_role_key,json=workItemRoleKey,proto3" json:"work_item_role_key,omitempty"`
	FlowNodeUuid        string `protobuf:"bytes,3,opt,name=flow_node_uuid,json=flowNodeUuid,proto3" json:"flow_node_uuid,omitempty"`                        //节点uuid
	FlowNodeStatus      int32  `protobuf:"zigzag32,4,opt,name=flow_node_status,json=flowNodeStatus,proto3" json:"flow_node_status,omitempty"`               //节点状态; 0:未定义 , 1: 未开启 2: 进行中 3:已完成
	FlowNodeCode        string `protobuf:"bytes,5,opt,name=flow_node_code,json=flowNodeCode,proto3" json:"flow_node_code,omitempty"`                        //节点编码
	FlowNodePassed      int32  `protobuf:"zigzag32,6,opt,name=flow_node_passed,json=flowNodePassed,proto3" json:"flow_node_passed,omitempty"`               //节点是否通过
	FlowNodeReached     int32  `protobuf:"zigzag32,7,opt,name=flow_node_reached,json=flowNodeReached,proto3" json:"flow_node_reached,omitempty"`            //节点是否到达
	FlowModeVersion     string `protobuf:"bytes,9,opt,name=flow_mode_version,json=flowModeVersion,proto3" json:"flow_mode_version,omitempty"`               //[已废弃]模式版本
	FlowModeCode        string `protobuf:"bytes,10,opt,name=flow_mode_code,json=flowModeCode,proto3" json:"flow_mode_code,omitempty"`                       //[已废弃]模式编码
	StartAt             int64  `protobuf:"zigzag64,11,opt,name=start_at,json=startAt,proto3" json:"start_at,omitempty"`                                     //开始时间
	FinishAt            int64  `protobuf:"zigzag64,12,opt,name=finish_at,json=finishAt,proto3" json:"finish_at,omitempty"`                                  //结束时间
	Directors           string `protobuf:"bytes,13,opt,name=directors,proto3" json:"directors,omitempty"`                                                   //节点负责人
	CreatedAt           int64  `protobuf:"zigzag64,14,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`                               //创建时间
	UpdatedAt           int64  `protobuf:"zigzag64,15,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`                               //更新时间
	DeletedAt           int64  `protobuf:"zigzag64,16,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`                               //删除时间
	SpaceId             int64  `protobuf:"zigzag64,17,opt,name=space_id,json=spaceId,proto3" json:"space_id,omitempty"`                                     //空间id
	PlanStartAt         int64  `protobuf:"zigzag64,18,opt,name=plan_start_at,json=planStartAt,proto3" json:"plan_start_at,omitempty"`                       //计 划开始时间
	PlanCompleteAt      int64  `protobuf:"zigzag64,19,opt,name=plan_complete_at,json=planCompleteAt,proto3" json:"plan_complete_at,omitempty"`              //计划结束时间
	FlowId              int64  `protobuf:"zigzag64,20,opt,name=flow_id,json=flowId,proto3" json:"flow_id,omitempty"`                                        //流程Id
	FlowKey             string `protobuf:"bytes,21,opt,name=flow_key,json=flowKey,proto3" json:"flow_key,omitempty"`                                        //流程key
	FlowMode            string `protobuf:"bytes,22,opt,name=flow_mode,json=flowMode,proto3" json:"flow_mode,omitempty"`                                     //模式类型
	FlowTemplateId      int64  `protobuf:"varint,23,opt,name=flow_template_id,json=flowTemplateId,proto3" json:"flow_template_id,omitempty"`                //模式版本
	FlowTemplateVersion int64  `protobuf:"varint,24,opt,name=flow_template_version,json=flowTemplateVersion,proto3" json:"flow_template_version,omitempty"` //模式编码
}
