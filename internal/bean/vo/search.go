package vo

import (
	"go-cs/api/comm"
	db "go-cs/internal/bean/biz"
	"strconv"

	"gorm.io/datatypes"
)

type SearchSpaceWorkItemGroupInfoVo struct {
	SpaceIds       []int64                        `json:"space_id"`
	UserId         int64                          `json:"user_id"`
	Pagination     *comm.Pagination               `json:"pagination"`
	ConditionGroup *comm.SearchOptConditionGroup  `json:"condition_group"`
	Sorts          []*comm.SearchOptConditionSort `json:"sorts"`
	Groups         []*comm.SearchOptGroup         `json:"groups"`
}

type SearchSpaceWorkItemGroupInfoResultItemVo struct {
	WorkItemId         int64  `json:"work_item_id,omitempty"`
	WorkObjectId       int64  `json:"work_object_id,omitempty"`
	WorkObjectName     string `json:"work_object_name,omitempty"`
	WorkItemName       string `json:"work_item_name,omitempty"`
	ParentWorkItemName string `json:"parent_work_item_name,omitempty"`
	ParentWorkItemId   int64  `json:"parent_work_item_id,omitempty"`
	IsMainWorkItem     int64  `json:"is_main_work_item"`
	SpaceId            int64  `json:"space_id,omitempty"`
	SpaceName          string `json:"space_name,omitempty"`
	DirectorName       string `json:"director_name,omitempty"`
	DirectorUserId     int64  `json:"director_user_id,omitempty"`
	DirectorId         int64  `json:"director_id,omitempty"`
	Priority           string `json:"priority,omitempty"`
	TagId              int64  `json:"tag_id,omitempty"`
	TagName            string `json:"tag_name,omitempty"`
	SortWeightGroup    string
}

func (p *SearchSpaceWorkItemGroupInfoResultItemVo) void() {

}

func (p *SearchSpaceWorkItemGroupInfoResultItemVo) GetFiledStrValue(filedName string) string {
	switch filedName {
	case "work_item_id":
		return strconv.FormatInt(p.WorkItemId, 10)
	case "work_object_id":
		return strconv.FormatInt(p.WorkObjectId, 10)
	case "parent_work_item_id":
		return strconv.FormatInt(p.ParentWorkItemId, 10)
	case "work_object_name":
		return p.WorkObjectName
	case "work_item_name":
		return p.WorkItemName
	case "parent_work_item_name":
		return p.ParentWorkItemName
	case "space_name":
		return p.SpaceName
	case "space_id":
		return strconv.FormatInt(p.SpaceId, 10)
	case "director_id":
		return strconv.FormatInt(p.DirectorId, 10)
	case "priority":
		return p.Priority
	case "director_name":
		return p.DirectorName
	case "director_user_id":
		return strconv.FormatInt(p.DirectorUserId, 10)
	case "is_main_work_item":
		return strconv.FormatInt(p.IsMainWorkItem, 10)
	default:
		return ""
	}
}

type SearchMySpaceWorkItemGroupInfoViewResultVo struct {
	TotalNum int64
	List     []*SearchMySpaceWorkItemGroupInfoViewResultListItemVo
}

type SearchMySpaceWorkItemGroupInfoViewResultListItemVo struct {
	GroupInfo []*SearchMySpaceWorkItemGroupInfoViewResultListItemGroupInfoVo
	WorkItems []*SearchMySpaceWorkItemGroupInfoViewResultListItemWorkItemVo
}

type SearchMySpaceWorkItemGroupInfoViewResultListItemGroupInfoVo struct {
	FieldKey    string
	DisplayName string
}

type SearchMySpaceWorkItemGroupInfoViewResultListItemWorkItemVo struct {
	WorkItemId int64
	WorkItems  []*SearchMySpaceWorkItemGroupInfoViewResultListItemWorkItemVo
}

type SearchSpaceWorkTaskSimpleInfoModel struct {
	Id             int64  `query:"work_item_id" db:"id" gorm:"column:id" json:"id,omitempty"`
	Pid            int64  `query:"pid" db:"pid" gorm:"column:pid" json:"pid,omitempty"`
	SpaceId        int64  `query:"space_id" db:"space_id" gorm:"column:space_id" json:"space_id,omitempty"`
	UserId         int64  `query:"user_id" db:"user_id" gorm:"column:user_id" json:"user_id,omitempty"`
	WorkItemType   int64  `query:"work_item_type" db:"work_item_type" gorm:"column:work_item_type" json:"work_item_type,omitempty"`
	WorkObjectId   int64  `query:"work_object_id" db:"work_object_id" gorm:"column:work_object_id" json:"work_object_id,omitempty"`
	WorkItemGuid   string `query:"work_item_guid" db:"work_item_guid" gorm:"column:work_item_guid" json:"work_item_guid,omitempty"`
	WorkItemName   string `query:"work_item_name" db:"work_item_name" gorm:"column:work_item_name" json:"work_item_name,omitempty"`
	WorkItemStatus int32  `query:"work_item_status" db:"work_item_status" gorm:"column:work_item_status" json:"work_item_status,omitempty"`
	CreatedAt      int64  `query:"created_at" db:"created_at" gorm:"column:created_at" dt:"date" json:"created_at,omitempty"`
	UpdatedAt      int64  `query:"updated_at" db:"updated_at" gorm:"column:updated_at" dt:"date" json:"updated_at,omitempty"`
	DeletedAt      int64  `query:"deleted_at" db:"deleted_at" gorm:"column:deleted_at" dt:"date" json:"deleted_at,omitempty"`

	// 以下为doc内字段
	//Doc            datatypes.JSON              `query:"doc" db:"doc" gorm:"column:doc"` //这个字段不能直接读取
	ProcessRate    int32                       `query:"process_rate" db:"doc->'$.process_rate'" gorm:"column:process_rate" json:"process_rate,omitempty"`
	Priority       string                      `query:"priority" db:"doc->>'$.priority'" gorm:"column:priority" json:"priority,omitempty"`
	PlanStartAt    int64                       `query:"plan_start_at" db:"doc->'$.plan_start_at'" gorm:"column:plan_start_at"  dt:"date" json:"plan_start_at,omitempty"`
	PlanCompleteAt int64                       `query:"plan_complete_at" db:"doc->'$.plan_complete_at'" gorm:"column:plan_complete_at"  dt:"date" json:"plan_complete_at,omitempty"`
	Directors      datatypes.JSONSlice[string] `query:"directors" db:"doc->'$.directors'" gorm:"column:directors" dt:"multi-user" json:"directors,omitempty"` //当前负责人
}

type OpLogPaginationSearchVo struct {
	Pos  int64
	Size int

	SpaceIds          []int64
	OperatorType      int
	ModuleType        int
	ModuleId          int64
	OperId            int64
	IncludeModuleType []int
}

type SpaceWorkItemDb struct {
	db.SpaceWorkItemV2
	db.SpaceWorkItemDocV2
}
