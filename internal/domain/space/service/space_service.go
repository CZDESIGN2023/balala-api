package service

import (
	"context"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/space"
	"go-cs/internal/domain/space/repo"
	"go-cs/internal/pkg/biz_id"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils/errs"
)

type SpaceService struct {
	repo repo.SpaceRepo

	idService *biz_id.BusinessIdService
}

func NewSpaceService(
	repo repo.SpaceRepo,
	idService *biz_id.BusinessIdService,
) *SpaceService {
	return &SpaceService{
		repo:      repo,
		idService: idService,
	}
}

func (s *SpaceService) CreateSpace(ctx context.Context, uid int64, name string, describe string, oper shared.Oper) (*domain.Space, error) {

	//idService 不是业务领域服务，只是一个工具，用于生成业务ID
	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_Space)
	if bizId == nil {
		return nil, errs.Business(ctx, "空间ID分配失败")
	}

	space := domain.NewSpace(bizId.Id, uid, name, describe, oper)
	return space, nil
}

func (s *SpaceService) UpdateSpaceInfo(ctx context.Context, space *domain.Space, newName string, newDescribe string, oper shared.Oper) error {

	space.UpdateDescribe(newDescribe, oper)
	space.UpdateName(newName, oper)
	return nil
}

func (s *SpaceService) UpdateSpaceName(ctx context.Context, space *domain.Space, newName string, oper shared.Oper) error {
	space.UpdateName(newName, oper)
	return nil
}
