package repo

import (
	"context"
	"go-cs/internal/bean/vo/rsp"
	domain "go-cs/internal/domain/space_member"
)

type SpaceMemberRepo interface {
	SpaceMemberQueryRepo

	//新增项目空间成员
	AddSpaceMember(context.Context, *domain.SpaceMember) error
	AddSpaceMembers(context.Context, []*domain.SpaceMember) error

	DelSpaceMember(ctx context.Context, spaceId int64, userId int64) error

	GetSpaceMember(ctx context.Context, spaceId int64, userId int64) (*domain.SpaceMember, error)

	GetSpaceMemberByUserIds(ctx context.Context, spaceId int64, userId []int64) ([]*domain.SpaceMember, error)
	GetSpaceMemberBySpaceId(ctx context.Context, spaceId int64) ([]*domain.SpaceMember, error)

	SaveSpaceMember(ctx context.Context, member *domain.SpaceMember) error

	IsExistSpaceMember(ctx context.Context, spaceId int64, userId int64) (bool, error)

	AllIsMember(ctx context.Context, spaceId int64, userIds ...int64) (bool, error)
	AnyOneIsMember(ctx context.Context, spaceId int64, userIds ...int64) (bool, error)
	UserInAllSpace(ctx context.Context, userId int64, spaceIds ...int64) (bool, error)
	GetUserSpaceIdList(ctx context.Context, userId int64) ([]int64, error)
	UserSpaceRoleMap(ctx context.Context, userId int64, spaceIds []int64) (map[int64]int64, error)
	UserSpaceMemberMap(ctx context.Context, userId int64, spaceIds []int64) (map[int64]*domain.SpaceMember, error)
	UserSpaceMemberMapFromDB(ctx context.Context, userId int64, spaceIds []int64) (map[int64]*domain.SpaceMember, error)
	GetUserSpaceMemberList(ctx context.Context, userId int64, spaceIds []int64) ([]*domain.SpaceMember, error)
	GetSpaceAllMemberIds(ctx context.Context, spaceIds ...int64) ([]int64, error)
	GetUserSpaceMemberBySpaceId(ctx context.Context, userId int64, spaceIds []int64) (map[int64]*domain.SpaceMember, error)
	UpdateUserSpaceOrder(ctx context.Context, userId int64, formIdx, toIdx int64) error
	GetUserAllSpaceMember(ctx context.Context, userId int64) ([]*domain.SpaceMember, error)
	GetUserSpaceMemberListFromDB(ctx context.Context, userId int64, spaceId []int64) ([]*domain.SpaceMember, error)
	GetMultiUserSpaceMemberMap(ctx context.Context, userIds []int64) (map[int64][]*domain.SpaceMember, error)
	GetManagerIds(ctx context.Context, spaceId int64) ([]int64, error)
	GetSuperManagerIds(ctx context.Context, spaceId int64) ([]int64, error)

	//管理员列表
	SpaceMemberMapByUserIds(ctx context.Context, spaceId int64, userIds []int64) (map[int64]*domain.SpaceMember, error)
	SpaceMemberByUserIdsFromDB(ctx context.Context, spaceId int64, userIds []int64) ([]*domain.SpaceMember, error)

	SpaceMemberMapBySpaceIds(ctx context.Context, spaceIds []int64) (map[int64][]*domain.SpaceMember, error)
	SpaceMemberNumMapBySpaceIds(ctx context.Context, spaceIds []int64) (map[int64]int64, error)

	//删除空间下所有的成员
	DelSpaceMemberBySpaceId(ctx context.Context, spaceId int64) (int64, error)
}

type SpaceMemberQueryRepo interface {
	QSpaceManagerList(ctx context.Context, spaceId int64) ([]*rsp.SpaceMemberInfo, error)
	QSpaceMemberList(ctx context.Context, spaceId int64, userName string) ([]*rsp.SpaceMemberInfo, error)
	QSpaceMemberByUids(ctx context.Context, spaceId int64, userIds []int64) ([]*rsp.SpaceMemberInfo, error)
}
