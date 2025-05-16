package search2

import (
	"cmp"
	"go-cs/internal/bean/vo"
	"go-cs/pkg/pprint"
	"slices"
	"testing"
)

var testCmpFuncMap = map[string]func(a, b *vo.SearchSpaceWorkItemGroupInfoResultItemVo) int{
	"space_id": func(a, b *vo.SearchSpaceWorkItemGroupInfoResultItemVo) int {
		return cmp.Compare(a.SpaceId, b.SpaceId)
	},
	"work_item_object": func(a, b *vo.SearchSpaceWorkItemGroupInfoResultItemVo) int {
		return cmp.Compare(a.WorkObjectId, b.WorkObjectId)
	},
	"work_item_id": func(a, b *vo.SearchSpaceWorkItemGroupInfoResultItemVo) int {
		return cmp.Compare(a.WorkItemId, b.WorkItemId)
	},
	"priority": func(a, b *vo.SearchSpaceWorkItemGroupInfoResultItemVo) int {
		return cmp.Compare(a.Priority, b.Priority)
	},
}

func TestSort(t *testing.T) {
	arr := []*vo.SearchSpaceWorkItemGroupInfoResultItemVo{
		{SpaceId: 1, WorkItemId: 1, WorkObjectId: 1, Priority: "p4"},
		{SpaceId: 1, WorkItemId: 2, WorkObjectId: 2, Priority: "p7"},
		{SpaceId: 2, WorkItemId: 3, WorkObjectId: 3, Priority: "p4"},
		{SpaceId: 2, WorkItemId: 4, WorkObjectId: 4, Priority: "p2"},
	}

	cmpFunc := MergeCmpFunc(
		testCmpFuncMap["space_id"],
		testCmpFuncMap["priority"],
	)

	slices.SortFunc(arr, cmpFunc)

	pprint.Println(arr)
}
