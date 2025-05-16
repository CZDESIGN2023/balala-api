package biz

import (
	"cmp"
	"context"
	v1 "go-cs/api/search/v1"
	"go-cs/internal/consts"
	search22 "go-cs/internal/domain/search/search2"
	member "go-cs/internal/domain/space_member"
	"go-cs/internal/utils/errs"
	group2 "go-cs/internal/utils/group"
	"go-cs/pkg/stream"
)

func (uc *SearchUsecase) SearchGroupInfoByDB(
	ctx context.Context,
	uid int64,
	req *v1.SearchSpaceWorkItemGroupInfoRequestV2,
	finalSorts []Sort,
	groupSorts []Sort,
) (*v1.SearchSpaceWorkItemGroupInfoReplyV2Data, error) {

	condition, err := search22.ConvertCondition(req.ConditionGroup)
	if err != nil {
		return nil, errs.Param(ctx, err.Error())
	}

	// 仅选择需要的字段
	selectFields := buildSelectFields(req)

	// 没传spaceId，则获取用户全部的spaceId
	var spaceIds []int64
	if req.SpaceId != 0 {
		spaceIds = append(spaceIds, req.SpaceId)
	} else {
		ids, err := uc.spaceRepo.GetUserSpaceIds(ctx, uid)
		if err != nil {
			return nil, err
		}
		spaceIds = ids
	}

	// 获取数据
	list, err := uc.repo.QueryWorkItem(ctx, spaceIds, condition, selectFields)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	// 提取模块id
	workObjectIds := stream.Map(list, func(item *search22.Model) int64 {
		return item.WorkObjectId
	})

	// 提取版本id
	versionIds := stream.Map(list, func(item *search22.Model) int64 {
		return item.VersionId
	})

	// 提取流程id
	flowIds := stream.Map(list, func(item *search22.Model) int64 {
		return item.FlowId
	})

	getWorkItemSortFunc := func(field string) func(a, b *search22.Model) int {
		switch field {
		case "space_id":
			return uc.workSpaceSortCmpFunc(ctx, uid, spaceIds)
		case "work_object_id":
			return uc.workObjectSortCmpFunc(ctx, workObjectIds)
		case "version_id":
			return uc.versionSortCmpFunc(ctx, versionIds)
		case "priority":
			return func(a, b *search22.Model) int {
				ranking := consts.GetPriorityRankingMap()
				return cmp.Compare(ranking[a.Priority], ranking[b.Priority])
			}
		case "work_item_id":
			return func(a, b *search22.Model) int {
				return cmp.Compare(a.Id, b.Id)
			}
		case "work_item_flow_id":
			return uc.flowSortCmpFunc(ctx, flowIds)
		default:
			return nil
		}
	}

	var cmpFunctions []func(a *search22.Model, b *search22.Model) int
	for _, s := range finalSorts {
		fn := getWorkItemSortFunc(s.Field)
		if fn == nil {
			continue
		}

		if s.Order == "DESC" {
			fn = search22.ReverseCmpFunc(fn)
		}

		cmpFunctions = append(cmpFunctions, fn)
	}

	search22.Sort(list, cmpFunctions...)

	// 分组
	g := group2.New(list)
	groups := req.Groups
	for _, group := range groups {
		g.GroupByToMulti(group.Field, func(item *search22.Model) []string {
			return search22.GroupValue(group.Field, item)
		})
	}

	sortOrderMap := stream.ToMap(finalSorts, func(i int, t Sort) (string, string) {
		return t.Field, t.Order
	})

	searchCtx := &DataCtx{
		ctx:            ctx,
		Data:           list,
		userId:         uid,
		SpaceIds:       spaceIds,
		childOrderMap:  sortOrderMap,
		spaceRepo:      uc.spaceRepo,
		userRepo:       uc.userRepo,
		workObjectRepo: uc.spaceWorkObjectRepo,
		workItemRepo:   uc.spaceWorkItemRepo,
		memberRepo:     uc.spaceMemberRepo,
		versionRepo:    uc.workVersionRepo,
		flowRepo:       uc.workFlowRepo,
		statusRepo:     uc.workItemStatusRepo,
	}

	//将分组转换为响应结构
	ret := searchCtx.parseGroup(g)

	if len(groupSorts) > 0 {
		ret = searchCtx.sortGroup(groupSorts, ret)
	}

	return &v1.SearchSpaceWorkItemGroupInfoReplyV2Data{
		List:     ret,
		TotalNum: int32(len(list)),
	}, nil
}

func buildSelectFields(req *v1.SearchSpaceWorkItemGroupInfoRequestV2) string {
	var fields = []string{"space_id", "pid", "work_item_id", "work_object_id", "version_id", "priority"} //必须的字段
	for _, v := range req.Sorts {
		if v.Field != "" {
			fields = append(fields, v.Field)
		}
	}
	for _, v := range req.Groups {
		fields = append(fields, v.Field)
	}

	fields = stream.Unique(fields) //去重

	return search22.SelectByQuery(fields...)
}

func (uc *SearchUsecase) workSpaceSortCmpFunc(ctx context.Context, userId int64, spaceIds []int64) func(a, b *search22.Model) int {
	var m map[int64]*member.SpaceMember
	if len(spaceIds) > 1 {
		m, _ = uc.spaceMemberRepo.GetUserSpaceMemberBySpaceId(ctx, userId, spaceIds)
	}

	return func(a, b *search22.Model) int {
		ao := m[a.SpaceId]
		bo := m[b.SpaceId]

		if ao == nil || bo == nil {
			return cmp.Compare(a.SpaceId, b.SpaceId)
		}

		return cmp.Compare(ao.Ranking, bo.Ranking)
	}
}

func (uc *SearchUsecase) workObjectSortCmpFunc(ctx context.Context, objectIds []int64) func(a, b *search22.Model) int {
	objectIds = stream.Unique(objectIds)

	m, _ := uc.spaceWorkObjectRepo.SpaceWorkObjectMapByObjectIds(ctx, objectIds)

	return func(a, b *search22.Model) int {
		ao := m[a.WorkObjectId]
		bo := m[b.WorkObjectId]

		if ao == nil || bo == nil {
			return cmp.Compare(a.WorkObjectId, b.WorkObjectId)
		}

		return cmp.Compare(ao.Ranking, bo.Ranking)
	}
}

func (uc *SearchUsecase) versionSortCmpFunc(ctx context.Context, versionIds []int64) func(a, b *search22.Model) int {
	versionIds = stream.Unique(versionIds)

	m, _ := uc.workVersionRepo.SpaceWorkVersionMapByVersionIds(ctx, versionIds)

	return func(a, b *search22.Model) int {
		ao := m[a.VersionId]
		bo := m[b.VersionId]

		if ao == nil || bo == nil {
			return cmp.Compare(a.VersionId, b.VersionId)
		}

		return cmp.Compare(ao.Ranking, bo.Ranking)
	}
}

func (uc *SearchUsecase) flowSortCmpFunc(ctx context.Context, flowIds []int64) func(a, b *search22.Model) int {
	flowIds = stream.Unique(flowIds)

	m, _ := uc.workFlowRepo.WorkFlowMap(ctx, flowIds)

	return func(a, b *search22.Model) int {
		ao := m[a.FlowId]
		bo := m[b.FlowId]

		if ao == nil || bo == nil {
			return cmp.Compare(a.FlowId, b.FlowId)
		}

		return cmp.Compare(ao.Ranking, bo.Ranking)
	}
}
