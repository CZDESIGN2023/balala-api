package biz

import (
	"context"
	v1 "go-cs/api/admin/v1"
	"go-cs/api/comm"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	user_repo "go-cs/internal/domain/user/repo"
	user_service "go-cs/internal/domain/user/service"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/rand"
	"go-cs/pkg/stream"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type AdminRepo interface {
	SearchUser(ctx context.Context, keyword string, pos, size int, sorts []vo.OrderBy) (list []*db.User, nextPos int, hasNext bool, err error)
	TotalUser(ctx context.Context, keyword string) (int64, error)
	SetAdminNeedResetPwdStatus(ctx context.Context, userId int64) error
	GetAdminNeedResetPwdStatus(ctx context.Context, userId int64) (bool, error)
	DelAdminNeedResetPwdStatus(ctx context.Context, userId int64) error
}

type AdminUsecase struct {
	tm       trans.Transaction
	log      *log.Helper
	repo     AdminRepo
	userRepo user_repo.UserRepo

	userService *user_service.UserService

	domainMessageProducer *domain_message.DomainMessageProducer
}

func NewAdminUsecase(repo AdminRepo, tm trans.Transaction,
	userRepo user_repo.UserRepo,
	userService *user_service.UserService,
	domainMessageProducer *domain_message.DomainMessageProducer,

	logger log.Logger) *AdminUsecase {
	moduleName := "AdminUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &AdminUsecase{
		tm:                    tm,
		log:                   hlog,
		repo:                  repo,
		userRepo:              userRepo,
		domainMessageProducer: domainMessageProducer,
		userService:           userService,
	}
}

func (uc *AdminUsecase) SearchUser(ctx context.Context, uid int64, req *v1.SearchUserRequest) (*v1.SearchUserReplyData, error) {
	user, err := uc.userRepo.GetUserByUserId(ctx, uid)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}
	if user.Role == 0 {
		return nil, errs.NoPerm(ctx)
	}

	orderBy := stream.Map(req.Sorts, func(orderBy *v1.SearchUserRequest_OrderBy) vo.OrderBy {
		return vo.OrderBy{
			Field: orderBy.Field,
			Order: orderBy.Order,
		}
	})

	list, nextPos, hasNext, err := uc.repo.SearchUser(ctx, req.Keyword, int(req.Pos), int(req.Size), orderBy)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var items []*v1.SearchUserReplyData_User
	for _, v := range list {
		items = append(items, &v1.SearchUserReplyData_User{
			Id:            v.Id,
			UserName:      v.UserName,
			UserNickname:  v.UserNickname,
			Avatar:        v.Avatar,
			Role:          v.Role,
			LastLoginTime: v.LastLoginTime,
			CreatedAt:     v.CreatedAt,
		})
	}

	total, err := uc.repo.TotalUser(ctx, req.Keyword)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	return &v1.SearchUserReplyData{
		Items:   items,
		HasNext: hasNext,
		NextPos: int64(nextPos),
		Total:   total,
	}, nil
}

func (uc *AdminUsecase) ResetPwd(ctx context.Context, uid int64, userId int64) (*v1.ResetPwdReplyData, error) {
	userMap, err := uc.userRepo.UserMap(ctx, []int64{userId, uid})
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_USER_GET_FAIL)
	}

	operator := userMap[uid]
	user := userMap[userId]

	if uid != userId && operator.Role <= user.Role {
		return nil, errs.NoPerm(ctx)
	}

	newPwd := rand.Digits(8)
	user.ResetPwd(newPwd, shared.SysOper(uid))

	err = uc.userRepo.SaveUser(ctx, user)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	msg := &domain_message.AdminResetUserPassword{
		DomainMessageBase: shared.DomainMessageBase{
			Oper:     shared.SysOper(uid),
			OperTime: time.Now(),
		},
		UserId: userId,
	}

	uc.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})

	return &v1.ResetPwdReplyData{NewPwd: newPwd}, nil
}

func (uc *AdminUsecase) SetRole(ctx context.Context, uid int64, userId int64, roleId consts.SystemRole) error {
	if uid == userId {
		return errs.NoPerm(ctx)
	}

	if roleId >= consts.SystemRole_SuperAdmin {
		return errs.NoPerm(ctx)
	}

	userMap, err := uc.userRepo.UserMap(ctx, []int64{userId, uid})
	if err != nil {
		return errs.New(ctx, comm.ErrorCode_USER_GET_FAIL)
	}

	operator := userMap[uid]
	if operator.Role != consts.SystemRole_SuperAdmin {
		return errs.NoPerm(ctx)
	}

	user := userMap[userId]
	if user == nil {
		return errs.NoPerm(ctx)
	}

	oldValue := user.Role

	//设置角色
	user.UpdateRole(roleId)

	err = uc.userRepo.SaveUser(ctx, user)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	switch {
	case roleId >= consts.SystemRole_Admin:
		uc.repo.SetAdminNeedResetPwdStatus(ctx, userId)
	case roleId == consts.SystemRole_Normal:
		uc.repo.DelAdminNeedResetPwdStatus(ctx, userId)
	}

	msg := &domain_message.AdminChangeUserRole{
		DomainMessageBase: shared.DomainMessageBase{
			Oper:     shared.SysOper(uid),
			OperTime: time.Now(),
		},
		UserId:   userId,
		OldValue: oldValue,
		NewValue: user.Role,
	}

	uc.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})

	return nil
}

func (uc *AdminUsecase) SetNickname(ctx context.Context, uid int64, userId int64, nickname string) error {
	userMap, err := uc.userRepo.UserMap(ctx, []int64{userId, uid})
	if err != nil {
		return errs.New(ctx, comm.ErrorCode_USER_GET_FAIL)
	}

	operator := userMap[uid]
	user := userMap[userId]

	if uid != userId && operator.Role <= user.Role {
		return errs.NoPerm(ctx)
	}

	oldValue := user.UserNickname

	//设置新昵称
	err = user.ChangeNickName(nickname, nil)
	if err != nil {
		return err
	}

	err = uc.userRepo.SaveUser(ctx, user)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	msg := &domain_message.AdminChangeUserNickname{
		DomainMessageBase: shared.DomainMessageBase{
			Oper:     shared.SysOper(uid),
			OperTime: time.Now(),
		},
		UserId:   userId,
		OldValue: oldValue,
		NewValue: nickname,
	}

	uc.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})

	return nil
}

func (uc *AdminUsecase) AddUser(ctx context.Context, uid int64, req *v1.AddUserRequest) error {
	userMap, err := uc.userRepo.UserMap(ctx, []int64{uid})
	if err != nil {
		return errs.New(ctx, comm.ErrorCode_USER_GET_FAIL)
	}

	operator := userMap[uid]
	if operator.Role < consts.SystemRole_Admin {
		return errs.NoPerm(ctx)
	}

	newUser, err := uc.userService.NewUser(ctx, req.Username, req.Nickname, req.Password, "", 0, nil)
	if err != nil {
		return err
	}

	err = uc.tm.InTx(ctx, func(ctx context.Context) error {
		err = uc.userRepo.AddUser(ctx, newUser)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = uc.userRepo.InitUserConfig(ctx, newUser.Id)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	uc.domainMessageProducer.Send(ctx, shared.DomainMessages{
		&domain_message.AdminAddUser{
			DomainMessageBase: shared.DomainMessageBase{
				Oper:     shared.SysOper(uid),
				OperTime: time.Now(),
			},
			UserId: newUser.Id,
		},
	})

	return nil
}
