package search2

import (
	"fmt"
	v1 "go-cs/api/search/v1"
	"slices"
)

var sortFields = []QueryField{
	UserId,
	WorkItemId,
	SpaceId,
	WorkObjectId,
	WorkItemStatusId,
	Priority,
	WorkItemFlowId,
	VersionId,
}

func CheckSortField(list []*v1.Sort) error {
	for _, v := range list {
		if !slices.Contains(sortFields, QueryField(v.Field)) {
			return fmt.Errorf("排序参数不支持 %v", v.Field)
		}
	}
	return nil
}
