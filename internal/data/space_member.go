package data

import (
	"context"
	"encoding/json"
	"errors"
	"go-cs/api/comm"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/consts"
	"go-cs/internal/data/convert"
	domain "go-cs/internal/domain/space_member"
	repo "go-cs/internal/domain/space_member/repo"
	"go-cs/internal/utils"
	"go-cs/pkg/bus"
	"go-cs/pkg/stream"
	"go-cs/pkg/stream/tuple"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type spaceMemberRepo struct {
	baseRepo
}

func NewSpaceMemberRepo(data *Data, logger log.Logger) repo.SpaceMemberRepo {
	moduleName := "SpaceMemberRepo"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	repo := &spaceMemberRepo{
		baseRepo: baseRepo{
			data: data,
			log:  hlog,
		},
	}

	bus.On(comm.CanalEvent_ce_SpaceMember, "", repo.clearCache2)
	return repo
}

func (p *spaceMemberRepo) AddSpaceMember(ctx context.Context, member *domain.SpaceMember) error {

	po := convert.SpaceMemberEntityToPo(member)
	err := p.data.DB(ctx).Save(po).Error
	return err
}

func (p *spaceMemberRepo) AddSpaceMembers(ctx context.Context, members []*domain.SpaceMember) error {

	inserts := stream.Map(members, func(v *domain.SpaceMember) *db.SpaceMember {
		return convert.SpaceMemberEntityToPo(v)
	})

	err := p.data.DB(ctx).CreateInBatches(inserts, 100).Error

	return err
}

func (p *spaceMemberRepo) DelSpaceMember(ctx context.Context, spaceId int64, userId int64) error {
	var opValue = make(map[string]interface{})
	err := p.data.DB(ctx).Model(&db.SpaceMember{}).Where("user_id=? and space_id=?", userId, spaceId).Unscoped().Delete(&opValue).Error
	if err != nil {
		return err
	}
	//p.clearCache(ctx, userId, spaceId)
	return err
}

func (c *spaceMemberRepo) DelSpaceMemberBySpaceId(ctx context.Context, spaceId int64) (int64, error) {
	res := c.data.DB(ctx).
		Where("space_id = ?", spaceId).
		Unscoped().
		Delete(&db.SpaceMember{})

	return res.RowsAffected, res.Error
}

func (p *spaceMemberRepo) SaveSpaceMember(ctx context.Context, member *domain.SpaceMember) error {

	diffs := member.GetDiffs()
	if len(diffs) == 0 {
		return nil
	}

	m := &db.SpaceMember{}
	mCloums := m.Cloumns()

	updateCloums := make(map[string]interface{})
	for _, v := range diffs {
		switch v {
		case domain.Diff_RoleId:
			updateCloums[mCloums.RoleId] = member.RoleId
			updateCloums[mCloums.HistoryRoleId] = member.HistoryRoleId
		case domain.Diff_Ranking:
			updateCloums[mCloums.Ranking] = member.Ranking
		case domain.Diff_Notify:
			updateCloums[mCloums.Notify] = member.Notify
		}
	}

	if len(updateCloums) == 0 {
		return nil
	}

	updateCloums[mCloums.UpdatedAt] = time.Now().Unix()
	err := p.data.DB(ctx).Model(m).Where("id=?", member.Id).UpdateColumns(updateCloums).Error
	if err != nil {
		return err
	}

	return err
}

func (p *spaceMemberRepo) IsExistSpaceMember(ctx context.Context, spaceId int64, userId int64) (bool, error) {

	var count int64
	res := p.data.RoDB(ctx).Model(&db.SpaceMember{}).Where("space_id = ? AND user_id = ?", spaceId, userId).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}

	return count > 0, nil
}

func (p *spaceMemberRepo) AllIsMember(ctx context.Context, spaceId int64, userIds ...int64) (bool, error) {
	if len(userIds) == 0 {
		return true, nil
	}

	userIds = stream.Unique(userIds)

	var count int64
	res := p.data.RoDB(ctx).Model(&db.SpaceMember{}).Where("space_id = ? AND user_id in ?", spaceId, userIds).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}

	if int(count) == len(userIds) {
		return true, nil
	}

	return false, nil
}

func (p *spaceMemberRepo) AnyOneIsMember(ctx context.Context, spaceId int64, userIds ...int64) (bool, error) {
	if len(userIds) == 0 {
		return true, nil
	}

	userIds = stream.Unique(userIds)

	var count int64
	res := p.data.RoDB(ctx).Model(&db.SpaceMember{}).Where("space_id = ? AND user_id in ?", spaceId, userIds).Count(&count)
	if res.Error != nil {
		return false, res.Error
	}

	return count > 0, nil
}

func (c *spaceMemberRepo) UserInAllSpace(ctx context.Context, userId int64, spaceIds ...int64) (bool, error) {
	spaceIds = stream.Unique(spaceIds)

	var count int64
	err := c.data.RoDB(ctx).Model(&db.SpaceMember{}).
		Where("user_id = ?", userId).
		Where("space_id in ?", spaceIds).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return len(spaceIds) == int(count), nil
}

func (c *spaceMemberRepo) GetUserSpaceIdList(ctx context.Context, userId int64) ([]int64, error) {
	var list []int64
	err := c.data.RoDB(ctx).Model(&db.SpaceMember{}).
		Where("user_id = ?", userId).
		Pluck("space_id", &list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (p *spaceMemberRepo) UserSpaceRoleMap(ctx context.Context, userId int64, spaceIds []int64) (map[int64]int64, error) {
	list, err := p.GetUserSpaceMemberList(ctx, userId, spaceIds)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(i int, t *domain.SpaceMember) (int64, int64) {
		return t.SpaceId, t.RoleId
	})

	return m, nil
}

func (p *spaceMemberRepo) GetSpaceMember(ctx context.Context, spaceId int64, userId int64) (*domain.SpaceMember, error) {
	list, err := p.GetUserSpaceMemberList(ctx, userId, []int64{spaceId})
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		var user db.User
		err := p.data.RoDB(ctx).Where("id = ?", userId).First(&user).Error
		if err != nil {
			return nil, err
		}

		if user.Role == int64(consts.SystemRole_Enterprise) {
			return &domain.SpaceMember{
				UserId:  userId,
				SpaceId: spaceId,
				RoleId:  consts.MEMBER_ROLE_WATCHER,
			}, nil
		}

		return nil, gorm.ErrRecordNotFound
	}

	return list[0], nil
}

func (p *spaceMemberRepo) GetSpaceMemberByUserIds(ctx context.Context, spaceId int64, userIds []int64) ([]*domain.SpaceMember, error) {
	var rows []*db.SpaceMember
	err := p.data.DB(ctx).Model(&db.SpaceMember{}).
		Where("space_id = ? AND user_id in ?", spaceId, userIds).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	list := convert.SpaceMemberPosToEntity(rows)
	return list, nil
}

func (p *spaceMemberRepo) GetSpaceMemberBySpaceId(ctx context.Context, spaceId int64) ([]*domain.SpaceMember, error) {
	var rows []*db.SpaceMember
	err := p.data.DB(ctx).Model(&db.SpaceMember{}).
		Where("space_id = ? ", spaceId).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}

	list := convert.SpaceMemberPosToEntity(rows)
	return list, nil
}

//func (p *spaceMemberRepo) clearCache(ctx context.Context, userId int64, spaceId ...int64) {
//	keys := stream.Map(spaceId, func(t int64) string {
//		return NewSpaceKey(t).Member(userId).Key()
//	})
//	p.data.rdb.Del(ctx, keys...)
//}

func (p *spaceMemberRepo) clearCache2(list []tuple.Pair[int64, int64]) {
	keys := stream.Map(list, func(t tuple.Pair[int64, int64]) string {
		return NewSpaceKey(t.First).Member(t.Second).Key()
	})
	p.data.rdb.Del(context.Background(), keys...)
}

func (p *spaceMemberRepo) UserSpaceMemberMap(ctx context.Context, userId int64, spaceIds []int64) (map[int64]*domain.SpaceMember, error) {
	spaceIds = stream.Unique(spaceIds)

	list, err := p.GetUserSpaceMemberList(ctx, userId, spaceIds)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(i int, v *domain.SpaceMember) (int64, *domain.SpaceMember) {
		return v.SpaceId, v
	})

	if len(spaceIds) != len(m) {
		for _, id := range spaceIds {
			if _, ok := m[id]; !ok {
				m[id] = &domain.SpaceMember{
					UserId: userId,
					RoleId: consts.MEMBER_ROLE_WATCHER,
				}
			}
		}
	}

	return m, nil
}

func (p *spaceMemberRepo) UserSpaceMemberMapFromDB(ctx context.Context, userId int64, spaceIds []int64) (map[int64]*domain.SpaceMember, error) {
	list, err := p.GetUserSpaceMemberListFromDB(ctx, userId, spaceIds)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(i int, v *domain.SpaceMember) (int64, *domain.SpaceMember) {
		return v.SpaceId, v
	})

	return m, nil
}

func (p *spaceMemberRepo) GetUserSpaceMemberList(ctx context.Context, userId int64, spaceIds []int64) ([]*domain.SpaceMember, error) {
	spaceIds = stream.Unique(spaceIds)
	spaceIds = stream.Filter(spaceIds, func(id int64) bool {
		return id != 0
	})

	if len(spaceIds) == 0 {
		return nil, nil
	}

	fromRedis, err := p.getUserSpaceMemberListFromCache(ctx, userId, spaceIds)
	if err != nil {
		p.log.Error(err)
	}

	if len(spaceIds) == len(fromRedis) {
		return fromRedis, nil
	}

	redisUserIds := stream.Map(fromRedis, func(v *domain.SpaceMember) int64 {
		return v.Id
	})

	var noCachedIds []int64
	for _, v := range spaceIds {
		if !stream.Contains(redisUserIds, v) {
			noCachedIds = append(noCachedIds, v)
		}
	}

	fromDB, err := p.GetUserSpaceMemberListFromDB(ctx, userId, noCachedIds)
	if err != nil {
		return nil, err
	}

	kv := map[string]string{}
	for _, v := range fromDB {
		k := NewSpaceKey(v.SpaceId).Member(v.UserId).Key()
		v := utils.ToJSON(v)
		kv[k] = v
	}

	_, _ = p.data.rdb.Pipelined(ctx, func(pipeline redis.Pipeliner) error {
		for k, v := range kv {
			pipeline.Set(ctx, k, v, time.Hour*24*3)
		}
		return nil
	})

	return append(fromRedis, fromDB...), nil
}

func (p *spaceMemberRepo) GetUserSpaceMemberListFromDB(ctx context.Context, userId int64, spaceId []int64) ([]*domain.SpaceMember, error) {
	var rows []*db.SpaceMember
	err := p.data.DB(ctx).Where("user_id=? and space_id in ?", userId, spaceId).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	var list []*domain.SpaceMember
	for _, v := range rows {
		list = append(list, convert.SpaceMemberPoToEntity(v))
	}

	return list, nil
}

func (p *spaceMemberRepo) getUserSpaceMemberListFromCache(ctx context.Context, userId int64, spaceIds []int64) ([]*domain.SpaceMember, error) {

	return nil, nil

	// keys := stream.Map(spaceIds, func(v int64) string {
	// 	return NewSpaceKey(v).Member(userId).Key()
	// })

	// result, err := p.data.rdb.MGet(ctx, keys...).Result()
	// if err != nil {
	// 	return nil, err
	// }

	// result = stream.Filter(result, func(e interface{}) bool {
	// 	return e != nil
	// })

	// if len(result) == 0 {
	// 	return nil, nil
	// }

	// list := stream.Map(result, func(v any) *domain.SpaceMember {
	// 	var u domain.SpaceMember
	// 	json.Unmarshal([]byte(v.(string)), &u)
	// 	return &u
	// })

	// return list, nil
}

func (p *spaceMemberRepo) GetSpaceAllMemberIds(ctx context.Context, spaceIds ...int64) ([]int64, error) {
	var ids []int64
	err := p.data.RoDB(ctx).Model(&db.SpaceMember{}).
		Where("space_id IN ?", spaceIds).
		Pluck("user_id", &ids).Error
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (p *spaceMemberRepo) GetUserSpaceMemberBySpaceId(ctx context.Context, userId int64, spaceIds []int64) (map[int64]*domain.SpaceMember, error) {
	var list []*db.SpaceMember
	err := p.data.DB(ctx).Model(&db.SpaceMember{}).
		Where("user_id = ? AND space_id in ?", userId, spaceIds).
		Find(&list).Error

	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(i int, t *db.SpaceMember) (int64, *domain.SpaceMember) {
		return t.SpaceId, convert.SpaceMemberPoToEntity(t)
	})

	return m, nil
}

func (p *spaceMemberRepo) UpdateUserSpaceOrder(ctx context.Context, userId int64, formIdx, toIdx int64) error {
	member, err := p.GetUserAllSpaceMember(ctx, userId)
	if err != nil {
		return err
	}

	if formIdx < 0 || formIdx >= int64(len(member)) {
		return errors.New("wrong index")
	}

	stream.Move(member, int(formIdx), int(toIdx))

	for i, v := range member {
		newRanking := len(member) - i
		if v.Ranking != int64(newRanking) {
			err = p.data.DB(ctx).Model(&db.SpaceMember{}).
				Where("id = ?", v.Id).
				Update("ranking", newRanking).Error
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (p *spaceMemberRepo) GetUserAllSpaceMember(ctx context.Context, userId int64) ([]*domain.SpaceMember, error) {
	var rows []*db.SpaceMember
	err := p.data.DB(ctx).Model(&db.SpaceMember{}).
		Where("user_id = ?", userId).
		Order("ranking DESC, id DESC").
		Find(&rows).Error

	var list []*domain.SpaceMember
	for _, v := range rows {
		list = append(list, convert.SpaceMemberPoToEntity(v))
	}
	return list, err
}

func (p *spaceMemberRepo) GetMultiUserSpaceMemberMap(ctx context.Context, userIds []int64) (map[int64][]*domain.SpaceMember, error) {
	var rows []*db.SpaceMember
	err := p.data.DB(ctx).Model(&db.SpaceMember{}).
		Where("user_id in ?", userIds).
		Find(&rows).Error

	var list []*domain.SpaceMember
	for _, v := range rows {
		list = append(list, convert.SpaceMemberPoToEntity(v))
	}

	m := stream.GroupBy(list, func(t *domain.SpaceMember) int64 {
		return t.UserId
	})
	return m, err
}

func (p *spaceMemberRepo) GetManagerIds(ctx context.Context, spaceId int64) ([]int64, error) {
	var list []int64
	err := p.data.DB(ctx).Model(&db.SpaceMember{}).
		Where("space_id = ?", spaceId).
		Where("role_id not in ?", []int64{consts.MEMBER_ROLE_EDITOR, consts.MEMBER_ROLE_WATCHER}).
		Pluck("user_id", &list).Error

	return list, err
}

func (p *spaceMemberRepo) GetSuperManagerIds(ctx context.Context, spaceId int64) ([]int64, error) {
	var list []int64
	err := p.data.DB(ctx).Model(&db.SpaceMember{}).
		Where("space_id = ?", spaceId).
		Where("role_id not in ?", []int64{consts.MEMBER_ROLE_EDITOR, consts.MEMBER_ROLE_WATCHER, consts.MEMBER_ROLE_MANAGER}).
		Pluck("user_id", &list).Error

	return list, err
}

func (p *spaceMemberRepo) SpaceMemberMapBySpaceIds(ctx context.Context, spaceIds []int64) (map[int64][]*domain.SpaceMember, error) {
	var list []*db.SpaceMember
	err := p.data.RoDB(ctx).Model(&db.SpaceMember{}).
		Where("space_id in ?", spaceIds).
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	domainEntities := convert.SpaceMemberPosToEntity(list)

	m := stream.GroupBy(domainEntities, func(v *domain.SpaceMember) int64 {
		return v.SpaceId
	})

	return m, nil
}

func (p *spaceMemberRepo) SpaceMemberNumMapBySpaceIds(ctx context.Context, spaceIds []int64) (map[int64]int64, error) {
	type count struct {
		SpaceId int64
		Count   int64
	}
	var list []count
	err := p.data.RoDB(ctx).Model(&db.SpaceMember{}).
		Select("space_id, count(*) count").
		Where("space_id in ?", spaceIds).
		Group("space_id").
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(i int, t count) (int64, int64) {
		return t.SpaceId, t.Count
	})

	return m, nil
}
func (p *spaceMemberRepo) SpaceMemberMapByUserIds(ctx context.Context, spaceId int64, userIds []int64) (map[int64]*domain.SpaceMember, error) {
	list, err := p.SpaceMemberListByUserIds(ctx, spaceId, userIds)
	if err != nil {
		return nil, err
	}

	m := stream.ToMap(list, func(_ int, v *domain.SpaceMember) (int64, *domain.SpaceMember) {
		return v.UserId, v
	})

	return m, nil
}

func (p *spaceMemberRepo) SpaceMemberListByUserIds(ctx context.Context, spaceId int64, userIds []int64) ([]*domain.SpaceMember, error) {
	userIds = stream.Unique(userIds)
	userIds = stream.Filter(userIds, func(id int64) bool {
		return id != 0
	})

	if len(userIds) == 0 {
		return nil, nil
	}

	fromRedis, err := p.SpaceMemberListByUserIdsFromRedis(ctx, spaceId, userIds)
	if err != nil {
		p.log.Error(err)
	}

	if len(userIds) == len(fromRedis) {
		return fromRedis, nil
	}

	redisUserIds := stream.Map(fromRedis, func(v *domain.SpaceMember) int64 {
		return v.Id
	})

	var noCachedIds []int64
	for _, v := range userIds {
		if !stream.Contains(redisUserIds, v) {
			noCachedIds = append(noCachedIds, v)
		}
	}

	fromDB, err := p.SpaceMemberByUserIdsFromDB(ctx, spaceId, noCachedIds)
	if err != nil {
		return nil, err
	}

	kv := map[string]string{}
	for _, v := range fromDB {
		k := NewSpaceKey(v.SpaceId).Member(v.UserId).Key()
		v := utils.ToJSON(v)
		kv[k] = v
	}

	_, _ = p.data.rdb.Pipelined(ctx, func(pipeline redis.Pipeliner) error {
		for k, v := range kv {
			pipeline.Set(ctx, k, v, time.Hour*24*3)
		}
		return nil
	})

	return append(fromRedis, fromDB...), nil
}

func (p *spaceMemberRepo) SpaceMemberByUserIdsFromDB(ctx context.Context, spaceId int64, userIds []int64) ([]*domain.SpaceMember, error) {
	var rows []*db.SpaceMember
	err := p.data.DB(ctx).Where("space_id=? and user_id in ?", spaceId, userIds).Find(&rows).Error
	if err != nil {
		return nil, err
	}

	list := convert.SpaceMemberPosToEntity(rows)
	return list, nil
}

func (p *spaceMemberRepo) SpaceMemberListByUserIdsFromRedis(ctx context.Context, spaceId int64, userIds []int64) ([]*domain.SpaceMember, error) {
	keys := stream.Map(userIds, func(v int64) string {
		return NewSpaceKey(spaceId).Member(v).Key()
	})

	result, err := p.data.rdb.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	result = stream.Filter(result, func(e interface{}) bool {
		return e != nil
	})

	if len(result) == 0 {
		return nil, nil
	}

	list := stream.Map(result, func(v any) *domain.SpaceMember {
		var u domain.SpaceMember
		json.Unmarshal([]byte(v.(string)), &u)
		return &u
	})

	return list, nil
}

func (p *spaceMemberRepo) QSpaceMemberByUids(ctx context.Context, spaceId int64, userIds []int64) ([]*rsp.SpaceMemberInfo, error) {
	var rows []*rsp.SpaceMemberInfo

	tbUser := (&db.User{}).TableName()
	tbSpace := (&db.Space{}).TableName()
	tbSpaceMember := (&db.SpaceMember{}).TableName()

	err := p.data.RoDB(ctx).Table(tbSpaceMember+" m").
		Select("m.id, m.id member_id, m.space_id, m.user_id, m.role_id, u.user_name, u.mobile, u.user_nickname, u.user_pinyin, u.user_status, u.avatar").
		Joins("LEFT JOIN "+tbUser+" u on m.user_id=u.id LEFT JOIN "+tbSpace+" s on m.space_id = s.id").
		Where("m.space_id=? and m.user_id in ?", spaceId, userIds).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (p *spaceMemberRepo) QSpaceManagerList(ctx context.Context, spaceId int64) ([]*rsp.SpaceMemberInfo, error) {
	var list []*rsp.SpaceMemberInfo

	tbUser := (&db.User{}).TableName()
	tbSpace := (&db.Space{}).TableName()
	tbSpaceMember := (&db.SpaceMember{}).TableName()

	tx := p.data.RoDB(ctx).Table(tbSpaceMember + " m")

	err := tx.Select("m.id, m.id member_id, m.space_id, m.user_id, m.role_id, u.user_name, u.mobile, u.user_nickname, u.user_pinyin, u.user_status, u.avatar").
		Joins("LEFT JOIN "+tbUser+" u on m.user_id=u.id LEFT JOIN "+tbSpace+" s on m.space_id = s.id").
		Where("m.space_id=? and m.role_id in ?", spaceId, []int64{consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER, consts.MEMBER_ROLE_SPACE_CREATOR}).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (p *spaceMemberRepo) QSpaceMemberList(ctx context.Context, spaceId int64, userName string) ([]*rsp.SpaceMemberInfo, error) {
	var list []*rsp.SpaceMemberInfo

	tbUser := (&db.User{}).TableName()
	tbSpace := (&db.Space{}).TableName()
	tbSpaceMember := (&db.SpaceMember{}).TableName()

	tx := p.data.RoDB(ctx).Table(tbSpaceMember + " m")

	if userName != "" {
		like := "%" + userName + "%"
		tx.Where("u.user_nickname like ? OR u.user_pinyin like ? OR u.user_name like ?", like, like, like)
	}

	err := tx.Select("m.id, m.id member_id, m.space_id, m.user_id, m.role_id, u.user_name, u.mobile, u.user_nickname, u.user_pinyin, u.user_status, u.avatar").
		Joins("LEFT JOIN "+tbUser+" u on m.user_id=u.id LEFT JOIN "+tbSpace+" s on m.space_id = s.id").
		Where("m.space_id=?", spaceId).
		Find(&list).Error
	if err != nil {
		return nil, err
	}
	return list, nil
}
