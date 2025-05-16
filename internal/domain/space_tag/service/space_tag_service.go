package service

import (
	"context"
	"go-cs/api/comm"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/space_tag"
	"go-cs/internal/domain/space_tag/repo"
	"go-cs/internal/pkg/biz_id"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/errs"
)

type SpaceTagService struct {
	repo      repo.SpaceTagRepo
	idService *biz_id.BusinessIdService
}

func NewSpaceTagService(
	repo repo.SpaceTagRepo,
	idService *biz_id.BusinessIdService,
) *SpaceTagService {
	return &SpaceTagService{
		repo:      repo,
		idService: idService,
	}
}

func (s *SpaceTagService) CreateSpaceTag(ctx context.Context, spaceId int64, name string, oper shared.Oper) (*domain.SpaceTag, error) {

	isExist, err := s.repo.CheckTagNameIsExist(ctx, spaceId, name)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_DB_QUERY_FAIL)
	}

	if isExist {
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_TAG_NAME_IS_EXIST)
	}

	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_SpaceTag)
	if bizId == nil {
		return nil, errs.Business(ctx, "生成标签ID失败")
	}

	return domain.NewSpaceTag(bizId.Id, spaceId, name, oper), nil
}

func (s *SpaceTagService) FilterExistSpaceTagIds(ctx context.Context, spaceId int64, tagIds []int64) ([]int64, error) {
	// 查询标签是否存在
	existTagIds, err := s.repo.FilterExistSpaceTagIds(ctx, spaceId, tagIds)
	if err != nil {
		return nil, err
	}
	return existTagIds, nil
}

func (s *SpaceTagService) UpdateTagName(ctx context.Context, tag *domain.SpaceTag, newName string, oper shared.Oper) error {
	isExist, err := s.repo.CheckTagNameIsExist(ctx, tag.SpaceId, newName)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if isExist {
		//errInfo := errs.New(ctx, comm.ErrorCode_SPACE_TAG_NAME_IS_EXIST)
		//return nil, errInfo
		return errs.Business(ctx, "标签名重复")
	}

	tag.ChangeName(newName, oper)
	return nil
}
