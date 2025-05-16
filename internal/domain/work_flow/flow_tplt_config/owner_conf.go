package config

import "github.com/mitchellh/mapstructure"

const (
	UsageMode_None      UsageMode = "none"
	UsageMode_Appointed UsageMode = "appointed"

	FillOwnerType_User FillOwnerType = "user" //用户id
	FillOwnerType_Role FillOwnerType = "role" //角色
)

// 节点 负责人分配方式配置
type OwnerConf struct {
	//是否必须需要负责人
	ForceOwner bool `json:"forceOwner"`
	//负责人分配方式 none 不指定 appointed 指定负责人
	UsageMode UsageMode `json:"usageMode"`
	//存储一些特别条件的信息，比如默认指定负责人
	Value any `json:"value"`
	//负责人角色
	OwnerRole []*OwnerConf_Role `json:"ownerRole"`
}

func (w *OwnerConf) IsAppointedUsageMode() bool {
	return w.UsageMode == UsageMode_Appointed
}

func (w *OwnerConf) CheckOwnerRole(role string) bool {
	for _, v := range w.OwnerRole {
		if role == v.Key {
			return true
		}
	}
	return false
}

func (w *OwnerConf) GetAppointedUsageModeVal() *OwnerConf_UsageMode_Appointed {
	val := w.GetUsageModeVal()
	v, _ := val.(*OwnerConf_UsageMode_Appointed)
	return v
}

func (w *OwnerConf) GetNoneUsageModeVal() *OwnerConf_UsageMode_None {
	val := w.GetUsageModeVal()
	v, _ := val.(*OwnerConf_UsageMode_None)
	return v
}

func (w *OwnerConf) GetUsageModeVal() any {
	var v any
	switch w.UsageMode {
	case UsageMode_None:
		v = &OwnerConf_UsageMode_None{}
	case UsageMode_Appointed:
		v = &OwnerConf_UsageMode_Appointed{}
	}

	mapstructure.Decode(w.Value, v)
	return v
}

type OwnerConf_Role struct {
	Id   string `json:"id"`
	Key  string `json:"key"`
	Uuid string `json:"uuid"`
}

type OwnerConf_UsageMode_FillOwner struct {
	// type是user，值为用户id
	// type为role，任务创建人: _creator
	Type  FillOwnerType `json:"type"`
	Value any           `json:"value"`
}

func (w *OwnerConf_UsageMode_FillOwner) IsUserType() bool {
	return w.Type == FillOwnerType_User
}

// 不指定，可以通过规则填充，也可以通过指定用户填充
type OwnerConf_UsageMode_None struct {
	//自动填充指定的负责人 一组用户id
	FillOwner []*OwnerConf_UsageMode_FillOwner `json:"fillOwner"`
}

func (w *OwnerConf_UsageMode_None) Contains(typ FillOwnerType, val string) bool {
	for _, v := range w.FillOwner {
		if v.Type == typ && v.Value == val {
			return true
		}
	}

	return false
}

// 指定可选用户，默认填充的也只能从指定用户里面做选择
type OwnerConf_UsageMode_Appointed struct {
	//自动填充指定的负责人 一组用户id
	FillOwner []*OwnerConf_UsageMode_FillOwner `json:"fillOwner"`
	//指定可选择的负责人
	AppointedOwner []*OwnerConf_UsageMode_FillOwner `json:"appointedOwner"`
}

func (w *OwnerConf_UsageMode_Appointed) Contains(typ FillOwnerType, val string) bool {
	for _, v := range w.FillOwner {
		if v.Type == typ && v.Value == val {
			return true
		}
	}

	for _, v := range w.AppointedOwner {
		if v.Type == typ && v.Value == val {
			return true
		}
	}

	return false
}
