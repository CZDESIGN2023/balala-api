package search2

import (
	"fmt"
	"github.com/spf13/cast"
	v1 "go-cs/api/search/v1"
)

func CheckGroupField(list []*v1.GroupBy) error {
	for _, v := range list {
		switch QueryField(v.Field) {
		case UserId, WorkObjectId, SpaceId, Priority, Directors, WorkItemFlowId, WorkItemStatusId, VersionId:
		default:
			return fmt.Errorf("分组参数不支持 %v", v.Field)
		}
	}
	return nil
}

func GroupValue(field string, item *Model) []string {
	switch field {
	case "user_id":
		return []string{cast.ToString(item.UserId)}
	case "work_object_id":
		return []string{cast.ToString(item.WorkObjectId)}
	case "space_id":
		return []string{cast.ToString(item.SpaceId)}
	case "priority":
		return []string{cast.ToString(item.Priority)}
	case "directors":
		return item.Directors
	case "work_item_flow_id":
		return []string{cast.ToString(item.WorkItemFlowId)}
	case "version_id":
		return []string{cast.ToString(item.VersionId)}
	case "work_item_status_id":
		return []string{cast.ToString(item.WorkItemStatusId)}
	}
	return nil
}
