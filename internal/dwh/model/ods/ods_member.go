package ods

import (
	"go-cs/internal/dwh/pkg/model"
)

type OdsMember struct {
	model.OdsModel

	Id            int64  `protobuf:"zigzag64,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        int64  `protobuf:"zigzag64,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`                         //创建用户id
	SpaceId       int64  `protobuf:"zigzag64,3,opt,name=space_id,json=spaceId,proto3" json:"space_id,omitempty"`                      //空间Guid
	RoleId        int64  `protobuf:"zigzag64,4,opt,name=role_id,json=roleId,proto3" json:"role_id,omitempty"`                         //空间名称
	Remark        string `protobuf:"bytes,5,opt,name=remark,proto3" json:"remark,omitempty"`                                          //备注
	Ranking       int64  `protobuf:"varint,6,opt,name=ranking,proto3" json:"ranking,omitempty"`                                       //排序值
	Notify        int32  `protobuf:"varint,7,opt,name=notify,proto3" json:"notify,omitempty"`                                         //是否通知;0:否,1:是
	CreatedAt     int64  `protobuf:"zigzag64,8,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`                //创建时间
	UpdatedAt     int64  `protobuf:"zigzag64,9,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`                //更新时间
	DeletedAt     int64  `protobuf:"zigzag64,10,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`               //删除时间
	HistoryRoleId int64  `protobuf:"zigzag64,11,opt,name=history_role_id,json=historyRoleId,proto3" json:"history_role_id,omitempty"` //空间名称
}
