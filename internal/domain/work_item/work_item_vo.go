package work_item

import (
	"encoding/json"
	"go-cs/pkg/stream"
	"slices"

	"github.com/spf13/cast"
)

type Tags []string

func (t Tags) ToJsonString() string {
	if t == nil {
		return "[]"
	}
	v, _ := json.Marshal(t)
	return string(v)
}

func (t Tags) ToInt64s() []int64 {
	var vals []int64
	for _, v := range t {
		vals = append(vals, cast.ToInt64(v))
	}
	return vals
}

type Directors []string

func (t Directors) Clone() Directors {
	return slices.Clone(t)
}

func (t Directors) FormJsonString(s string) Directors {
	_ = json.Unmarshal([]byte(s), &t)
	return t
}

func (t Directors) ToStrings() []string {
	return t
}

func (t Directors) ToInt64s() []int64 {
	var vals []int64
	for _, v := range t {
		vals = append(vals, cast.ToInt64(v))
	}
	return vals
}

func (t Directors) ToJsonString() string {
	if t == nil {
		return "[]"
	}

	v, _ := json.Marshal(t)
	return string(v)
}

func (t Directors) Contains(userIds ...string) bool {
	return stream.ContainsArr(t, userIds)
}

type PlanTime struct {
	StartAt    int64 `json:"plan_start_at,omitempty"`
	CompleteAt int64 `json:"plan_complete_at,omitempty"`
}

type Restart struct {
	IsRestart     int32 ` json:"is_restart,omitempty"` //是否为重启任务
	RestartAt     int64 ` json:"restart_at,omitempty"` //重启时间
	RestartUserId int64 ` json:"restart_user_id,omitempty"`
}

type Resume struct {
	ResumeAt     int64 ` json:"resume_at,omitempty"`      //恢复时间
	ResumeUserId int64 ` json:"resume_user_id,omitempty"` //恢复操作人
}

type WorkItemStatus struct {
	Val string ` json:"work_item_status,omitempty"`     //任务状态val
	Key string ` json:"work_item_status_key,omitempty"` //工作项状态key
	Id  int64  ` json:"work_item_status_id,omitempty"`  //工作项状态id
}

type LastWorkItemStatus struct {
	LastAt int64  ` json:"last_status_at,omitempty"`  //历史状态更新时间
	Val    string ` json:"last_status,omitempty"`     //历史状态
	Key    string ` json:"last_status_key,omitempty"` //最后一次状态key
	Id     int64  ` json:"last_status_id,omitempty"`  //最后一次状态id
}

const (
	ICON_FLAG_Pic  = 1 << iota //有图片
	ICON_FLAG_Url  = 1 << iota //有图片
	ICON_FLAG_File = 1 << iota //有附件
)

const (
	ICON_FLAG_TYPE_Pic  = iota + 1 //有图片
	ICON_FLAG_TYPE_Url             //有链接
	ICON_FLAG_TYPE_File            //有附件
)

type IconFlag uint32 //图标标记
type IconFlagUpdate struct {
	Flag uint32
	Val  uint32
}

func (i *IconFlag) HasFlag(flag uint32) bool {
	return (uint32(*i) & flag) != 0
}

func (i *IconFlag) SetFlag(updates ...*IconFlagUpdate) {
	var v = uint32(*i)
	for _, update := range updates {
		if update.Val == 1 {
			v |= update.Flag
		} else {
			v &= ^update.Flag
		}
	}
	*i = IconFlag(v)
}

func (i *IconFlag) AddFlag(flags ...uint32) {
	var v = uint32(*i)
	for _, flag := range flags {
		v |= flag
	}
	*i = IconFlag(v)
}

func (i *IconFlag) RemoveFlag(flags ...uint32) {
	var v = uint32(*i)
	for _, flag := range flags {
		v &= (^flag)
	}
	*i = IconFlag(v)
}

func (i *IconFlag) ToFlags() []uint32 {
	bits := uint32(*i)
	var flags []uint32
	for _, flag := range i.all() {
		if flag&bits != 0 {
			flags = append(flags, flag)
		}
	}
	return flags
}

func (i *IconFlag) IsValidFlags(flags ...uint32) bool {
	allIconFlags := i.all()
	for _, flag := range flags {
		if !slices.Contains(allIconFlags, flag) {
			return false
		}
	}
	return true
}

func (i *IconFlag) all() []uint32 {
	return []uint32{
		ICON_FLAG_Pic,
		ICON_FLAG_Url,
		ICON_FLAG_File,
	}
}
