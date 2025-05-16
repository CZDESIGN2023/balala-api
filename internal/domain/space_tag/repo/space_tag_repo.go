package repo

import (
	"context"
	domain "go-cs/internal/domain/space_tag"
)

type SpaceTagRepo interface {
	SpaceTagQueryRepo

	CreateTag(ctx context.Context, spaceTag *domain.SpaceTag) error
	SaveTag(ctx context.Context, spaceTag *domain.SpaceTag) error

	FilterExistSpaceTagIds(ctx context.Context, spaceId int64, tagIds []int64) ([]int64, error)
	CheckTagNameIsExist(ctx context.Context, spaceId int64, tagName string) (bool, error)

	GetSpaceTag(ctx context.Context, spaceId int64, tagId int64) (*domain.SpaceTag, error)
	GetSpaceTags(ctx context.Context, spaceId int64, tagIds []int64) ([]*domain.SpaceTag, error)

	GetTagByIds(ctx context.Context, tagIds []int64) ([]*domain.SpaceTag, error)
	TagMap(ctx context.Context, ids []int64) (map[int64]*domain.SpaceTag, error)

	DelSpaceTag(ctx context.Context, id int64) error
	DelSpaceTagBySpaceId(ctx context.Context, spaceId int64) error

	GetSpaceTagCount(ctx context.Context, spaceId int64) (int64, error)
}

type SpaceTagQueryRepo interface {
	QSpaceTagList(ctx context.Context, spaceId int64) ([]*domain.SpaceTag, error)
}
