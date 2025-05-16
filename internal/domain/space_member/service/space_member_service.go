package service

import (
	"context"
	domain "go-cs/internal/domain/space_member"
	repo "go-cs/internal/domain/space_member/repo"
	"go-cs/internal/pkg/biz_id"
	shared "go-cs/internal/pkg/domain"

	"go-cs/internal/utils/errs"
)

type SpaceMemberService struct {
	repo      repo.SpaceMemberRepo
	idService *biz_id.BusinessIdService
}

func NewSpaceMemberService(
	repo repo.SpaceMemberRepo,
	idService *biz_id.BusinessIdService,
) *SpaceMemberService {

	return &SpaceMemberService{
		repo:      repo,
		idService: idService,
	}
}

func (s *SpaceMemberService) NewSpaceMember(ctx context.Context, spaceId int64, userId int64, roleId int64, oper shared.Oper) (*domain.SpaceMember, error) {

	//检查是否重复添加了
	exist, err := s.repo.IsExistSpaceMember(ctx, spaceId, userId)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errs.Business(ctx, "成员已存在")
	}

	member := domain.NewSpaceMember(spaceId, userId, roleId, oper)

	return member, nil
}

func (s *SpaceMemberService) CheckAllIsMember(ctx context.Context, spaceId int64, userIds []int64) (bool, error) {

	isAll, err := s.repo.AllIsMember(ctx, spaceId, userIds...)
	if err != nil {
		return false, err
	}

	if !isAll {
		return false, nil
	}

	return true, nil
}
