package repo

import (
	"context"
	domain "go-cs/internal/domain/work_item"
)

type WorkItemEsRepo interface {
	CreateWorkItemEs(ctx context.Context, workItem *domain.WorkItem) error
}

type WorkItemRepo interface {
	// WorkItemQueryRepo
	WorkItemFlowNodeRepo
	WorkItemFlowRoleRepo

	CreateWorkItem(ctx context.Context, workItem *domain.WorkItem) error
	SaveWorkItem(ctx context.Context, workItem *domain.WorkItem) error

	GetWorkItem(context.Context, int64, *WithDocOption, *WithOption) (*domain.WorkItem, error)
	GetWorkItemByIds(context.Context, []int64, *WithDocOption, *WithOption) (domain.WorkItems, error)

	GetWorkItemByPid(context.Context, int64, *WithDocOption, *WithOption) (domain.WorkItems, error)
	GetWorkItemByPids(context.Context, []int64, *WithDocOption, *WithOption) (domain.WorkItems, error)

	CreateWorkItemFiles(ctx context.Context, workItemFiles domain.WorkItemFiles) error

	DelSpaceWorkItem(ctx context.Context, workItemIds ...int64) (int64, error)
	DelSpaceWorkItemByWorkObjectId(ctx context.Context, workObjectId int64) (int64, error)

	DelSpaceWorkItemBySpaceId(ctx context.Context, spaceId int64) (int64, error)

	// GetSpaceWorkItemDetail(ctx context.Context, workItemId int64) (*domain.WorkItem, error)
	// GetSpaceWorkItemCountByPid(ctx context.Context, workItemPid int64) (int64, error)

	// // 字段更新
	UpdateSpaceAllWorkItemCreator(ctx context.Context, spaceId, oldUserId, newUserId int64) (int64, error)
	// UpdateWorkItemsFields(ctx context.Context, workItemIds []int64, kvs map[string]any) error
	// UpdateWorkItemFields(ctx context.Context, workItemId int64, fields map[string]interface{}) error
	// UpdateWorkItemFieldsByPid(ctx context.Context, workItemPid int64, kvs map[string]any) error

	// // 标签操作
	// WorkItemTagAdd(ctx context.Context, workItemId int64, tagIds []int64) error
	// WorkItemTagRemove(ctx context.Context, workItemId int64, tagIds []int64) error
	// WorkItemTagSet(ctx context.Context, workItemId int64, tagIds []int64) error
	RemoveTagFromAllWorkItem(ctx context.Context, spaceId, tagId int64) (int64, error)

	// // 当前负责人
	// ClearDirectorForWorkItem(ctx context.Context, workItemId int64) (int64, error)
	// // 任务查询

	// WorkItemMapByPid(ctx context.Context, ids []int64, fields ...string) (map[int64]*domain.WorkItem, error)

	// 获取任务id列表
	GetSpaceWorkItemIdsForUpgradeFlow(ctx context.Context, spaceId int64, flowId int64, lastVersion int64) ([]int64, error)
	GetSpaceWorkItemIdsByParticipators(ctx context.Context, spaceId, userId int64) ([]int64, error)
	GetSpaceWorkItemIdsByWorkObject(ctx context.Context, workObjectId int64) ([]int64, error)
	GetSpaceAllWorkItemIdsHasComment(ctx context.Context, spaceId int64) ([]int64, error)
	GetWorkItemIdsByFollower(ctx context.Context, userId int64, spaceId int64) ([]int64, error)
	GetSpaceWorkItemIdsByPid(ctx context.Context, pid int64) ([]int64, error)

	//替换负责人
	ReplaceDirectorForWorkItemBySpace(ctx context.Context, spaceId int64, userId int64, newUserId int64) (int64, error)
	ReplaceParticipatorsForWorkItemBySpace(ctx context.Context, spaceId int64, userId int64, newUserId int64) (int64, error)
	ReplaceNodeDirectorsForWorkItemBySpace(ctx context.Context, spaceId int64, userId int64, newUserId int64) (int64, error)
	ReplaceDirectorForWorkItemFlowRolesBySpace(ctx context.Context, spaceId int64, userId int64, newUserId int64) (int64, error)
	ReplaceDirectorForWorkItemFlowBySpace(ctx context.Context, spaceId int64, userId int64, newUserId int64) (int64, error)

	// GetSpaceAllWorkItemIds(ctx context.Context, spaceId int64) ([]int64, error)
	// MoveToNewWorkObject(ctx context.Context, workItemId int64, objectId int64) error

	// // 修改任务状态
	// ChangeSpaceWorkItemStatus(ctx context.Context, workItemId int64, newStatus *vo.WorkItemStatusVo) (int64, error)
	// ChangeSpaceWorkItemStatusByPid(ctx context.Context, workItemPid int64, newStatus *vo.WorkItemStatusVo) error
	// ChangeWorkItemStatusByPidIfNotEqual(ctx context.Context, workItemPid int64, newStatus *vo.WorkItemStatusVo, oldStatus *vo.WorkItemStatusVo) error

	// // 设置重启任务标志
	// SetRestartWorkItemFlag(ctx context.Context, workItemId, restartUserId int64) error
	// SetRestartWorkItemFlagByPid(ctx context.Context, workItemPid, restartUserId int64) error

	// // 恢复任务
	// ResumeWorkItem(ctx context.Context, workItemId int64) error
	// ResumeWorkItemByPid(ctx context.Context, workItemPid int64, newStatus *vo.WorkItemStatusVo) error

	// // 图标icon
	// AddWorkItemIconFlags(ctx context.Context, workItemId int64, flag uint32) error
	// RemoveWorkItemIconFlags(ctx context.Context, workItemId int64, flag uint32) error

	// // 统计用户关联的任务
	CountUserRelatedSpaceWorkItem(ctx context.Context, spaceId, userId int64) (int64, error)
	CountUserRelatedSpaceWorkItemBySpaceIds(ctx context.Context, userId int64, spaceIds []int64) (map[int64]int64, error)
	// 获取工作流相关的数量
	CountWorkFlowRelatedSpaceWorkItem(ctx context.Context, spaceId, workFlowId int64, excludeStatusKeys []string) (int64, error)
	// 获取工作状态相关的数量
	CountWorkItemStatusRelatedSpaceWorkItem(ctx context.Context, spaceId, workItemStatusId int64) (int64, error)

	// 增加评论数
	IncrCommentNum(ctx context.Context, workItemId int64, num int) (int64, error)
	AllCommentNumMap(ctx context.Context, spaceId int64) (map[int64]int64, error)

	// 重置子任务统计数量
	ResetChildTaskNum(ctx context.Context, workItemId int64) error

	// //工作流信息更新
	// UpdateWorkItemTypeInfo(ctx context.Context, in vo.UpdateWorkItemTypeInfoVo) error

	// 取消关注
	Unfollow(ctx context.Context, userId int64, workItemIds []int64) error

	// //设置工作项版本至默认版本
	ResetVersion(ctx context.Context, oldVersionId int64, newVersionId int64) error
	// SetWorkVersion(ctx context.Context, workItemId int64, workVersionId int64) error
	// SetWorkVersionByPid(ctx context.Context, workItemPid int64, workVersionId int64) error

	// // 同步计算字段: 当前负责人，参与人
	// SyncComputeField(ctx context.Context, workItemIds ...int64)

	UpdateWorkItemWorkObjectIdByIds(ctx context.Context, workItemIds []int64, newWorkObjectId int64) error
	UpdateNewWorkObject(ctx context.Context, workItemId int64, objectId int64) error

	FilterChildWorkItemIds(ctx context.Context, workItemIds []int64) ([]int64, error)

	HasWorkItemRelateFlow(ctx context.Context, spaceId, workFlowId int64) (bool, error)
}

type WorkItemQueryRepo interface {
	// GetSpaceWorkItemByDirector(ctx context.Context, spaceId, userId int64) ([]*search2.Model, error)
	// GetSpaceWorkItemList(ctx context.Context, workItemIds []int64) ([]*search2.Model, error)
	// GetSpaceWorkItemListByPid(ctx context.Context, workItemPid int64) ([]*search2.Model, error)
	// GetWorkItemByTag(ctx context.Context, tagId int64) ([]*search2.Model, error)
	// WorkItemSearchModelMapByPid(ctx context.Context, ids []int64, fields ...string) (map[int64]*search2.Model, error)
	// SpaceWorkItemMap(ctx context.Context, workItemIds []int64) (map[int64]*search2.Model, error)
	// SearchByName(ctx context.Context, spaceId int64, keyword string) ([]*search2.Model, error)
	// GetSpaceWorkItemWithColumns(ctx context.Context, workItemId int64, columns ...string) (*search2.Model, error)
	// GetSpaceWorkItem(ctx context.Context, workItemId int64) (*search2.Model, error)
	// GetSpaceWorkItemFromDB(ctx context.Context, workItemId int64) (*search2.Model, error)
}

type WorkItemFlowNodeRepo interface {
	// GetWorkItemFlowNodesByIds(ctx context.Context, workItemIds []int64) ([]*domain.WorkItemFlowNode, error)

	CreateWorkItemFlowNode(ctx context.Context, workItemFlowNode *domain.WorkItemFlowNode) error
	CreateWorkItemFlowNodes(ctx context.Context, workItemFlowNode ...*domain.WorkItemFlowNode) error
	SaveWorkItemFlowNode(ctx context.Context, workItemFlowNode *domain.WorkItemFlowNode) error

	GetWorkItemFlowNodes(ctx context.Context, workItemId int64) (domain.WorkItemFlowNodes, error)

	// GetWorkItemFlowNode(ctx context.Context, id int64) (*domain.WorkItemFlowNode, error)
	GetWorkItemFlowNodeByNodeCode(ctx context.Context, workItemId int64, workFlowNodeCode string) (*domain.WorkItemFlowNode, error)
	ProcessingNodeMapByWorkItemIds(ctx context.Context, workItemIds []int64) (map[int64][]*domain.WorkItemFlowNode, error)

	DelWorkItemFlowNodeByWorkItemIds(ctx context.Context, workItemIds ...int64) (int64, error)
	DelWorkItemFlowNodeByIds(ctx context.Context, ids ...int64) (int64, error)

	AddDirectorForWorkItemFlows(ctx context.Context, workFlowIds []int64, userId int64) (int64, error)
	RemoveDirectorForWorkItemFlows(ctx context.Context, workFlowIds []int64, userId int64) (int64, error)

	// SetDirectorForWorkItemFlow(ctx context.Context, workFlowId int64, userIds []int64) (int64, error)
	// ResetWorkItemFlowNodeInfoByWorkItemId(ctx context.Context, workItemId int64) error
	// ResetWorkItemFlowNodeInfoByNodeCode(ctx context.Context, workItemId int64, flowNodeCode string) error

	//GetWorkItemProgressingFlowNode(ctx context.Context, workItemId int64) ([]*domain.WorkItemFlowNode, error)
	// WorkItemFlowNodeMap(ctx context.Context, workItemIds []int64) (map[int64][]*domain.WorkItemFlowNode, error)
	// NodeDirectorsMapByWorkItemIds(ctx context.Context, workItemIds []int64) (map[int64][]int64, error)
	// GetWorkItemFlowIdsBySpaceUserId(ctx context.Context, spaceId int64, userId int64) ([]int64, error)
	// UpdateFields(ctx context.Context, workFlowId int64, kvs map[string]any) error

	DelSpaceWorkItemFlowBySpaceId(ctx context.Context, spaceId int64) (int64, error)

	GetSpaceWorkItemFlowIdsBySpaceUserId(ctx context.Context, spaceId int64, userId int64) ([]int64, error)
}

type WorkItemFlowRoleRepo interface {
	CreateWorkItemFlowRole(ctx context.Context, workItemFlowNode *domain.WorkItemFlowRole) error
	CreateWorkItemFlowRoles(ctx context.Context, workItemFlowNode ...*domain.WorkItemFlowRole) error
	SaveWorkItemFlowRole(ctx context.Context, workItemFlowNode *domain.WorkItemFlowRole) error

	GetWorkItemFlowRoles(ctx context.Context, workItemId int64) (domain.WorkItemFlowRoles, error)

	DelWorkItemFlowRoleByWorkItemIds(ctx context.Context, workItemIds ...int64) (int64, error)
	DelWorkItemFlowRoleByIds(ctx context.Context, ids ...int64) (int64, error)

	// 获取角色关联的数量
	CountWorkFlowRoleRelatedSpaceWorkItem(ctx context.Context, spaceId, flowRoleId int64) (int64, error)

	DelSpaceWorkItemFlowRoleBySpaceId(ctx context.Context, spaceId int64) (int64, error)
}
