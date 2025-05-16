package work_item

import shared "go-cs/internal/pkg/domain"

const (
	Diff_Name             shared.PropDiff = "workItemName"
	Diff_Priority         shared.PropDiff = "priority"
	Diff_Remark           shared.PropDiff = "remark"
	Diff_PlanTime         shared.PropDiff = "planTime"
	Diff_Describe         shared.PropDiff = "describe"
	Diff_IconFlags        shared.PropDiff = "iconFlags"
	Diff_Status           shared.PropDiff = "status"
	Diff_Tags             shared.PropDiff = "tags"
	Diff_Directors        shared.PropDiff = "directors"
	Diff_Followers        shared.PropDiff = "followers"
	Diff_Participators    shared.PropDiff = "participators"
	Diff_NodeDirectors    shared.PropDiff = "node_directors"
	Diff_Restart          shared.PropDiff = "restart"
	Diff_Resume           shared.PropDiff = "resume"
	Diff_CommentNum       shared.PropDiff = "commentNum"
	Diff_WorkItemType     shared.PropDiff = "workItemType"
	Diff_ProcessRate      shared.PropDiff = "processRate"
	Diff_ChildNum         shared.PropDiff = "childNum"
	Diff_VersionId        shared.PropDiff = "versionId"
	Diff_ObjectId         shared.PropDiff = "objectId"
	Diff_WorkFlowTemplate shared.PropDiff = "workFlowTemplate"
	Diff_Reason           shared.PropDiff = "reason"
	Diff_CountAt          shared.PropDiff = "countAt"

	Diff_WorkItemFlowNode_Status    shared.PropDiff = "workItemFlowNode.Status"
	Diff_WorkItemFlowNode_Directors shared.PropDiff = "workItemFlowNode.Directors"
	Diff_WorkItemFlowNode_PlanTime  shared.PropDiff = "workItemFlowNode.PlanTime"
	Diff_WorkItemFlowNode_UpdatedAt shared.PropDiff = "workItemFlowNode.UpdatedAt"

	Diff_WorkItemFlowRole_Directors shared.PropDiff = "workItemFlowRole.Directors"
)
