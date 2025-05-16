package biz

import (
	"cmp"
	"context"
	"github.com/spf13/cast"
	"go-cs/api/comm"
	v1 "go-cs/api/user/v1"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/consts"
	notify_snapshot_repo "go-cs/internal/domain/notify_snapshot/repo"
	domain_message "go-cs/internal/domain/pkg/message"
	search_repo "go-cs/internal/domain/search/repo"
	space_repo "go-cs/internal/domain/space/repo"
	member "go-cs/internal/domain/space_member"
	member_repo "go-cs/internal/domain/space_member/repo"
	user_domain "go-cs/internal/domain/user"
	user_repo "go-cs/internal/domain/user/repo"
	user_service "go-cs/internal/domain/user/service"
	witem_repo "go-cs/internal/domain/work_item/repo"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/internal/utils/rand"
	"go-cs/internal/utils/third_platform"
	"go-cs/pkg/stream"
	"slices"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type UserUsecase struct {
	tpClient *third_platform.Client
	log      *log.Helper
	tm       trans.Transaction

	repo         user_repo.UserRepo
	loginRepo    LoginRepo
	memberRepo   member_repo.SpaceMemberRepo
	searchRepo   search_repo.SearchRepo
	spaceRepo    space_repo.SpaceRepo
	notifyRepo   notify_snapshot_repo.NotifySnapshotRepo
	adminRepo    AdminRepo
	configRepo   ConfigRepo
	workItemRepo witem_repo.WorkItemRepo

	userService *user_service.UserService

	domainMessageProducer *domain_message.DomainMessageProducer
}

func NewUserUsecase(

	tpClient *third_platform.Client,
	logger log.Logger,
	tm trans.Transaction,

	repo user_repo.UserRepo,
	loginRepo LoginRepo,
	memberRepo member_repo.SpaceMemberRepo,
	searchRepo search_repo.SearchRepo,
	spaceRepo space_repo.SpaceRepo,
	notifyRepo notify_snapshot_repo.NotifySnapshotRepo,
	adminRepo AdminRepo,
	configRepo ConfigRepo,
	workItemRepo witem_repo.WorkItemRepo,

	userService *user_service.UserService,

	domainMessageProducer *domain_message.DomainMessageProducer,

) *UserUsecase {
	moduleName := "UserUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &UserUsecase{
		repo:                  repo,
		loginRepo:             loginRepo,
		memberRepo:            memberRepo,
		spaceRepo:             spaceRepo,
		searchRepo:            searchRepo,
		notifyRepo:            notifyRepo,
		adminRepo:             adminRepo,
		configRepo:            configRepo,
		workItemRepo:          workItemRepo,
		log:                   hlog,
		tm:                    tm,
		tpClient:              tpClient,
		userService:           userService,
		domainMessageProducer: domainMessageProducer,
	}
}

func (uc *UserUsecase) MyInfo(ctx context.Context, userId int64) (*rsp.MyInfo, error) {
	user, err := uc.repo.GetUserByUserId(ctx, userId)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_DB_QUERY_FAIL)
	}

	account, err := uc.repo.GetAllThirdPfAccount(ctx, userId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var bindAccounts []*rsp.MyInfo_ThirdPartyInfo
	for _, v := range account {
		bindAccounts = append(bindAccounts, &rsp.MyInfo_ThirdPartyInfo{
			Id:            v.Id,
			PfCode:        v.PfInfo.PfCode,
			PfName:        v.PfInfo.PfName,
			PfUserName:    v.PfInfo.PfUserName,
			PfUserId:      v.PfInfo.PfUserId,
			PfUserAccount: v.PfInfo.PfUserAccount,
			Notify:        v.Notify,
		})
	}

	config, err := uc.repo.GetUserAllConfig(ctx, userId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	configMap := stream.MapValue(config, func(item *user_domain.UserConfig) string {
		return item.Value
	})

	return &rsp.MyInfo{
		Id:             user.Id,
		UserName:       user.UserName,
		Mobile:         user.Mobile,
		UserNickname:   user.UserNickname,
		UserPinyin:     user.UserPinyin,
		UserStatus:     user.UserStatus,
		UserEmail:      user.UserEmail,
		Sex:            user.Sex,
		Avatar:         user.Avatar,
		LastLoginIp:    user.LastLoginIp,
		LastLoginTime:  user.LastLoginTime,
		CanSetNickName: true,
		BindAccounts:   bindAccounts,
		Role:           int64(user.Role),
		Config:         configMap,
	}, nil
}

func (uc *UserUsecase) CheckUserName(ctx context.Context, userName string) (bool, error) {
	return uc.repo.IsExistByUserName(ctx, userName)
}

func (uc *UserUsecase) Register(ctx context.Context, req *v1.RegUserRequest) (*user_domain.User, error) {
	if req.PfToken == "" && !uc.configRepo.CanRegister(ctx) {
		return nil, errs.Business(ctx, "注册功能已关闭，请联系管理员进行添加")
	}

	var user *user_domain.User
	var err error

	switch req.Way {
	case v1.RegUserRequest_Pwd:
		user, err = uc.RegisterByPWD(ctx, req)
	case v1.RegUserRequest_IM, v1.RegUserRequest_QL, v1.RegUserRequest_Halala:
		user, err = uc.RegisterByIM(ctx, req)
	}

	switch req.Way {
	case v1.RegUserRequest_IM, v1.RegUserRequest_QL, v1.RegUserRequest_Halala:
		platformCode := comm.ThirdPlatformCode(req.Way)

		chatToken := req.PfToken

		client := uc.tpClient.ByPfCode(platformCode)
		if client == nil {
			return nil, errs.Business(ctx, "暂不支持该平台")
		}

		pfUserInfo, err := client.GetUserInfo(chatToken)
		if err != nil {
			uc.log.Error(err)
			return nil, errs.Internal(ctx, err)
		}

		pfAccount := uc.userService.NewPfAccount(ctx, user, user_domain.ThirdPfInfo{
			PfCode:        int32(platformCode),
			PfName:        req.PfName,
			PfUserKey:     req.PfToken,
			PfUserName:    pfUserInfo.NickName,
			PfUserAccount: pfUserInfo.UserName,
			PfUserId:      pfUserInfo.Id,
		})

		err = uc.repo.BindThirdPlatform(ctx, pfAccount)
		if err != nil {
			return nil, errs.Business(ctx, "账号已经被绑定过了")
		}

		utils.Go(func() {
			client.Bind(chatToken)
		})
	}

	err = uc.tm.InTx(ctx, func(ctx context.Context) error {
		err = uc.repo.AddUser(ctx, user)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = uc.repo.InitUserConfig(ctx, user.Id)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUsecase) RegisterByPWD(ctx context.Context, req *v1.RegUserRequest) (*user_domain.User, error) {
	//检查用户名是否已存在
	isExist, err := uc.repo.IsExistByUserName(ctx, req.UserName)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}
	if isExist {
		return nil, errs.New(ctx, comm.ErrorCode_USER_NAME_IS_EXIST)
	}

	user, err := uc.userService.NewUser(ctx, req.UserName, req.NickName, req.Password, req.Avatar, consts.SystemRole_Normal, nil)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	return user, nil
}

func (uc *UserUsecase) RegisterByIM(ctx context.Context, req *v1.RegUserRequest) (*user_domain.User, error) {
	if req.UserName == "" {
		for {
			userName := utils.GenerateUserName("im_")
			isExist, err := uc.repo.IsExistByUserName(ctx, userName)
			if err != nil {
				return nil, errs.Internal(ctx, err)
			}
			if !isExist {
				req.UserName = userName
				break
			}
		}
	}

	if req.Password == "" {
		req.Password = rand.Digits(8)
	}

	if req.NickName == "" {
		req.NickName = req.UserName
		if len(req.NickName) > 14 {
			req.NickName = req.NickName[:14]
		}
	}

	user, err := uc.userService.NewUser(ctx, req.UserName, req.NickName, req.Password, req.Avatar, consts.SystemRole_Normal, nil)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	return user, nil
}

func (uc *UserUsecase) SearchUserList(ctx context.Context, spaceIds []int64, py string, userIds []int64) (*v1.SearchUserListReplyData, error) {

	var list []*rsp.ViewUserWithSpaceInfo
	var err error

	// 转义特殊字符
	py = utils.EscapeSqlSpecialCharacters(py)
	if py == "," { //因为 , 是拼音字段的分隔符，所以不能直接使用
		return nil, nil
	}

	// 获取全部命中的用户
	if len(spaceIds) == 0 {
		list, err = uc.repo.SearchUser(ctx, py, userIds)
	} else {
		list, err = uc.repo.SearchSpaceMember(ctx, py, spaceIds)
	}

	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	list = stream.UniqueFunc(list, func(a *rsp.ViewUserWithSpaceInfo) int64 {
		return a.Id
	})

	slices.SortFunc(list, func(a, b *rsp.ViewUserWithSpaceInfo) int {
		return cmp.Compare(strings.ToLower(a.UserPinyin), strings.ToLower(b.UserPinyin))
	})

	var items []*rsp.SimpleUserInfo
	for _, v := range list {
		items = append(items, &rsp.SimpleUserInfo{
			Id:           v.Id,
			UserId:       v.Id,
			UserName:     v.UserName,
			UserNickname: v.UserNickname,
			Avatar:       v.Avatar,
		})
	}

	return &v1.SearchUserListReplyData{List: items}, nil
}

func (uc *UserUsecase) SetUserAvatar(ctx context.Context, userId int64, avatar string) error {

	user, err := uc.repo.GetUserByUserId(ctx, userId)
	if err != nil {
		return errs.New(ctx, comm.ErrorCode_DB_QUERY_FAIL)
	}

	user.UpdateAvatar(avatar, shared.UserOper(userId))

	err = uc.repo.SaveUser(ctx, user)
	if err != nil {
		return errs.New(ctx, comm.ErrorCode_DB_MODIFY_FAIL)
	}

	uc.domainMessageProducer.Send(ctx, user.GetMessages())

	return nil
}

func (uc *UserUsecase) ChangeUserPwd(ctx context.Context, uid int64, req *v1.ChangeMyPwdRequest) error {
	var user *user_domain.User
	var err error

	if uid != 0 {
		user, err = uc.repo.GetUserByUserId(ctx, uid)
	} else {
		user, err = uc.repo.GetUserByUserName(ctx, req.Username)
	}
	if err != nil {
		return errs.New(ctx, comm.ErrorCode_USER_GET_FAIL)
	}

	uid = user.Id

	// 系统管理员强制要求强密码
	if user.Role >= consts.SystemRole_Admin {
		validate := utils.NewValidator()
		if err := validate.Var(req.NewPwd, "required,utf8Len=5-20,stronger_password"); err != nil {
			return errs.New(ctx, comm.ErrorCode_USER_WRONG_RULE_PASSWORD)
		}
	}

	switch req.Type {
	case v1.ChangeMyPwdRequest_by_old_pwd:
		//占位
	case v1.ChangeMyPwdRequest_by_forceUpdate:
		//强制修改密码
	default:
		return errs.Business(ctx, "不支持的修改密码方式")
	}

	if !user.ValidatePwd(req.Pwd) {
		return errs.Business(ctx, "原密码不正确")
	}

	err = user.ChangePwd(req.Pwd, req.NewPwd, &utils.LoginUserInfo{UserId: user.Id})
	if err != nil {
		return err
	}

	//设置新密码
	err = uc.repo.SaveUser(ctx, user)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	// 清理掉强制更新标记
	if req.Type == v1.ChangeMyPwdRequest_by_forceUpdate {
		uc.adminRepo.DelAdminNeedResetPwdStatus(ctx, uid)
	}

	uc.domainMessageProducer.Send(ctx, user.GetMessages())
	return nil
}

func (uc *UserUsecase) ChangeUserNickName(ctx context.Context, userId int64, nickName string) error {

	user, err := uc.repo.GetUserByUserId(ctx, userId)
	if err != nil {
		return errs.New(ctx, comm.ErrorCode_USER_GET_FAIL)
	}

	err = user.ChangeNickName(nickName, &utils.LoginUserInfo{UserId: userId})
	if err != nil {
		return err
	}

	err = uc.repo.SaveUser(ctx, user)
	if err != nil /**/ {
		return errs.New(ctx, comm.ErrorCode_DB_MODIFY_FAIL)
	}

	uc.domainMessageProducer.Send(ctx, user.GetMessages())

	return nil
}

func (uc *UserUsecase) GetMySpaceMemberInfo(ctx context.Context, userId int64, spaceId int64) (*rsp.SpaceMemberInfo, error) {

	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, userId)
	if err != nil || member == nil {
		return nil, errs.New(ctx, comm.ErrorCode_USER_GET_FAIL)
	}

	userInfo, err := uc.repo.GetUserByUserId(ctx, userId)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_USER_GET_FAIL)
	}

	return &rsp.SpaceMemberInfo{
		Id:           member.Id,
		MemberId:     member.Id,
		SpaceId:      member.SpaceId,
		UserId:       member.UserId,
		RoleId:       int32(member.RoleId),
		UserName:     userInfo.UserName,
		Mobile:       userInfo.Mobile,
		UserNickname: userInfo.UserNickname,
		UserPinyin:   userInfo.UserPinyin,
		UserStatus:   userInfo.UserStatus,
		Avatar:       userInfo.Avatar,
	}, nil
}

func (uc *UserUsecase) Bind(ctx context.Context, uid int64, req *v1.BindRequest) error {
	//  判断绑定的账号是否已经被其他人绑定了
	AlreadyBindErr := errs.Business(ctx, "该账号已被其它三方账号关联")

	switch req.Type {
	default:
		return errs.Business(ctx, "不支持的绑定方式")
	case v1.BindRequest_IM, v1.BindRequest_Ql, v1.BindRequest_Halala:
		_, err := uc.repo.GetThirdPfAccountByPfUserKey(ctx, req.Key, int32(req.Type))
		if err == nil {
			return AlreadyBindErr
		}

		_, err = uc.repo.GetThirdPfAccount(ctx, uid, int32(req.Type))
		if err == nil {
			return AlreadyBindErr
		}

		if err != nil && !errs.IsDbRecordNotFoundErr(err) {
			return errs.Internal(ctx, err)
		}
	}

	user, err := uc.repo.GetUserByUserId(ctx, uid)
	if err != nil {
		return err
	}

	platformCode := comm.ThirdPlatformCode(req.Type)

	switch req.Type {
	case v1.BindRequest_IM, v1.BindRequest_Ql, v1.BindRequest_Halala:
		client := uc.tpClient.ByPfCode(platformCode)
		if client == nil {
			return errs.Business(ctx, "不支持的第三方平台")
		}

		pfUserInfo, err := client.GetUserInfo(req.Key)
		if err != nil {
			uc.log.Error(err)
			return errs.Internal(ctx, err)
		}

		var pfName = req.PfName
		if pfName == "" {
			pfName = consts.GetPlatformName(platformCode)
		}

		newAccount := uc.userService.NewPfAccount(ctx, user, user_domain.ThirdPfInfo{
			PfCode:        int32(platformCode),
			PfName:        pfName,
			PfUserKey:     req.Key,
			PfUserName:    pfUserInfo.NickName,
			PfUserId:      pfUserInfo.Id,
			PfUserAccount: pfUserInfo.UserName,
		})

		err = uc.repo.BindThirdPlatform(ctx, newAccount)
		if err != nil {
			return AlreadyBindErr
		}

		// 通知添加会话
		utils.Go(func() {
			client.Bind(req.Key)
		})
	}

	msg := &domain_message.PersonalBindThirdPlatform{
		DomainMessageBase: shared.DomainMessageBase{
			Oper:     shared.UserOper(user.Id),
			OperTime: time.Now(),
		},
		PlatformName: consts.GetPlatformName(platformCode),
	}
	uc.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})

	return nil
}

func (uc *UserUsecase) Unbind(ctx context.Context, uid int64, req *v1.UnbindRequest) error {
	platformCode := comm.ThirdPlatformCode(req.Type)

	switch req.Type {
	case v1.UnbindRequest_IM, v1.UnbindRequest_Ql, v1.UnbindRequest_Halala:

		pfAccount, err := uc.repo.GetThirdPfAccount(ctx, uid, int32(platformCode))
		if errs.IsDbRecordNotFoundErr(err) {
			return nil
		}
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = uc.repo.UnbindThirdPlatform(ctx, pfAccount)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		client := uc.tpClient.ByPfCode(platformCode)
		if err != nil {
			uc.log.Error(err)
			return nil
		}

		err = client.Unbind(pfAccount.PfInfo.PfUserKey)
		if err != nil {
			uc.log.Error(err)
		}
	default:
		return errs.Business(ctx, "不支持的类型")
	}

	msg := &domain_message.PersonalUnBindThirdPlatform{
		DomainMessageBase: shared.DomainMessageBase{
			Oper:     shared.UserOper(uid),
			OperTime: time.Now(),
		},
		PlatformName: consts.GetPlatformName(platformCode),
	}
	uc.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})

	return nil
}

func (uc *UserUsecase) Cancel(ctx context.Context, uid int64, targetId int64) error {
	userMap, err := uc.repo.UserMap(ctx, []int64{uid, targetId})
	if err != nil {
		return err
	}

	curUser := userMap[uid]
	target := userMap[targetId]

	if target.IsSystemSuperAdmin() {
		return errs.Business(ctx, "无法注销超级管理员")
	}

	if uid != targetId && !curUser.IsSystemAdmin() {
		return errs.NoPerm(ctx)
	}

	list, err := uc.memberRepo.GetUserSpaceIdList(ctx, targetId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if len(list) > 0 {
		return errs.Business(ctx, "无法注销，需退出所有项目")
	}

	pfAccount, _ := uc.repo.GetAllThirdPfAccount(ctx, targetId)

	oldTarget := *target

	target.Cancel()

	err = uc.tm.InTx(ctx, func(ctx context.Context) error {
		err = uc.repo.SaveUser(ctx, target)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		// 清理绑定的第三方账号
		err := uc.repo.RemoveAllThirdPlatform(ctx, targetId)
		if err != nil {
			return errs.Internal(ctx, err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// 解绑所有第三方账号
	if pfAccount != nil {
		utils.Go(func() {
			for _, account := range pfAccount {
				client := uc.tpClient.ByPfCode(comm.ThirdPlatformCode(account.PfInfo.PfCode))
				if client != nil {
					client.Unbind(account.PfInfo.PfUserKey)
				}
			}
		})
	}

	// 注销后清理登录token
	uc.loginRepo.ClearAllJwtToken(ctx, targetId)

	// 系统管理员注销日志
	if uid != targetId {
		msg := &domain_message.AdminCancelUser{
			DomainMessageBase: shared.DomainMessageBase{
				Oper:     shared.SysOper(uid),
				OperTime: time.Now(),
			},
			UserId:   targetId,
			Username: oldTarget.UserName,
			Nickname: oldTarget.UserNickname,
		}

		uc.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})
	}

	return err
}

func (uc *UserUsecase) SetSpaceNotify(ctx context.Context, uid int64, req *v1.SetSpaceNotifyRequest) error {
	_, err := uc.spaceRepo.GetSpace(ctx, req.SpaceId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	member, err := uc.memberRepo.GetSpaceMember(ctx, req.SpaceId, uid)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	member.SetNotify(req.Notify, shared.UserOper(uid))

	err = uc.memberRepo.SaveSpaceMember(ctx, member)
	if err != nil {
		uc.log.Error(err)
	}

	uc.domainMessageProducer.Send(ctx, member.GetMessages())

	return nil
}

func (uc *UserUsecase) SetSpaceOrder(ctx context.Context, uid int64, req *v1.SetSpaceOrderRequest) error {
	err := uc.memberRepo.UpdateUserSpaceOrder(ctx, uid, req.FromIdx, req.ToIdx)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	return nil
}

func (uc *UserUsecase) MyPendingWorkItem(ctx context.Context, uid int64) (*v1.MyPendingWorkItemReplyData, error) {
	spaceIds, err := uc.spaceRepo.GetUserSpaceIds(ctx, uid)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	items, err := uc.searchRepo.PendingWorkItem(ctx, uid, spaceIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	return &v1.MyPendingWorkItemReplyData{Items: items}, nil
}

func (uc *UserUsecase) MyRelatedComment(ctx context.Context, uid int64, req *v1.MyRelatedCommentRequest) (*v1.MyRelatedCommentReplyData, error) {
	items, nextPos, hasNext, err := uc.notifyRepo.GetUserRelatedCommentIds(ctx, uid, int(req.Pos), int(req.Size))
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	return &v1.MyRelatedCommentReplyData{
		Items:   items,
		NextPos: nextPos,
		HasNext: hasNext,
	}, nil
}

func (uc *UserUsecase) MyRelatedCommentByIds(ctx context.Context, uid int64, ids []int64) (*v1.MyRelatedCommentByIdsReplyData, error) {
	list, err := uc.notifyRepo.GetNotifyByIds(ctx, uid, ids)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var items []*v1.MyRelatedCommentByIdsReplyData_Item
	for _, v := range list {
		items = append(items, &v1.MyRelatedCommentByIdsReplyData_Item{
			Id:  v.Id,
			Doc: v.Doc,
		})
	}

	return &v1.MyRelatedCommentByIdsReplyData{
		Items: items,
	}, nil
}

func (uc *UserUsecase) NotifyCount(ctx context.Context, uid int64) (*v1.NotifyCountReplyData, error) {
	spaceIds, err := uc.spaceRepo.GetUserSpaceIds(ctx, uid)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	items, err := uc.searchRepo.PendingWorkItem(ctx, uid, spaceIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	pending := int64(len(items))

	return &v1.NotifyCountReplyData{
		Pending: pending,
	}, nil
}

func (uc *UserUsecase) GetTempConfig(ctx context.Context, uid int64, keys []string) (map[string]string, error) {
	config := uc.repo.GetTempConfig(ctx, uid, keys...)
	return config, nil
}

func (uc *UserUsecase) SetTempConfig(ctx context.Context, uid int64, configs map[string]string) error {
	oldConfigs, _ := uc.GetTempConfig(ctx, uid, stream.Keys(configs))

	err := uc.repo.SetTempConfig(ctx, uid, configs)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	// 删除没有变化的
	for key, newVal := range configs {
		oldVal := oldConfigs[key]
		if oldVal == newVal {
			delete(configs, key)
			delete(oldConfigs, key)
		}
	}

	msg := &domain_message.PersonalSetTempConfig{
		DomainMessageBase: shared.DomainMessageBase{
			Oper:     utils.GetLoginUser(ctx),
			OperTime: time.Now(),
		},
		OldValues: oldConfigs,
		NewValues: configs,
	}

	uc.domainMessageProducer.Send(ctx, []shared.DomainMessage{
		msg,
	})

	return nil
}

func (uc *UserUsecase) DelTempConfig(ctx context.Context, uid int64, keys []string) error {

	err := uc.repo.DelTempConfig(ctx, uid, keys...)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	return nil
}

func (uc *UserUsecase) AllSpaceInfo(ctx context.Context, uid int64, userId int64) (*v1.AllSpaceProfileReply_Data, error) {

	members, err := uc.memberRepo.GetUserAllSpaceMember(ctx, userId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}
	userMemberMap := stream.ToMap(members, func(i int, t *member.SpaceMember) (int64, *member.SpaceMember) {
		return t.SpaceId, t
	})

	spaceIds := stream.Map(members, func(item *member.SpaceMember) int64 {
		return item.SpaceId
	})
	spaceMap, err := uc.spaceRepo.SpaceMap(ctx, spaceIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	spaceMemberMap, err := uc.memberRepo.SpaceMemberMapBySpaceIds(ctx, spaceIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	countMap, err := uc.workItemRepo.CountUserRelatedSpaceWorkItemBySpaceIds(ctx, userId, spaceIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var list []*v1.AllSpaceProfileReply_Data_SpaceProfile

	for _, v := range members {
		space := spaceMap[v.SpaceId]
		if space == nil {
			continue
		}

		userMember := userMemberMap[v.SpaceId]
		if userMember == nil {
			continue
		}

		list = append(list, &v1.AllSpaceProfileReply_Data_SpaceProfile{
			SpaceId:            space.Id,
			SpaceName:          space.SpaceName,
			RelatedWorkItemNum: countMap[v.SpaceId],
			SpaceRole:          userMember.RoleId,
			MemberNum:          int64(len(spaceMemberMap[v.SpaceId])),
		})
	}

	slices.SortFunc(list, func(a, b *v1.AllSpaceProfileReply_Data_SpaceProfile) int {
		if v := cmp.Compare(a.SpaceRole, b.SpaceRole); v != 0 {
			return -v
		}

		return cmp.Compare(a.SpaceId, b.SpaceId)
	})

	return &v1.AllSpaceProfileReply_Data{
		List: list,
	}, nil
}

func (uc *UserUsecase) SetThirdPlatformNotify(ctx context.Context, uid int64, req *v1.SetThirdPlatformNotifyRequest) error {
	account, err := uc.repo.GetThirdPfAccount(ctx, uid, int32(req.PlatformCode))
	if err != nil {
		return errs.Internal(ctx, err)
	}

	account.SetNotify(req.Notify, shared.UserOper(uid))

	err = uc.repo.SaveThirdPfAccount(ctx, account)
	if err != nil {
		uc.log.Error(err)
	}

	uc.domainMessageProducer.Send(ctx, account.GetMessages())

	return nil
}

func (uc *UserUsecase) SetUserConfig(ctx context.Context, uid int64, req *v1.SetUserConfigRequest) error {
	config, err := uc.repo.GetUserConfig(ctx, uid, req.Key)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if config == nil {
		return errs.Business(ctx, "没有该配置")
	}

	switch req.Key {
	case consts.UserConfigKey_NotifySwitchGlobal, consts.UserConfigKey_NotifySwitchThirdPlatform, consts.UserConfigKey_NotifySwitchSpace:
		req.Value = cast.ToString(cast.ToInt(req.Value)) // 保证是数字
	default:
		return errs.Business(ctx, "没有该配置")
	}

	config.SetValue(req.Value, shared.UserOper(uid))

	err = uc.repo.SaveUserConfig(ctx, config)
	if err != nil {
		uc.log.Error(err)
	}

	uc.domainMessageProducer.Send(ctx, config.GetMessages())

	return nil
}
