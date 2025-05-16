package service

import (
	"context"
	domain "go-cs/internal/domain/space_file_info"
	"go-cs/internal/domain/space_file_info/repo"
	"go-cs/internal/pkg/biz_id"
)

type SpaceFileInfoService struct {
	repo      repo.SpaceFileInfoRepo
	idService *biz_id.BusinessIdService
}

func NewSpaceFileInfoService(
	repo repo.SpaceFileInfoRepo,
	idService *biz_id.BusinessIdService,
) *SpaceFileInfoService {
	return &SpaceFileInfoService{
		repo:      repo,
		idService: idService,
	}
}

func (s *SpaceFileInfoService) DeleteSpaceFileInfoByWorkItemIds(ctx context.Context, spaceId int64, workItemIds []int64) (domain.SpaceFileInfos, error) {
	spaceFiles, err := s.repo.GetSpaceWorkItemFileInfoByWorkItemIds(ctx, spaceId, workItemIds)
	if err != nil {
		return nil, err
	}

	for _, v := range spaceFiles {
		v.OnDelete()
	}

	return spaceFiles, nil
}
