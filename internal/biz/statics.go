package biz

import (
	"cmp"
	"context"
	"go-cs/api/comm"
	v12 "go-cs/api/search/v1"
	v1 "go-cs/api/workbench/v1"
	"go-cs/internal/bean/vo"
	"go-cs/internal/consts"
	"go-cs/internal/domain/search/condition_translater"
	search_es2 "go-cs/internal/domain/search/search_es"
	member_repo "go-cs/internal/domain/space_member/repo"
	statics_repo "go-cs/internal/domain/statics/repo"
	"go-cs/internal/domain/work_item_status"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/date"
	"go-cs/internal/utils/errs"
	esV8 "go-cs/internal/utils/es/v8"
	"go-cs/pkg/stream"
	"slices"

	witem_status_repo "go-cs/internal/domain/work_item_status/repo"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type StaticsUsecase struct {
	log                *log.Helper
	repo               statics_repo.StaticsRepo
	spaceMemberRepo    member_repo.SpaceMemberRepo
	tm                 trans.Transaction
	workItemStatusRepo witem_status_repo.WorkItemStatusRepo
}

func NewStaticsUsecase(repo statics_repo.StaticsRepo, tm trans.Transaction, spaceMemberRepo member_repo.SpaceMemberRepo, workItemStatusRepo witem_status_repo.WorkItemStatusRepo, logger log.Logger) *StaticsUsecase {
	moduleName := "StaticsUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &StaticsUsecase{
		log:                hlog,
		repo:               repo,
		spaceMemberRepo:    spaceMemberRepo,
		tm:                 tm,
		workItemStatusRepo: workItemStatusRepo,
	}
}

func (uc *StaticsUsecase) GetWorkbenchCount(ctx context.Context, userId int64) (*v1.GetWorkBenchCountReplyData, error) {

	spaceIds, _ := uc.spaceMemberRepo.GetUserSpaceIdList(ctx, userId)
	data, err := uc.repo.GetWorkbenchCount(ctx, userId, spaceIds)

	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	statusList, err := uc.workItemStatusRepo.GetWorkItemStatusItemsBySpaceIds(ctx, spaceIds)

	completeStatusItems := statusList.GetCompleteStatus()
	archiveStatusItems := statusList.GetArchivedStatus()

	completeStatusList := stream.Map(completeStatusItems, func(item *work_item_status.WorkItemStatusItem) string {
		return cast.ToString(item.Id)
	})
	archiveStatusList := stream.Map(archiveStatusItems, func(item *work_item_status.WorkItemStatusItem) string {
		return cast.ToString(item.Id)
	})

	uidStr := cast.ToString(userId)

	begin, end := date.ThisWeekBeginEnd()

	var weekStart = date.Format(begin)
	var weekEnd = date.Format(end)

	weekCond := []*v12.Condition{
		{Field: "plan_time", Operator: "BETWEEN", Values: []string{weekStart, weekEnd}},
	}

	weekCompleteCond := []*v12.Condition{
		{Field: "last_status_at", Operator: "BETWEEN", Values: []string{weekStart, weekEnd}},
	}

	participatorCond := []*v12.Condition{
		{Field: "participators", Operator: "IN", Values: []string{uidStr}},
	}

	directorCond := []*v12.Condition{
		{Field: "directors", Operator: "IN", Values: []string{uidStr}},
	}

	workItemDoneStatusCond := []*v12.Condition{
		{Field: "work_item_status_id", Operator: "IN", Values: completeStatusList},
	}

	workItemProcessingCond := []*v12.Condition{
		{Field: "work_item_status_id", Operator: "NOT_IN", Values: archiveStatusList},
	}

	weekDoneCond := &v12.ConditionGroup{
		Conjunction: "AND",
		Conditions: stream.Concat(
			participatorCond,
			weekCompleteCond,
			workItemDoneStatusCond,
		),
	}

	weekProcessingCond := &v12.ConditionGroup{
		Conjunction: "AND",
		Conditions: stream.Concat(
			directorCond,
			weekCond,
			workItemProcessingCond,
		),
	}

	doneCond := &v12.ConditionGroup{
		Conjunction: "AND",
		Conditions: stream.Concat(
			participatorCond,
			workItemDoneStatusCond,
		),
	}

	processingCond := &v12.ConditionGroup{
		Conjunction: "AND",
		Conditions: stream.Concat(
			directorCond,
			workItemProcessingCond,
		),
	}

	participatedProcessingCond := &v12.ConditionGroup{
		Conjunction: "AND",
		Conditions: stream.Concat(
			participatorCond,
			workItemProcessingCond,
		),
	}

	followedCond := &v12.ConditionGroup{
		Conjunction: "AND",
		Conditions: []*v12.Condition{
			{Field: "followers", Operator: "IN", Values: []string{uidStr}},
		},
	}

	createdCond := &v12.ConditionGroup{
		Conjunction: "AND",
		Conditions: []*v12.Condition{
			{Field: "user_id", Operator: "EQ", Values: []string{uidStr}},
		},
	}

	followedCount, _ := uc.repo.GetUserFollowCount(ctx, userId, spaceIds)

	return &v1.GetWorkBenchCountReplyData{
		Processing:             data.ProcessingCount,
		WeekProcessing:         data.WeekProcessingCount,
		ParticipatedProcessing: data.ParticipatedProcessingCount,
		Completed:              data.CompleteCount,
		WeekCompleted:          data.WeekCompleteCount,
		Followed:               followedCount,
		Created:                data.CreatedCount,

		ProcessingConditionGroup:             processingCond,
		WeekProcessingConditionGroup:         weekProcessingCond,
		ParticipatedProcessingConditionGroup: participatedProcessingCond,
		DoneConditionGroup:                   doneCond,
		WeekDoneConditionGroup:               weekDoneCond,
		FollowedConditionGroup:               followedCond,
		CreatedConditionGroup:                createdCond,
	}, nil
}

func (uc *StaticsUsecase) GetSpaceWorkbenchCount2(ctx context.Context, uid int64, spaceId int64, conditionGroups []*v12.ConditionGroup) (*v1.GetSpaceWorkBenchCountReply2_Data, error) {
	buildSource := func(cg *v12.ConditionGroup) *esV8.SearchSource {
		translateCtx := condition_translater.Ctx{
			Ctx:        ctx,
			StatusRepo: uc.workItemStatusRepo,
			SpaceIds:   []int64{spaceId},
		}
		cg = translateCtx.TranslateCond(cg)

		condition, err := search_es2.ConvertCondition(cg)
		if err != nil {
			uc.log.Error(err)
			return nil
		}

		var spaceIds = []int64{spaceId}
		var spaceIdsAny = utils.ToAnyArray(spaceIds)

		// 构建filter条件
		sourceQuery := search_es2.BuildCondition(condition)
		if sourceQuery == nil {
			sourceQuery = esV8.NewBoolQuery()
		}
		// 添加spaceId过滤条件
		sourceQuery.Must(esV8.NewTermsQuery(search_es2.SpaceIdField.EsKeyword(), spaceIdsAny...))

		source := esV8.NewSearchSource().
			Query(sourceQuery).
			TrackTotalHits(true).
			Size(0)

		return source
	}

	sources := stream.Map(conditionGroups, func(item *v12.ConditionGroup) *esV8.SearchSource {
		return buildSource(item)
	})

	// 获取数据
	result, err := uc.repo.MultiCountByEs(ctx, sources)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var items []*v1.GetSpaceWorkBenchCountReply2_Data_Item
	for _, v := range result {
		items = append(items, &v1.GetSpaceWorkBenchCountReply2_Data_Item{
			Value: v.Value,
		})
	}

	return &v1.GetSpaceWorkBenchCountReply2_Data{Results: items}, nil
}

func (uc *StaticsUsecase) GetSpaceWorkbenchCount(ctx context.Context, userId, spaceId int64) (*v1.GetSpaceWorkBenchCountReplyData, error) {
	statusList, err := uc.workItemStatusRepo.GetWorkItemStatusItemsBySpace(ctx, spaceId)

	completeStatusItems := statusList.GetCompleteStatus()
	archiveStatusItems := statusList.GetArchivedStatus()

	completeStatusList := stream.Map(completeStatusItems, func(item *work_item_status.WorkItemStatusItem) string {
		return cast.ToString(item.Id)
	})
	archiveStatusList := stream.Map(archiveStatusItems, func(item *work_item_status.WorkItemStatusItem) string {
		return cast.ToString(item.Id)
	})

	completeCond := &v12.ConditionGroup{
		Conjunction: "AND",
		Conditions: []*v12.Condition{
			{Field: "work_item_status_id", Operator: "IN", Values: completeStatusList},
		},
	}

	processingCond := &v12.ConditionGroup{
		Conjunction: "AND",
		Conditions: []*v12.Condition{
			{Field: "work_item_status_id", Operator: "NOT_IN", Values: archiveStatusList},
		},
	}

	_, err = uc.spaceMemberRepo.GetSpaceMember(ctx, spaceId, userId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	data, err := uc.repo.GetSpaceWorkbenchCount(ctx, userId, spaceId)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_DB_QUERY_FAIL)
	}

	return &v1.GetSpaceWorkBenchCountReplyData{
		Processing:                data.ProcessingCount,
		ProcessingAndExpiredCount: data.ProcessingAndExpiredCount,
		ProcessingAndExpiredRate:  data.ProcessingAndExpiredRate,
		Completed:                 data.CompleteCount,
		CompleteRate:              data.CompleteRate,
		Total:                     data.TotalCount,
		Followed:                  data.FollowCount,
		ProcessingConditionGroup:  processingCond,
		DoneConditionGroup:        completeCond,
		//FollowedConditionGroup:    followedCond,
	}, nil
}

func (uc *StaticsUsecase) GetSpaceWorkObjectCountByIds(ctx context.Context, userId int64, spaceId int64, workObjectIds []int64, startTime, endTime int64) (*v1.GetSpaceWorkObjectCountReplyData, error) {
	_, err := uc.spaceMemberRepo.GetSpaceMember(ctx, spaceId, userId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	countMap, err := uc.repo.GetSpaceWorkObjectCountByIds(ctx, spaceId, workObjectIds, startTime, endTime)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var items []*v1.GetSpaceWorkObjectCountReplyData_Item

	allPriority := consts.GetAllPriority()
	for _, id := range workObjectIds {
		info := countMap[id]
		if info == nil {
			info = &vo.SpaceWorkObjectCountInfo{
				PriorityInfo: map[string]int64{},
			}
		}

		// 填充关联任务数为0的优先级
		for _, v := range allPriority {
			if _, ok := info.PriorityInfo[v]; !ok {
				info.PriorityInfo[v] = 0
			}
		}

		priorityInfos := stream.ToSlice(info.PriorityInfo, func(priority string, count int64) *v1.GetSpaceWorkObjectCountReplyData_PriorityInfo {
			return &v1.GetSpaceWorkObjectCountReplyData_PriorityInfo{
				Priority: priority,
				Count:    count,
			}
		})

		// 对优先级信息排序
		slices.SortFunc(priorityInfos, func(a, b *v1.GetSpaceWorkObjectCountReplyData_PriorityInfo) int {
			av := consts.GetWorkItemPriorityPower(a.Priority)
			bv := consts.GetWorkItemPriorityPower(b.Priority)
			return -cmp.Compare(av, bv)
		})

		items = append(items, &v1.GetSpaceWorkObjectCountReplyData_Item{
			Id:                 id,
			Processing:         info.ProcessingCount,
			WeekProcessing:     info.WeekProcessingCount,
			Completed:          info.CompletedCount,
			Expired:            info.ExpiredCount,
			ClosedOrTerminated: info.ClosedOrTerminatedCount,
			Total:              info.TotalCount,
			CompleteRate:       info.CompleteRate,
			PriorityInfos:      priorityInfos,
		})
	}

	return &v1.GetSpaceWorkObjectCountReplyData{
		List: items,
	}, nil
}

func (uc *StaticsUsecase) GetSpaceUserCount(ctx context.Context, uid int64, spaceId int64, startTime, endTime int64) (*v1.GetSpaceUserCountReplyData, error) {
	_, err := uc.spaceMemberRepo.GetSpaceMember(ctx, spaceId, uid)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	countMap, err := uc.repo.GetSpaceUserCount(ctx, spaceId, startTime, endTime)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	allPriority := consts.GetAllPriority()
	var items []*v1.GetSpaceUserCountReplyData_Item
	for id, v := range countMap {
		info := countMap[id]
		if info == nil {
			info = &vo.SpaceUserCountInfo{
				PriorityInfo: map[string]int64{},
			}
		}

		// 填充关联任务数为0的优先级
		for _, v := range allPriority {
			if _, ok := info.PriorityInfo[v]; !ok {
				info.PriorityInfo[v] = 0
			}
		}

		priorityInfos := stream.ToSlice(info.PriorityInfo, func(priority string, count int64) *v1.GetSpaceUserCountReplyData_PriorityInfo {
			return &v1.GetSpaceUserCountReplyData_PriorityInfo{
				Priority: priority,
				Count:    count,
			}
		})

		// 对优先级信息排序
		slices.SortFunc(priorityInfos, func(a, b *v1.GetSpaceUserCountReplyData_PriorityInfo) int {
			av := consts.GetWorkItemPriorityPower(a.Priority)
			bv := consts.GetWorkItemPriorityPower(b.Priority)
			return -cmp.Compare(av, bv)
		})

		items = append(items, &v1.GetSpaceUserCountReplyData_Item{
			Id:                 id,
			Processing:         v.ProcessingCount,
			WeekProcessing:     v.WeekProcessingCount,
			ClosedOrTerminated: v.ClosedOrTerminatedCount,
			Expired:            v.ExpiredCount,
			Completed:          v.CompletedCount,
			CompleteRate:       v.CompleteRate,
			Total:              v.TotalCount,
			PriorityInfos:      priorityInfos,
		})
	}

	return &v1.GetSpaceUserCountReplyData{
		List: items,
	}, nil
}

func (uc *StaticsUsecase) GetSpaceVersionCount(ctx context.Context, uid int64, spaceId int64, startTime, endTime int64) (*v1.GetSpaceVersionCountReplyData, error) {
	_, err := uc.spaceMemberRepo.GetSpaceMember(ctx, spaceId, uid)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	countMap, err := uc.repo.GetSpaceWorkVersionCountByIds(ctx, spaceId, nil, startTime, endTime)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	allPriority := consts.GetAllPriority()
	var items []*v1.GetSpaceVersionCountReplyData_Item
	for id, v := range countMap {
		info := countMap[id]
		if info == nil {
			info = &vo.SpaceWorkVersionCountInfo{
				PriorityInfo: map[string]int64{},
			}
		}

		// 填充关联任务数为0的优先级
		for _, v := range allPriority {
			if _, ok := info.PriorityInfo[v]; !ok {
				info.PriorityInfo[v] = 0
			}
		}

		priorityInfos := stream.ToSlice(info.PriorityInfo, func(priority string, count int64) *v1.GetSpaceVersionCountReplyData_PriorityInfo {
			return &v1.GetSpaceVersionCountReplyData_PriorityInfo{
				Priority: priority,
				Count:    count,
			}
		})

		// 对优先级信息排序
		slices.SortFunc(priorityInfos, func(a, b *v1.GetSpaceVersionCountReplyData_PriorityInfo) int {
			av := consts.GetWorkItemPriorityPower(a.Priority)
			bv := consts.GetWorkItemPriorityPower(b.Priority)
			return -cmp.Compare(av, bv)
		})

		items = append(items, &v1.GetSpaceVersionCountReplyData_Item{
			Id:                 id,
			Processing:         v.ProcessingCount,
			WeekProcessing:     v.WeekProcessingCount,
			ClosedOrTerminated: v.ClosedOrTerminatedCount,
			Completed:          v.CompletedCount,
			CompleteRate:       v.CompleteRate,
			Expired:            v.ExpiredCount,
			Total:              v.TotalCount,
			PriorityInfos:      priorityInfos,
		})
	}

	return &v1.GetSpaceVersionCountReplyData{
		List: items,
	}, nil
}
