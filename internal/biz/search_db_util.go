package biz

import (
	"cmp"
	"context"
	"fmt"
	v1 "go-cs/api/search/v1"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/consts"
	"go-cs/internal/domain/search/search2"
	"go-cs/internal/domain/work_flow"
	work_flow_repo "go-cs/internal/domain/work_flow/repo"
	"go-cs/internal/domain/work_item_status"
	"go-cs/internal/utils"
	group2 "go-cs/internal/utils/group"
	"go-cs/pkg/stream"
	"go-cs/pkg/stream/tuple"
	"golang.org/x/exp/maps"
	"slices"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"

	"go-cs/internal/domain/space"
	space_repo "go-cs/internal/domain/space/repo"
	member "go-cs/internal/domain/space_member"
	member_repo "go-cs/internal/domain/space_member/repo"
	workObj "go-cs/internal/domain/space_work_object"
	workObj_repo "go-cs/internal/domain/space_work_object/repo"
	"go-cs/internal/domain/user"
	user_repo "go-cs/internal/domain/user/repo"
	witem_repo "go-cs/internal/domain/work_item/repo"

	witem "go-cs/internal/domain/work_item"

	workVersion "go-cs/internal/domain/space_work_version"
	workVersion_repo "go-cs/internal/domain/space_work_version/repo"
	workItemStatus_repo "go-cs/internal/domain/work_item_status/repo"
)

type Sort struct {
	Field string
	Order string
}

type DataCtx struct {
	ctx context.Context

	userId        int64
	SpaceIds      []int64
	WorkObjectIds []int64

	curUserSpaceRankingMap map[string]int64

	childOrderMap map[string]string

	Data        []*search2.Model
	ObjectMap   map[int64]*workObj.SpaceWorkObject
	UserMap     map[int64]*user.User
	SpaceMap    map[int64]*space.Space
	workItemMap map[int64]*search2.Model
	versionMap  map[int64]*workVersion.SpaceWorkVersion
	statusMap   map[int64]*work_item_status.WorkItemStatusItem // val -> item
	flowMap     map[int64]*work_flow.WorkFlow

	workItemChainMap map[int64]*v1.WorkItem

	spaceRepo      space_repo.SpaceRepo
	userRepo       user_repo.UserRepo
	workObjectRepo workObj_repo.SpaceWorkObjectRepo
	workItemRepo   witem_repo.WorkItemRepo
	versionRepo    workVersion_repo.SpaceWorkVersionRepo
	memberRepo     member_repo.SpaceMemberRepo
	statusRepo     workItemStatus_repo.WorkItemStatusRepo // val -> item
	flowRepo       work_flow_repo.WorkFlowRepo
}

func (d *DataCtx) parseGroup(g *group2.Group[*search2.Model]) []*v1.GroupItem {
	if g == nil || len(g.Values) == 0 {
		return nil
	}

	groupInfo := &v1.GroupInfo{
		FieldKey:    g.FieldName,
		DisplayName: d.GetDisplayName(g.FieldName, g.Key),
		FieldId:     g.Key,
	}

	if len(g.NextNodes) == 0 {
		return []*v1.GroupItem{{
			GroupInfo: []*v1.GroupInfo{groupInfo},
			WorkItems: d.buildChain(g.Values),
		}}
	}

	var gs []*v1.GroupItem
	for _, v := range g.NextNodes {
		items := d.parseGroup(v)
		if g.Key != "" {
			for _, item := range items {
				item.GroupInfo = append([]*v1.GroupInfo{groupInfo}, item.GroupInfo...)
			}
		}
		gs = append(gs, items...)
	}
	return gs
}

// 对分组排序
func (d *DataCtx) sortGroup(groupOrder []Sort, groups []*v1.GroupItem) []*v1.GroupItem {
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
			rankingMap = d.WorkObjectRankingMap()
		case "version_id":
			rankingMap = d.VersionRankingMap()
		case "space_id":
			rankingMap = d.SpaceRankingMap()
		case "priority":
			rankingMap = d.PriorityRankingMap()
		case "user_id":
			rankingMap = d.MemberRankingMap()
		case "directors":
			rankingMap = d.MemberRankingMap()
		case "work_item_flow_id":
			rankingMap = d.WorkItemFlowRankingMap()
		case "work_item_status_id":
			rankingMap = d.WorkItemStatusRankingMap()
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

func (d *DataCtx) buildChain(list []*search2.Model) []*v1.WorkItem {
	// 记录下顺序值
	var rankMap = map[int64]int{}
	for i, v := range list {
		rankMap[v.Id] = i
		if _, ok := rankMap[v.Pid]; !ok {
			rankMap[v.Pid] = i
		}
	}

	// 分为两组 有父级 没有父级
	hasParent := stream.GroupBy(list, func(v *search2.Model) bool {
		return v.Pid != 0
	})

	var pMap = map[int64]*v1.WorkItem{}
	for _, v := range hasParent[true] {
		if pMap[v.Pid] == nil {
			pMap[v.Pid] = &v1.WorkItem{WorkItemId: v.Pid}
		}
		pMap[v.Pid].WorkItems = append(pMap[v.Pid].WorkItems, &v1.WorkItem{WorkItemId: v.Id})
	}
	for _, v := range hasParent[false] {
		//pMap[v.Id] = ctx.GetWorkItemTree(v.Id)
		if pMap[v.Id] == nil {
			pMap[v.Id] = &v1.WorkItem{WorkItemId: v.Id}
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

func (d *DataCtx) getWorkObjectIds() []int64 {
	if d.WorkObjectIds != nil {
		return d.WorkObjectIds
	}

	WorkObjectIds := stream.Map(d.Data, func(v *search2.Model) int64 {
		return v.WorkObjectId
	})

	d.WorkObjectIds = stream.Unique(WorkObjectIds)

	return d.WorkObjectIds
}

func (d *DataCtx) GetWorkItemTree(workItemId int64) *v1.WorkItem {
	if d.workItemChainMap != nil {
		return d.workItemChainMap[workItemId]
	}

	var chainMap = map[int64]*v1.WorkItem{}
	var hasChild []int64
	for _, v := range d.Data {
		if v.ChildNum != 0 {
			hasChild = append(hasChild, v.Id)
		}
	}

	values, _ := d.workItemRepo.GetWorkItemByPids(d.ctx, hasChild, nil, nil)

	// 根据主任务id分组
	subTaskMap := stream.GroupBy(values, func(v *witem.WorkItem) int64 {
		return v.Pid
	})

	// 对子任务进行排序
	stream.MapValue(subTaskMap, func(v []*witem.WorkItem) []*witem.WorkItem {
		if d.childOrderMap["work_item_id"] == "DESC" {
			slices.SortFunc(v, func(a, b *witem.WorkItem) int {
				return cmp.Compare(b.Id, a.Id)
			})
		} else {
			slices.SortFunc(v, func(a, b *witem.WorkItem) int {
				return cmp.Compare(a.Id, b.Id)
			})
		}

		return v
	})

	for _, v := range d.Data {
		if v.Pid == 0 {
			subTask := subTaskMap[v.Id]
			var apiItems []*v1.WorkItem
			for _, v := range subTask {
				apiItems = append(apiItems, &v1.WorkItem{
					WorkItemId: v.Id,
				})
			}
			chainMap[v.Id] = &v1.WorkItem{
				WorkItemId: v.Id,
				WorkItems:  apiItems,
			}
		} else {
			chainMap[v.Id] = &v1.WorkItem{
				WorkItemId: v.Pid,
				WorkItems: []*v1.WorkItem{
					{WorkItemId: v.Id},
				},
			}
		}
	}

	d.workItemChainMap = chainMap

	return d.workItemChainMap[workItemId]
}

func (d *DataCtx) getObjectMap() map[int64]*workObj.SpaceWorkObject {
	if d.ObjectMap != nil {
		return d.ObjectMap
	}

	var ids []int64
	for _, v := range d.Data {
		ids = append(ids, v.WorkObjectId)
	}

	ids = stream.Unique(ids) //去重

	d.ObjectMap, _ = d.workObjectRepo.SpaceWorkObjectMapByObjectIds(d.ctx, ids)

	return d.ObjectMap
}

func (d *DataCtx) getFlowMap() map[int64]*work_flow.WorkFlow {
	if d.flowMap != nil {
		return d.flowMap
	}

	var ids []int64
	for _, v := range d.Data {
		ids = append(ids, v.WorkItemFlowId)
	}

	ids = stream.Unique(ids) //去重

	d.flowMap, _ = d.flowRepo.GetWorkFlowMap(d.ctx, ids)

	return d.flowMap
}

func (d *DataCtx) getStatusMap() map[int64]*work_item_status.WorkItemStatusItem {
	if d.statusMap != nil {
		return d.statusMap
	}

	var ids []int64
	for _, v := range d.Data {
		ids = append(ids, v.Id)
	}

	ids = stream.Unique(ids) //去重

	if len(d.SpaceIds) == 0 {
		return nil
	}

	statusList, _ := d.statusRepo.GetWorkItemStatusItemsBySpace(d.ctx, d.SpaceIds[0])

	d.statusMap = stream.ToMap(statusList, func(_ int, v *work_item_status.WorkItemStatusItem) (int64, *work_item_status.WorkItemStatusItem) {
		return v.Id, v
	})

	return d.statusMap
}

func (d *DataCtx) getUserMap() map[int64]*user.User {
	if d.UserMap != nil {
		return d.UserMap
	}

	var ids []int64
	for _, v := range d.Data {
		switch {
		case len(v.Participators) != 0:
			ids = append(ids, utils.ToInt64Array(v.Participators)...)
		case len(v.Directors) != 0:
			ids = append(ids, utils.ToInt64Array(v.Directors)...)
		default:
			ids = append(ids, v.UserId)
		}
	}
	slices.Sort(ids)
	ids = slices.Compact(ids) //去重

	d.UserMap, _ = d.userRepo.UserMap(d.ctx, ids)

	return d.UserMap
}

func (d *DataCtx) getSpaceMap() map[int64]*space.Space {
	if d.SpaceMap != nil {
		return d.SpaceMap
	}

	d.SpaceMap, _ = d.spaceRepo.SpaceMap(d.ctx, d.SpaceIds)

	return d.SpaceMap
}

func (d *DataCtx) getVersionMap() map[int64]*workVersion.SpaceWorkVersion {
	if d.versionMap != nil {
		return d.versionMap
	}

	var ids []int64
	for _, v := range d.Data {
		ids = append(ids, v.VersionId)
	}
	slices.Sort(ids)
	ids = slices.Compact(ids) //去重

	d.versionMap, _ = d.versionRepo.SpaceWorkVersionMapByVersionIds(d.ctx, ids)

	return d.versionMap
}

func (d *DataCtx) GetDisplayName(field, fieldValue string) string {
	switch field {
	case "work_item_flow_id":
		typ := d.getFlowMap()[cast.ToInt64(fieldValue)]
		if typ == nil {
			return ""
		}
		return typ.Name
	case "work_item_status_id":
		return d.getStatusMap()[cast.ToInt64(fieldValue)].Name
	case "work_object_id":
		object := d.getObjectMap()[cast.ToInt64(fieldValue)]
		if object == nil {
			return ""
		}
		return object.WorkObjectName
	case "user_id", "directors":
		object := d.getUserMap()[cast.ToInt64(fieldValue)]
		if object == nil {
			return ""
		}
		return object.UserNickname
	case "space_id":
		object := d.getSpaceMap()[cast.ToInt64(fieldValue)]
		if object == nil {
			return ""
		}
		return object.SpaceName
	case "priority":
		return fieldValue
	case "version_id":
		object := d.getVersionMap()[cast.ToInt64(fieldValue)]
		if object == nil {
			return ""
		}
		return object.VersionName
	}

	return fmt.Sprintf("%v %v", field, fieldValue)
}

func (d *DataCtx) VersionRankingMap() map[string]int64 {
	spaceIds := d.SpaceIds

	if len(spaceIds) != 1 {
		return nil
	}

	versionMap, err := d.versionRepo.SpaceWorkVersionMap(context.Background(), spaceIds[0])
	if err != nil {
		return nil
	}

	m := stream.MapKV(versionMap, func(_ int64, v *workVersion.SpaceWorkVersion) (string, int64) {
		return cast.ToString(v.Id), v.Ranking
	})

	return m
}

func (d *DataCtx) WorkObjectRankingMap() map[string]int64 {
	spaceIds := d.SpaceIds

	if len(spaceIds) != 1 {
		return nil
	}

	m := stream.MapKV(d.getObjectMap(), func(_ int64, v *workObj.SpaceWorkObject) (string, int64) {
		return cast.ToString(v.Id), v.Ranking
	})

	return m
}

func (d *DataCtx) SpaceRankingMap() map[string]int64 {
	userId := d.userId
	spaceIds := d.SpaceIds

	if len(spaceIds) <= 1 {
		return nil
	}

	objectMap, err := d.memberRepo.UserSpaceMemberMap(context.Background(), userId, spaceIds)
	if err != nil {
		return nil
	}

	m := stream.MapKV(objectMap, func(_ int64, v *member.SpaceMember) (string, int64) {
		return cast.ToString(v.Id), v.Ranking
	})

	return m
}

func (d *DataCtx) PriorityRankingMap() map[string]int64 {
	return consts.GetPriorityRankingMap()
}

func (d *DataCtx) WorkItemStatusRankingMap() map[string]int64 {
	if len(d.SpaceIds) != 1 {
		return nil
	}

	return stream.MapKV(d.getStatusMap(), func(k int64, v *work_item_status.WorkItemStatusItem) (string, int64) {
		return cast.ToString(k), v.Ranking
	})
}

func (d *DataCtx) MemberRankingMap() map[string]int64 {
	spaceIds := d.SpaceIds

	if len(spaceIds) != 1 {
		return nil
	}

	spaceId := spaceIds[0]

	list, err := d.memberRepo.QSpaceMemberList(d.ctx, spaceId, "")
	if err != nil {
		log.Error(err)
		return nil
	}

	slices.SortFunc(list, func(a, b *rsp.SpaceMemberInfo) int {
		if a.RoleId != b.RoleId {
			return cmp.Compare(consts.GetMemberRoleRank(a.RoleId), consts.GetMemberRoleRank(b.RoleId))
		}

		return cmp.Compare(strings.ToLower(a.UserPinyin), strings.ToLower(b.UserPinyin))
	})

	m := stream.ToMap(list, func(i int, v *rsp.SpaceMemberInfo) (string, int64) {
		return cast.ToString(v.UserId), int64(len(list) - i)
	})

	return m
}

func (d *DataCtx) WorkItemFlowRankingMap() map[string]int64 {
	spaceIds := d.SpaceIds

	if len(spaceIds) != 1 {
		return nil
	}

	m := stream.MapKV(d.getFlowMap(), func(_ int64, v *work_flow.WorkFlow) (string, int64) {
		return cast.ToString(v.Id), v.Ranking
	})

	return m
}

func combineToSortableInt64(a, b int64) int64 {
	// 将a的高位移动到int64的高32位
	// 将b的低位保持在int64的低32位
	// 由于a和b已经有符号，这里不需要额外的处理
	return a<<32 | b
}
