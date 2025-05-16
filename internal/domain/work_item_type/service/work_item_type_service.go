package service

import (
	"context"
	domain "go-cs/internal/domain/work_item_type"
	"go-cs/internal/domain/work_item_type/repo"
	"go-cs/internal/pkg/biz_id"

	"go-cs/internal/consts"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/rand"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type WorkItemTypeService struct {
	log       *log.Helper
	idService *biz_id.BusinessIdService
	repo      repo.WorkItemTypeRepo
}

func NewWorkItemTypeService(
	idService *biz_id.BusinessIdService,
	repo repo.WorkItemTypeRepo,
	logger log.Logger,
) *WorkItemTypeService {

	moduleName := "WorkItemTypeService"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &WorkItemTypeService{
		repo:      repo,
		log:       hlog,
		idService: idService,
	}
}

func (s *WorkItemTypeService) CreateSpaceDefaultWorkItemTypes(ctx context.Context, spaceId int64, creatorUid int64) (map[string]*domain.WorkItemType, error) {

	task, err := s.newWorkItemType(ctx, domain.SpaceId(spaceId), "任务", string(consts.WorkItemTypeKey_Task), consts.FlowMode_WorkFlow, 1, creatorUid)
	if err != nil {
		return nil, errs.Business(ctx, "创建系统默认工作项类型失败")
	}

	subTask, err := s.newWorkItemType(ctx, domain.SpaceId(spaceId), "子任务", string(consts.WorkItemTypeKey_SubTask), consts.FlowMode_StateFlow, 1, creatorUid)
	if err != nil {
		return nil, errs.Business(ctx, "创建系统默认工作项类型失败")
	}

	stateTask, err := s.newWorkItemType(ctx, domain.SpaceId(spaceId), "状态任务", string(consts.WorkItemTypeKey_StateTask), consts.FlowMode_StateFlow, 1, creatorUid)
	if err != nil {
		return nil, errs.Business(ctx, "创建系统默认工作项类型失败")
	}

	return map[string]*domain.WorkItemType{
		string(consts.WorkItemTypeKey_Task):      task,
		string(consts.WorkItemTypeKey_SubTask):   subTask,
		string(consts.WorkItemTypeKey_StateTask): stateTask,
	}, nil
}

func (s *WorkItemTypeService) CreateWorkFlowTaskWorkType(ctx context.Context, spaceId int64, uid int64) (*domain.WorkItemType, error) {
	task, err := s.newWorkItemType(ctx, domain.SpaceId(spaceId), "任务", string(consts.WorkItemTypeKey_Task), consts.FlowMode_WorkFlow, 1, uid)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *WorkItemTypeService) CreateStateFlowTaskWorkType(ctx context.Context, spaceId int64, uid int64) (*domain.WorkItemType, error) {
	stateTask, err := s.newWorkItemType(ctx, domain.SpaceId(spaceId), "状态任务", string(consts.WorkItemTypeKey_StateTask), consts.FlowMode_StateFlow, 1, uid)
	if err != nil {
		return nil, err
	}

	return stateTask, nil
}

func (s *WorkItemTypeService) CreateSubTaskWorkType(ctx context.Context, spaceId int64, uid int64) (*domain.WorkItemType, error) {
	task, err := s.newWorkItemType(ctx, domain.SpaceId(spaceId), "子任务", string(consts.WorkItemTypeKey_SubTask), consts.FlowMode_StateFlow, 1, uid)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *WorkItemTypeService) newWorkItemType(ctx context.Context, spaceId domain.SpaceId, name string, key string, flowMode consts.WorkFlowMode, isSys int32, uid int64) (*domain.WorkItemType, error) {

	if name == "" || key == "" {
		return nil, errs.Business(ctx, "名称和键不能为空")
	}

	isExistName, err := s.repo.IsExistByName(ctx, cast.ToInt64(spaceId), name)
	if err != nil {
		return nil, err
	}
	if isExistName {
		return nil, errs.Business(ctx, "名称已存在")
	}

	isExistKey, err := s.repo.IsExistByKey(ctx, cast.ToInt64(spaceId), key)
	if err != nil {
		return nil, err
	}
	if isExistKey {
		return nil, errs.Business(ctx, "键已存在")
	}

	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_WorkItemType)
	if bizId == nil {
		return nil, errs.Business(ctx, "分配Id失败")
	}

	if key == "" {
		key = "key_" + rand.S(5)
	}

	itemType := domain.NewWorkItemType(bizId.Id, spaceId, name, key, flowMode, isSys, uid)
	return itemType, nil
}
