package service

import (
	"context"
	"go-cs/api/comm"
	"go-cs/internal/consts"
	domain "go-cs/internal/domain/user"
	"go-cs/internal/domain/user/repo"
	"go-cs/internal/pkg/biz_id"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/rand"
	"time"
)

type UserService struct {
	repo      repo.UserRepo
	idService *biz_id.BusinessIdService
}

func NewUserService(
	repo repo.UserRepo,
	idService *biz_id.BusinessIdService,
) *UserService {
	return &UserService{
		repo:      repo,
		idService: idService,
	}
}

func (s *UserService) NewUser(ctx context.Context, name string, nickName string, pwd string, avatar string, role consts.SystemRole, oper shared.Oper) (*domain.User, error) {

	//检查用户名是否已存在
	isExist, err := s.repo.IsExistByUserName(ctx, name)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}
	if isExist {
		return nil, errs.New(ctx, comm.ErrorCode_USER_NAME_IS_EXIST)
	}

	bizId := s.idService.NewId(ctx, consts.BusinessId_Type_User)
	if bizId == nil {
		return nil, errs.Business(ctx, "生成用户ID失败")
	}

	salt := rand.S(10)
	password := utils.EncryptUserPassword(pwd, salt)
	nickNamePy := utils.Pinyin(nickName)

	user := domain.NewUser(bizId.Id, name, nickName, nickNamePy, password, salt, avatar, role, oper)
	return user, nil
}

func (s *UserService) NewPfAccount(ctx context.Context, user *domain.User, pfInfo domain.ThirdPfInfo) *domain.ThirdPfAccount {
	account := &domain.ThirdPfAccount{
		PfInfo:    pfInfo,
		Notify:    1,
		CreatedAt: time.Now().Unix(),
	}

	if user != nil {
		account.UserId = user.Id
	}

	return account
}
