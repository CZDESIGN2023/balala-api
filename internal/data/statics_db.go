package data

import (
	"context"
	"fmt"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo"
	"go-cs/internal/domain/search/search2"
	"go-cs/internal/domain/work_item_status"
	status_repo "go-cs/internal/domain/work_item_status/repo"
	"go-cs/internal/utils"
	"go-cs/internal/utils/date"
	esV8 "go-cs/internal/utils/es/v8"
	"go-cs/pkg/stream"
	"slices"
	"strconv"
	"time"

	"go-cs/internal/domain/statics/repo"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type staticsRepo struct {
	baseRepo
	statusRepo status_repo.WorkItemStatusRepo
}

func NewStaticsRepo(data *Data, statusRepo status_repo.WorkItemStatusRepo, logger log.Logger) repo.StaticsRepo {
	moduleName := "StaticsRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	repo := &staticsRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
		statusRepo: statusRepo,
	}
	return repo
}

// 获取tag的被使用关联数量
func (c *staticsRepo) GetTagRelationCounts(ctx context.Context, tagIds []int64) ([]*vo.TagRelationCountInfo, error) {

	if len(tagIds) == 0 {
		return nil, nil
	}

	var rows []*vo.TagRelationCountInfo

	sql := fmt.Sprintf(`SELECT count(dat.id) total_num, tags.id FROM %v dat
		CROSS JOIN json_table(concat(doc->>'$.tags'), '$[*]' COLUMNS( id VARCHAR(50) PATH '$')) AS tags 
		where tags.id in (?) group by  tags.id`, (&db.SpaceWorkItemV2{}).TableName())

	err := c.data.RoDB(ctx).Raw(sql, tagIds).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *staticsRepo) GetSpaceTagRelationCounts(ctx context.Context, spaceId int64, tagIds []int64) ([]*vo.TagRelationCountInfo, error) {
	if len(tagIds) == 0 {
		return nil, nil
	}

	sql := `
SELECT
	count( dat.id ) total_num,
	tags.id 
FROM
	space_work_item_v2 dat
	CROSS JOIN json_table (
		concat( doc ->> '$.tags' ),
		'$[*]' COLUMNS ( id VARCHAR ( 50 ) PATH '$' )
	) AS tags 
WHERE
	dat.space_id = ? and
	tags.id IN ? 
GROUP BY
	tags.id
`

	var rows []*vo.TagRelationCountInfo

	err := c.data.RoDB(ctx).Raw(sql, spaceId, tagIds).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *staticsRepo) TagCountMap(ctx context.Context, tagIds []int64) (map[int64]int64, error) {
	counts, err := c.GetTagRelationCounts(ctx, tagIds)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(counts, func(i int, t *vo.TagRelationCountInfo) (int64, int64) {
		return t.Id, int64(t.TotalNum)
	})

	return m, nil
}

func (c *staticsRepo) SpaceTagCountMap(ctx context.Context, spaceId int64, tagIds []int64) (map[int64]int64, error) {
	counts, err := c.GetSpaceTagRelationCounts(ctx, spaceId, tagIds)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(counts, func(i int, t *vo.TagRelationCountInfo) (int64, int64) {
		return t.Id, int64(t.TotalNum)
	})

	return m, nil
}

func (c *staticsRepo) GetWorkbenchCount(ctx context.Context, userId int64, spaceIds []int64) (*vo.UserWorkbenchCountInfo, error) {
	if len(spaceIds) == 0 {
		return nil, nil
	}

	//获取对应的任务状态
	status := c.spaceStatus(ctx, spaceIds...)

	completeStatusList := status.CompleteStatusList
	endStatusList := status.EndStatusList

	sql := `
select
	id,
	user_id,
	work_item_status,
	last_status_at,
	doc->'$.plan_start_at' plan_start_at,
	doc->'$.plan_complete_at' plan_complete_at,
	doc->'$.directors' directors,
	doc->'$.participators' participators
from
	space_work_item_v2 swiv
where
	space_id in (?)
	and ? member of(doc->'$.participators')
`
	type entity struct {
		Id               int64    `json:"id"`
		UserId           int64    `json:"user_id"`
		WorkItemStatusId int64    `json:"work_item_status_id"`
		LastStatusAt     int64    `json:"last_status_at"`
		PlanStartAt      int64    `json:"plan_start_at"`
		PlanCompleteAt   int64    `json:"plan_complete_at"`
		Directors        []string `gorm:"serializer:json"`
		Participators    []string `gorm:"serializer:json"`
	}

	userIdStr := cast.ToString(userId)

	var list []entity
	err := c.data.RoDB(ctx).Raw(sql, spaceIds, userIdStr).Find(&list).Error
	if err != nil {
		return nil, err
	}

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

func (c *staticsRepo) GetSpaceWorkObjectCountByIds(ctx context.Context, spaceId int64, workObjectIds []int64, startTime, endTime int64) (map[int64]*vo.SpaceWorkObjectCountInfo, error) {

	//获取对应的任务状态
	status := c.spaceStatus(ctx, spaceId)

	completeStatusList := status.CompleteStatusList
	endStatusList := status.EndStatusList

	type item struct {
		WorkObjectId     int64
		WorkItemStatusId int64
		Priority         string
		PlanStartAt      int64
		PlanCompleteAt   int64
	}

	tx := c.data.RoDB(ctx).
		Model(&db.SpaceWorkItemV2{}).
		Select(search2.SelectByColumn("work_object_id", "work_item_status", "priority", "plan_start_at", "plan_complete_at"))

	if len(workObjectIds) == 0 {
		tx = tx.Where("space_id = ?", spaceId)
	} else {
		tx = tx.Where("work_object_id in ?", spaceId, workObjectIds)
	}

	if startTime > 0 && endTime > 0 {
		tx = tx.Where("plan_start_at >= ? and plan_complete_at <= ?", startTime, endTime)
	}

	var list []item
	err := tx.Find(&list).Error
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	weekBegin, weekEnd := date.ThisWeekBeginEnd()
	weekBeginUnix, weekEndUnix := weekBegin.Unix(), weekEnd.Unix()

	m := map[int64]*vo.SpaceWorkObjectCountInfo{}

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

			// 周待办计数
			if date.HasInter(v.PlanStartAt, v.PlanCompleteAt, weekBeginUnix, weekEndUnix) { //是否存在交集
				info.WeekProcessingCount++
			}

			// 优先级计数
			info.PriorityInfo[v.Priority] += 1

			// 已过期
			if v.PlanCompleteAt < now {
				info.ExpiredCount += 1
			}
		}

		// 已完成
		if slices.Contains(completeStatusList, v.WorkItemStatusId) {
			info.CompletedCount += 1
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

func (c *staticsRepo) GetSpaceWorkVersionCountByIds(ctx context.Context, spaceId int64, ids []int64, startTime, endTime int64) (map[int64]*vo.SpaceWorkVersionCountInfo, error) {

	//获取对应的任务状态
	status := c.spaceStatus(ctx, spaceId)

	completeStatusList := status.CompleteStatusList
	endStatusList := status.EndStatusList

	type item struct {
		VersionId        int64
		WorkItemStatusId int64
		PlanStartAt      int64
		PlanCompleteAt   int64
	}

	tx := c.data.RoDB(ctx).Model(&db.SpaceWorkItemV2{}).
		Select(search2.SelectByColumn("version_id", "work_item_status", "priority", "plan_start_at", "plan_complete_at")).
		Where("space_id = ?", spaceId)

	if len(ids) > 0 {
		tx = tx.Where("version_id in ?", ids)
	}

	if startTime > 0 && endTime > 0 {
		tx = tx.Where("plan_start_at >= ? and plan_complete_at <= ?", startTime, endTime)
	}

	var list []item
	err := tx.Find(&list).Error
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()

	weekBegin, weekEnd := date.ThisWeekBeginEnd()
	weekBeginUnix, weekEndUnix := weekBegin.Unix(), weekEnd.Unix()

	m := map[int64]*vo.SpaceWorkVersionCountInfo{}

	for _, v := range list {
		info := m[v.VersionId]
		if info == nil {
			info = &vo.SpaceWorkVersionCountInfo{}
			m[v.VersionId] = info
		}

		info.TotalCount += 1

		isProcessing := !slices.Contains(endStatusList, v.WorkItemStatusId)

		// 进行中
		if isProcessing {
			// 待办计数
			info.ProcessingCount += 1

			// 周待办计数
			if date.HasInter(v.PlanStartAt, v.PlanCompleteAt, weekBeginUnix, weekEndUnix) { //是否存在交集
				info.WeekProcessingCount++
			}

			// 逾期计数
			if v.PlanCompleteAt < now {
				info.ExpiredCount++
			}
		}

		// 已完成
		if slices.Contains(completeStatusList, v.WorkItemStatusId) {
			info.CompletedCount += 1
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

func (c *staticsRepo) GetSpaceUserCount(ctx context.Context, spaceId int64, startTime, endTime int64) (map[int64]*vo.SpaceUserCountInfo, error) {

	//获取对应的任务状态
	status := c.spaceStatus(ctx, spaceId)

	endStatusList := status.EndStatusList

	type item struct {
		Id               int64
		UserId           int64
		WorkItemStatusId int64
		PlanStartAt      int64
		PlanCompleteAt   int64
		Directors        []string `gorm:"serializer:json"`
	}

	tx := c.data.RoDB(ctx).
		Model(&db.SpaceWorkItemV2{}).
		Select(search2.SelectByColumn("id", "user_id", "work_item_status_id", "plan_start_at", "plan_complete_at", "directors"))

	if startTime > 0 && endTime > 0 {
		tx = tx.Where("plan_start_at >= ? and plan_complete_at <= ?", startTime, endTime)
	}

	var itemList []*item
	err := tx.Find(&itemList).Error
	if err != nil {
		return nil, err
	}

	type node struct {
		WorkItemId int64
		Directors  []string `gorm:"serializer:json"`
		directors  []int64
	}
	var nodeList []*node
	err = c.data.RoDB(ctx).Model(&db.SpaceWorkItemFlowV2{}).
		Where("space_id = ?", spaceId).
		Select("work_item_id", "directors").
		Find(&nodeList).Error
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(nodeList); i++ {
		nodeList[i].directors = utils.ToInt64Array(nodeList[i].Directors)
	}

	itemId2NodeListMap := stream.GroupBy(nodeList, func(t *node) int64 {
		return t.WorkItemId
	})

	itemId2UserIds := stream.MapValue(itemId2NodeListMap, func(nodeList []*node) []int64 {
		var userIds []int64

		for _, node := range nodeList {
			userIds = append(userIds, node.directors...)
		}

		return userIds
	})

	userCountMap := map[int64]*vo.SpaceUserCountInfo{}

	now := time.Now().Unix()

	weekBegin, weekEnd := date.ThisWeekBeginEnd()
	weekBeginUnix, weekEndUnix := weekBegin.Unix(), weekEnd.Unix()

	for _, v := range itemList {
		userIds := stream.Unique(append(itemId2UserIds[v.Id], v.UserId)) //节点负责人+创建人

		isProcessing := !slices.Contains(endStatusList, v.WorkItemStatusId)
		expired := v.PlanCompleteAt < now
		inWeek := date.HasInter(v.PlanStartAt, v.PlanCompleteAt, weekBeginUnix, weekEndUnix)

		for _, userId := range userIds {
			isDirector := slices.Contains(v.Directors, cast.ToString(userId))

			info := userCountMap[userId]
			if info == nil {
				info = &vo.SpaceUserCountInfo{}
				userCountMap[userId] = info
			}

			info.TotalCount++
			// 进行中
			if isProcessing && isDirector {
				// 待办计数
				info.ProcessingCount += 1

				// 周待办计数
				if inWeek { //是否存在交集
					info.WeekProcessingCount++
				}

				// 逾期计数
				if expired {
					info.ExpiredCount++
				}
			}
		}
	}

	for _, v := range userCountMap {
		if v.TotalCount == 0 {
			continue
		}
		v.CompleteRate = int64(100 - (float64(v.ProcessingCount) * 100 / float64(v.TotalCount)))
	}

	return userCountMap, nil
}

func (c *staticsRepo) spaceStatus(ctx context.Context, spaceIds ...int64) vo.SpaceStatus {
	//获取对应的任务状态
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

func (c *staticsRepo) GetSpaceWorkbenchCount(ctx context.Context, uid int64, spaceId int64) (*vo.SpaceWorkbenchCountInfo, error) {
	if spaceId <= 0 {
		return nil, nil
	}

	status := c.spaceStatus(ctx, spaceId)

	sql := `
select
	id,
	work_item_status,
	doc->'$.plan_start_at' plan_start_at,
	doc->'$.plan_complete_at' plan_complete_at
from
	space_work_item_v2 swiv
where
	space_id = ?
`
	type entity struct {
		Id               int64 `json:"id"`
		WorkItemStatusId int64 `json:"work_item_status"`
		PlanStartAt      int64 `json:"plan_start_at"`
		PlanCompleteAt   int64 `json:"plan_complete_at"`
	}

	var list []entity
	err := c.data.RoDB(ctx).Raw(sql, spaceId).Find(&list).Error
	if err != nil {
		return nil, err
	}

	var total = int64(len(list))
	var progressingCount int64
	var completeCount int64
	var progressingAndExpiredCount int64

	beginOfToday := date.TodayBegin().Unix()

	for _, v := range list {
		if !slices.Contains(status.EndStatusList, v.WorkItemStatusId) { //待办
			progressingCount++
			if v.PlanCompleteAt <= beginOfToday {
				progressingAndExpiredCount++
			}
		} else if slices.Contains(status.CompleteStatusList, v.WorkItemStatusId) { //已完成
			completeCount++
		}
	}

	var completeRate int64
	var progressingAndExpiredRate int64
	if total > 0 {
		completeRate = int64(100 - float64(progressingCount)*100/float64(total))
		progressingAndExpiredRate = progressingAndExpiredCount * 100 / total
	}

	return &vo.SpaceWorkbenchCountInfo{
		ProcessingCount:           progressingCount,
		ProcessingAndExpiredCount: progressingAndExpiredCount,
		ProcessingAndExpiredRate:  progressingAndExpiredRate,
		CompleteCount:             completeCount,
		CompleteRate:              completeRate,
		TotalCount:                total,
	}, nil
}

func (c *staticsRepo) GetSpaceWorkItemTagCount(ctx context.Context, spaceId int64) (map[int64]int64, error) {
	type item struct {
		Tags []string `gorm:"serializer:json"`
	}

	var list []item

	err := c.data.RoDB(ctx).Model(&db.SpaceWorkItemV2{}).
		Where("space_id = ?", spaceId).
		Where("JSON_LENGTH(doc->'$.tags') > 0").
		Select(search2.SelectByColumn("id", "tags")).
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	var countMap = map[string]int64{}
	for _, v := range list {
		for _, tag := range v.Tags {
			countMap[tag] += 1
		}
	}

	ret := stream.MapKey(countMap, func(k string) int64 {
		i, _ := strconv.ParseInt(k, 10, 64)
		return i
	})

	return ret, nil
}

func (c *staticsRepo) GetUserFollowCount(ctx context.Context, userId int64, spaceIds []int64) (int64, error) {
	var count int64
	err := c.data.RoDB(ctx).
		Model(&db.SpaceWorkItemV2{}).
		Where("space_id in (?)", spaceIds).
		Where("? MEMBER OF(doc->'$.followers')", cast.ToString(userId)).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *staticsRepo) GetWorkItemCountBySpaceFlowId(ctx context.Context, spaceId int64, workFlowId int64) (int64, error) {
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

func (c *staticsRepo) MultiCountByEs(ctx context.Context, sources []*esV8.SearchSource) ([]vo.CountInfo, error) {
	//TODO implement me
	panic("implement me")
}
