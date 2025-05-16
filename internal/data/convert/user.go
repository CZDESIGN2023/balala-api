package convert

import (
	db "go-cs/internal/bean/biz"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/user"
)

func UserEntityToPo(user *domain.User) *db.User {
	return &db.User{
		Id:            user.Id,
		UserName:      user.UserName,
		Mobile:        user.Mobile,
		UserNickname:  user.UserNickname,
		UserPassword:  user.UserPassword,
		UserPinyin:    user.UserPinyin,
		UserSalt:      user.UserSalt,
		UserStatus:    user.UserStatus,
		UserEmail:     user.UserEmail,
		Sex:           user.Sex,
		Avatar:        user.Avatar,
		Remark:        user.Remark,
		Describe:      user.Describe,
		LastLoginIp:   user.LastLoginIp,
		LastLoginTime: user.LastLoginTime,
		Role:          int64(user.Role),
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		DeletedAt:     user.DeletedAt,
	}
}

func UserPoToEntity(user *db.User) *domain.User {
	return &domain.User{
		Id:            user.Id,
		UserName:      user.UserName,
		Mobile:        user.Mobile,
		UserNickname:  user.UserNickname,
		UserPassword:  user.UserPassword,
		UserPinyin:    user.UserPinyin,
		UserSalt:      user.UserSalt,
		UserStatus:    user.UserStatus,
		UserEmail:     user.UserEmail,
		Sex:           user.Sex,
		Avatar:        user.Avatar,
		Remark:        user.Remark,
		Describe:      user.Describe,
		LastLoginIp:   user.LastLoginIp,
		LastLoginTime: user.LastLoginTime,
		Role:          consts.SystemRole(user.Role),
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		DeletedAt:     user.DeletedAt,
	}
}

func ThirdPfAccountEntityToPo(account *domain.ThirdPfAccount) *db.ThirdPfAccount {
	return &db.ThirdPfAccount{
		Id:            account.Id,
		UserId:        account.UserId,
		PfCode:        account.PfInfo.PfCode,
		PfName:        account.PfInfo.PfName,
		PfUserKey:     account.PfInfo.PfUserKey,
		PfUserName:    account.PfInfo.PfUserName,
		PfUserId:      account.PfInfo.PfUserId,
		PfUserAccount: account.PfInfo.PfUserAccount,
		Notify:        account.Notify,
		CreatedAt:     account.CreatedAt,
		UpdatedAt:     account.UpdatedAt,
		DeletedAt:     account.DeletedAt,
	}
}

func ThirdPfAccountPoToEntity(account *db.ThirdPfAccount) *domain.ThirdPfAccount {
	ent := &domain.ThirdPfAccount{
		Id:        account.Id,
		UserId:    account.UserId,
		Notify:    account.Notify,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
		DeletedAt: account.DeletedAt,
	}
	ent.PfInfo = domain.ThirdPfInfo{
		PfCode:        account.PfCode,
		PfName:        account.PfName,
		PfUserKey:     account.PfUserKey,
		PfUserName:    account.PfUserName,
		PfUserId:      account.PfUserId,
		PfUserAccount: account.PfUserAccount,
	}

	return ent
}

func UserConfigPoToEntity(e *db.UserConfig) *domain.UserConfig {
	return &domain.UserConfig{
		Id:        e.Id,
		UserId:    e.UserId,
		Key:       e.Key,
		Value:     e.Value,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func UserConfigEntityToPo(e *domain.UserConfig) *db.UserConfig {
	return &db.UserConfig{
		Id:        e.Id,
		UserId:    e.UserId,
		Key:       e.Key,
		Value:     e.Value,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

func UserConfigPoToEntities(list []*db.UserConfig) []*domain.UserConfig {
	entities := make([]*domain.UserConfig, 0, len(list))
	for _, v := range list {
		entities = append(entities, UserConfigPoToEntity(v))
	}

	return entities
}
