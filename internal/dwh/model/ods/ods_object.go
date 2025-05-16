package ods

import "go-cs/internal/dwh/pkg/model"

type OdsObject struct {
	model.OdsModel

	Id               int64  `protobuf:"zigzag64,1,opt,name=id,proto3" json:"id,omitempty"`
	SpaceId          int64  `protobuf:"zigzag64,2,opt,name=space_id,json=spaceId,proto3" json:"space_id,omitempty"`                              //空间id
	UserId           int64  `protobuf:"zigzag64,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`                                 //创建用户id
	WorkObjectGuid   string `protobuf:"bytes,4,opt,name=work_object_guid,json=workObjectGuid,proto3" json:"work_object_guid,omitempty"`          //空间Guid
	WorkObjectName   string `protobuf:"bytes,5,opt,name=work_object_name,json=workObjectName,proto3" json:"work_object_name,omitempty"`          //空间名称
	WorkObjectStatus int32  `protobuf:"zigzag32,6,opt,name=work_object_status,json=workObjectStatus,proto3" json:"work_object_status,omitempty"` //空间状态;0:禁用,1:正常,2:未验证
	Remark           string `protobuf:"bytes,7,opt,name=remark,proto3" json:"remark,omitempty"`                                                  //备注
	Describe         string `protobuf:"bytes,8,opt,name=describe,proto3" json:"describe,omitempty"`                                              //描述信息
	Ranking          int64  `protobuf:"varint,9,opt,name=ranking,proto3" json:"ranking,omitempty"`                                               //排序值
	CreatedAt        int64  `protobuf:"zigzag64,10,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`                       //创建时间
	UpdatedAt        int64  `protobuf:"zigzag64,11,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`                       //更新时间
	DeletedAt        int64  `protobuf:"zigzag64,12,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`                       //删除时间
}
