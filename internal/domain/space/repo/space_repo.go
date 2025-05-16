package repo

import (
	"context"
	domain "go-cs/internal/domain/space"
)

type SpaceRepo interface {
	SpaceRepoQuery

	CreateSpace(ctx context.Context, space *domain.Space) error

	GetSpace(ctx context.Context, id int64) (*domain.Space, error)
	GetSpaceDetail(ctx context.Context, spaceId int64) (*domain.Space, error)

	IsExistBySpaceName(context.Context, int64, string) (bool, error)
	GetSpaceByCreator(ctx context.Context, userId int64, spaceId int64) (*domain.Space, error)
	GetUserSpaceIds(ctx context.Context, userId int64) ([]int64, error)
	GetSpaceByIds(ctx context.Context, ids []int64) ([]*domain.Space, error)
	SpaceMap(ctx context.Context, ids []int64) (map[int64]*domain.Space, error)

	SaveSpace(ctx context.Context, space *domain.Space) error

	DelSpace(ctx context.Context, spaceId int64) error

	CreateSpaceConfig(ctx context.Context, space *domain.SpaceConfig) error
	SaveSpaceConfig(ctx context.Context, space *domain.SpaceConfig) error
	GetSpaceConfig(ctx context.Context, spaceId int64) (*domain.SpaceConfig, error)
	DelSpaceConfig(ctx context.Context, spaceId int64) error
	SpaceConfigMap(ctx context.Context, spaceIds []int64) (map[int64]*domain.SpaceConfig, error)

	SetTempConfig(ctx context.Context, userId int64, confMap map[string]string) error
	GetTempConfig(ctx context.Context, userId int64, keys ...string) map[string]string
	DelTempConfig(ctx context.Context, userId int64, keys ...string) error

	GetAllSpaceIds() ([]int64, error)
	GetAllSpace(ctx context.Context) ([]*domain.Space, error)
}

type SpaceRepoQuery interface {
	GetUserSpaceList(ctx context.Context, userId int64) ([]*domain.Space, error)
}
