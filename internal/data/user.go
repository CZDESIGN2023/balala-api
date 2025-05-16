package data

import (
	"context"
	"encoding/json"
	"fmt"
	"go-cs/api/comm"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/consts"
	"go-cs/internal/data/convert"
	domain "go-cs/internal/domain/user"
	"go-cs/internal/domain/user/repo"
	"go-cs/internal/utils"
	"go-cs/pkg/bus"
	"go-cs/pkg/stream"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	baseRepo
}

func RegisterHandler(key string) {

}

// --缓存失效时间
const kCacheTime = 60 * 60 * 24

func NewUserRepo(data *Data, logger log.Logger) repo.UserRepo {
	moduleName := "UserRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	repo := &userRepo{
		baseRepo{
			data: data,
			log:  hlog,
		},
	}

	bus.On(comm.CanalEvent_ce_User, "", repo.clearCache2)

	return repo
}

func (c *userRepo) GetUserByUserName(ctx context.Context, userName string) (*domain.User, error) {

	row := &db.User{}
	err := c.data.RoDB(ctx).Model(&db.User{}).Where("BINARY user_name = ?", userName).Take(row).Error
	if err != nil {
		return nil, err
	}
	return convert.UserPoToEntity(row), err
}

func (c *userRepo) IsExistByUserName(ctx context.Context, userName string) (bool, error) {

	var rowNum int64
	err := c.data.RoDB(ctx).Model(&db.User{}).Where("BINARY user_name=?", userName).Count(&rowNum).Error
	if err != nil {
		return false, err
	}

	if rowNum > 0 {
		return true, nil
	}
	return false, nil
}

func (c *userRepo) ExistEmail(ctx context.Context, email string) (bool, error) {

	var rowNum int64
	err := c.data.RoDB(ctx).Model(&db.User{}).Where("BINARY user_email=?", email).Count(&rowNum).Error
	if err != nil {
		return false, err
	}

	if rowNum > 0 {
		return true, nil
	}
	return false, nil
}

func (c *userRepo) ExistMobile(ctx context.Context, mobile string) (bool, error) {

	var rowNum int64
	err := c.data.RoDB(ctx).Model(&db.User{}).Where("mobile = ?", mobile).Count(&rowNum).Error
	if err != nil {
		return false, err
	}

	if rowNum > 0 {
		return true, nil
	}
	return false, nil
}

func (c *userRepo) IsExistByUserNickName(ctx context.Context, userNickName string) (bool, error) {
	var rowNum int64
	err := c.data.RoDB(ctx).Model(&db.User{}).Where("BINARY user_nickname=?", userNickName).Count(&rowNum).Error
	if err != nil {
		return false, err
	}

	if rowNum > 0 {
		return true, nil
	}
	return false, nil
}

func (c *userRepo) AddUser(ctx context.Context, user *domain.User) error {
	var omitFields []string
	if user.UserEmail == "" {
		omitFields = append(omitFields, "user_email")
	}
	if user.Mobile == "" {
		omitFields = append(omitFields, "mobile")
	}

	err := c.data.DB(ctx).Omit(omitFields...).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (c *userRepo) InitUserConfig(ctx context.Context, userId int64) error {
	// 默认配置
	err := c.data.db.Create([]*db.UserConfig{
		{
			UserId: userId,
			Key:    consts.UserConfigKey_NotifySwitchGlobal,
			Value:  "1",
		},
		{
			UserId: userId,
			Key:    consts.UserConfigKey_NotifySwitchThirdPlatform,
			Value:  "1",
		},
		{
			UserId: userId,
			Key:    consts.UserConfigKey_NotifySwitchSpace,
			Value:  "1",
		},
	}).Error

	return err
}

func (c *userRepo) SaveUser(ctx context.Context, user *domain.User) error {

	diffs := user.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.User{}
	mColumns := m.Cloumns()

	updateColumns := make(map[string]interface{})

	for _, v := range diffs {
		switch v {
		case domain.Diff_Avatar:
			updateColumns[mColumns.Avatar] = user.Avatar
		case domain.Diff_Password:
			updateColumns[mColumns.UserPassword] = user.UserPassword
			updateColumns[mColumns.UserSalt] = user.UserSalt

		case domain.Diff_NickName:
			updateColumns[mColumns.UserNickname] = user.UserNickname
			updateColumns[mColumns.UserPinyin] = user.UserPinyin

		case domain.Diff_Status:
			updateColumns[mColumns.UserStatus] = user.UserStatus
		case domain.Diff_Name:
			updateColumns[mColumns.UserName] = user.UserName
			if user.UserName == "" {
				updateColumns[mColumns.UserName] = nil
			}
		case domain.Diff_Role:
			updateColumns[mColumns.Role] = user.Role
		case domain.Diff_LastLogin:
			updateColumns[mColumns.LastLoginIp] = user.LastLoginIp
			updateColumns[mColumns.LastLoginTime] = user.LastLoginTime
		}
	}

	if len(updateColumns) == 0 {
		return nil
	}

	updateColumns[mColumns.UpdatedAt] = time.Now().Unix()
	err := c.data.DB(ctx).Model(&db.User{}).Where("id=?", user.Id).UpdateColumns(updateColumns).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *userRepo) clearCache2(userId []int64) {
	keys := stream.Map(userId, func(t int64) string {
		return NewUserKey(t).Key()
	})

	c.data.rdb.Del(context.Background(), keys...)
}

//func (c *userRepo) clearCache(ctx context.Context, userId ...int64) {
//	return
//	keys := stream.Map(userId, func(t int64) string {
//		return NewUserKey(t).Key()
//	})
//
//	c.data.rdb.Del(ctx, keys...)
//}

func (c *userRepo) UserMap(ctx context.Context, ids []int64) (map[int64]*domain.User, error) {
	list, err := c.GetUserByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(i int, t *domain.User) (int64, *domain.User) {
		return t.Id, t
	})

	return m, nil
}

func (c *userRepo) GetUserByUserId(ctx context.Context, id int64) (*domain.User, error) {
	ids, err := c.GetUserByIds(ctx, []int64{id})
	if err != nil {
		return nil, err
	}
	if len(ids) > 0 {
		return ids[0], nil
	}

	return nil, gorm.ErrRecordNotFound
}

func (c *userRepo) GetUserByIds(ctx context.Context, userIds []int64) ([]*domain.User, error) {
	userIds = stream.Unique(userIds)
	userIds = stream.Filter(userIds, func(id int64) bool {
		return id != 0
	})

	if len(userIds) == 0 {
		return nil, nil
	}

	fromRedis, err := c.GetUserByIdsFromRedis(ctx, userIds)
	if err != nil {
		c.log.Error(err)
	}

	if len(userIds) == len(fromRedis) {
		return fromRedis, nil
	}

	noCachedIds := stream.Filter(userIds, func(id int64) bool {
		return !stream.ContainsFunc(fromRedis, func(user *domain.User) bool {
			return id == user.Id
		})
	})

	fromDB, err := c.GetUserByIdsFromDB(ctx, noCachedIds)
	if err != nil {
		return nil, err
	}

	kv := map[string]string{}
	for _, v := range fromDB {
		k := NewUserKey(v.Id).Key()
		v := utils.ToJSON(v)
		kv[k] = v
	}

	_, _ = c.data.rdb.Pipelined(ctx, func(pipeline redis.Pipeliner) error {
		for k, v := range kv {
			pipeline.Set(ctx, k, v, time.Hour*24*3)
		}
		return nil
	})

	return append(fromRedis, fromDB...), nil
}

func (c *userRepo) GetUserByIdsFromRedis(ctx context.Context, userIds []int64) ([]*domain.User, error) {

	keys := stream.Map(userIds, func(v int64) string {
		return NewUserKey(v).Key()
	})

	result, err := c.data.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	result = stream.Filter(result, func(e interface{}) bool {
		return e != nil
	})

	list := stream.Map(result, func(v any) *domain.User {
		var u domain.User
		json.Unmarshal([]byte(v.(string)), &u)
		return &u
	})

	return list, err
}

func (c *userRepo) GetUserByIdsFromDB(ctx context.Context, userIds []int64) ([]*domain.User, error) {
	if len(userIds) == 0 {
		return nil, nil
	}

	var rows []*db.User
	err := c.data.RoDB(ctx).Model(&db.User{}).Where("id in (?)", userIds).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	var users []*domain.User
	for _, v := range rows {
		users = append(users, convert.UserPoToEntity(v))
	}

	return users, err
}

func (c *userRepo) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := &db.User{}
	err := c.data.RoDB(ctx).Where("user_email = ?", email).Take(row).Error
	if err != nil {
		return nil, err
	}
	return convert.UserPoToEntity(row), err
}

func (c *userRepo) GetUserByMobile(ctx context.Context, mobile string) (*domain.User, error) {
	row := &db.User{}
	err := c.data.RoDB(ctx).Where("mobile = ?", mobile).Take(row).Error
	if err != nil {
		return nil, err
	}
	return convert.UserPoToEntity(row), err
}

func (c *userRepo) BindEmail(ctx context.Context, userId int64, email string) error {
	err := c.data.RoDB(ctx).
		Model(&db.User{}).
		Where("id = ?", userId).
		UpdateColumn("user_email", email).Error
	if err != nil {
		return err
	}

	//c.clearCache(ctx, userId)
	return nil
}

func (c *userRepo) UnbindEmail(ctx context.Context, userId int64) error {
	err := c.data.RoDB(ctx).
		Model(&db.User{}).
		Where("id = ?", userId).
		Update("user_email", nil).Error
	if err != nil {
		return err
	}

	//c.clearCache(ctx, userId)
	return nil
}

func (c *userRepo) BindMobile(ctx context.Context, userId int64, mobile string) error {
	err := c.data.RoDB(ctx).
		Model(&db.User{}).
		Where("id = ?", userId).
		UpdateColumn("mobile", mobile).Error
	if err != nil {
		return err
	}

	//c.clearCache(ctx, userId)
	return nil
}

func (c *userRepo) UnbindMobile(ctx context.Context, userId int64) error {
	err := c.data.RoDB(ctx).
		Model(&db.User{}).
		Where("id = ?", userId).
		Update("mobile", nil).Error
	if err != nil {
		return err
	}

	//c.clearCache(ctx, userId)
	return nil
}

func (c *userRepo) BindThirdPlatform(ctx context.Context, account *domain.ThirdPfAccount) error {

	po := convert.ThirdPfAccountEntityToPo(account)
	err := c.data.DB(ctx).Create(po).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *userRepo) UnbindThirdPlatform(ctx context.Context, account *domain.ThirdPfAccount) error {
	err := c.data.DB(ctx).
		Where("user_id = ? AND pf_code = ?", account.UserId, account.PfInfo.PfCode).
		Delete(&db.ThirdPfAccount{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *userRepo) GetThirdPfAccount(ctx context.Context, userId int64, platformCode int32) (*domain.ThirdPfAccount, error) {
	var v *db.ThirdPfAccount
	err := c.data.DB(ctx).
		Where("user_id = ? AND pf_code = ?", userId, platformCode).
		Take(&v).Error
	if err != nil {
		return nil, err
	}

	return convert.ThirdPfAccountPoToEntity(v), nil
}

func (c *userRepo) GetThirdPfAccountByPfUserKey(ctx context.Context, pfUserKey string, platformCode int32) (*domain.ThirdPfAccount, error) {
	var v *db.ThirdPfAccount
	err := c.data.DB(ctx).
		Where("pf_user_key = ? AND pf_code = ?", pfUserKey, platformCode).
		Take(&v).Error
	if err != nil {
		return nil, err
	}
	return convert.ThirdPfAccountPoToEntity(v), nil
}

func (c *userRepo) GetAllThirdPfAccount(ctx context.Context, userId int64) ([]*domain.ThirdPfAccount, error) {
	var list []*db.ThirdPfAccount
	err := c.data.DB(ctx).
		Where("user_id = ?", userId).
		Order("id DESC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	var accounts []*domain.ThirdPfAccount
	for _, v := range list {
		accounts = append(accounts, convert.ThirdPfAccountPoToEntity(v))
	}
	return accounts, nil
}

func (c *userRepo) GetThirdPfAccountByUserIds(ctx context.Context, userIds []int64) ([]*domain.ThirdPfAccount, error) {
	var list []*db.ThirdPfAccount
	err := c.data.DB(ctx).
		Where("user_id IN ?", userIds).
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	var accounts []*domain.ThirdPfAccount
	for _, v := range list {
		accounts = append(accounts, convert.ThirdPfAccountPoToEntity(v))
	}
	return accounts, nil
}

func (c *userRepo) RemoveAllThirdPlatform(ctx context.Context, userId int64) error {
	err := c.data.DB(ctx).Where("user_id = ?", userId).Delete(&db.ThirdPfAccount{}).Error
	return err
}

func (c *userRepo) SaveThirdPfAccount(ctx context.Context, account *domain.ThirdPfAccount) error {
	diffs := account.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.ThirdPfAccount{}
	mColumns := m.Cloumns()

	updateColumns := make(map[string]interface{})

	for _, v := range diffs {
		switch v {
		case domain.Diff_ThirdPlatform_Notify:
			updateColumns[mColumns.Notify] = account.Notify
		}
	}

	if len(updateColumns) == 0 {
		return nil
	}

	updateColumns[mColumns.UpdatedAt] = time.Now().Unix()
	err := c.data.DB(ctx).Model(m).Where("id=?", account.Id).UpdateColumns(updateColumns).Error
	if err != nil {
		return err
	}

	return nil

}

func (c *userRepo) FieldAllowChange(ctx context.Context, userId int64, field string) bool {
	key := fmt.Sprintf("balala:user_field:%d:%s", userId, field)
	return c.data.rdb.Exists(ctx, key).Val() == 0
}

func (c *userRepo) SetFieldChangeTime(ctx context.Context, userId int64, field string) {
	key := fmt.Sprintf("balala:user_field:%d:%s", userId, field)
	c.data.rdb.Set(ctx, key, time.Now().UnixMilli(), time.Hour*24*7)
}

func (c *userRepo) SearchSpaceMember(ctx context.Context, py string, spaceIds []int64) ([]*rsp.ViewUserWithSpaceInfo, error) {

	tbSpaceMember := (&db.SpaceMember{}).TableName()
	tbUser := (&db.User{}).TableName()

	var rows []*rsp.ViewUserWithSpaceInfo
	var err error

	model := c.data.RoDB(ctx).Table(tbUser + " u").Select("u.*, m.space_id, m.id space_memeber_id")
	model = model.Joins("LEFT JOIN " + tbSpaceMember + " m on u.id = m.user_id")
	model = model.Where("u.user_status = 1")

	if py != "" {
		pyLike := "%" + py + "%"
		model = model.Where("u.user_pinyin like ? or u.user_nickname like ? or u.user_name like ?", pyLike, pyLike, pyLike)
	}

	model = model.Where("m.space_id in ?", spaceIds)

	model = model.Find(&rows)
	err = model.Error

	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *userRepo) SearchUser(ctx context.Context, py string, useIds []int64) ([]*rsp.ViewUserWithSpaceInfo, error) {
	var rows []*rsp.ViewUserWithSpaceInfo
	var err error

	model := c.data.RoDB(ctx).
		Model(&db.User{}).Select("*").
		Where("user_status = 1")

	if py != "" {
		pyLike := "%" + py + "%"
		model = model.Where("user_pinyin like ? or user_nickname like ? or user_name like ?", pyLike, pyLike, pyLike)
	}

	if len(useIds) != 0 {
		model = model.Where("id in ? ", useIds)
	} else {
		model = model.Where("role != ?", consts.SystemRole_SuperAdmin)
	}

	err = model.Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *userRepo) SetTempConfig(ctx context.Context, userId int64, confMap map[string]string) error {
	for key, val := range confMap {
		key = fmt.Sprintf("balala:user_temp_config:%d:%s", userId, key)
		c.data.rdb.Set(ctx, key, val, time.Hour*24*90)
	}

	return nil
}

func (c *userRepo) GetTempConfig(ctx context.Context, userId int64, userKeys ...string) map[string]string {
	keys := stream.Map(userKeys, func(v string) string {
		return fmt.Sprintf("balala:user_temp_config:%d:%s", userId, v)
	})

	ret := c.data.rdb.MGet(ctx, keys...)

	values := stream.Map(ret.Val(), func(v any) string {
		return cast.ToString(v)
	})

	for _, key := range keys {
		c.data.rdb.Expire(ctx, key, time.Hour*24*90)
	}

	return stream.Zip(userKeys, values)
}

func (c *userRepo) DelTempConfig(ctx context.Context, userId int64, userKeys ...string) error {
	keys := stream.Map(userKeys, func(v string) string {
		return fmt.Sprintf("balala:user_temp_config:%d:%s", userId, v)
	})

	return c.data.rdb.Del(ctx, keys...).Err()
}

func (c *userRepo) GetUserConfig(ctx context.Context, userId int64, key string) (*domain.UserConfig, error) {
	var config db.UserConfig
	err := c.data.DB(ctx).Where("user_id = ? AND `key` = ?", userId, key).Take(&config).Error
	if err != nil {
		return nil, err
	}

	return convert.UserConfigPoToEntity(&config), nil
}

func (c *userRepo) GetUserAllConfig(ctx context.Context, userId int64) (map[string]*domain.UserConfig, error) {
	var configs []*db.UserConfig
	err := c.data.DB(ctx).Where("user_id = ?", userId).Find(&configs).Error
	if err != nil {
		return nil, err
	}

	configMap := stream.ToMap(convert.UserConfigPoToEntities(configs), func(_ int, v *domain.UserConfig) (string, *domain.UserConfig) {
		return v.Key, v
	})

	return configMap, nil
}

func (c *userRepo) GetUserConfigMapByUserIdsAndKey(ctx context.Context, userIds []int64, key string) (map[int64]*domain.UserConfig, error) {
	var configs []*db.UserConfig
	err := c.data.DB(ctx).Where("user_id IN ? AND `key` = ?", userIds, key).Find(&configs).Error
	if err != nil {
		return nil, err
	}

	configMap := stream.ToMap(convert.UserConfigPoToEntities(configs), func(_ int, v *domain.UserConfig) (int64, *domain.UserConfig) {
		return v.UserId, v
	})

	return configMap, nil
}
func (c *userRepo) GetUserConfigMapByUserIdsAndKeys(ctx context.Context, userIds []int64, keys []string) (map[int64]map[string]*domain.UserConfig, error) {
	var configs []*db.UserConfig
	err := c.data.DB(ctx).Where("user_id IN ? AND `key` IN ?", userIds, keys).Find(&configs).Error
	if err != nil {
		return nil, err
	}

	configMap := stream.GroupBy(convert.UserConfigPoToEntities(configs), func(v *domain.UserConfig) int64 {
		return v.UserId
	})

	finalMap := stream.MapValue(configMap, func(v []*domain.UserConfig) map[string]*domain.UserConfig {
		return stream.ToMap(v, func(_ int, v *domain.UserConfig) (string, *domain.UserConfig) {
			return v.Key, v
		})
	})

	return finalMap, nil
}

func (c *userRepo) SaveUserConfig(ctx context.Context, config *domain.UserConfig) error {
	po := convert.UserConfigEntityToPo(config)

	err := c.data.DB(ctx).Model(&db.UserConfig{}).Where("user_id = ? AND `key` = ?", po.UserId, po.Key).
		Updates(po).Error

	if err != nil {
		return err
	}

	return nil
}

func (c *userRepo) IsEnterpriseAdmin(ctx context.Context, userId int64) bool {

	info, err := c.GetUserByUserId(ctx, userId)
	if err != nil {
		return false
	}

	return info.IsEnterpriseAdmin()
}
