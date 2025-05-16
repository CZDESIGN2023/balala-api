package biz

import (
	"context"
	v1 "go-cs/api/search/v1"
	search_es2 "go-cs/internal/domain/search/search_es"
	"go-cs/internal/utils/errs"
	esV8 "go-cs/internal/utils/es/v8"
	"go-cs/pkg/stream"
)

func (uc *SearchUsecase) SearchGroupInfoByEs(
	ctx context.Context,
	uid int64,
	spaceIds []int64,
	req *v1.SearchSpaceWorkItemGroupInfoRequestV2,
	finalSorts []Sort,
	groupSorts []Sort,
) (*v1.SearchSpaceWorkItemGroupInfoReplyV2Data, error) {

	groups := req.Groups

	condition, err := search_es2.ConvertCondition(req.ConditionGroup)
	if err != nil {
		return nil, errs.Param(ctx, err.Error())
	}

	spaceIdsAny := stream.Map(spaceIds, func(id int64) any {
		return id
	})

	sortOrderMap := stream.ToMap(finalSorts, func(i int, t Sort) (string, string) {
		return t.Field, t.Order
	})

	searchCtx := &SearchEsDataCtx{
		ctx:            ctx,
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

	// 构建filter条件
	sourceQuery := search_es2.BuildCondition(condition)
	if sourceQuery == nil {
		sourceQuery = esV8.NewBoolQuery()
	}
	// 添加spaceId过滤条件
	sourceQuery.Must(esV8.NewTermsQuery(search_es2.SpaceIdField.EsKeyword(), spaceIdsAny...))

	var source *esV8.SearchSource
	// 分组 按层级来 扩展的话需要实现每个对应查询字段的aggregation
	if len(groups) > 0 {
		var firstAgg *esV8.TermsAggregation
		var firstField string
		var curAgg *esV8.TermsAggregation
		for _, g := range groups {
			agg := esV8.
				NewTermsAggregation().
				Field(search_es2.GetFieldEsKeyword(g.Field)).
				Size(10000)
			if curAgg != nil {
				curAgg.SubAggregation(g.Field, agg)
			} else {
				firstField = g.Field
				firstAgg = agg
			}
			curAgg = agg
		}

		//最后一层，加入hit
		sorters := searchCtx.BuildSorter(finalSorts)
		hitsSource := esV8.NewSearchSource().
			FetchSource(false).
			Fields(
				search_es2.PidField.Es(),
				search_es2.ChildNumField.Es(),
			).
			Size(10000).
			SortBy(sorters...)

		hitAgg := esV8.NewTopHitsAggregation().SearchSource(hitsSource)
		curAgg.SubAggregation("top_hits", hitAgg)

		source = esV8.NewSearchSource().
			Query(sourceQuery).
			Aggregation(firstField, firstAgg).
			Size(0) //这里不需要hits信息。 分组top_hits中有

	} else {
		sorters := searchCtx.BuildSorter(finalSorts)
		source = esV8.NewSearchSource().
			Query(sourceQuery).
			FetchSource(false).
			Fields(
				search_es2.PidField.Es(),
				search_es2.ChildNumField.Es(),
			).
			Size(10000).
			SortBy(sorters...)
	}

	// 获取数据
	result, err := uc.repo.QueryWorkItemEs(ctx, source)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}
	if result == nil {
		return nil, nil
	}

	var esGroup [][]*search_es2.Group
	if len(groups) > 0 {
		// 获取分组信息, 任务id已经在top_hits中存在了
		esGroup = search_es2.ConvertAgg(result.Aggregations)
	} else {
		// 没有分组，直接获取hits
		esGroup = search_es2.ConvertHits(result.Hits.Hits)
	}

	// 扁平化
	flattenEsGroup := search_es2.FlattenGroup(esGroup)

	//分组的名称需要补全
	searchCtx.flattedGroups = flattenEsGroup

	// 将分组转换为响应结构
	ret := searchCtx.parseEsGroup()

	// 对分组进行排序
	ret = searchCtx.sortGroup(groupSorts, ret)

	return &v1.SearchSpaceWorkItemGroupInfoReplyV2Data{
		List:     ret,
		TotalNum: int32(result.TotalHits()),
	}, nil
}
