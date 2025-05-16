package ods

import "go-cs/internal/dwh/pkg/model"

type OdsSpace struct {
	model.OdsModel

	Id          int64  `protobuf:"zigzag64,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId      int64  `protobuf:"zigzag64,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`                //创建用户id
	SpaceGuid   string `protobuf:"bytes,3,opt,name=space_guid,json=spaceGuid,proto3" json:"space_guid,omitempty"`          //空间Guid
	SpaceName   string `protobuf:"bytes,4,opt,name=space_name,json=spaceName,proto3" json:"space_name,omitempty"`          //空间名称
	SpaceStatus int32  `protobuf:"zigzag32,5,opt,name=space_status,json=spaceStatus,proto3" json:"space_status,omitempty"` //空间状态;0:禁用,1:正常,2:未验证
	Remark      string `protobuf:"bytes,6,opt,name=remark,proto3" json:"remark,omitempty"`                                 //备注
	Describe    string `protobuf:"bytes,7,opt,name=describe,proto3" json:"describe,omitempty"`                             //描述信息
	Notify      int64  `protobuf:"zigzag64,8,opt,name=notify,proto3" json:"notify,omitempty"`                              //是否开启通知;0:否,1:是
	CreatedAt   int64  `protobuf:"zigzag64,9,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`       //创建时间
	UpdatedAt   int64  `protobuf:"zigzag64,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`      //更新时间
	DeletedAt   int64  `protobuf:"zigzag64,11,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`      //删除时间
}
