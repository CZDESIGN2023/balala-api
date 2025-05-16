package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"go-cs/api/comm"
	"go-cs/api/notify"
	"go-cs/internal/bean/vo/message"
	msg "go-cs/internal/bean/vo/message/notify_message"
	"go-cs/internal/biz"
	"go-cs/internal/consts"
	notify_snapshot_repo "go-cs/internal/domain/notify_snapshot/repo"
	notify_snapshot_service "go-cs/internal/domain/notify_snapshot/service"
	domain_message "go-cs/internal/domain/pkg/message"
	comment_repo "go-cs/internal/domain/space_work_item_comment/repo"
	"go-cs/internal/domain/user"
	"go-cs/internal/domain/work_item"
	wiz_repo "go-cs/internal/domain/work_item_status/repo"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/server/websock"
	"go-cs/internal/utils"
	"go-cs/internal/utils/third_platform"
	"go-cs/internal/utils/third_platform/raw_sdk/tea_im"
	"go-cs/pkg/bus"
	"go-cs/pkg/stream"
	"reflect"
	"runtime/debug"

	space_repo "go-cs/internal/domain/space/repo"
	member_repo "go-cs/internal/domain/space_member/repo"
	tag_repo "go-cs/internal/domain/space_tag/repo"
	workObj_repo "go-cs/internal/domain/space_work_object/repo"
	workVersion_repo "go-cs/internal/domain/space_work_version/repo"
	user_repo "go-cs/internal/domain/user/repo"
	work_flow_repo "go-cs/internal/domain/work_flow/repo"
	witem_repo "go-cs/internal/domain/work_item/repo"

	"github.com/go-kratos/kratos/v2/log"
)

type Notify struct {
	log                *log.Helper
	workItemRepo       witem_repo.WorkItemRepo
	userRepo           user_repo.UserRepo
	spaceRepo          space_repo.SpaceRepo
	memberRepo         member_repo.SpaceMemberRepo
	tagRepo            tag_repo.SpaceTagRepo
	workObjectRepo     workObj_repo.SpaceWorkObjectRepo
	notifySnapShotRepo notify_snapshot_repo.NotifySnapshotRepo
	configRepo         biz.ConfigRepo
	workVersionRepo    workVersion_repo.SpaceWorkVersionRepo
	workItemStatusRepo wiz_repo.WorkItemStatusRepo
	workFlowRepo       work_flow_repo.WorkFlowRepo
	commentRepo        comment_repo.SpaceWorkItemCommentRepo

	notifySnapShotService *notify_snapshot_service.NotifySnapShotService

	domainMessageConsumer *domain_message.DomainMessageConsumer

	ws       *websock.Server
	tpClient *third_platform.Client

	domainMessageRouter map[shared.MessageType]any
}

func NewNotify(
	spaceRepo space_repo.SpaceRepo,
	userRepo user_repo.UserRepo,
	workVersionRepo workVersion_repo.SpaceWorkVersionRepo,
	workItemRepo witem_repo.WorkItemRepo,
	memberRepo member_repo.SpaceMemberRepo,
	tagRepo tag_repo.SpaceTagRepo,
	notifySnapShotRepo notify_snapshot_repo.NotifySnapshotRepo,
	workObjectRepo workObj_repo.SpaceWorkObjectRepo,
	configRepo biz.ConfigRepo,
	workItemStatusRepo wiz_repo.WorkItemStatusRepo,
	workFlowRepo work_flow_repo.WorkFlowRepo,
	commentRepo comment_repo.SpaceWorkItemCommentRepo,

	notifySnapShotService *notify_snapshot_service.NotifySnapShotService,

	domainMessageConsumer *domain_message.DomainMessageConsumer,

	ws *websock.Server, tpClient *third_platform.Client, logger log.Logger,
) *Notify {

	moduleName := "Notify"
	_, l := utils.InitModuleLogger(logger, moduleName)

	n := &Notify{

		userRepo:           userRepo,
		spaceRepo:          spaceRepo,
		workItemRepo:       workItemRepo,
		memberRepo:         memberRepo,
		tagRepo:            tagRepo,
		notifySnapShotRepo: notifySnapShotRepo,
		workObjectRepo:     workObjectRepo,
		configRepo:         configRepo,
		workVersionRepo:    workVersionRepo,
		workItemStatusRepo: workItemStatusRepo,
		workFlowRepo:       workFlowRepo,
		commentRepo:        commentRepo,

		notifySnapShotService: notifySnapShotService,

		domainMessageConsumer: domainMessageConsumer,

		ws:       ws,
		tpClient: tpClient,
		log:      l,

		domainMessageRouter: map[shared.MessageType]any{},
	}

	return n
}

func (s *Notify) Init() {
	s.initHandler()
	s.initInternalHandler()
	s.initDomainMessageHandler()
}

func (s *Notify) busOn(event any, id bus.ID, handler any) {
	bus.On(event, id, func(args ...any) {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					s.log.Errorf("handle internal event %s panic: %v, %s", event, err, debug.Stack())
				}
			}()

			var rvs []reflect.Value
			for _, v := range args {
				rvs = append(rvs, reflect.ValueOf(v))
			}
			reflect.ValueOf(handler).Call(rvs)
		}()
	})
}

func (s *Notify) handle(event notify.Event, handler any) {
	s.busOn(event, "notify", handler)
}

func (s *Notify) handleInternal(event notify.InternalEvent, handler any) {
	s.busOn(event, "notify_internal", handler)
}

func (s *Notify) initInternalHandler() {
	s.handleInternal(notify.InternalEvent_WsConnected, s.WsConnected)
}

func (s *Notify) initHandler() {
	s.handle(notify.Event_DeleteSpace, s.deleteSpace)
	s.handle(notify.Event_SpaceAbnormal, s.SpaceAbnormal)
	s.handle(notify.Event_DeleteWorkItem, s.DeleteWorkItem)
	s.handle(notify.Event_WorkItemExpired, s.WorkItemExpired)
	s.handle(notify.Event_WorkItemFlowNodeExpired, s.WorkItemFlowNodeExpired)
	s.handle(notify.Event_TransferWorkItem, s.TransferWorkItem)
	s.handle(notify.Event_RemindWork, s.RemindWork)
}

func (s *Notify) initDomainMessageHandler() {
	s.domainMessageConsumer.SetMessageListener("notify_handler", s.domainMessageHandler)

	// 空间
	s.handDomainMessage(domain_message.Message_Type_Space_Create, s.createSpaceByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Space_Modify, s.modifySpaceByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Space_SetNotify, s.setSpaceNotifyByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Space_Transfer, s.transferSpaceByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Space_SetWorkingDay, s.setSpaceWorkingDayByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Space_Quit, s.QuitSpaceByDomainMessage)

	// 空间视图
	s.handDomainMessage(domain_message.Message_Type_SpaceView_Update, s.updateSpaceViewByDomainMessage)

	// 空间成员
	s.handDomainMessage(domain_message.Message_Type_Space_AddMember, s.addMemberByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Space_RemoveMember, s.removeMemberByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Space_SetMemberRole, s.setSpaceMemberRoleByDomainMessage)

	// 标签
	s.handDomainMessage(domain_message.Message_Type_SpaceTag_Delete, s.deleteSpaceTagByDomainMessage)

	// 空间配置
	s.handDomainMessage(domain_message.Message_Type_Space_SetCommentDeletable, s.setCommentDeletableByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Space_SetCommentDeletableWhenArchived, s.setCommentDeletableWhenArchivedByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Space_SetCommentShowPos, s.setCommentShowPosByDomainMessage)

	// 任务
	s.handDomainMessage(domain_message.Message_Type_WorkItem_Create, s.createWorkItemByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_WorkItem_SubTask_Create, s.createChildWorkItemByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_WorkItem_Modify, s.changeWorkItemFieldByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_WorkItem_Director_Change, s.changeWorkItemDirectorByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_WorkItem_Status_Change, s.changeWorkItemStatusByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_WorkItem_Tag_Change, s.changeWorkItemTagByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_WorkItem_File_Change, s.changeWorkItemFileByDomainMessage)

	// 节点
	s.handDomainMessage(domain_message.Message_Type_WorkItem_FlowNode_Reach, s.reachWorkItemFlowNodeByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_WorkItem_FlowNode_Rollback, s.rollbackWorkItemFlowNodeByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_WorkItem_FlowNode_Modify, s.modifymWorkItemFlowNodeByDomainMessage)

	// 工作流
	s.handDomainMessage(domain_message.Message_Type_WorkFlow_Create, s.createWorkFlowByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_WorkFlow_Modify, s.modifyWorkFlowByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_WorkFlow_Delete, s.deleteWorkFlowByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_WorkFlow_Template_Save, s.saveWorkFlowTemplateByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Task_WorkFlow_Upgrade, s.upgradeWorkFlowByDomainMessage)

	// 超管
	s.handDomainMessage(domain_message.Message_Type_Admin_Change_User_Nickname, s.adminChangeUserNicknameByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Admin_Change_User_Role, s.adminChangeUserRoleByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Admin_Cancel_User, s.adminCancelUserByDomainMessage)

	// 评论
	s.handDomainMessage(domain_message.Message_Type_Comment_Create, s.addCommentByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Comment_Delete, s.deleteCommentByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Comment_Update, s.updateCommentByDomainMessage)
	s.handDomainMessage(domain_message.Message_Type_Comment_EmojiAdd, s.addCommentEmojiByDomainMessage)
}

func (s *Notify) handDomainMessage(msgType shared.MessageType, handler any) {
	if _, ok := s.domainMessageRouter[msgType]; ok {
		panic(fmt.Sprintf("notify domain message handler already exists: %s", msgType))
	}
	s.domainMessageRouter[msgType] = handler
}

func (s *Notify) domainMessageHandler(evt domain_message.DomainMessagePublishEvent) {
	defer func() {
		if err := recover(); err != nil {
			s.log.Errorf("handle domainMessageHandler panic: %v %s", err, debug.Stack())
		}
	}()

	for _, v := range evt.Messages {
		handler, ok := s.domainMessageRouter[v.MessageType()]
		if !ok {
			s.log.Errorf("notify domain message handler not found: %s", v.MessageType())
			continue
		}

		var rvs = []reflect.Value{
			reflect.ValueOf(v),
		}

		reflect.ValueOf(handler).Call(rvs)
	}
}

func (s *Notify) Send(msg *msg.Message, userIds ...int64) {
	if len(userIds) == 0 {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	m := utils.ToJSONBytes(&message.Msg{
		Type: notify.MsgType_mt_Notify,
		Data: json.RawMessage(data),
	})

	for _, userId := range userIds {
		ok := s.ws.SendData2User(m, userId)
		if !ok {
			// 保存到redis
			s.notifySnapShotRepo.SaveOfflineNotify(context.Background(), userId, data)
		}
	}
}

func (s *Notify) Send2(ctx *notifyCtx, msg *msg.Message, userIds ...int64) {
	if len(userIds) == 0 {
		return
	}

	// 过滤掉不接收通知的用户
	userIds = stream.Filter(userIds, func(userId int64) bool {
		if ctx.forceNotify {
			return true
		}

		space := ctx.space
		member := ctx.memberMap[userId]
		globalNotify := ctx.userNotifySwitchGlobalMap[userId]
		spaceNotify := ctx.userNotifySwitchSpaceMap[userId]

		if globalNotify == nil || spaceNotify == nil {
			return false
		}

		return globalNotify.Value == "1" && spaceNotify.Value == "1" &&
			(space == nil || space.Notify == 1) &&
			(member == nil || member.Notify == 1)
	})

	s.Send(msg, userIds...)
}

func (s *Notify) SendAllMember(msg any, spaceId int64) {
	if spaceId == 0 {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	m := utils.ToJSONBytes(&message.Msg{
		Type: notify.MsgType_mt_Cooperation,
		Data: json.RawMessage(data),
	})

	ids, err := s.memberRepo.GetSpaceAllMemberIds(context.Background(), spaceId)
	if err != nil {
		s.log.Error(err)
		return
	}

	s.ws.SendData(m, ids...)
}

func (s *Notify) SendToMember(msg any, userId int64) {
	if userId == 0 {
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	m := utils.ToJSONBytes(&message.Msg{
		Type: notify.MsgType_mt_Cooperation,
		Data: json.RawMessage(data),
	})

	s.ws.SendData(m, userId)
}

func (s *Notify) pushThirdPlatformMessage(message *tea_im.RobotMessage, userIds []int64) {
	if len(userIds) == 0 {
		return
	}

	accounts, err := s.userRepo.GetThirdPfAccountByUserIds(context.Background(), userIds)
	if err != nil {
		s.log.Error(err)
		return
	}

	pfMap := stream.GroupBy(accounts, func(account *user.ThirdPfAccount) int32 {
		return account.PfInfo.PfCode
	})

	for pfCode, pfAccounts := range pfMap {
		code := comm.ThirdPlatformCode(pfCode)

		client := s.tpClient.ByPfCode(code)
		if client == nil {
			continue
		}

		pfAccounts := stream.Filter(pfAccounts, func(account *user.ThirdPfAccount) bool {
			return account.Notify == 1
		})

		chatTokens := stream.Map(pfAccounts, func(account *user.ThirdPfAccount) string {
			return account.PfInfo.PfUserKey
		})

		err := client.Push(message, chatTokens)
		if err != nil {
			s.log.Error(err)
		}
	}
}

func (s *Notify) pushThirdPlatformMessage2(ctx *notifyCtx, message *tea_im.RobotMessage, userIds []int64) {
	if len(userIds) == 0 {
		return
	}

	// 过滤掉不接收通知的用户
	userIds = stream.Filter(userIds, func(userId int64) bool {
		if ctx.forceNotify {
			return true
		}

		space := ctx.space
		member := ctx.memberMap[userId]
		globalNotify := ctx.userNotifySwitchGlobalMap[userId]
		thirdPlatformNotify := ctx.userNotifySwitchThirdPlatformMap[userId]
		spaceNotify := ctx.userNotifySwitchSpaceMap[userId]

		return (globalNotify == nil || globalNotify.Value == "1") &&
			(thirdPlatformNotify == nil || thirdPlatformNotify.Value == "1") &&
			(spaceNotify == nil || spaceNotify.Value == "1") &&
			(space == nil || space.Notify == 1) &&
			(member == nil || member.Notify == 1)
	})

	s.pushThirdPlatformMessage(message, userIds)
}

func (s *Notify) getNotifyRedirectDomain() string {
	configRow, err := s.configRepo.GetByKey(context.Background(), consts.CONFIG_NOTIFY_REDIRECT_DOMAIN)
	if err != nil {
		return ""
	}
	return configRow.ConfigValue
}

func (s *Notify) getBalalaAssectDomain() string {
	configRow, err := s.configRepo.GetByKey(context.Background(), consts.CONFIG_BALALA_ASSECT_DOMIN)
	if err != nil {
		return ""
	}
	return configRow.ConfigValue
}

func (s *Notify) makeIconResLink(linkDomain string, res IconRes) string {
	if linkDomain == "" {
		linkDomain = s.getBalalaAssectDomain()
	}
	link := fmt.Sprintf("/notice/%v", res)
	return linkDomain + link
}

func (s *Notify) makeSpaceRedirectLink(linkDomain string, spaceId int64) string {
	if linkDomain == "" {
		linkDomain = s.getNotifyRedirectDomain()
	}
	link := fmt.Sprintf("/#/project/%v?type=intro", spaceId)
	return linkDomain + link
}

func (s *Notify) makeTableRedirectLink(linkDomain string, spaceId int64) string {
	if linkDomain == "" {
		linkDomain = s.getNotifyRedirectDomain()
	}
	link := fmt.Sprintf("/#/project/%v?type=table", spaceId)
	return linkDomain + link
}

func (s *Notify) makeWorkFlowRedirectLink(linkDomain string, spaceId int64) string {
	if linkDomain == "" {
		linkDomain = s.getNotifyRedirectDomain()
	}
	link := fmt.Sprintf("/#/project/%v?type=flow", spaceId)
	return linkDomain + link
}

func (s *Notify) makeViewRedirectLink(linkDomain string, spaceId, outerId int64) string {
	if linkDomain == "" {
		linkDomain = s.getNotifyRedirectDomain()
	}
	link := fmt.Sprintf("/#/project/%v?type=view&outer_id=%v", spaceId, outerId)
	return linkDomain + link
}

func (s *Notify) makeWorkItemRedirectLink(linkDomain string, spaceId int64, workItemId int64, pid int64) string {
	if linkDomain == "" {
		linkDomain = s.getNotifyRedirectDomain()
	}
	link := fmt.Sprintf("/#/project/%v?workItemId=%v&tabKey=base&spaceId=%v&pid=%v", spaceId, workItemId, spaceId, pid)
	return linkDomain + link
}

func (s *Notify) makeWorkItemCommentRedirectLink(linkDomain string, spaceId int64, workItemId, pid int64) string {
	if linkDomain == "" {
		linkDomain = s.getNotifyRedirectDomain()
	}
	link := fmt.Sprintf("/#/project/%v?workItemId=%v&tabKey=comment&spaceId=%v&pid=%v", spaceId, workItemId, spaceId, pid)
	return linkDomain + link
}

func (s *Notify) makeRedirectLink(message *msg.Message) string {

	notifyRedirectDomain := s.getNotifyRedirectDomain()

	var link string

	obj := message.Notification.Object
	subObj := message.Notification.SubObject

	//评论链接特殊处理
	if subObj != nil && subObj.Type == msg.ObjectType_workItemComment {
		workItem := obj.Data.(*msg.WorkItemData)
		link = s.makeWorkItemCommentRedirectLink(notifyRedirectDomain, message.Space.SpaceId, workItem.Id, workItem.Pid)
		return link
	}

	switch obj.Type {
	case msg.ObjectType_workItem:
		workItem := obj.Data.(*msg.WorkItemData)
		link = s.makeWorkItemRedirectLink(notifyRedirectDomain, message.Space.SpaceId, workItem.Id, workItem.Pid)
	case msg.ObjectType_space:
		space := obj.Data.(*msg.SpaceData)
		link = s.makeSpaceRedirectLink(notifyRedirectDomain, space.Id)
	case msg.ObjectType_workFlow:
		link = s.makeWorkFlowRedirectLink(notifyRedirectDomain, message.Space.SpaceId)
	case msg.ObjectType_View:
		viewData := obj.Data.(*msg.ViewData)
		link = s.makeViewRedirectLink(notifyRedirectDomain, message.Space.SpaceId, viewData.Id)
	case msg.ObjectType_user:
	case msg.ObjectType_workItemComment:
	}

	if link == "" {
		return notifyRedirectDomain
	}

	return link
}

func (s *Notify) isTodoWorkItem(workItem *work_item.WorkItem) bool {
	statusInfo, _ := s.workItemStatusRepo.GetWorkItemStatusInfo(context.Background(), workItem.SpaceId)
	for _, v := range statusInfo.Items {
		if v.IsProcessingTypeState() && workItem.WorkItemStatus.Val == v.Val {
			return true
		}
	}
	return false
}
