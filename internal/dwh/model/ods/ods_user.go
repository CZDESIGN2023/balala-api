package ods

import "go-cs/internal/dwh/pkg/model"

type OdsUser struct {
	model.OdsModel

	Id            int64  `protobuf:"zigzag64,1,opt,name=id,proto3" json:"id,omitempty"`
	UserName      string `protobuf:"bytes,2,opt,name=user_name,json=userName,proto3" json:"user_name,omitempty"`                      //用户名
	Mobile        string `protobuf:"bytes,3,opt,name=mobile,proto3" json:"mobile,omitempty"`                                          //中国手机不带国家代码，国际手机号格式为：国家代码-手机号
	UserNickname  string `protobuf:"bytes,4,opt,name=user_nickname,json=userNickname,proto3" json:"user_nickname,omitempty"`          //用户昵称
	UserPinyin    string `protobuf:"bytes,5,opt,name=user_pinyin,json=userPinyin,proto3" json:"user_pinyin,omitempty"`                //用户昵称拼音
	UserPassword  string `protobuf:"bytes,6,opt,name=user_password,json=userPassword,proto3" json:"user_password,omitempty"`          //登录密码;cmf_password加密
	UserSalt      string `protobuf:"bytes,7,opt,name=user_salt,json=userSalt,proto3" json:"user_salt,omitempty"`                      //加密盐
	UserStatus    int32  `protobuf:"zigzag32,8,opt,name=user_status,json=userStatus,proto3" json:"user_status,omitempty"`             //用户状态;0:禁用,1:正常,2:未验证
	UserEmail     string `protobuf:"bytes,9,opt,name=user_email,json=userEmail,proto3" json:"user_email,omitempty"`                   //用户登录邮箱
	Sex           int32  `protobuf:"zigzag32,10,opt,name=sex,proto3" json:"sex,omitempty"`                                            //性别;0:保密,1:男,2:女
	Avatar        string `protobuf:"bytes,11,opt,name=avatar,proto3" json:"avatar,omitempty"`                                         //用户头像
	Remark        string `protobuf:"bytes,12,opt,name=remark,proto3" json:"remark,omitempty"`                                         //备注
	Describe      string `protobuf:"bytes,13,opt,name=describe,proto3" json:"describe,omitempty"`                                     //描述信息
	LastLoginIp   string `protobuf:"bytes,14,opt,name=last_login_ip,json=lastLoginIp,proto3" json:"last_login_ip,omitempty"`          //最后登录ip
	LastLoginTime int64  `protobuf:"zigzag64,15,opt,name=last_login_time,json=lastLoginTime,proto3" json:"last_login_time,omitempty"` //最后登录时间
	Role          int64  `protobuf:"zigzag64,16,opt,name=role,proto3" json:"role,omitempty"`                                          //用户角色;0:普通用户,50:管理员,100:超级管理员
	CreatedAt     int64  `protobuf:"zigzag64,17,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`               //创建时间
	UpdatedAt     int64  `protobuf:"zigzag64,18,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`               //更新时间
	DeletedAt     int64  `protobuf:"zigzag64,19,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`               //删除时间
}
