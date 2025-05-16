package biz

import (
	"context"
	"go-cs/api/comm"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/conf"
	"go-cs/internal/consts"
	domain_message "go-cs/internal/domain/pkg/message"
	user_repo "go-cs/internal/domain/user/repo"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"slices"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type ConfigRepo interface {
	GetById(context.Context, int64) (*db.Config, error)
	GetByKey(context.Context, string) (*db.Config, error)
	List(ctx context.Context) ([]*db.Config, error)
	Map(ctx context.Context) (map[string]*db.Config, error)
	UpdateByKey(ctx context.Context, key string, value string) error

	CanRegister(ctx context.Context) bool
	AttachSize(ctx context.Context) int64
}

type ConfigUsecase struct {
	configRepo ConfigRepo
	userRepo   user_repo.UserRepo
	log        *log.Helper
	tm         trans.Transaction
	conf       *conf.Bootstrap

	domainMessageProducer *domain_message.DomainMessageProducer
}

func NewConfigUsecase(
	tm trans.Transaction, conf *conf.Bootstrap, logger log.Logger,
	configRepo ConfigRepo, userRepo user_repo.UserRepo,
	domainMessageProducer *domain_message.DomainMessageProducer,
) *ConfigUsecase {
	moduleName := "ConfigUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &ConfigUsecase{
		configRepo:            configRepo,
		userRepo:              userRepo,
		log:                   hlog,
		tm:                    tm,
		conf:                  conf,
		domainMessageProducer: domainMessageProducer,
	}
}

func (uc *ConfigUsecase) GetById(ctx context.Context, id int64) (*db.Config, error) {
	return uc.configRepo.GetById(ctx, id)
}

func (uc *ConfigUsecase) GetByKey(ctx context.Context, ConfigKey string) (*db.Config, error) {
	return uc.configRepo.GetByKey(ctx, ConfigKey)
}

func (uc *ConfigUsecase) List(ctx context.Context) ([]*rsp.ViewConfigInfo, error) {
	m, err := uc.configRepo.Map(ctx)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var items []*rsp.ViewConfigInfo
	for _, v := range m {
		items = append(items, &rsp.ViewConfigInfo{
			ConfigKey:   v.ConfigKey,
			ConfigValue: v.ConfigValue,
		})
	}

	// 返回im的平台码
	if uc.conf.TeaIm != nil {
		items = append(items, &rsp.ViewConfigInfo{
			ConfigKey:   consts.CONFIG_BALALA_THIRD_IM_CODE,
			ConfigValue: uc.conf.TeaIm.PlatformCode,
		})
	}

	// ql
	if uc.conf.Ql != nil {
		items = append(items, &rsp.ViewConfigInfo{
			ConfigKey:   consts.CONFIG_BALALA_THIRD_QL_CODE,
			ConfigValue: uc.conf.Ql.PlatformCode,
		})
	}

	// halala
	if uc.conf.Halala != nil {
		items = append(items, &rsp.ViewConfigInfo{
			ConfigKey:   consts.CONFIG_BALALA_THIRD_HALALA_CODE,
			ConfigValue: uc.conf.Halala.PlatformCode,
		})
	}

	return items, nil
}

func (uc *ConfigUsecase) UpdateByKey(ctx context.Context, uid int64, key string, value string) error {
	userMap, err := uc.userRepo.UserMap(ctx, []int64{uid})
	if err != nil {
		return errs.New(ctx, comm.ErrorCode_USER_GET_FAIL)
	}

	operator := userMap[uid]
	if !operator.IsSystemSuperAdmin() {
		return errs.NoPerm(ctx)
	}

	if !slices.Contains(consts.MutableConfigKeyList(), key) {
		return errs.Business(ctx, "invalid key")
	}

	oldConf, err := uc.configRepo.GetByKey(ctx, key)
	if oldConf == nil {
		oldConf = &db.Config{}
	}

	err = uc.configRepo.UpdateByKey(ctx, key, value)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	var msg shared.DomainMessage
	var oldValue = oldConf.ConfigValue
	switch key {
	case consts.CONFIG_BALALA_LOGO:
		msg = &domain_message.AdminChangeSystemLogo{}
	case consts.CONFIG_BALALA_TITLE:
		msg = &domain_message.AdminChangeSystemTitle{
			OldValue: oldValue,
			NewValue: value,
		}
	case consts.CONFIG_BALALA_BG:
		msg = &domain_message.AdminChangeSystemLoginBg{}
	case consts.CONFIG_BALALA_ATTACH:
		msg = &domain_message.AdminChangeSystemAttachSize{
			OldValue: oldValue,
			NewValue: value,
		}
	case consts.CONFIG_BALALA_REGISTER_ENTRY:
		msg = &domain_message.AdminChangeSystemRegisterEntry{
			OldValue: oldValue,
			NewValue: value,
		}
	case consts.CONFIG_NOTIFY_REDIRECT_DOMAIN:
		msg = &domain_message.AdminChangeSystemAccessUrl{
			OldValue: oldValue,
			NewValue: value,
		}
	default:
		return nil
	}

	msg.SetOper(shared.SysOper(uid), time.Now())
	uc.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})
	return nil
}
