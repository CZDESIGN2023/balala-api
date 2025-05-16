package ods

import "go-cs/internal/dwh/pkg/model"

type OdsVersion struct {
	model.OdsModel

	Id            int64  `protobuf:"zigzag64,1,opt,name=id,proto3" json:"id,omitempty"`
	SpaceId       int64  `protobuf:"zigzag64,2,opt,name=space_id,json=spaceId,proto3" json:"space_id,omitempty"`                 //空间id
	VersionKey    string `protobuf:"bytes,10,opt,name=version_key,json=versionKey,proto3" json:"version_key,omitempty"`          //版本KEY
	VersionName   string `protobuf:"bytes,3,opt,name=version_name,json=versionName,proto3" json:"version_name,omitempty"`        //版本名称
	VersionStatus int64  `protobuf:"varint,4,opt,name=version_status,json=versionStatus,proto3" json:"version_status,omitempty"` //版本状态
	Remark        string `protobuf:"bytes,5,opt,name=remark,proto3" json:"remark,omitempty"`                                     //备注
	Ranking       int64  `protobuf:"zigzag64,6,opt,name=ranking,proto3" json:"ranking,omitempty"`                                //排序
	CreatedAt     int64  `protobuf:"zigzag64,7,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`           //创建时间
	UpdatedAt     int64  `protobuf:"zigzag64,8,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`           //更新时间
	DeletedAt     int64  `protobuf:"zigzag64,9,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`           //删除时间
}
