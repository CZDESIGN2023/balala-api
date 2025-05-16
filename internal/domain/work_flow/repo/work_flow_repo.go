package repo

import (
	"context"
	"go-cs/internal/bean/vo/query"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/work_flow"
)

type WorkFlowRepo interface {
	WorkFlowQueryRepo
	WorkFlowMemoryCacheRepo

	CreateWorkFlow(ctx context.Context, wf *domain.WorkFlow) error
	CreateWorkFlows(ctx context.Context, wf []*domain.WorkFlow) error
	SaveWorkFlow(ctx context.Context, wf *domain.WorkFlow) error

	GetWorkFlowBySpaceWorkItemTypeId(ctx context.Context, spaceId int64, workItemTypeId int64) ([]*domain.WorkFlow, error)
	GetWorkFlow(ctx context.Context, id int64) (*domain.WorkFlow, error)
	GetWorkFlowBySpaceId(ctx context.Context, spaceId int64) ([]*domain.WorkFlow, error)
	GetWorkFlowBySpaceIds(ctx context.Context, spaceIds []int64) ([]*domain.WorkFlow, error)
	GetWorkFlowMap(ctx context.Context, ids []int64) (map[int64]*domain.WorkFlow, error)

	CreateWorkFlowTemplate(ctx context.Context, template *domain.WorkFlowTemplate) error
	CreateWorkFlowTemplates(ctx context.Context, templates []*domain.WorkFlowTemplate) error
	SaveWorkFlowTemplate(ctx context.Context, template *domain.WorkFlowTemplate) error

	GetFlowTemplate(ctx context.Context, id int64) (*domain.WorkFlowTemplate, error)
	GetFlowTemplateByIds(ctx context.Context, ids []int64) ([]*domain.WorkFlowTemplate, error)
	FlowTemplateMap(ctx context.Context, ids []int64) (map[int64]*domain.WorkFlowTemplate, error)

	DelWorkFlow(ctx context.Context, id int64) error
	DelWorkFlowTemplate(ctx context.Context, id int64) error
	DelWorkFlowTemplateByFlowId(ctx context.Context, id int64) error
	DelWorkFlowBySpaceId(ctx context.Context, spaceId int64) error
	DelWorkFlowTemplateBySpaceId(ctx context.Context, spaceId int64) error

	GetAllWorkFlowNameBySpaceId(ctx context.Context, spaceId int64) ([]string, error)

	IsExistByWorkFlowName(context.Context, int64, string) (bool, error)
	WorkFlowMap(ctx context.Context, ids []int64) (map[int64]*domain.WorkFlow, error)

	GetMaxRanking(ctx context.Context, spaceId int64) (int64, error)
	ClearHistoryTemplate(ctx context.Context, flowId int64) error
}

type WorkFlowQueryRepo interface {
	QTaskWorkFlowViewList(ctx context.Context, req *query.TaskWorkFlowListQuery) (*query.TaskWorkFlowListQueryResult, error)
	QTaskWorkFlowList(ctx context.Context, spaceId int64, workItemTypeIds []int64, flowMode consts.WorkFlowMode) ([]*domain.WorkFlow, error)
	QTaskWorkFlowById(ctx context.Context, spaceId int64, workFlowIds []int64) ([]*domain.WorkFlow, error)
	QWorkFlowInfo(ctx context.Context, req *query.WorkFlowInfoQuery) (*query.WorkFlowInfoQueryResult, error)

	SearchTaskWorkFlowLastTemplateByOwnerRule(ctx context.Context, spaceId int64, owner string) ([]int64, error)
	SearchTaskWorkFlowTemplateByNodeStateEvent(ctx context.Context, spaceId int64, subStateId string) ([]int64, error)
	SearchTaskWorkFlowTemplateByOwnerRoleRule(ctx context.Context, spaceId int64, flowRoleId string) ([]int64, error)

	SearchHistoryTaskWorkFlowTemplateByOwnerRule(ctx context.Context, spaceId int64, owner string) ([]*domain.WorkFlowTemplate, error)
}

type WorkFlowMemoryCacheRepo interface {
	GetWorkFlowTemplateFormMemoryCache(ctx context.Context, id int64) (*domain.WorkFlowTemplate, error)
}
