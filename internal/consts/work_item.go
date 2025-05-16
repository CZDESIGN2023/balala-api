package consts

type WorkItemStatusType int32
type WorkItemStatusDefaultVal string
type WorkItemStatusDefaultKey string
type WorkItemTypeKey string

type WorkItemStatusDefault struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

const (

	//工作项类型美剧，系统预设
	WorkItemTypeKey_Task      WorkItemTypeKey = "task"       //工作项类型-任务
	WorkItemTypeKey_SubTask   WorkItemTypeKey = "sub_task"   //工作项类型-子任务
	WorkItemTypeKey_StateTask WorkItemTypeKey = "state_task" //工作项类型-状态任务

	//任务状态类型
	WorkItemStatusType_Init     WorkItemStatusType = 1 //起始
	WorkItemStatusType_Process  WorkItemStatusType = 2 //过程
	WorkItemStatusType_Archived WorkItemStatusType = 3 //归档

	//工作任务状态枚举，系统级别
	WorkItemStatus_TerminatedKey WorkItemStatusDefaultKey = "terminated" //终止
	WorkItemStatus_TerminatedVal WorkItemStatusDefaultVal = "3"          //终止

	//节点模式工作流 扩展 任务状态默认值, 仅提供给初始化项目/流程模版等使用
	WorkItemStatus_WorkFlowProgressingDefaultKey WorkItemStatusDefaultKey = "progressing" //进行中
	WorkItemStatus_WorkFlowProgressingDefaultVal WorkItemStatusDefaultVal = "1"           //进行中

	WorkItemStatus_WorkFlowCompletedDefaultKey WorkItemStatusDefaultKey = "completed" //完成
	WorkItemStatus_WorkFlowCompletedDefaultVal WorkItemStatusDefaultVal = "2"         //完成

	WorkItemStatus_WorkFlowCloseDefaultKey WorkItemStatusDefaultKey = "close" //已关闭
	WorkItemStatus_WorkFlowCloseDefaultVal WorkItemStatusDefaultVal = "8"     //已关闭

	WorkItemStatus_WorkflowTestingDefaultKey WorkItemStatusDefaultKey = "testing" //测试中
	WorkItemStatus_WorkflowTestingDefaultVal WorkItemStatusDefaultVal = "4"       //测试中

	WorkItemStatus_WorkflowCheckingDefaultKey WorkItemStatusDefaultKey = "checking" //验收中
	WorkItemStatus_WorkflowCheckingDefaultVal WorkItemStatusDefaultVal = "5"        //验收中

	WorkItemStatus_WorkflowWaitConfirmDefaultKey WorkItemStatusDefaultKey = "wait_confirm" //待确认
	WorkItemStatus_WorkflowWaitConfirmDefaultVal WorkItemStatusDefaultVal = "6"            //待确认

	WorkItemStatus_WorkflowPlanningDefaultKey WorkItemStatusDefaultKey = "planning" //策划中
	WorkItemStatus_WorkflowPlanningDefaultVal WorkItemStatusDefaultVal = "9"        //策划中

	WorkItemStatus_WorkflowDesigningDefaultKey WorkItemStatusDefaultKey = "designing" //设计中
	WorkItemStatus_WorkflowDesigningDefaultVal WorkItemStatusDefaultVal = "10"        //设计中

	WorkItemStatus_WorkflowEvaluatingDefaultKey WorkItemStatusDefaultKey = "evaluating" //评审中
	WorkItemStatus_WorkflowEvaluatingDefaultVal WorkItemStatusDefaultVal = "11"         //评审中

	//状态模式工作流-扩展任务状态, 仅提供给初始化项目/流程模版等使用
	WorkItemStatus_StateflowPendingDefaultKey WorkItemStatusDefaultKey = "st_pending" //待处理
	WorkItemStatus_StateflowPendingDefaultVal WorkItemStatusDefaultVal = "st_pending" //待处理

	WorkItemStatus_StateflowFixingDefaultKey WorkItemStatusDefaultKey = "st_fixing" //修复中
	WorkItemStatus_StateflowFixingDefaultVal WorkItemStatusDefaultVal = "st_fixing" //修复中

	WorkItemStatus_StateflowPendingVerificationDefaultKey WorkItemStatusDefaultKey = "st_pending_verification" //待验证
	WorkItemStatus_StateflowPendingVerificationDefaultVal WorkItemStatusDefaultVal = "st_pending_verification" //待验证

	WorkItemStatus_StateflowRestartDefaultKey WorkItemStatusDefaultKey = "st_restart" //重新打开
	WorkItemStatus_StateflowRestartDefaultVal WorkItemStatusDefaultVal = "st_restart" //重新打开

	WorkItemStatus_StateflowConvertToStoryDefaultKey WorkItemStatusDefaultKey = "st_convert_to_story" //转成故事
	WorkItemStatus_StateflowConvertToStoryDefaultVal WorkItemStatusDefaultVal = "st_convert_to_story" //转成故事

	WorkItemStatus_StateflowDoNotProcessDefaultKey WorkItemStatusDefaultKey = "st_do_not_process" //不处理
	WorkItemStatus_StateflowDoNotProcessDefaultVal WorkItemStatusDefaultVal = "st_do_not_process" //不处理

	WorkItemStatus_StateflowProgressingDefaultKey WorkItemStatusDefaultKey = "st_progressing" //进行中
	WorkItemStatus_StateflowProgressingDefaultVal WorkItemStatusDefaultVal = "st_progressing" //进行中

	WorkItemStatus_StateflowClosedDefaultKey WorkItemStatusDefaultKey = "st_closed" //已关闭
	WorkItemStatus_StateflowClosedDefaultVal WorkItemStatusDefaultVal = "st_closed" //已关闭

)

var sysStatusKeyList = []WorkItemStatusDefaultKey{
	WorkItemStatus_WorkFlowProgressingDefaultKey, // 顺序影响状态排序
	WorkItemStatus_WorkFlowCompletedDefaultKey,
	WorkItemStatus_WorkFlowCloseDefaultKey,
	WorkItemStatus_TerminatedKey,
}

func SysStatusKeyList() []WorkItemStatusDefaultKey {
	return sysStatusKeyList
}
