package service

import (
	"context"
	"go-cs/api/comm"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/space_work_object"
	"go-cs/internal/domain/space_work_object/repo"
	"go-cs/internal/pkg/biz_id"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/errs"
)

type SpaceWorkObjectService struct {
	repo      repo.SpaceWorkObjectRepo
	idService *biz_id.BusinessIdService
}

func NewSpaceWorkObjectService(
	repo repo.SpaceWorkObjectRepo,
	idService *biz_id.BusinessIdService,
) *SpaceWorkObjectService {
	return &SpaceWorkObjectService{
		repo:      repo,
		idService: idService,
	}
}

func (s *SpaceWorkObjectService) CreateSpaceDefaultWorkObject(ctx context.Context, spaceId int64, oper shared.Oper) (*domain.SpaceWorkObject, error) {
	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_SpaceWorkObject)
	if bizId == nil {
		return nil, errs.Business(ctx, "生成模块ID失败")
	}

	ins := domain.NewSpaceWorkObject(bizId.Id, spaceId, consts.DefaultWorkObjectName, 0, oper.GetId(), oper)
	return ins, nil
}

func (s *SpaceWorkObjectService) CreateSpaceWorkObject(ctx context.Context, spaceId int64, name string, ranking int64, uid int64, oper shared.Oper) (*domain.SpaceWorkObject, error) {

	// 检查空间名称是否存在
	exist, err := s.repo.CheckSpaceWorkObjectName(ctx, spaceId, name)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_DB_QUERY_FAIL)
	}

	if exist {
		return nil, errs.Business(ctx, "模块名称重复")
	}

	if ranking == 0 {
		ranking, err = s.repo.GetMaxRanking(ctx, spaceId)
		if err != nil {
			return nil, errs.New(ctx, comm.ErrorCode_DB_QUERY_FAIL)
		}
		ranking = ranking + 100
	}

	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_SpaceWorkObject)
	if bizId == nil {
		return nil, errs.Business(ctx, "生成模块ID失败")
	}

	return domain.NewSpaceWorkObject(bizId.Id, spaceId, name, ranking, uid, oper), nil
}

func (s *SpaceWorkObjectService) UpdateSpaceWorkObjectName(ctx context.Context, spaceWorkObject *domain.SpaceWorkObject, newName string, oper shared.Oper) error {
	isExist, err := s.repo.CheckSpaceWorkObjectName(ctx, spaceWorkObject.SpaceId, newName)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if isExist {
		return errs.Business(ctx, "模块名称重复")
	}

	spaceWorkObject.UpdateName(newName, oper)
	return nil
}

func (s *SpaceWorkObjectService) CheckTransfer(ctx context.Context, fromWorkObj *domain.SpaceWorkObject, toWorkObjId int64, oper shared.Oper) error {
	//如果项目只有最后一个模块，就不让删除了
	wCount, err := s.repo.GetSpaceWorkObjectCount(ctx, fromWorkObj.SpaceId)
	if err != nil {
		return err
	}

	if wCount <= 1 {
		return errs.Business(ctx, "项目下至少需要一个模块")
	}

	if toWorkObjId != 0 {
		toWorkObj, err := s.repo.GetSpaceWorkObject(ctx, fromWorkObj.SpaceId, toWorkObjId)
		if err != nil {
			return err
		}

		if toWorkObj.Id == fromWorkObj.Id {
			return errs.Business(ctx, "不能移动到自身")
		}
	}

	return nil
}
