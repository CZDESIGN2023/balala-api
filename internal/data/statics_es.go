package data

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo"
	"go-cs/internal/consts"
	search_es2 "go-cs/internal/domain/search/search_es"
	"go-cs/internal/domain/work_item_status"
	"go-cs/internal/utils"
	"go-cs/internal/utils/date"
	esV8 "go-cs/internal/utils/es/v8"
	"go-cs/pkg/stream"
	"net/http"
	"slices"
	"strings"
	"time"

	"go-cs/internal/domain/statics/repo"
	status_repo "go-cs/internal/domain/work_item_status/repo"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type staticsEsRepo struct {
	baseRepo
	statusRepo status_repo.WorkItemStatusRepo
}

func NewStaticsEsRepo(data *Data, statusRepo status_repo.WorkItemStatusRepo, logger log.Logger) repo.StaticsRepo {
	moduleName := "staticsEsRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	repo := &staticsEsRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
		statusRepo: statusRepo,
	}
	return repo
}

func (c *staticsEsRepo) esDoSearch(source *esV8.SearchSource) (*esV8.SearchResult, error) {
	marshal, _ := source.MarshalJSON()

	//pprint.Println(source)

	es := c.data.es
	indexName := c.data.conf.Es.Index

	esRes, err := es.Search(
		es.Search.WithIndex(indexName),
		es.Search.WithBody(bytes.NewReader(marshal)),
	)

	if err != nil {
		return nil, err
	}
	defer esRes.Body.Close()

	r := &esV8.SearchResult{}
	// 解析查询结果
	if err := json.NewDecoder(esRes.Body).Decode(&r); err != nil {
		c.log.Errorf("Error parsing the response body: %s", err)
		return nil, err
	}

	if esRes.StatusCode == http.StatusNotFound {
		return &esV8.SearchResult{Hits: &esV8.SearchHits{}}, nil
	}

	if esRes.IsError() {
		c.log.Errorf("Search request returned an error: %s", esRes.String())
		return nil, errors.New(esRes.String())
	}

	if esRes.StatusCode != http.StatusOK {
		c.log.Debug("Search request failed: %s", esRes.String())
		return nil, errors.New(esRes.String())
	}

	return r, nil
}

func (c *staticsEsRepo) GetWorkbenchCount(ctx context.Context, userId int64, spaceIds []int64) (*vo.UserWorkbenchCountInfo, error) {
	if len(spaceIds) == 0 {
		return nil, nil
	}

	//获取对应的任务状态
	status := c.spaceStatus(ctx, spaceIds...)

	completeStatusList := status.CompleteStatusList
	endStatusList := status.EndStatusList

	source := esV8.NewSearchSource().
		FetchSourceIncludeExclude([]string{
			"id",
			"user_id",
			"work_item_status_id",
			"last_status_at",
			"doc.plan_start_at",
			"doc.plan_complete_at",
			"doc.directors",
			"doc.participators",
		}, nil).
		Query(esV8.NewBoolQuery().Must(
			esV8.NewTermsQuery("space_id", utils.ToAnyArray(spaceIds)...),
			esV8.NewTermQuery("doc.participators", userId),
		)).Size(10000)

	res, err := c.esDoSearch(source)
	if err != nil {
		return nil, err
	}

	type entity struct {
		Id               int64    `json:"id"`
		UserId           int64    `json:"user_id"`
		WorkItemStatusId int64    `json:"work_item_status_id"`
		LastStatusAt     int64    `json:"last_status_at"`
		PlanStartAt      int64    `json:"doc.plan_start_at"`
		PlanCompleteAt   int64    `json:"doc.plan_complete_at"`
		Directors        []string `json:"doc.directors"`
		Participators    []string `json:"doc.participators"`

		Doc struct {
			PlanStartAt    int64    `json:"plan_start_at"`
			PlanCompleteAt int64    `json:"plan_complete_at"`
			Directors      []string `json:"directors"`
			Participators  []string `json:"participators"`
		} `json:"doc"`
	}

	userIdStr := cast.ToString(userId)

	var list []entity
	list = stream.Map(res.Hits.Hits, func(v *esV8.SearchHit) entity {
		var e entity
		json.Unmarshal(v.Source, &e)

		e.PlanStartAt = e.Doc.PlanStartAt
		e.PlanCompleteAt = e.Doc.PlanCompleteAt
		e.Directors = e.Doc.Directors
		e.Participators = e.Doc.Participators
		return e
	})

	var progressingCount int64
	var weekProcessingCount int64
	var participatedProcessingCount int64

	var completeCount int64
	var weekCompleteCount int64

	var createdCount int64

	weekBegin, weekEnd := date.ThisWeekBeginEnd()
	weekBeginUnix, weekEndUnix := weekBegin.Unix(), weekEnd.Unix()

	for _, v := range list {
		// 是当前负责人
		isDirector := slices.Contains(v.Directors, userIdStr)
		// 是创建人
		isCreator := v.UserId == userId
		// 任务是待办
		isProcessing := !slices.Contains(endStatusList, v.WorkItemStatusId)
		// 任务是完成
		isComplete := slices.Contains(completeStatusList, v.WorkItemStatusId)

		// 统计我创建的
		if isCreator {
			createdCount++
		}

		// 统计待办相关
		if isProcessing {
			if isDirector { // 我的当前负责的待办
				progressingCount++

				// 我的本周待办
				if date.HasInter(v.PlanStartAt, v.PlanCompleteAt, weekBeginUnix, weekEndUnix) { //是否存在交集
					weekProcessingCount++
				}
			}

			// 我参与的待办
			if slices.Contains(v.Participators, userIdStr) {
				participatedProcessingCount++
			}

			continue
		}

		// 统计已完成相关
		if isComplete { //已完成
			completeCount++
			if weekBeginUnix <= v.LastStatusAt && v.LastStatusAt <= weekEndUnix { //是否存在交集
				weekCompleteCount++
			}
		}
	}

	return &vo.UserWorkbenchCountInfo{
		ProcessingCount:             progressingCount,
		WeekProcessingCount:         weekProcessingCount,
		ParticipatedProcessingCount: participatedProcessingCount,
		CompleteCount:               completeCount,
		WeekCompleteCount:           weekCompleteCount,
		CreatedCount:                createdCount,
	}, nil
}

func (c *staticsEsRepo) GetSpaceWorkObjectCountByIds(ctx context.Context, spaceId int64, workObjectIds []int64, startTime, endTime int64) (map[int64]*vo.SpaceWorkObjectCountInfo, error) {

	//获取对应的任务状态
	status := c.spaceStatus(ctx, spaceId)

	completeStatusList := status.CompleteStatusList
	endStatusList := status.EndStatusList
	endAndNotCompleteStatusList := status.EndAndNotCompleteStatusList

	queryList := []esV8.Query{
		esV8.NewTermQuery("space_id", spaceId),
	}

	if len(workObjectIds) != 0 {
		queryList = append(queryList, esV8.NewTermsQuery("work_object_id", utils.ToAnyArray(workObjectIds)...))
	}

	if startTime != 0 && endTime != 0 {
		queryList = append(queryList, esV8.NewRangeQuery("created_at").Gte(startTime).Lte(endTime))
	}

	query := esV8.NewBoolQuery().Must(queryList...)
	source := esV8.NewSearchSource().
		FetchSourceIncludeExclude([]string{
			"work_object_id",
			"work_item_status_id",
			"last_status_at",
			"doc.plan_start_at",
			"doc.plan_complete_at",
			"doc.priority",
		}, nil).
		Query(query).
		Size(10000)

	res, err := c.esDoSearch(source)
	if err != nil {
		return nil, err
	}

	type item struct {
		WorkObjectId     int64  `json:"work_object_id"`
		WorkItemStatusId int64  `json:"work_item_status_id"`
		Priority         string `json:"priority"`
		PlanStartAt      int64  `json:"plan_start_at"`
		PlanCompleteAt   int64  `json:"plan_complete_at"`

		Doc struct {
			Priority       string `json:"priority"`
			PlanStartAt    int64  `json:"plan_start_at"`
			PlanCompleteAt int64  `json:"plan_complete_at"`
		} `json:"doc"`
	}

	var list []item
	list = stream.Map(res.Hits.Hits, func(v *esV8.SearchHit) item {
		var e item
		json.Unmarshal(v.Source, &e)
		e.PlanStartAt = e.Doc.PlanStartAt
		e.PlanCompleteAt = e.Doc.PlanCompleteAt
		e.Priority = e.Doc.Priority
		return e
	})

	weekBegin, weekEnd := date.ThisWeekBeginEnd()
	weekBeginUnix, weekEndUnix := weekBegin.Unix(), weekEnd.Unix()

	m := map[int64]*vo.SpaceWorkObjectCountInfo{}
	now := time.Now().Unix()
	for _, v := range list {
		info := m[v.WorkObjectId]
		if info == nil {
			info = &vo.SpaceWorkObjectCountInfo{
				PriorityInfo: map[string]int64{},
			}
			m[v.WorkObjectId] = info
		}

		info.TotalCount += 1

		// 进行中
		if !slices.Contains(endStatusList, v.WorkItemStatusId) {
			// 待办计数
			info.ProcessingCount += 1

			// 优先级计数
			info.PriorityInfo[v.Priority] += 1

			// 逾期计数
			if v.PlanCompleteAt < now {
				info.ExpiredCount++
			}

			// 本周待办计数
			if date.HasInter(v.PlanStartAt, v.PlanCompleteAt, weekBeginUnix, weekEndUnix) { //是否存在交集
				info.WeekProcessingCount++
			}
		} else if slices.Contains(completeStatusList, v.WorkItemStatusId) { //完成
			info.CompletedCount++
		} else if slices.Contains(endAndNotCompleteStatusList, v.WorkItemStatusId) { //已关闭或已终止
			info.ClosedOrTerminatedCount++
		}
	}

	// 计算完成度
	for _, v := range m {
		if v.TotalCount == 0 {
			continue
		}
		v.CompleteRate = int64(100 - (float64(v.ProcessingCount) * 100 / float64(v.TotalCount)))
	}

	return m, nil
}

func (c *staticsEsRepo) GetSpaceWorkVersionCountByIds(ctx context.Context, spaceId int64, versionIds []int64, startTime, endTime int64) (map[int64]*vo.SpaceWorkVersionCountInfo, error) {

	//获取对应的任务状态
	status := c.spaceStatus(ctx, spaceId)

	completeStatusList := status.CompleteStatusList
	endStatusList := status.EndStatusList
	endAndNotCompleteStatusList := status.EndAndNotCompleteStatusList

	weekBegin, weekEnd := date.ThisWeekBeginEnd()
	weekBeginUnix, weekEndUnix := weekBegin.Unix(), weekEnd.Unix()

	queryList := []esV8.Query{
		esV8.NewTermQuery("space_id", spaceId),
	}

	if len(versionIds) != 0 {
		queryList = append(queryList, esV8.NewTermsQuery("version_id", utils.ToAnyArray(versionIds)...))
	}

	if startTime != 0 && endTime != 0 {
		queryList = append(queryList, esV8.NewRangeQuery("created_at").Gte(startTime).Lte(endTime))
	}

	query := esV8.NewBoolQuery().Must(queryList...)
	source := esV8.NewSearchSource().
		FetchSourceIncludeExclude([]string{
			"version_id",
			"work_item_status_id",
			"last_status_at",
			"doc.plan_start_at",
			"doc.plan_complete_at",
			"doc.priority",
		}, nil).
		Query(query).
		Size(10000)

	res, err := c.esDoSearch(source)
	if err != nil {
		return nil, err
	}

	type item struct {
		VersionId        int64  `json:"version_id"`
		WorkItemStatusId int64  `json:"work_item_status_id"`
		Priority         string `json:"priority"`
		PlanStartAt      int64  `json:"plan_start_at"`
		PlanCompleteAt   int64  `json:"plan_complete_at"`

		Doc struct {
			PlanStartAt    int64  `json:"plan_start_at"`
			PlanCompleteAt int64  `json:"plan_complete_at"`
			Priority       string `json:"priority"`
		} `json:"doc"`
	}

	var list []item
	list = stream.Map(res.Hits.Hits, func(v *esV8.SearchHit) item {
		var e item
		json.Unmarshal(v.Source, &e)
		e.PlanStartAt = e.Doc.PlanStartAt
		e.PlanCompleteAt = e.Doc.PlanCompleteAt
		e.Priority = e.Doc.Priority
		return e
	})

	now := time.Now().Unix()

	m := map[int64]*vo.SpaceWorkVersionCountInfo{}

	for _, v := range list {
		info := m[v.VersionId]
		if info == nil {
			info = &vo.SpaceWorkVersionCountInfo{
				PriorityInfo: map[string]int64{},
			}
			m[v.VersionId] = info
		}

		info.TotalCount += 1

		isProcessing := !slices.Contains(endStatusList, v.WorkItemStatusId)

		// 进行中
		if isProcessing {
			// 待办计数
			info.ProcessingCount += 1

			// 逾期计数
			if v.PlanCompleteAt < now {
				info.ExpiredCount++
			}

			// 优先级计数
			info.PriorityInfo[v.Priority] += 1

			// 本周待办计数
			if date.HasInter(v.PlanStartAt, v.PlanCompleteAt, weekBeginUnix, weekEndUnix) { //是否存在交集
				info.WeekProcessingCount++
			}
		} else if slices.Contains(completeStatusList, v.WorkItemStatusId) { //完成
			info.CompletedCount++
		} else if slices.Contains(endAndNotCompleteStatusList, v.WorkItemStatusId) { //已关闭或已终止
			info.ClosedOrTerminatedCount++
		}
	}

	// 计算完成度
	for _, v := range m {
		if v.TotalCount == 0 {
			continue
		}
		v.CompleteRate = int64(100 - (float64(v.ProcessingCount) * 100 / float64(v.TotalCount)))
	}

	return m, nil
}

func (c *staticsEsRepo) GetSpaceUserCount(ctx context.Context, spaceId int64, startTime, endTime int64) (map[int64]*vo.SpaceUserCountInfo, error) {

	//获取对应的任务状态
	status := c.spaceStatus(ctx, spaceId)

	endStatusList := status.EndStatusList
	completeStatusList := status.CompleteStatusList
	endAndNotCompleteStatusList := status.EndAndNotCompleteStatusList

	weekBegin, weekEnd := date.ThisWeekBeginEnd()
	weekBeginUnix, weekEndUnix := weekBegin.Unix(), weekEnd.Unix()

	queryList := []esV8.Query{
		esV8.NewTermQuery("space_id", spaceId),
	}

	if startTime != 0 && endTime != 0 {
		queryList = append(queryList, esV8.NewRangeQuery("created_at").Gte(startTime).Lte(endTime))
	}

	query := esV8.NewBoolQuery().Must(queryList...)
	source := esV8.NewSearchSource().
		FetchSourceIncludeExclude([]string{
			"id",
			"work_item_status_id",
			"work_item_type_key",
			"doc.plan_start_at",
			"doc.plan_complete_at",
			"doc.directors",
			"doc.participators",
			"doc.node_directors",
			"doc.priority",
		}, nil).
		Query(query).
		Size(10000)

	res, err := c.esDoSearch(source)
	if err != nil {
		return nil, err
	}

	type entity struct {
		Id               int64  `json:"id"`
		WorkItemStatusId int64  `json:"work_item_status_id"`
		WorkItemTypeKey  string `json:"work_item_type_key"`

		PlanStartAt    int64  `json:"doc.plan_start_at"`
		PlanCompleteAt int64  `json:"doc.plan_complete_at"`
		Priority       string `json:"doc.priority"`
		Directors      []int64
		NodeDirectors  []int64
		Participators  []int64

		Doc struct {
			PlanStartAt    int64    `json:"plan_start_at"`
			PlanCompleteAt int64    `json:"plan_complete_at"`
			Directors      []string `json:"directors"`
			Participators  []string `json:"participators"`
			NodeDirectors  []string `json:"node_directors"`
			Priority       string   `json:"priority"`
		} `json:"doc"`
	}

	var list []entity
	list = stream.Map(res.Hits.Hits, func(v *esV8.SearchHit) entity {
		var e entity
		json.Unmarshal(v.Source, &e)

		e.PlanStartAt = e.Doc.PlanStartAt
		e.PlanCompleteAt = e.Doc.PlanCompleteAt
		e.Directors = utils.ToInt64Array(e.Doc.Directors)
		e.Participators = utils.ToInt64Array(e.Doc.Participators)
		e.NodeDirectors = utils.ToInt64Array(e.Doc.NodeDirectors)
		e.Priority = e.Doc.Priority
		return e
	})

	userCountMap := map[int64]*vo.SpaceUserCountInfo{}

	now := time.Now().Unix()

	for _, v := range list {
		userIds := v.Participators

		isProcessing := !slices.Contains(endStatusList, v.WorkItemStatusId)
		isCompleted := slices.Contains(completeStatusList, v.WorkItemStatusId)
		isClosedOrTerminated := slices.Contains(endAndNotCompleteStatusList, v.WorkItemStatusId)
		expired := v.PlanCompleteAt < now

		for _, userId := range userIds {
			isDirector := slices.Contains(v.Directors, userId)
			isNodeDirector := slices.Contains(v.NodeDirectors, userId)

			info := userCountMap[userId]
			if info == nil {
				info = &vo.SpaceUserCountInfo{
					PriorityInfo: map[string]int64{},
				}
				userCountMap[userId] = info
			}

			info.TotalCount++
			// 进行中
			if isProcessing {
				if isDirector {
					// 待办计数
					info.ProcessingCount += 1

					// 优先级计数
					info.PriorityInfo[v.Priority] += 1

					// 本周待办计数
					if date.HasInter(v.PlanStartAt, v.PlanCompleteAt, weekBeginUnix, weekEndUnix) { //是否存在交集
						info.WeekProcessingCount++
					}
				}

				// 逾期
				if expired {
					if isNodeDirector && (v.WorkItemTypeKey == string(consts.WorkItemTypeKey_Task) || v.WorkItemTypeKey == string(consts.WorkItemTypeKey_SubTask)) {
						info.ExpiredCount++
					} else if isDirector && (v.WorkItemTypeKey == string(consts.WorkItemTypeKey_StateTask)) {
						info.ExpiredCount++
					}
				}
			}

			if isCompleted {
				info.CompletedCount++
			}

			if isClosedOrTerminated {
				info.ClosedOrTerminatedCount++
			}
		}
	}

	for _, v := range userCountMap {
		if v.TotalCount == 0 {
			continue
		}
		v.CompleteRate = int64(v.CompletedCount+v.ClosedOrTerminatedCount) * 100 / v.TotalCount
	}

	return userCountMap, nil
}

func (c *staticsEsRepo) spaceStatus(ctx context.Context, spaceIds ...int64) vo.SpaceStatus {

	statusList, err := c.statusRepo.GetWorkItemStatusItemsBySpaceIds(ctx, spaceIds)
	if err != nil {
		return vo.SpaceStatus{}
	}

	completeStatusItems := statusList.GetCompleteStatus()
	archiveStatusItems := statusList.GetArchivedStatus()

	completeStatusList := stream.Map(completeStatusItems, func(item *work_item_status.WorkItemStatusItem) int64 {
		return item.Id
	})
	archiveStatusList := stream.Map(archiveStatusItems, func(item *work_item_status.WorkItemStatusItem) int64 {
		return item.Id
	})

	//获取对应的任务状态
	return vo.SpaceStatus{
		EndStatusList:               archiveStatusList,
		EndAndNotCompleteStatusList: stream.Diff(archiveStatusList, completeStatusList),
		CompleteStatusList:          completeStatusList,
	}
}

func (c *staticsEsRepo) GetSpaceWorkbenchCount(ctx context.Context, uid int64, spaceId int64) (*vo.SpaceWorkbenchCountInfo, error) {
	if spaceId <= 0 {
		return nil, nil
	}

	status := c.spaceStatus(ctx, spaceId)

	// 获取总数
	totalSource := esV8.NewSearchSource().
		Query(
			esV8.NewBoolQuery().Must(
				esV8.NewTermQuery(search_es2.SpaceIdField.EsKeyword(), spaceId),
			),
		).
		TrackTotalHits(true).
		Size(0)

	// 获取进行中任务
	processingSource := esV8.NewSearchSource().
		Query(
			esV8.NewBoolQuery().Must(
				esV8.NewTermQuery(search_es2.SpaceIdField.EsKeyword(), spaceId),
			).MustNot(
				esV8.NewTermsQuery(search_es2.WorkItemStatusIdField.EsKeyword(), utils.ToAnyArray(status.EndStatusList)...),
			),
		).
		TrackTotalHits(true).
		Size(0)
	// 获取已完成任务
	completedSource := esV8.NewSearchSource().
		Query(
			esV8.NewBoolQuery().Must(
				esV8.NewTermQuery(search_es2.SpaceIdField.EsKeyword(), spaceId),
				esV8.NewTermsQuery(search_es2.WorkItemStatusIdField.EsKeyword(), utils.ToAnyArray(status.CompleteStatusList)...),
			),
		).
		TrackTotalHits(true).
		Size(0)

	// 获取逾期任务
	expiredSource := esV8.NewSearchSource().
		Query(
			esV8.NewBoolQuery().Must(
				esV8.NewTermQuery(search_es2.SpaceIdField.EsKeyword(), spaceId),
				esV8.NewRangeQuery(search_es2.PlanCompleteAtField.EsKeyword()).Lt(time.Now().Unix()),
			).MustNot(
				esV8.NewTermsQuery(search_es2.WorkItemStatusIdField.EsKeyword(), utils.ToAnyArray(status.EndStatusList)...),
			),
		).
		TrackTotalHits(true).
		Size(0)

	// 获取关注任务
	followedSource := esV8.NewSearchSource().
		Query(
			esV8.NewBoolQuery().Must(
				esV8.NewTermQuery(search_es2.SpaceIdField.EsKeyword(), spaceId),
				esV8.NewTermQuery(search_es2.FollowersField.EsKeyword(), cast.ToString(uid)),
			),
		).
		TrackTotalHits(true).
		Size(0)

	es, err := c.MultiCountByEs(ctx, []*esV8.SearchSource{totalSource, processingSource, completedSource, expiredSource, followedSource})
	if err != nil {
		return nil, err
	}

	var total = es[0].Value
	var progressingCount = es[1].Value
	var completedCount = es[2].Value
	var progressingAndExpiredCount = es[3].Value
	var followedCount = es[4].Value

	var completeRate int64
	var progressingAndExpiredRate int64
	if total > 0 {
		completeRate = int64(100 - float64(progressingCount)*100/float64(total))
		progressingAndExpiredRate = progressingAndExpiredCount * 100 / total
	}

	return &vo.SpaceWorkbenchCountInfo{
		ProcessingCount:           progressingCount,
		CompleteCount:             completedCount,
		ProcessingAndExpiredCount: progressingAndExpiredCount,
		ProcessingAndExpiredRate:  progressingAndExpiredRate,
		CompleteRate:              completeRate,
		FollowCount:               followedCount,
		TotalCount:                total,
	}, nil
}

func (c *staticsEsRepo) GetUserFollowCount(ctx context.Context, userId int64, spaceIds []int64) (int64, error) {
	source := esV8.NewSearchSource().
		FetchSource(false).
		Query(esV8.NewBoolQuery().Must(
			esV8.NewTermsQuery("space_id", utils.ToAnyArray(spaceIds)...),
			esV8.NewTermQuery("doc.followers", userId),
		)).Size(0)

	res, err := c.esDoSearch(source)
	if err != nil {
		return 0, err
	}

	return res.TotalHits(), nil
}

func (c *staticsEsRepo) GetWorkItemCountBySpaceFlowId(ctx context.Context, spaceId int64, workFlowId int64) (int64, error) {
	var count int64
	err := c.data.RoDB(ctx).
		Model(&db.SpaceWorkItemV2{}).
		Where("space_id =? and work_flow_id = ?", spaceId, workFlowId).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *staticsEsRepo) MultiCountByEs(ctx context.Context, sources []*esV8.SearchSource) ([]vo.CountInfo, error) {
	es := c.data.es
	indexName := c.data.conf.Es.Index

	var params []any
	for _, v := range sources {
		params = append(params, json.RawMessage(`{}`), v)
	}

	queries1 := stream.Map(params, func(v any) []string {
		s, _ := json.Marshal(v)
		return []string{string(s), "\n"}
	})

	queries2 := stream.Flat(queries1)

	body := strings.Join(queries2, "")

	esRes, err := es.Msearch(
		strings.NewReader(body),
		es.Msearch.WithContext(ctx),
		es.Msearch.WithIndex(indexName),
	)

	c.log.Info(esRes.String())

	if err != nil {
		c.log.Error(err)
		return nil, err
	}

	defer esRes.Body.Close()

	if esRes.StatusCode == http.StatusNotFound {
		return make([]vo.CountInfo, len(sources)), nil
	}

	if esRes.IsError() {
		c.log.Error(esRes.String())
		return nil, errors.New(esRes.String())
	}
	if esRes.StatusCode != http.StatusOK {
		c.log.Error(esRes.String())
		return nil, errors.New(esRes.String())
	}

	type esResp struct {
		Responses []*esV8.SearchResult `json:"responses"`
	}

	var esRespData esResp
	err = json.NewDecoder(esRes.Body).Decode(&esRespData)
	if err != nil {
		c.log.Error(err)
		return nil, err
	}

	return stream.Map(esRespData.Responses, func(v *esV8.SearchResult) vo.CountInfo {
		return vo.CountInfo{
			Value: v.TotalHits(),
		}
	}), nil
}
