package repo

import (
	"context"
	"go-cs/internal/domain/space_view"
)

type SpaceViewRepo interface {
	CreateUserView(ctx context.Context, item *space_view.SpaceUserView) error
	CreateUserViews(ctx context.Context, items []*space_view.SpaceUserView) error
	CreateGlobalView(ctx context.Context, item *space_view.SpaceGlobalView) error
	CreateGlobalViews(ctx context.Context, items []*space_view.SpaceGlobalView) error

	GetUserViewById(ctx context.Context, id int64) (*space_view.SpaceUserView, error)
	GetGlobalViewById(ctx context.Context, id int64) (*space_view.SpaceGlobalView, error)
	GetGlobalViewBySpaceIds(ctx context.Context, spaceIds []int64) ([]*space_view.SpaceGlobalView, error)

	GetGlobalViewMap(ctx context.Context, spaceId int64) (map[int64]*space_view.SpaceGlobalView, error)
	GetGlobalViewList(ctx context.Context, spaceId int64) ([]*space_view.SpaceGlobalView, error)
	GetUserViewMap(ctx context.Context, userId, spaceId int64) (map[int64]*space_view.SpaceUserView, error)
	GetUserViewMapByIds(ctx context.Context, ids []int64) (map[int64]*space_view.SpaceUserView, error)

	UserViewList(ctx context.Context, userId int64, spaceIds []int64, key string) ([]*space_view.SpaceUserView, error)

	SaveSpaceUserView(ctx context.Context, workVersion *space_view.SpaceUserView) error
	SaveSpaceGlobalView(ctx context.Context, workVersion *space_view.SpaceGlobalView) error

	DeleteUserViewById(ctx context.Context, id int64) error
	DeleteUserViewByOuterId(ctx context.Context, outerId int64) error
	DeleteUserViewBySpaceId(ctx context.Context, spaceId int64) error
	DeleteUserViewByUserId(ctx context.Context, userId, spaceId int64) error

	DeleteGlobalViewById(ctx context.Context, id int64) error
	DeleteGlobalViewBySpaceId(ctx context.Context, id int64) error

	GetMaxRanking(ctx context.Context, spaceId int64) (int64, error)
}
