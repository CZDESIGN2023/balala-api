package vo

type TagRelationCountInfo struct {
	Id       int64 `gorm:"column:id" bson:"id" json:"id"`
	TotalNum int32 `gorm:"column:total_num" bson:"total_num" json:"total_num"`
}

type UserWorkbenchCountInfo struct {
	ProcessingCount             int64
	WeekProcessingCount         int64
	ParticipatedProcessingCount int64
	CreatedCount                int64
	CompleteCount               int64
	WeekCompleteCount           int64
	FollowCount                 int64
}

type SpaceWorkbenchCountInfo struct {
	ProcessingCount           int64
	ProcessingAndExpiredCount int64
	ProcessingAndExpiredRate  int64
	CompleteCount             int64
	CompleteRate              int64
	TotalCount                int64
	FollowCount               int64
}

type SpaceWorkObjectCountInfo struct {
	CompletedCount          int64 // 已完成任务数
	ProcessingCount         int64 // 待办任务数
	WeekProcessingCount     int64 // 本周待办任务数
	ExpiredCount            int64 // 过期任务数
	CompleteRate            int64 // 完成率
	TotalCount              int64 // 总任务数
	ClosedOrTerminatedCount int64 // 关闭/终止

	PriorityInfo map[string]int64 //待办任务优先级统计
}

type SpaceWorkVersionCountInfo struct {
	CompletedCount          int64 // 已完成任务数
	ProcessingCount         int64 // 待办任务数
	WeekProcessingCount     int64 // 本周待办任务数
	ExpiredCount            int64 // 过期任务数
	CompleteRate            int64 // 完成率
	TotalCount              int64 // 总任务数
	ClosedOrTerminatedCount int64 // 关闭/终止

	PriorityInfo map[string]int64 //待办任务优先级统计
}

type SpaceUserCountInfo struct {
	CompletedCount          int64 // 已完成任务数
	ProcessingCount         int64 // 待办任务数
	WeekProcessingCount     int64 // 本周待办任务数
	ExpiredCount            int64 // 过期任务数
	CompleteRate            int64 // 完成率
	TotalCount              int64 // 总任务数
	ClosedOrTerminatedCount int64 // 关闭/终止

	PriorityInfo map[string]int64 //待办任务优先级统计
}

type SpaceStatus struct {
	EndStatusList               []int64
	EndAndNotCompleteStatusList []int64
	CompleteStatusList          []int64
}

type CountInfo struct {
	Value int64
}
