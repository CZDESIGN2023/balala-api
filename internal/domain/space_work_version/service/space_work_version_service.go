package service

import (
	"context"
	"go-cs/api/comm"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/space_work_version"
	"go-cs/internal/domain/space_work_version/repo"
	"go-cs/internal/pkg/biz_id"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/rand"
	"time"
)

type SpaceWorkVersionService struct {
	repo      repo.SpaceWorkVersionRepo
	idService *biz_id.BusinessIdService
}

func NewSpaceWorkVersionService(
	repo repo.SpaceWorkVersionRepo,
	idService *biz_id.BusinessIdService,
) *SpaceWorkVersionService {
	// 初始化逻辑
	return &SpaceWorkVersionService{
		repo:      repo,
		idService: idService,
	}
}

func (s *SpaceWorkVersionService) CreateSpaceDefaultWorkVersion(ctx context.Context, spaceId int64, oper shared.Oper) (*domain.SpaceWorkVersion, error) {

	isExist, err := s.repo.CheckSpaceWorkVersionName(ctx, spaceId, consts.DefaultWorkItemVersionName)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_DB_QUERY_FAIL)
	}

	if isExist {
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_TAG_NAME_IS_EXIST)
	}

	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_SpaceWorkVersion)
	if bizId == nil {
		return nil, errs.Business(ctx, "生成版本ID失败")
	}

	ins := domain.NewSpaceWorkVersion(bizId.Id, spaceId, consts.DefaultWorkItemVersionKey, consts.DefaultWorkItemVersionName, 100, nil)
	return ins, nil
}

func (s *SpaceWorkVersionService) CreateSpaceWorkVersion(ctx context.Context, spaceId int64, name string, ranking int64, oper shared.Oper) (*domain.SpaceWorkVersion, error) {

	isExist, err := s.repo.CheckSpaceWorkVersionName(ctx, spaceId, name)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_DB_QUERY_FAIL)
	}

	if isExist {
		return nil, errs.Business(ctx, "版本名称重复")
	}

	if ranking == 0 {
		ranking, _ = s.repo.GetMaxRanking(ctx, spaceId)
		if ranking == 0 {
			ranking = time.Now().Unix()
		}
	}

	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_SpaceWorkVersion)
	if bizId == nil {
		return nil, errs.Business(ctx, "生成版本ID失败")
	}

	versionKey := "spVersion_" + rand.S(5)
	ins := domain.NewSpaceWorkVersion(bizId.Id, spaceId, versionKey, name, ranking+100, oper)
	return ins, nil
}

func (s *SpaceWorkVersionService) UpdateSpaceWorkVersionName(ctx context.Context, workVersion *domain.SpaceWorkVersion, newName string, oper shared.Oper) error {
	isExist, err := s.repo.CheckSpaceWorkVersionName(ctx, workVersion.SpaceId, newName)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if isExist {
		return errs.Business(ctx, "版本名称重复")
	}

	workVersion.UpdateName(newName, oper)
	return nil
}
