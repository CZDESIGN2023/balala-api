package repo

import (
	"context"
	"go-cs/internal/bean/vo/rsp"
	domain "go-cs/internal/domain/user"
)

type UserRepo interface {
	UserQueryRepo

	GetUserByIds(ctx context.Context, userIds []int64) ([]*domain.User, error)
	// GetUserByUserId 根据Id获得用户信息
	GetUserByUserId(context.Context, int64) (*domain.User, error)
	//根据用户名获取用户信息
	GetUserByUserName(context.Context, string) (*domain.User, error)
	//判断用户名是否存在
	IsExistByUserName(context.Context, string) (bool, error)
	//判断用户昵称是否存在
	IsExistByUserNickName(context.Context, string) (bool, error)

	//新增用户
	AddUser(context.Context, *domain.User) error
	//初始化用户配置
	InitUserConfig(ctx context.Context, userId int64) error
	//搜索系统所有用户
	//更新用户字段
	SaveUser(ctx context.Context, user *domain.User) error
	//更新用户缓存
	UserMap(ctx context.Context, ids []int64) (map[int64]*domain.User, error)

	BindThirdPlatform(ctx context.Context, account *domain.ThirdPfAccount) error
	UnbindThirdPlatform(ctx context.Context, account *domain.ThirdPfAccount) error
	GetThirdPfAccount(ctx context.Context, userId int64, platformCode int32) (*domain.ThirdPfAccount, error)
	GetThirdPfAccountByPfUserKey(ctx context.Context, pfUserKey string, platformCode int32) (*domain.ThirdPfAccount, error)
	GetAllThirdPfAccount(ctx context.Context, userId int64) ([]*domain.ThirdPfAccount, error)
	GetThirdPfAccountByUserIds(ctx context.Context, userIds []int64) ([]*domain.ThirdPfAccount, error)
	RemoveAllThirdPlatform(ctx context.Context, userId int64) error
	SaveThirdPfAccount(ctx context.Context, account *domain.ThirdPfAccount) error

	FieldAllowChange(ctx context.Context, userId int64, field string) bool
	SetFieldChangeTime(ctx context.Context, userId int64, field string)

	SetTempConfig(ctx context.Context, userId int64, confMap map[string]string) error
	GetTempConfig(ctx context.Context, userId int64, userKeys ...string) map[string]string
	DelTempConfig(ctx context.Context, userId int64, userKeys ...string) error

	GetUserConfig(ctx context.Context, userId int64, key string) (*domain.UserConfig, error)
	GetUserAllConfig(ctx context.Context, userId int64) (map[string]*domain.UserConfig, error)
	GetUserConfigMapByUserIdsAndKey(ctx context.Context, userIds []int64, key string) (map[int64]*domain.UserConfig, error)
	GetUserConfigMapByUserIdsAndKeys(ctx context.Context, userIds []int64, keys []string) (map[int64]map[string]*domain.UserConfig, error)

	SaveUserConfig(ctx context.Context, config *domain.UserConfig) error

	IsEnterpriseAdmin(ctx context.Context, userId int64) bool
}

type UserQueryRepo interface {
	//拼音搜索项目内用户
	SearchUser(ctx context.Context, py string, useIds []int64) ([]*rsp.ViewUserWithSpaceInfo, error)
	SearchSpaceMember(ctx context.Context, py string, spaceIds []int64) ([]*rsp.ViewUserWithSpaceInfo, error)
}
