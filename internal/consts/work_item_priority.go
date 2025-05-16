package consts

import "slices"

const (
	//工作任务优先级

	WorkItemPriority_P0      = "P0"      //P0
	WorkItemPriority_P1      = "P1"      //P1
	WorkItemPriority_P2      = "P2"      //P2
	WorkItemPriority_P3      = "P3"      //P3
	WorkItemPriority_P4      = "P4"      //P4
	WorkItemPriority_Pending = "PENDING" //PENDING 暂停
	WorkItemPriority_Suspend = "SUSPEND" //SUSPEND 中断
)

var allPriority = []string{
	WorkItemPriority_P0,
	WorkItemPriority_P1,
	WorkItemPriority_P2,
	WorkItemPriority_P3,
	WorkItemPriority_P4,
	WorkItemPriority_Pending,
	WorkItemPriority_Suspend,
}

func GetAllPriority() []string {
	return slices.Clone(allPriority)
}

func IsOneOfWorkItemPriority(in string) bool {
	return slices.Contains(allPriority, in)
}

func GetWorkItemPriorityPower(in string) int32 {
	return int32(GetPriorityRankingMap()[in])
}

func GetWorkItemPriorityName(in string) string {
	switch in {
	case WorkItemPriority_P0:
		return "P0"
	case WorkItemPriority_P1:
		return "P1"
	case WorkItemPriority_P2:
		return "P2"
	case WorkItemPriority_P3:
		return "P3"
	case WorkItemPriority_P4:
		return "P4"
	case WorkItemPriority_Pending:
		return "待定"
	case WorkItemPriority_Suspend:
		return "暂定"

	}
	return ""
}

var priorityRankingMap = map[string]int64{
	WorkItemPriority_P0:      10,
	WorkItemPriority_P1:      9,
	WorkItemPriority_P2:      8,
	WorkItemPriority_P3:      7,
	WorkItemPriority_P4:      6,
	WorkItemPriority_Pending: -1,
	WorkItemPriority_Suspend: -2,
}

func GetPriorityRankingMap() map[string]int64 {
	return priorityRankingMap
}
