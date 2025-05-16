package biz

import (
	"cmp"
	"context"
	"fmt"
	"github.com/spf13/cast"
	"go-cs/api/comm"
	"go-cs/api/notify"
	pb "go-cs/api/space_member/v1"
	db "go-cs/internal/bean/biz"
	vo "go-cs/internal/bean/vo"
	"go-cs/internal/bean/vo/event"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/consts"
	"go-cs/internal/domain/perm/facade"
	perm_service "go-cs/internal/domain/perm/service"
	domain_message "go-cs/internal/domain/pkg/message"
	space_repo "go-cs/internal/domain/space/repo"
	member "go-cs/internal/domain/space_member"
	member_repo "go-cs/internal/domain/space_member/repo"
	member_service "go-cs/internal/domain/space_member/service"
	space_view_repo "go-cs/internal/domain/space_view/repo"
	view_service "go-cs/internal/domain/space_view/service"
	comment_repo "go-cs/internal/domain/space_work_item_comment/repo"
	user_repo "go-cs/internal/domain/user/repo"
	flow_repo "go-cs/internal/domain/work_flow/repo"
	flow_service "go-cs/internal/domain/work_flow/service"
	witem_repo "go-cs/internal/domain/work_item/repo"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/pkg/bus"
	"go-cs/pkg/stream"
	"slices"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
)

type SpaceMemberUsecase struct {
	repo              member_repo.SpaceMemberRepo
	userRepo          user_repo.UserRepo
	spaceRepo         space_repo.SpaceRepo
	spaceWorkItemRepo witem_repo.WorkItemRepo
	commentRepo       comment_repo.SpaceWorkItemCommentRepo
	flowRepo          flow_repo.WorkFlowRepo
	viewRepo          space_view_repo.SpaceViewRepo

	permService   *perm_service.PermService
	memberService *member_service.SpaceMemberService
	flowService   *flow_service.WorkFlowService
	viewService   *view_service.SpaceViewService

	domainMessageProducer *domain_message.DomainMessageProducer

	log *log.Helper
	tm  trans.Transaction
}

func NewSpaceMemberUsecase(
	repo member_repo.SpaceMemberRepo,
	spaceRepo space_repo.SpaceRepo,
	userRepo user_repo.UserRepo,
	spaceWorkItemRepo witem_repo.WorkItemRepo,
	commentRepo comment_repo.SpaceWorkItemCommentRepo,
	flowRepo flow_repo.WorkFlowRepo,
	viewRepo space_view_repo.SpaceViewRepo,

	permService *perm_service.PermService,
	memberService *member_service.SpaceMemberService,
	flowService *flow_service.WorkFlowService,
	viewService *view_service.SpaceViewService,

	domainMessageProducer *domain_message.DomainMessageProducer,

	tm trans.Transaction,
	logger log.Logger,
) *SpaceMemberUsecase {
	moduleName := "SpaceMemberUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceMemberUsecase{
		tm:                tm,
		repo:              repo,
		userRepo:          userRepo,
		spaceRepo:         spaceRepo,
		spaceWorkItemRepo: spaceWorkItemRepo,
		commentRepo:       commentRepo,
		flowRepo:          flowRepo,
		viewRepo:          viewRepo,

		log: hlog,

		permService:   permService,
		memberService: memberService,
		flowService:   flowService,
		viewService:   viewService,

		domainMessageProducer: domainMessageProducer,
	}
}

func (s *SpaceMemberUsecase) AddMySpaceMembers(ctx context.Context, oper *utils.LoginUserInfo, in *vo.AddSpaceMembersVo) ([]*db.SpaceMember, error) {

	space, err := s.spaceRepo.GetSpace(ctx, in.SpaceId)
	if err != nil {
		//空间信息不存在
		err := errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
		return nil, err
	}

	//查看操作人是不是这个空间的成员，并且是否有相关添加成员的权限
	memberRelation, err := s.repo.GetSpaceMember(ctx, in.SpaceId, oper.UserId)
	if err != nil {
		//用户不是当前空间的成员，无权操作
		err := errs.New(ctx, comm.ErrorCode_SPACE_MEMBER_WRONG)
		return nil, err
	}

	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: facade.BuildSpaceMemberFacade(memberRelation),
		Perm:              consts.PERM_ADD_SPACE_MEMBER,
	})

	if err != nil {
		return nil, err
	}

	userMap, err := s.userRepo.UserMap(ctx, in.MemberUserIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var members member.SpaceMembers
	for i, id := range in.MemberUserIds {
		u := userMap[id]
		if u == nil {
			continue
		}

		newMember, err := s.memberService.NewSpaceMember(ctx, space.Id, u.Id, in.MemberUserRoleIds[i], oper)
		if err != nil {
			//忽略掉不符合条件的
			continue
		}
		members = append(members, newMember)
	}

	if len(members) == 0 {
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_MEMBER_ADD_FAIL)
	}

	err = s.tm.InTx(ctx, func(ctx context.Context) error {
		//添加成员到空间中
		err = s.repo.AddSpaceMembers(ctx, members)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		// 初始化全局视图
		userIds := stream.Map(members, func(m *member.SpaceMember) int64 {
			return m.UserId
		})
		userViews, err := s.viewService.InitUserGlobalView(ctx, space.Id, userIds)
		if err != nil {
			return errs.Internal(ctx, err)
		}
		err = s.viewRepo.CreateUserViews(ctx, userViews)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	s.domainMessageProducer.Send(ctx, members.GetMessages())

	//设置评论未读数
	numMap, _ := s.spaceWorkItemRepo.AllCommentNumMap(ctx, space.Id)
	for _, v := range members {
		s.commentRepo.SetUserUnreadNum(ctx, v.UserId, numMap)
	}

	// for _, member := range members {
	// 	bus.Emit(notify.Event_AddMember, &event.AddMember{
	// 		Event:    notify.Event_AddMember,
	// 		Operator: userId,
	// 		Space:    convert.SpaceEntityToPo(space),
	// 		TargetId: member.UserId,
	// 		RoleId:   member.RoleId,
	// 	})
	// }

	return nil, nil
}

func (s *SpaceMemberUsecase) GetMySpaceMemberList(ctx context.Context, userId int64, spaceId int64, userName string) ([]*rsp.SpaceMemberInfo, error) {
	user, _ := s.userRepo.GetUserByUserId(ctx, userId)

	//判断当前用户是否在要查询的项目空间内
	member, err := s.repo.GetSpaceMember(ctx, spaceId, userId)
	if !user.IsSystemAdmin() && (member == nil || err != nil) {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.New(ctx, comm.ErrorCode_PERMISSION_INSUFFICIENT_PERMISSIONS)
		return nil, err
	}

	list, err := s.repo.QSpaceMemberList(ctx, spaceId, userName)
	if err != nil {
		err := errs.New(ctx, comm.ErrorCode_DB_QUERY_FAIL)
		return nil, err
	}

	slices.SortFunc(list, func(a, b *rsp.SpaceMemberInfo) int {
		if a.RoleId != b.RoleId {
			return cmp.Compare(consts.GetMemberRoleRank(a.RoleId), consts.GetMemberRoleRank(b.RoleId))
		}

		return cmp.Compare(strings.ToLower(a.UserPinyin), strings.ToLower(b.UserPinyin))
	})

	return list, nil
}

func (s *SpaceMemberUsecase) SetMySpaceMemberRoleId(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, memberUserId int64, roleId int) error {

	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		//空间信息不存在
		err := errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
		return err
	}

	//判断当前用户是否在要查询的项目空间内
	uMember, err := s.repo.GetSpaceMember(ctx, space.Id, oper.UserId)
	if err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		return errs.NoPerm(ctx)
	}

	//不能设置自己
	if uMember.IsCurrentMemberUser(memberUserId) {
		return errs.NoPerm(ctx)
	}

	member, err := s.repo.GetSpaceMember(ctx, spaceId, memberUserId)
	if err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.New(ctx, comm.ErrorCode_SPACE_MEMBER_WRONG)
		return err
	}

	// 验证权限
	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: facade.BuildSpaceMemberFacade(uMember),
		Perm:              consts.PERM_SET_SPACE_MEMBER_ROLE,
	})
	if err != nil {
		return err
	}

	//是否拥有可以操作对方的权限
	err = s.permService.CheckRoleOperatePerm(ctx, uMember.GetRole(), member.GetRole())
	if err != nil {
		return err
	}

	//验证是否有切换成员至管理员的权限
	if roleId == consts.MEMBER_ROLE_MANAGER {
		err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
			SpaceMemberFacade: facade.BuildSpaceMemberFacade(uMember),
			Perm:              consts.PERM_SET_SPACE_MEMBER_ROLE_TO_MANAGER,
		})
		if err != nil {
			return err
		}
	}

	member.ChangeRoleId(int64(roleId), oper)

	err2 := s.repo.SaveSpaceMember(ctx, member)
	if err2 != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		return err2
	}

	s.domainMessageProducer.Send(ctx, member.GetMessages())

	return nil
}

const (
	QuitSpaceBySelf         = iota // 主动离开项目
	QuitSpaceBySystemAdmin         // 被系统管理员操作离开项目
	QuitSpaceBySpaceManager        // 被项目管理员操作离开项目
)

func (s *SpaceMemberUsecase) QuitMySpaceV2(ctx context.Context, uid, spaceId int64, userId, targetUserId int64) error {
	scene := QuitSpaceBySelf
	if uid != userId {
		scene = QuitSpaceBySystemAdmin
	}

	return s.removeMySpaceMember(ctx, uid, scene, spaceId, userId, targetUserId)
}

func (s *SpaceMemberUsecase) KickOut(ctx context.Context, uid, spaceId int64, userId, targetUserId int64) error {
	scene := QuitSpaceBySpaceManager
	return s.removeMySpaceMember(ctx, uid, scene, spaceId, userId, targetUserId)
}

func (s *SpaceMemberUsecase) removeMySpaceMember(ctx context.Context, uid int64, scene int, spaceId int64, oldUserId, newUserId int64) error {
	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		return errs.Custom(ctx, comm.ErrorCode_SPACE_NOT_EXIST, "项目不存在")
	}

	switch scene {
	case QuitSpaceBySelf:
		if space.IsCreator(oldUserId) {
			return errs.Business(ctx, "项目创建者不能离开自己的项目")
		}
	case QuitSpaceBySystemAdmin:
		userMap, _ := s.userRepo.UserMap(ctx, []int64{uid, oldUserId})

		curUser := userMap[uid]
		oldUser := userMap[oldUserId]

		if !curUser.IsSystemAdmin() || !curUser.RoleGreaterThan(oldUser.Role) {
			return errs.NoPerm(ctx)
		}

	case QuitSpaceBySpaceManager:
		if uid == oldUserId {
			return errs.NoPerm(ctx)
		}

		// 查看当前用户是否为空间成员
		memberMap, err := s.repo.SpaceMemberMapByUserIds(ctx, spaceId, []int64{uid, oldUserId})
		if err != nil {
			return errs.Internal(ctx, err)
		}

		curUserRelation := memberMap[uid]
		oldUserRelation := memberMap[oldUserId]

		// 判断当前用户项目权限是否大于待移除用户
		if !curUserRelation.RoleGreaterThan(oldUserRelation) {
			return errs.NoPerm(ctx)
		}

		// 不能操作自己 || 验证权限
		err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
			SpaceMemberFacade: facade.BuildSpaceMemberFacade(curUserRelation),
			Perm:              consts.PERM_REMOVE_SPACE_MEMBER,
		})
		if err != nil {
			return errs.NoPerm(ctx)
		}
	default:
		return errs.NoPerm(ctx)
	}

	// 检查任务
	workItemIds, err := s.spaceWorkItemRepo.GetSpaceWorkItemIdsByParticipators(ctx, spaceId, oldUserId)
	if err != nil {
		return errs.Internal(ctx, err)
	}
	if len(workItemIds) != 0 && newUserId == 0 {
		return errs.Business(ctx, "存在未转交的任务")
	}

	// 检查流程
	templates, err := s.flowRepo.SearchHistoryTaskWorkFlowTemplateByOwnerRule(ctx, spaceId, cast.ToString(oldUserId))
	if err != nil {
		return errs.Internal(ctx, err)
	}
	// 替换流程中的人员
	for _, template := range templates {
		template.RemoveOwner(oldUserId)
	}

	// 获取已关注的任务
	followedIds, err := s.spaceWorkItemRepo.GetWorkItemIdsByFollower(ctx, oldUserId, spaceId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	// 获取有评论的任务id
	allCommentWorkItemIds, err := s.spaceWorkItemRepo.GetSpaceAllWorkItemIdsHasComment(ctx, spaceId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	removeMember, err := s.repo.GetSpaceMember(ctx, spaceId, oldUserId)
	if err != nil { // 用户不是当前空间的成员，无权操作
		return errs.Business(ctx, "要移除的成员不存在")
	}

	var targetSpaceMember *member.SpaceMember
	if newUserId != 0 {
		targetSpaceMember, err = s.repo.GetSpaceMember(ctx, spaceId, newUserId)
		if err != nil {
			return errs.Internal(ctx, err)
		}
	}

	switch scene {
	case QuitSpaceBySelf:
		removeMember.OnQuit(len(workItemIds), newUserId, shared.UserOper(uid))
	case QuitSpaceBySystemAdmin:
		removeMember.OnQuit(len(workItemIds), newUserId, shared.UserOper(oldUserId))
	case QuitSpaceBySpaceManager:
		removeMember.OnKickOut(len(workItemIds), newUserId, shared.UserOper(uid))
	default:
		return errs.NoPerm(ctx)
	}

	//持久化
	err = s.tm.InTx(ctx, func(ctx context.Context) error {

		err = s.spaceRepo.SaveSpace(ctx, space)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		// 移除成员
		err = s.repo.DelSpaceMember(ctx, removeMember.SpaceId, removeMember.UserId)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		//移除关注
		err = s.spaceWorkItemRepo.Unfollow(ctx, oldUserId, followedIds)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		// 清理评论未读数
		if len(allCommentWorkItemIds) > 0 {
			err = s.commentRepo.RemoveUnreadNumForUser(ctx, oldUserId, allCommentWorkItemIds)
			if err != nil {
				return errs.Internal(ctx, err)
			}
		}

		// 替换流程模版人员
		for _, template := range templates {
			err = s.flowRepo.SaveWorkFlowTemplate(ctx, template)
			if err != nil {
				return errs.Internal(ctx, err)
			}
		}

		// 清理个人视图
		err = s.viewRepo.DeleteUserViewByUserId(ctx, oldUserId, spaceId)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		//--负责人关联转移
		if targetSpaceMember != nil {

			// 修改任务创建人
			_, err = s.spaceWorkItemRepo.UpdateSpaceAllWorkItemCreator(ctx, spaceId, oldUserId, newUserId)
			if err != nil {
				return errs.Internal(ctx, err)
			}

			// 当前负责人调整为目标负责人
			_, err = s.spaceWorkItemRepo.ReplaceDirectorForWorkItemBySpace(ctx, spaceId, oldUserId, newUserId)
			if err != nil {
				return errs.Internal(ctx, err)
			}

			// 任务参与人调整为目标负责人
			_, err = s.spaceWorkItemRepo.ReplaceParticipatorsForWorkItemBySpace(ctx, spaceId, oldUserId, newUserId)
			if err != nil {
				return errs.Internal(ctx, err)
			}

			// 节点负责人调整为目标负责人
			_, err = s.spaceWorkItemRepo.ReplaceNodeDirectorsForWorkItemBySpace(ctx, spaceId, oldUserId, newUserId)
			if err != nil {
				return errs.Internal(ctx, err)
			}

			// 节点负责人调整为目标负责人
			_, err = s.spaceWorkItemRepo.ReplaceDirectorForWorkItemFlowBySpace(ctx, spaceId, oldUserId, newUserId)
			if err != nil {
				return errs.Internal(ctx, err)
			}

			// 节点角色负责人调整为目标负责人
			_, err = s.spaceWorkItemRepo.ReplaceDirectorForWorkItemFlowRolesBySpace(ctx, spaceId, oldUserId, newUserId)
			if err != nil {
				return errs.Internal(ctx, err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 清理空间配置
	keys := []string{
		fmt.Sprintf("condition_%v", spaceId),
	}
	err = s.userRepo.DelTempConfig(ctx, oldUserId, keys...)

	s.domainMessageProducer.Send(ctx, removeMember.GetMessages())
	s.domainMessageProducer.Send(ctx, space.GetMessages())

	//这里的推送只能特殊处理了
	if newUserId != 0 {
		bus.Emit(notify.Event_TransferWorkItem, &event.TransferWorkItem{
			Event:    notify.Event_TransferWorkItem,
			Space:    space,
			Operator: oldUserId,
			TargetId: newUserId,
			Num:      len(workItemIds),
		})
	}

	return err
}

func (s *SpaceMemberUsecase) GetSpaceMemberWorkItemCountV2(ctx context.Context, uid int64, memberUserId int64, spaceId int64) (int64, error) {

	//判断当前用户是否在要查询的项目空间内
	member, err := s.repo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		s.log.Error(err)
		return 0, errs.New(ctx, comm.ErrorCode_PERMISSION_INSUFFICIENT_PERMISSIONS)
	}

	member2, err := s.repo.GetSpaceMember(ctx, spaceId, memberUserId)
	if member2 == nil || err != nil {
		return 0, errs.New(ctx, comm.ErrorCode_SPACE_MEMBER_WRONG)
	}

	counts, err := s.spaceWorkItemRepo.CountUserRelatedSpaceWorkItem(ctx, spaceId, memberUserId)
	if err != nil {
		return 0, errs.New(ctx, comm.ErrorCode_PERMISSION_INSUFFICIENT_PERMISSIONS)
	}

	return counts, nil
}

func (s *SpaceMemberUsecase) GetMySpaceManagerList(ctx context.Context, userId int64, spaceId int64) ([]*rsp.SpaceMemberInfo, error) {

	//判断当前用户是否在要查询的项目空间内
	member, err := s.repo.GetSpaceMember(ctx, spaceId, userId)
	if member == nil || err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.NoPerm(ctx)
		return nil, err
	}

	list, err := s.repo.QSpaceManagerList(ctx, spaceId)
	if err != nil {
		err := errs.New(ctx, comm.ErrorCode_DB_QUERY_FAIL)
		return nil, err
	}

	slices.SortFunc(list, func(a, b *rsp.SpaceMemberInfo) int {
		if a.RoleId == consts.MEMBER_ROLE_SPACE_CREATOR {
			return -1
		}

		if b.RoleId == consts.MEMBER_ROLE_SPACE_CREATOR {
			return 1
		}

		if a.RoleId == consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER {
			return -1
		}

		if b.RoleId == consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER {
			return 1
		}

		if a.RoleId != b.RoleId {
			return cmp.Compare(a.RoleId, b.RoleId)
		}

		return cmp.Compare(strings.ToLower(a.UserPinyin), strings.ToLower(b.UserPinyin))
	})

	return list, nil
}

func (s *SpaceMemberUsecase) AddMySpaceManager(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, managerUids []int64) ([]*db.SpaceMember, error) {

	managerUids = stream.Unique(managerUids)
	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		err := errs.NoPerm(ctx)
		return nil, err
	}

	//只有创建者才能添加管理员
	//查看操作人是不是这个空间的成员，并且是否有相关添加成员的权限
	operMember, err := s.repo.GetSpaceMember(ctx, space.Id, oper.UserId)
	if err != nil {
		err := errs.NoPerm(ctx)
		return nil, err
	}

	if operMember.RoleId != consts.MEMBER_ROLE_SPACE_CREATOR {
		err := errs.NoPerm(ctx)
		return nil, err
	}

	// 检查是否已经是成员
	allIsMember, err := s.repo.AllIsMember(ctx, spaceId, managerUids...)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	if !allIsMember {
		return nil, errs.Business(ctx, "要添加的成员不属于该项目")
	}

	allMembers, err := s.repo.GetSpaceMemberByUserIds(ctx, spaceId, managerUids)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	//过滤出管理员
	normalMembers := stream.Filter(allMembers, func(e *member.SpaceMember) bool {
		return !e.IsSpaceSupperManager() && !e.IsSpaceCreator()
	})

	for _, v := range normalMembers {
		v.SetSpaceManagerRole(oper)
	}

	//把对应的成员设置成管理员
	err = s.tm.InTx(ctx, func(ctx context.Context) (txErr error) {
		for _, v := range normalMembers {
			txErr = s.repo.SaveSpaceMember(ctx, v)
			if txErr != nil {
				return txErr
			}
		}
		return nil
	})

	if err != nil {
		return nil, errs.Business(ctx, "添加失败")
	}

	for _, v := range normalMembers {
		s.domainMessageProducer.Send(ctx, v.GetMessages())
	}

	return nil, nil
}

func (s *SpaceMemberUsecase) RemoveMySpaceManager(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, managerUids []int64) ([]*db.SpaceMember, error) {

	managerUids = stream.Unique(managerUids)

	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		err := errs.NoPerm(ctx)
		return nil, err
	}

	operMember, err := s.repo.GetSpaceMember(ctx, space.Id, oper.UserId)
	if err != nil {
		err := errs.NoPerm(ctx)
		return nil, err
	}

	//只有创建者才能添加管理员
	if !operMember.IsSpaceCreator() {
		err := errs.NoPerm(ctx)
		return nil, err
	}

	if slices.Contains(managerUids, oper.UserId) {
		return nil, errs.Business(ctx, "创建者不能移除自己")
	}

	allMembers, err := s.repo.GetSpaceMemberByUserIds(ctx, spaceId, managerUids)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	//过滤出管理员
	managerMembers := stream.Filter(allMembers, func(e *member.SpaceMember) bool {
		return e.IsSpaceSupperManager()
	})

	if len(managerMembers) == 0 {
		return nil, errs.Business(ctx, "要移除的管理员不属于该项目")
	}

	for _, v := range managerMembers {
		v.RemoveSupperManagerRole(oper)
	}

	//把对应的成员设置成管理员
	err = s.tm.InTx(ctx, func(ctx context.Context) (txErr error) {
		for _, v := range managerMembers {
			txErr = s.repo.SaveSpaceMember(ctx, v)
			if txErr != nil {
				return txErr
			}
		}
		return nil
	})

	if err != nil {
		return nil, errs.Business(ctx, "移除失败")
	}

	for _, v := range managerMembers {
		s.domainMessageProducer.Send(ctx, v.GetMessages())
	}

	return nil, nil
}

func (s *SpaceMemberUsecase) QSpaceMemberByIds(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, uids []int64) (*pb.SpaceMemberByIdReplyData, error) {
	//判断当前用户是否在要查询的项目空间内
	member, err := s.repo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if member == nil || err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.NoPerm(ctx)
		return nil, err
	}

	result := &pb.SpaceMemberByIdReplyData{
		List: make([]*rsp.SpaceMemberInfo, 0),
	}

	uids = stream.Unique(uids)
	if len(uids) == 0 {
		return result, nil
	}

	list, err := s.repo.QSpaceMemberByUids(ctx, spaceId, uids)
	if err != nil {
		err := errs.Business(ctx, "查询失败")
		return nil, err
	}

	result.List = append(result.List, list...)
	return result, nil
}
