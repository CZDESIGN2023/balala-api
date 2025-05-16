package biz

import (
	"cmp"
	"context"
	"fmt"
	v1 "go-cs/api/search/v1"
	"go-cs/internal/consts"
	search_repo "go-cs/internal/domain/search/repo"
	search_es2 "go-cs/internal/domain/search/search_es"
	"go-cs/internal/domain/space"
	space_repo "go-cs/internal/domain/space/repo"
	"go-cs/internal/domain/space_member"
	member_repo "go-cs/internal/domain/space_member/repo"
	workObj "go-cs/internal/domain/space_work_object"
	workObj_repo "go-cs/internal/domain/space_work_object/repo"
	workVersion "go-cs/internal/domain/space_work_version"
	workVersion_repo "go-cs/internal/domain/space_work_version/repo"
	"go-cs/internal/domain/user"
	user_repo "go-cs/internal/domain/user/repo"
	"go-cs/internal/domain/work_flow"
	work_flow_repo "go-cs/internal/domain/work_flow/repo"
	witem_repo "go-cs/internal/domain/work_item/repo"
	"go-cs/internal/domain/work_item_status"
	workItemStatus_repo "go-cs/internal/domain/work_item_status/repo"
	esV8 "go-cs/internal/utils/es/v8"
	"go-cs/pkg/stream"
	"go-cs/pkg/stream/tuple"
	"slices"
	"strconv"

	"github.com/spf13/cast"
	"golang.org/x/exp/maps"
)

type SearchEsDataCtx struct {
	ctx context.Context

	userId   int64
	SpaceIds []int64

	childOrderMap map[string]string

	flattedGroups [][]*search_es2.Group // 扁平化后的分组

	objectMap    map[int64]*workObj.SpaceWorkObject
	userMap      map[int64]*user.User
	spaceMap     map[int64]*space.Space
	versionMap   map[int64]*workVersion.SpaceWorkVersion
	statusMap    map[int64]*work_item_status.WorkItemStatusItem
	statusValMap map[string]*work_item_status.WorkItemStatusItem
	flowMap      map[int64]*work_flow.WorkFlow

	workItemChainMap map[int64]*v1.WorkItem

	spaceRepo      space_repo.SpaceRepo
	userRepo       user_repo.UserRepo
	workObjectRepo workObj_repo.SpaceWorkObjectRepo
	workItemRepo   witem_repo.WorkItemRepo
	versionRepo    workVersion_repo.SpaceWorkVersionRepo
	memberRepo     member_repo.SpaceMemberRepo
	statusRepo     workItemStatus_repo.WorkItemStatusRepo // val -> item
	flowRepo       work_flow_repo.WorkFlowRepo
	searchRepo     search_repo.SearchRepo
}

func (ctx *SearchEsDataCtx) parseEsGroup() []*v1.GroupItem {

	//做分组
	var gs []*v1.GroupItem
	for _, gvs := range ctx.flattedGroups {

		var groupInfo = make([]*v1.GroupInfo, 0, len(gvs))
		var workItems []*v1.WorkItem

		for _, gv := range gvs {
			groupInfo = append(groupInfo, &v1.GroupInfo{
				FieldKey:    gv.Field,
				DisplayName: ctx.GetDisplayName(gv.Field, gv.Key),
				FieldId:     gv.Key,
			})

			//这是最后一层的分组, 添加工作项
			if len(gv.Hits) > 0 {
				workItems = ctx.buildChain(gv.Hits)
			}
		}

		gs = append(gs, &v1.GroupItem{
			GroupInfo: groupInfo,
			WorkItems: workItems,
		})
	}

	return gs
}

// 对分组排序
func (ctx *SearchEsDataCtx) sortGroup(groupOrder []Sort, groups []*v1.GroupItem) []*v1.GroupItem {
	var tuples []tuple.Pair[map[string]string, *v1.GroupItem]
	for _, group := range groups {
		fieldKV := stream.ToMap(group.GroupInfo, func(i int, t *v1.GroupInfo) (string, string) {
			return t.FieldKey, t.FieldId
		})
		tuples = append(tuples, tuple.Pair[map[string]string, *v1.GroupItem]{First: fieldKV, Second: group})
	}

	var cmpFuncList []func(a, b tuple.Pair[map[string]string, *v1.GroupItem]) int
	for _, by := range groupOrder {

		var rankingMap map[string]int64
		var f func(a, b tuple.Pair[map[string]string, *v1.GroupItem]) int

		switch by.Field {
		case "work_object_id":
			rankingMap = ctx.WorkObjectRankingMap()
		case "version_id":
			rankingMap = ctx.VersionRankingMap()
		case "space_id":
			rankingMap = ctx.SpaceRankingMap()
		case "priority":
			rankingMap = ctx.PriorityRankingMap()
		case "user_id", "directors":
			rankingMap = ctx.UserRankingMap()
		case "work_item_flow_id":
			rankingMap = ctx.FlowRankingMap()
		case "work_item_status_id":
			rankingMap = ctx.StatusRankingMap()
		default:
			continue
		}

		f = buildCmpFunc(by.Field, rankingMap)

		if by.Order == "DESC" {
			f = stream.ReverseCmpFunc(f)
		}

		cmpFuncList = append(cmpFuncList, f)
	}

	cmpFunc := stream.MergeCmpFunc(cmpFuncList...)
	slices.SortFunc(tuples, cmpFunc)

	return stream.Map(tuples, func(t tuple.Pair[map[string]string, *v1.GroupItem]) *v1.GroupItem {
		return t.Second
	})
}

func (ctx *SearchEsDataCtx) buildChain(list []*search_es2.HitItem) []*v1.WorkItem {
	// 记录下顺序值
	var rankMap = map[int64]int{}
	for i, v := range list {
		id := cast.ToInt64(v.Id)
		pid := int64(v.Fields.Float64s("pid")[0])

		rankMap[id] = i
		if _, ok := rankMap[pid]; !ok {
			rankMap[pid] = i
		}
	}

	// 分为两组 有父级 没有父级
	hasParent := stream.GroupBy(list, func(v *search_es2.HitItem) bool {
		pid := int64(v.Fields.Float64s("pid")[0])
		return pid != 0
	})

	var pMap = map[int64]*v1.WorkItem{}

	// 有父级
	for _, v := range hasParent[true] {
		id := cast.ToInt64(v.Id)
		pid := int64(v.Fields.Float64s("pid")[0])

		if pMap[pid] == nil {
			pMap[pid] = &v1.WorkItem{WorkItemId: pid}
		}
		pMap[pid].WorkItems = append(pMap[pid].WorkItems, &v1.WorkItem{WorkItemId: id})
	}

	// 没有父级
	for _, v := range hasParent[false] {
		id := cast.ToInt64(v.Id)

		//pMap[id] = ctx.GetWorkItemTree(id) // 填充子任务

		if pMap[id] == nil { //不填充子任务
			pMap[id] = &v1.WorkItem{WorkItemId: id}
		}
	}

	// 排序
	var sort func([]*v1.WorkItem)
	sort = func(items []*v1.WorkItem) {
		if len(items) == 0 {
			return
		}

		slices.SortFunc(items, func(a, b *v1.WorkItem) int {
			return cmp.Compare(rankMap[a.WorkItemId], rankMap[b.WorkItemId])
		})
		for _, v := range items {
			sort(v.WorkItems)
		}
	}

	values := maps.Values(pMap)
	sort(values)

	return values
}

func (ctx *SearchEsDataCtx) GetDisplayName(field, fieldValue string) string {

	switch field {
	case "priority":
		return consts.GetWorkItemPriorityName(fieldValue)
	case "work_object_id":
		object := ctx.getObjectMap()[cast.ToInt64(fieldValue)]
		if object == nil {
			return ""
		}
		return object.WorkObjectName
	case "user_id", "directors", "participators":
		object := ctx.getUserMap()[cast.ToInt64(fieldValue)]
		if object == nil {
			return ""
		}
		return object.UserNickname
	case "space_id":
		object := ctx.getSpaceMap()[cast.ToInt64(fieldValue)]
		if object == nil {
			return ""
		}
		return object.SpaceName
	case "work_item_flow_id":
		object := ctx.getFlowMap()[cast.ToInt64(fieldValue)]
		if object == nil {
			return ""
		}
		return object.Name
	case "work_item_status_id":
		object := ctx.getStatusMap()[cast.ToInt64(fieldValue)]
		if object == nil {
			return ""
		}
		return object.Name
	case "version_id":
		object := ctx.getVersionMap()[cast.ToInt64(fieldValue)]
		if object == nil {
			return ""
		}
		return object.VersionName
	}

	return fmt.Sprintf("%v %v", field, fieldValue)
}

//func (ctx *SearchEsDataCtx) GetWorkItemTree(workItemId int64) *v1.WorkItem {
//	if ctx.workItemChainMap != nil {
//		return ctx.workItemChainMap[workItemId]
//	}
//
//	var chainMap = map[int64]*v1.WorkItem{}
//	var hasChildIds []int64
//	for _, group := range ctx.flattedGroups {
//		if len(group) == 0 || len(group[len(group)-1].Hits) == 0 {
//			continue
//		}
//
//		for _, v := range group[len(group)-1].Hits {
//			id := cast.ToInt64(v.Id)
//			childNum := int64(v.Fields.Float64s("child_num")[0])
//
//			if childNum > 0 { //存在子任务
//				hasChildIds = append(hasChildIds, id)
//			} else { // 没有子任务
//				chainMap[id] = &v1.WorkItem{
//					WorkItemId: id,
//				}
//			}
//		}
//	}
//
//	// 获取子任务
//	values, _ := ctx.searchRepo.QueryWorkItemEsByPid(ctx.ctx, hasChildIds)
//
//	// 根据主任务id分组
//	subTaskMap := stream.GroupBy(values, func(v *search2.Model) int64 {
//		return v.Pid
//	})
//
//	// 对子任务进行排序
//	stream.MapValue(subTaskMap, func(v []*search2.Model) []*search2.Model {
//		slices.SortFunc(v, ctx.childTaskCmpFunc())
//		return v
//	})
//
//	for _, pid := range hasChildIds {
//		subTask := subTaskMap[pid]
//		var apiItems []*v1.WorkItem
//		for _, v := range subTask {
//			apiItems = append(apiItems, &v1.WorkItem{
//				WorkItemId: v.Id,
//			})
//		}
//		chainMap[pid] = &v1.WorkItem{
//			WorkItemId: pid,
//			WorkItems:  apiItems,
//		}
//	}
//
//	ctx.workItemChainMap = chainMap
//
//	return ctx.workItemChainMap[workItemId]
//}

func (ctx *SearchEsDataCtx) BuildSorter(sorts []Sort) []esV8.Sorter {
	var sorters []esV8.Sorter

	for _, v := range sorts {
		sorter := ctx.newEsSortByQueryFiled(search_es2.QueryField(v.Field), v.Order == "ASC")
		sorters = append(sorters, sorter)
	}

	sorters = stream.Filter(sorters, func(e esV8.Sorter) bool {
		return e != nil
	})

	return sorters
}

func (ctx *SearchEsDataCtx) newEsSortByQueryFiled(queryField search_es2.QueryField, ascending bool) esV8.Sorter {
	keyword := queryField.EsKeyword()
	if keyword == "" {
		keyword = queryField.String()
	}

	switch queryField {
	case search_es2.SpaceIdField: // 空间id
		return sortByScript(keyword, ctx.SpaceRankingMap(), ascending)
	case search_es2.PriorityField: // 优先级
		return sortByScript(keyword, ctx.PriorityRankingMap(), ascending)
	case search_es2.VersionIdField: // 版本
		return sortByScript(keyword, ctx.VersionRankingMap(), ascending)
	case search_es2.WorkObjectIdField: // 模块
		return sortByScript(keyword, ctx.WorkObjectRankingMap(), ascending)
	case search_es2.WorkItemFlowIdField: // 流程
		return sortByScript(keyword, ctx.FlowRankingMap(), ascending)
	case search_es2.UserIdField: // 创建人
		return sortByScript(keyword, ctx.UserRankingMap(), ascending)
	case search_es2.WorkItemStatusIdField: // 状态
		return sortByScript(keyword, ctx.StatusRankingMap(), ascending)
	default:
		return esV8.NewFieldSort(keyword).Order(ascending)
	}
}

// sortByScript 脚本排序
func sortByScript(queryField string, rankingMap map[string]int64, ascending bool) esV8.Sorter {
	if rankingMap == nil {
		return nil
	}

	scriptSource := `def v=doc[params.field].value.toString(); if(params.ranking.containsKey(v)) { return params.ranking[v] } return 0`
	script := esV8.NewScript(scriptSource).
		Param("ranking", rankingMap).
		Param("field", queryField)

	return esV8.NewScriptSort(script, "number").Order(ascending)
}

func buildCmpFunc(field string, rankingMap map[string]int64) func(a, b tuple.Pair[map[string]string, *v1.GroupItem]) int {
	return func(a, b tuple.Pair[map[string]string, *v1.GroupItem]) int {
		aId := a.First[field]
		bId := b.First[field]

		rankingA := rankingMap[aId]
		rankingB := rankingMap[bId]

		return cmp.Compare(rankingA, rankingB)
	}
}

func (ctx *SearchEsDataCtx) VersionRankingMap() map[string]int64 {
	versionMap := ctx.getVersionMap()

	m := stream.MapKV(versionMap, func(_ int64, v *workVersion.SpaceWorkVersion) (string, int64) {
		return cast.ToString(v.Id), v.Ranking
	})

	return m
}

func (ctx *SearchEsDataCtx) WorkObjectRankingMap() map[string]int64 {
	objectMap := ctx.getObjectMap()

	m := stream.MapKV(objectMap, func(_ int64, v *workObj.SpaceWorkObject) (string, int64) {
		return cast.ToString(v.Id), v.Ranking
	})

	return m
}

func (ctx *SearchEsDataCtx) SpaceRankingMap() map[string]int64 {
	userId := ctx.userId
	spaceIds := ctx.SpaceIds

	if len(spaceIds) <= 1 {
		return nil
	}

	memberMap, err := ctx.memberRepo.UserSpaceMemberMap(context.Background(), userId, spaceIds)
	if err != nil {
		return nil
	}

	m := stream.MapKV(memberMap, func(_ int64, v *space_member.SpaceMember) (string, int64) {
		return cast.ToString(v.SpaceId), v.Ranking
	})

	return m
}

func (ctx *SearchEsDataCtx) FlowRankingMap() map[string]int64 {
	flowMap := ctx.getFlowMap()

	m := stream.MapKV(flowMap, func(k int64, v *work_flow.WorkFlow) (string, int64) {
		return cast.ToString(v.Id), v.Ranking
	})

	return m
}

func (ctx *SearchEsDataCtx) PriorityRankingMap() map[string]int64 {
	return consts.GetPriorityRankingMap()
}

func sortStatus(list []*work_item_status.WorkItemStatusItem) {
	slices.SortFunc(list, func(i, j *work_item_status.WorkItemStatusItem) int {
		return i.Compare(j)
	})
}

func (ctx *SearchEsDataCtx) StatusRankingMap() map[string]int64 {
	srcMap := ctx.getStatusMap()

	list := stream.Values(srcMap)
	sortStatus(list)

	m := stream.ToMap(list, func(i int, t *work_item_status.WorkItemStatusItem) (string, int64) {
		return strconv.FormatInt(t.Id, 10), int64(-i)
	})

	return m
}

func (ctx *SearchEsDataCtx) UserRankingMap() map[string]int64 {
	userMap := ctx.getUserMap()
	entries := stream.ToEntries(userMap)
	slices.SortFunc(entries, func(a, b stream.Entry[int64, *user.User]) int {
		return cmp.Compare(b.Val.UserPinyin, a.Val.UserPinyin)
	})

	return stream.ToMap(entries, func(i int, t stream.Entry[int64, *user.User]) (string, int64) {
		return cast.ToString(t.Key), int64(i)
	})
}

func (ctx *SearchEsDataCtx) getRelativeUserIds() []int64 {
	var groupIds = make(map[string][]int64)
	for _, gvs := range ctx.flattedGroups {
		for _, gv := range gvs {
			groupIds[gv.Field] = append(groupIds[gv.Field], cast.ToInt64(gv.Key))
		}
	}

	var userIds []int64

	for k, ids := range groupIds {
		switch k {
		case "user_id", "directors":
			userIds = append(userIds, ids...)
		}
	}

	return stream.Unique(userIds)
}

func (ctx *SearchEsDataCtx) getVersionMap() map[int64]*workVersion.SpaceWorkVersion {
	if ctx.versionMap != nil {
		return ctx.versionMap
	}

	list, _ := ctx.versionRepo.GetSpaceWorkVersionBySpaceIds(ctx.ctx, ctx.SpaceIds)

	ctx.versionMap = stream.ToMap(list, func(_ int, v *workVersion.SpaceWorkVersion) (int64, *workVersion.SpaceWorkVersion) {
		return v.Id, v
	})

	return ctx.versionMap
}

func (ctx *SearchEsDataCtx) getFlowMap() map[int64]*work_flow.WorkFlow {
	if ctx.flowMap != nil {
		return ctx.flowMap
	}

	list, _ := ctx.flowRepo.GetWorkFlowBySpaceIds(ctx.ctx, ctx.SpaceIds)

	ctx.flowMap = stream.ToMap(list, func(_ int, v *work_flow.WorkFlow) (int64, *work_flow.WorkFlow) {
		return v.Id, v
	})

	return ctx.flowMap
}

func (ctx *SearchEsDataCtx) getUserMap() map[int64]*user.User {
	if ctx.userMap != nil {
		return ctx.userMap
	}

	userIds, _ := ctx.memberRepo.GetSpaceAllMemberIds(ctx.ctx, ctx.SpaceIds...)

	userIds = stream.Unique(userIds)
	ctx.userMap, _ = ctx.userRepo.UserMap(ctx.ctx, userIds)

	return ctx.userMap
}

func (ctx *SearchEsDataCtx) getStatusMap() map[int64]*work_item_status.WorkItemStatusItem {
	if ctx.statusMap != nil {
		return ctx.statusMap
	}

	ctx.statusMap, _ = ctx.statusRepo.StatusMapBySpaceIds(ctx.ctx, ctx.SpaceIds)

	return ctx.statusMap
}

func (ctx *SearchEsDataCtx) getObjectMap() map[int64]*workObj.SpaceWorkObject {
	if ctx.objectMap != nil {
		return ctx.objectMap
	}

	list, _ := ctx.workObjectRepo.GetSpaceWorkObjectBySpaceIds(ctx.ctx, ctx.SpaceIds)

	ctx.objectMap = stream.ToMap(list, func(_ int, v *workObj.SpaceWorkObject) (int64, *workObj.SpaceWorkObject) {
		return v.Id, v
	})

	return ctx.objectMap
}

func (ctx *SearchEsDataCtx) getSpaceMap() map[int64]*space.Space {
	if ctx.spaceMap != nil {
		return ctx.spaceMap
	}

	ctx.spaceMap, _ = ctx.spaceRepo.SpaceMap(ctx.ctx, ctx.SpaceIds)

	return ctx.spaceMap
}
