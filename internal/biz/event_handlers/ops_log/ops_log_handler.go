package ops_log

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/panjf2000/ants/v2"
	db "go-cs/internal/bean/biz"
	oper_repo "go-cs/internal/domain/oper_log/repo"
	domain_message "go-cs/internal/domain/pkg/message"
	space_repo "go-cs/internal/domain/space/repo"
	space_tag_repo "go-cs/internal/domain/space_tag/repo"
	view_repo "go-cs/internal/domain/space_view/repo"
	comment_repo "go-cs/internal/domain/space_work_item_comment/repo"
	wObject_repo "go-cs/internal/domain/space_work_object/repo"
	wVersion_repo "go-cs/internal/domain/space_work_version/repo"
	user_repo "go-cs/internal/domain/user/repo"
	wFlow_repo "go-cs/internal/domain/work_flow/repo"
	witem_repo "go-cs/internal/domain/work_item/repo"
	witem_role_repo "go-cs/internal/domain/work_item_role/repo"
	witem_status_repo "go-cs/internal/domain/work_item_status/repo"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils"
	"go-cs/internal/utils/oper"
	"go-cs/pkg/stream"
)

type OpsLogEventHandlers struct {
	log    *log.Helper
	pool   *ants.Pool
	router map[shared.MessageType]Handler

	userRepo     user_repo.UserRepo
	operRepo     oper_repo.OperLogRepo
	spaceRepo    space_repo.SpaceRepo
	wRoleRepo    witem_role_repo.WorkItemRoleRepo
	wStatusRepo  witem_status_repo.WorkItemStatusRepo
	wObjectRepo  wObject_repo.SpaceWorkObjectRepo
	wVersionRepo wVersion_repo.SpaceWorkVersionRepo
	wItemRepo    witem_repo.WorkItemRepo
	tagRepo      space_tag_repo.SpaceTagRepo
	wFlowRepo    wFlow_repo.WorkFlowRepo
	commentRepo  comment_repo.SpaceWorkItemCommentRepo
	viewRepo     view_repo.SpaceViewRepo

	domainMessageConsumer *domain_message.DomainMessageConsumer
}

func NewOpsLogEventHandlers(
	logger log.Logger,
	userRepo user_repo.UserRepo,
	operRepo oper_repo.OperLogRepo,
	spaceRepo space_repo.SpaceRepo,
	wRoleRepo witem_role_repo.WorkItemRoleRepo,
	wStatusRepo witem_status_repo.WorkItemStatusRepo,
	wObjectRepo wObject_repo.SpaceWorkObjectRepo,
	wVersionRepo wVersion_repo.SpaceWorkVersionRepo,
	wItemRepo witem_repo.WorkItemRepo,
	tagRepo space_tag_repo.SpaceTagRepo,
	wFlowRepo wFlow_repo.WorkFlowRepo,
	commentRepo comment_repo.SpaceWorkItemCommentRepo,
	viewRepo view_repo.SpaceViewRepo,

	domainMessageConsumer *domain_message.DomainMessageConsumer,

) *OpsLogEventHandlers {

	moduleName := "OpsLogEventHandlers"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	pool, err := ants.NewPool(10000)
	if err != nil {
		panic("OpsLogEventHandlers init pool failed")
	}

	evtHandlers := &OpsLogEventHandlers{
		log:          hlog,
		userRepo:     userRepo,
		operRepo:     operRepo,
		spaceRepo:    spaceRepo,
		wRoleRepo:    wRoleRepo,
		wStatusRepo:  wStatusRepo,
		wObjectRepo:  wObjectRepo,
		wVersionRepo: wVersionRepo,
		wItemRepo:    wItemRepo,
		tagRepo:      tagRepo,
		wFlowRepo:    wFlowRepo,
		commentRepo:  commentRepo,
		viewRepo:     viewRepo,

		domainMessageConsumer: domainMessageConsumer,
		router:                map[shared.MessageType]Handler{},
		pool:                  pool,
	}

	return evtHandlers
}

func (s *OpsLogEventHandlers) Init() {
	s.domainMessageConsumer.SetMessageListener("ops_log_handler", s.domainMessagePublishEventHandler)

	// 空间
	s.handle(domain_message.Message_Type_Space_Quit, s.spaceQuitLogHandler)
	s.handle(domain_message.Message_Type_Space_AddMember, s.spaceAddMemberHandler)
	s.handle(domain_message.Message_Type_Space_RemoveMember, s.spaceRemoveMemberHandler)
	s.handle(domain_message.Message_Type_Space_SetMemberRole, s.spaceSetMemberRoleHandler)
	s.handle(domain_message.Message_Type_Space_Transfer, s.TransferSpaceHandler)
	s.handle(domain_message.Message_Type_Space_Create, s.spaceCreateHandler)
	s.handle(domain_message.Message_Type_Space_Modify, s.spaceModifyHandler)
	s.handle(domain_message.Message_Type_Space_Del, s.spaceDeleteHandler)
	s.handle(domain_message.Message_Type_Space_SetCommentDeletable, s.spaceSetCommentDeletableHandler)
	s.handle(domain_message.Message_Type_Space_SetCommentDeletableWhenArchived, s.spaceSetCommentDeletableWhenArchivedHandler)
	s.handle(domain_message.Message_Type_Space_SetCommentShowPos, s.spaceSetCommentShowPosHandler)
	s.handle(domain_message.Message_Type_Space_SetWorkingDay, s.spaceSetWorkingDayHandler)
	s.handle(domain_message.Message_Type_Space_SetNotify, s.spaceSetSpaceNotifyHandler)
	s.handle(domain_message.Message_Type_Space_SetTempConfig, s.spaceSetTempConfigHandler)

	// 工作项
	s.handle(domain_message.Message_Type_WorkItem_Create, s.workItemCreateHandler)
	s.handle(domain_message.Message_Type_WorkItem_SubTask_Create, s.workItemSubTaskCreateHandler)
	s.handle(domain_message.Message_Type_WorkItem_Modify, s.workItemModifyHandler)
	s.handle(domain_message.Message_Type_WorkItem_Status_Change, s.workItemStatusChangeHandler)
	s.handle(domain_message.Message_Type_WorkItem_Director_Change, s.workItemDirectorChangeHandler)
	s.handle(domain_message.Message_Type_WorkItem_Delete, s.workItemDeleteHandler)
	s.handle(domain_message.Message_Type_WorkItem_Tag_Change, s.workItemTagChangeHandler)
	s.handle(domain_message.Message_Type_WorkItem_File_Change, s.workItemFileChangeHandler)
	s.handle(domain_message.Message_Type_Task_WorkFlow_Upgrade, s.workItemFlowUpgradeHandler)

	// 标签
	s.handle(domain_message.Message_Type_SpaceTag_Create, s.spaceTagCreateHandler)
	s.handle(domain_message.Message_Type_SpaceTag_Modify, s.spaceTagModifyHandler)
	s.handle(domain_message.Message_Type_SpaceTag_Delete, s.spaceTagDeleteHandler)

	// 视图
	s.handle(domain_message.Message_Type_SpaceView_Create, s.spaceViewCreateHandler)
	s.handle(domain_message.Message_Type_SpaceView_Delete, s.spaceViewDeleteHandler)
	s.handle(domain_message.Message_Type_SpaceView_SetName, s.spaceViewSetNameHandler)
	s.handle(domain_message.Message_Type_SpaceView_SetStatus, s.spaceViewSetStatusHandler)
	s.handle(domain_message.Message_Type_SpaceView_Update, s.updateSpaceViewHandler)

	// 工作项-节点
	s.handle(domain_message.Message_Type_WorkItem_FlowNode_Modify, s.workItemFlowNodeModifyHandler)
	s.handle(domain_message.Message_Type_WorkItem_FlowNode_Confirm, s.workItemFlowNodeConfirmHandler)
	s.handle(domain_message.Message_Type_WorkItem_FlowNode_Rollback, s.workItemFlowNodeRollbackHandler)

	// 工作流
	s.handle(domain_message.Message_Type_WorkFlow_Create, s.workFlowCreateHandler)
	s.handle(domain_message.Message_Type_WorkFlow_Modify, s.workFlowModifyHandler)
	s.handle(domain_message.Message_Type_WorkFlow_Delete, s.workFlowDeleteHandler)
	s.handle(domain_message.Message_Type_WorkFlow_Template_Save, s.workFlowTemplateSaveHandler)

	// 流程角色
	s.handle(domain_message.Message_Type_WorkItemRole_Create, s.workItemRoleCreateHandler)
	s.handle(domain_message.Message_Type_WorkItemRole_Modify, s.workItemRoleModifyHandler)
	s.handle(domain_message.Message_Type_WorkItemRole_Delete, s.workItemRoleDeleteHandler)

	// 工作项状态
	s.handle(domain_message.Message_Type_WorkItemStatus_Create, s.workItemStatusCreateHandler)
	s.handle(domain_message.Message_Type_WorkItemStatus_Modify, s.workItemStatusModifyHandler)
	s.handle(domain_message.Message_Type_WorkItemStatus_Delete, s.workItemStatusDeleteHandler)

	// 工作模块
	s.handle(domain_message.Message_Type_WorkObject_Create, s.workObjectCreateHandler)
	s.handle(domain_message.Message_Type_WorkObject_Modify, s.workObjectModifyHandler)
	s.handle(domain_message.Message_Type_WorkObject_Delete, s.workObjectDeleteHandler)

	// 任务版本
	s.handle(domain_message.Message_Type_WorkVersion_Create, s.workVersionCreateHandler)
	s.handle(domain_message.Message_Type_WorkVersion_Modify, s.workVersionModifyHandler)
	s.handle(domain_message.Message_Type_WorkVersion_Delete, s.workVersionDeleteHandler)

	// 顺序
	s.handle(domain_message.Message_Type_SpaceAllOrder_WorkObject, s.changeWorkObjectOrderHandler)
	s.handle(domain_message.Message_Type_SpaceAllOrder_Version, s.changeVersionOrderHandler)
	s.handle(domain_message.Message_Type_SpaceAllOrder_Role, s.changeRoleOrderHandler)
	s.handle(domain_message.Message_Type_SpaceAllOrder_Status, s.changeStatusOrderHandler)
	s.handle(domain_message.Message_Type_SpaceAllOrder_WorkFlow, s.changeWorkFlowOrderHandler)
	s.handle(domain_message.Message_Type_SpaceAllOrder_OverviewDataItem, s.changeOverviewDataItemOrderHandler)
	s.handle(domain_message.Message_Type_SpaceAllOrder_OverviewBlock, s.changeOverviewBlockOrderHandler)
	s.handle(domain_message.Message_Type_SpaceAllOrder_View, s.changeViewOrderHandler)

	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

	// 个人
	s.handle(domain_message.Message_Type_Personal_ResetPwd, s.personalResetPwdHandler)
	s.handle(domain_message.Message_Type_Personal_ChangeNickname, s.personalChangeNicknameHandler)
	s.handle(domain_message.Message_Type_Personal_ChangeAvatar, s.personalChangeAvatarHandler)
	s.handle(domain_message.Message_Type_Personal_SetSpaceNotify, s.personalSetSpaceNotifyHandler)
	s.handle(domain_message.Message_Type_Personal_BindThirdPlatform, s.personalBindThirdPlatformHandler)
	s.handle(domain_message.Message_Type_Personal_UnBindThirdPlatform, s.personalUnBindThirdPlatformHandler)
	s.handle(domain_message.Message_Type_Personal_SetTempConfig, s.personalSetTempConfigHandler)
	s.handle(domain_message.Message_Type_Personal_SetThirdPlatformNotify, s.personalSetThirdPlatformNotifyHandler)
	s.handle(domain_message.Message_Type_Personal_SetUserConfig, s.personalSetUserConfigHandler)

	//--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

	// 系统管理员日志
	s.handle(domain_message.Message_Type_Admin_Add_User, s.AdminAddUserHandler)
	s.handle(domain_message.Message_Type_Admin_Cancel_User, s.AdminCancelUserHandler)
	s.handle(domain_message.Message_Type_Admin_Change_User_Nickname, s.AdminChangeUserNicknameHandler)
	s.handle(domain_message.Message_Type_Admin_Change_User_Role, s.AdminChangeUserRoleHandler)
	s.handle(domain_message.Message_Type_Admin_Reset_User_Password, s.AdminResetUserPasswordHandler)

	s.handle(domain_message.Message_Type_Admin_Change_System_Logo, s.AdminChangeSystemLogoHandler)
	s.handle(domain_message.Message_Type_Admin_Change_System_Title, s.AdminChangeSystemTitleHandler)
	s.handle(domain_message.Message_Type_Admin_Change_System_AccessUrl, s.AdminChangeSystemAccessUrlHandler)
	s.handle(domain_message.Message_Type_Admin_Change_System_LoginBg, s.AdminChangeSystemLoginBgHandler)
	s.handle(domain_message.Message_Type_Admin_Change_System_RegisterEntry, s.AdminChangeSystemRegisterEntryHandler)
	s.handle(domain_message.Message_Type_Admin_Change_System_AttachSize, s.AdminChangeSystemAttachSizeHandler)

	// 工作项评论
	s.handle(domain_message.Message_Type_Comment_Create, s.addCommentHandler)
	s.handle(domain_message.Message_Type_Comment_Delete, s.deleteCommentHandler)
	s.handle(domain_message.Message_Type_Comment_Update, s.updateCommentHandler)
	s.handle(domain_message.Message_Type_Comment_EmojiAdd, s.addCommentEmojiHandler)
	s.handle(domain_message.Message_Type_Comment_EmojiRemove, s.removeCommentEmojiHandler)

}

type Handler func(ctx context.Context, operLogger *oper.OperLogger, opsLogs shared.DomainMessages)

func (s *OpsLogEventHandlers) handle(msgType shared.MessageType, handler Handler) {
	if _, ok := s.router[msgType]; ok {
		panic("ops log message handler already exists: " + msgType)
	}

	s.router[msgType] = handler
}

func (s *OpsLogEventHandlers) domainMessagePublishEventHandler(evt domain_message.DomainMessagePublishEvent) {
	//合并日志操作
	operLogger := oper.GetOperLoggerFormCtx(evt.Ctx)
	if operLogger == nil {
		//上下文日志记录对象，就不处理了
		return
	}

	//把同类型的日志分组，然后处理, 要不要合并，实际业务再说， 合并的日志，以最早的那条操作时间为主
	opsLogGroup := make(map[shared.MessageType]shared.DomainMessages)
	for _, v := range evt.Messages {
		opsLogGroup[v.MessageType()] = append(opsLogGroup[v.MessageType()], v)
	}

	for k, v := range opsLogGroup {
		handler, ok := s.router[k]
		if !ok {
			s.log.Errorf("no handler for %s", k)
			continue
		}

		handler(evt.Ctx, operLogger, v)
	}
}

func (s *OpsLogEventHandlers) getOperUserInfo(ctx context.Context, opsOper shared.Oper) *oper.OperUser {
	userInfo, err := s.userRepo.GetUserByUserId(ctx, opsOper.GetId())
	if err != nil {
		return nil
	}

	return &oper.OperUser{
		OperType:         int(opsOper.GetType()),
		OperUid:          userInfo.Id,
		OperUname:        userInfo.UserName,
		OperUserNickName: userInfo.UserNickname,
	}

}

func (s *OpsLogEventHandlers) invokeOperLog(ctx context.Context, logger *oper.OperLogger, result *oper.OperResultInfo) {
	s.pool.Submit(func() {
		var param = logger.RequestInfo.OperParam
		if len(logger.RequestInfo.OperParam) >= 2048 { //太长直接忽略
			param = ""
		}

		r := result
		operLog := &db.OperLog{
			Title:         r.ModuleTitle,
			BusinessType:  int32(r.BusinessType),
			ShowType:      int32(r.ShowType),
			RequestMethod: logger.RequestInfo.RequestMethod,
			ModuleType:    int32(r.ModuleType),
			ModuleId:      int64(r.ModuleId),
			OperatorType:  int32(logger.Operator.OperType),
			OperId:        logger.Operator.OperUid,
			OperName:      logger.Operator.OperUname,
			OperNickname:  logger.Operator.OperUserNickName,
			OperIp:        logger.RequestInfo.OperIp,
			OperLocation:  logger.RequestInfo.OperLocation,
			OperUrl:       logger.RequestInfo.OperUrl,
			OperParam:     param,
			OperMsg:       r.OperMsg,
			SpaceId:       r.SpaceId,
			SpaceName:     r.SpaceName,
			ModuleFlag: int64(stream.Reduce(r.ModuleFlags, func(acc, v oper.ModuleFlag) oper.ModuleFlag {
				return acc | v
			})),
		}

		s.operRepo.AddOperLog(ctx, operLog)
	})
}
