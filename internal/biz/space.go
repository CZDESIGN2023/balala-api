package biz

import (
	"cmp"
	"context"
	"encoding/json"
	"go-cs/api/comm"
	"go-cs/api/notify"
	pb "go-cs/api/space/v1"
	v1 "go-cs/api/space/v1"
	db "go-cs/internal/bean/biz"
	"go-cs/internal/bean/vo/event"
	"go-cs/internal/bean/vo/query"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/biz/space_temp_config"
	"go-cs/internal/consts"
	"go-cs/internal/data/convert"
	perm_facade "go-cs/internal/domain/perm/facade"
	perm_service "go-cs/internal/domain/perm/service"
	domain_message "go-cs/internal/domain/pkg/message"
	search_repo "go-cs/internal/domain/search/repo"
	"go-cs/internal/domain/search/search_es"
	"go-cs/internal/domain/space"
	space_repo "go-cs/internal/domain/space/repo"
	space_service "go-cs/internal/domain/space/service"
	space_file_repo "go-cs/internal/domain/space_file_info/repo"
	member_domain "go-cs/internal/domain/space_member"
	member_repo "go-cs/internal/domain/space_member/repo"
	member_service "go-cs/internal/domain/space_member/service"
	tag_repo "go-cs/internal/domain/space_tag/repo"
	space_view_repo "go-cs/internal/domain/space_view/repo"
	space_view_service "go-cs/internal/domain/space_view/service"
	workObj_repo "go-cs/internal/domain/space_work_object/repo"
	user_repo "go-cs/internal/domain/user/repo"
	wf "go-cs/internal/domain/work_flow"
	work_flow_facade "go-cs/internal/domain/work_flow/facade"
	wf_repo "go-cs/internal/domain/work_flow/repo"
	wf_service "go-cs/internal/domain/work_flow/service"
	witem_repo "go-cs/internal/domain/work_item/repo"
	"go-cs/internal/domain/work_item_role"
	role_repo "go-cs/internal/domain/work_item_role/repo"
	role_service "go-cs/internal/domain/work_item_role/service"
	witem_status "go-cs/internal/domain/work_item_status"
	witem_status_repo "go-cs/internal/domain/work_item_status/repo"
	witem_status_service "go-cs/internal/domain/work_item_status/service"
	witem_type "go-cs/internal/domain/work_item_type"
	witem_type_repo "go-cs/internal/domain/work_item_type/repo"
	witem_type_service "go-cs/internal/domain/work_item_type/service"
	shared "go-cs/internal/pkg/domain"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"

	work_object_service "go-cs/internal/domain/space_work_object/service"
	workVersion_repo "go-cs/internal/domain/space_work_version/repo"
	workVersion_service "go-cs/internal/domain/space_work_version/service"

	"go-cs/internal/pkg/biz_id"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"go-cs/pkg/bus"
	"go-cs/pkg/stream"
	"slices"

	"github.com/go-kratos/kratos/v2/log"
)

type SpaceUsecase struct {
	log *log.Helper
	tm  trans.Transaction

	repo               space_repo.SpaceRepo
	userRepo           user_repo.UserRepo
	memberRepo         member_repo.SpaceMemberRepo
	workObjectRepo     workObj_repo.SpaceWorkObjectRepo
	workVersionRepo    workVersion_repo.SpaceWorkVersionRepo
	workItemRepo       witem_repo.WorkItemRepo
	fileInfoRepo       space_file_repo.SpaceFileInfoRepo
	tagRepo            tag_repo.SpaceTagRepo
	workItemTypeRepo   witem_type_repo.WorkItemTypeRepo
	workFlowRepo       wf_repo.WorkFlowRepo
	workItemRoleRepo   role_repo.WorkItemRoleRepo
	workItemStatusRepo witem_status_repo.WorkItemStatusRepo
	searchRepo         search_repo.SearchRepo
	viewRepo           space_view_repo.SpaceViewRepo

	workFlowDomainService *wf_service.WorkFlowService
	idService             *biz_id.BusinessIdService
	wizService            *witem_status_service.WorkItemStatusService
	roleService           *role_service.WorkItemRoleService
	spaceServie           *space_service.SpaceService
	permService           *perm_service.PermService
	memberService         *member_service.SpaceMemberService
	witemTypeService      *witem_type_service.WorkItemTypeService
	workVersionService    *workVersion_service.SpaceWorkVersionService
	workObjectService     *work_object_service.SpaceWorkObjectService
	viewService           *space_view_service.SpaceViewService

	domainMessageProducer *domain_message.DomainMessageProducer
}

func NewSpaceUsecase(

	tm trans.Transaction,
	logger log.Logger,

	repo space_repo.SpaceRepo,
	workVersionRepo workVersion_repo.SpaceWorkVersionRepo,
	memberRepo member_repo.SpaceMemberRepo,
	workObjectRepo workObj_repo.SpaceWorkObjectRepo,
	userRepo user_repo.UserRepo,
	workItemRepo witem_repo.WorkItemRepo,
	fileInfoRepo space_file_repo.SpaceFileInfoRepo,
	tagRepo tag_repo.SpaceTagRepo,
	workItemTypeRepo witem_type_repo.WorkItemTypeRepo,
	workFlowRepo wf_repo.WorkFlowRepo,
	workItemRoleRepo role_repo.WorkItemRoleRepo,
	workItemStatusRepo witem_status_repo.WorkItemStatusRepo,
	searchRepo search_repo.SearchRepo,
	viewRepo space_view_repo.SpaceViewRepo,

	workFlowDomainService *wf_service.WorkFlowService,
	idService *biz_id.BusinessIdService,
	wizService *witem_status_service.WorkItemStatusService,
	roleService *role_service.WorkItemRoleService,
	spaceServie *space_service.SpaceService,
	permService *perm_service.PermService,
	memberService *member_service.SpaceMemberService,
	witemTypeService *witem_type_service.WorkItemTypeService,
	workVersionService *workVersion_service.SpaceWorkVersionService,
	workObjectService *work_object_service.SpaceWorkObjectService,
	viewService *space_view_service.SpaceViewService,

	domainMessageProducer *domain_message.DomainMessageProducer,

) *SpaceUsecase {
	moduleName := "SpaceUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceUsecase{
		log: hlog,
		tm:  tm,

		repo:               repo,
		userRepo:           userRepo,
		memberRepo:         memberRepo,
		workObjectRepo:     workObjectRepo,
		workVersionRepo:    workVersionRepo,
		workItemRepo:       workItemRepo,
		tagRepo:            tagRepo,
		workItemTypeRepo:   workItemTypeRepo,
		workFlowRepo:       workFlowRepo,
		fileInfoRepo:       fileInfoRepo,
		workItemRoleRepo:   workItemRoleRepo,
		workItemStatusRepo: workItemStatusRepo,
		searchRepo:         searchRepo,
		viewRepo:           viewRepo,

		workFlowDomainService: workFlowDomainService,
		idService:             idService,
		wizService:            wizService,
		roleService:           roleService,
		spaceServie:           spaceServie,
		memberService:         memberService,
		witemTypeService:      witemTypeService,
		workVersionService:    workVersionService,
		workObjectService:     workObjectService,
		viewService:           viewService,

		domainMessageProducer: domainMessageProducer,
	}
}

func (uc *SpaceUsecase) CreateMySpace(ctx context.Context, oper *utils.LoginUserInfo, spaceName string, spaceDescribe string, inMember []*db.SpaceMember) (*db.Space, error) {

	//检查创建空间的用户是否存在
	_, err := uc.userRepo.GetUserByUserId(ctx, oper.UserId)
	if err != nil {
		//创建空间失败-创建人的信息不存在
		err := errs.New(ctx, comm.ErrorCode_SPACE_CREATOR_WRONG)
		return nil, err
	}

	// 新增空间
	space, err := uc.spaceServie.CreateSpace(ctx, oper.UserId, spaceName, spaceDescribe, oper)
	if err != nil {
		return nil, err
	}

	var members []*member_domain.SpaceMember
	ownerMember, err := uc.memberService.NewSpaceMember(ctx, space.Id, oper.UserId, consts.MEMBER_ROLE_SPACE_CREATOR, oper)
	if err != nil {
		return nil, err
	}

	members = append(members, ownerMember)
	for _, v := range inMember {

		//忽略创建者的ID
		if v.UserId == oper.UserId {
			continue
		}

		if !slices.Contains([]int64{consts.MEMBER_ROLE_MANAGER, consts.MEMBER_ROLE_EDITOR, consts.MEMBER_ROLE_WATCHER}, v.RoleId) {
			v.RoleId = consts.MEMBER_ROLE_EDITOR
		}

		newMember, _ := uc.memberService.NewSpaceMember(ctx, space.Id, v.UserId, v.RoleId, oper)
		if newMember != nil {
			members = append(members, newMember)
		}
	}

	// 添加默认的工作状态
	stateInfo := uc.wizService.CreateSpaceDefaultStatusInfo(
		ctx, space.Id,
		[]consts.FlowScope{consts.FlowScope_All, consts.FlowScope_Workflow, consts.FlowScope_Stateflow},
		space.UserId,
	)
	//添加工作类型
	var taskType *witem_type.WorkItemType
	var subTaskType *witem_type.WorkItemType
	var stateTaskType *witem_type.WorkItemType

	itemTypes, _ := uc.witemTypeService.CreateSpaceDefaultWorkItemTypes(ctx, space.Id, space.UserId)
	if itemTypes != nil {
		taskType = itemTypes[string(consts.WorkItemTypeKey_Task)]
		subTaskType = itemTypes[string(consts.WorkItemTypeKey_SubTask)]
		stateTaskType = itemTypes[string(consts.WorkItemTypeKey_StateTask)]
	}

	//添加默认的工作项类型-角色
	roles := uc.roleService.CreateSpaceDefaultRoles(ctx, role_service.CreateSpaceDefaultRolesReq{
		SpaceId:              space.Id,
		OperUid:              space.UserId,
		WorkFlowWitemTypeId:  taskType.Id,
		StateFlowWitemTypeId: stateTaskType.Id,
	})

	//默认的流程模版
	workFlowResult, _ := uc.createSpaceWorkFlow(&createSpaceWorkFlowCtx{
		Ctx:                   ctx,
		SpaceId:               space.Id,
		WorkItemTypeId:        taskType.Id,
		SubTaskWorkItemTypeId: subTaskType.Id,
		StateWorkItemTypeId:   stateTaskType.Id,
		WorkItemStatusInfo:    stateInfo,
		Roles:                 roles,
	})

	//添加默认版本
	workVersion, err := uc.workVersionService.CreateSpaceDefaultWorkVersion(ctx, space.Id, oper)
	if err != nil {
		return nil, err
	}

	//添加默认模块
	workObject, err := uc.workObjectService.CreateSpaceDefaultWorkObject(ctx, space.Id, oper)
	if err != nil {
		return nil, err
	}

	userIds := stream.Map(members, func(member *member_domain.SpaceMember) int64 {
		return member.UserId
	})

	globalViews, err := uc.viewService.InitSpacePublicView(ctx, space.Id, stateInfo.Items)
	if err != nil {
		return nil, err
	}

	err = uc.tm.InTx(ctx, func(ctx context.Context) error {

		err = uc.memberRepo.AddSpaceMembers(ctx, members)
		if err != nil {
			return err
		}

		// 添加工作项默认版本
		err = uc.workVersionRepo.CreateSpaceWorkVersion(ctx, workVersion)
		if err != nil {
			return err
		}

		// 添加工作项默认版本
		err = uc.workObjectRepo.CreateSpaceWorkObject(ctx, workObject)
		if err != nil {
			return err
		}

		// 添加默认配置
		err = uc.repo.CreateSpaceConfig(ctx, space.SpaceConfig)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		// 添加默认工作状态
		err = uc.workItemStatusRepo.CreateWorkItemStatusItems(ctx, space.Id, stateInfo.Items)
		if err != nil {
			return err
		}

		// 添加默认工作项类型以及对应工作流程模版

		err = uc.workItemTypeRepo.CreateWorkItemType(ctx, taskType)
		if err != nil {
			return err
		}

		err = uc.workItemTypeRepo.CreateWorkItemType(ctx, subTaskType)
		if err != nil {
			return err
		}

		err = uc.workItemTypeRepo.CreateWorkItemType(ctx, stateTaskType)
		if err != nil {
			return err
		}

		//添加默认工作项类型-角色
		err = uc.workItemRoleRepo.CreateWorkItemRoles(ctx, roles)
		if err != nil {
			return err
		}

		// 保存space
		err = uc.repo.CreateSpace(ctx, space)
		if err != nil {
			return err
		}

		//保存工作流信息
		err = uc.workFlowRepo.CreateWorkFlows(ctx, workFlowResult.WorkFlow)
		if err != nil {
			return err
		}

		//保存工作流模版信息
		err = uc.workFlowRepo.CreateWorkFlowTemplates(ctx, workFlowResult.WorkFlowTemplate)
		if err != nil {
			return err
		}

		// 创建公共视图
		err = uc.viewRepo.CreateGlobalViews(ctx, globalViews)
		// 创建个人视图
		userViews, err := uc.viewService.InitUserGlobalView(ctx, space.Id, userIds)
		if err != nil {
			return err
		}
		err = uc.viewRepo.CreateUserViews(ctx, userViews)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	uc.domainMessageProducer.Send(ctx, space.GetMessages())

	return convert.SpaceEntityToPo(space), nil
}

type createSpaceWorkFlowCtx struct {
	Ctx                   context.Context
	SpaceId               int64
	WorkItemTypeId        int64
	SubTaskWorkItemTypeId int64
	StateWorkItemTypeId   int64
	WorkItemStatusInfo    *witem_status.WorkItemStatusInfo
	Roles                 []*work_item_role.WorkItemRole
}

type createSpaceWorkFlowResult struct {
	WorkFlow         []*wf.WorkFlow
	WorkFlowTemplate []*wf.WorkFlowTemplate
}

func (uc *SpaceUsecase) createSpaceWorkFlow(ctx *createSpaceWorkFlowCtx) (*createSpaceWorkFlowResult, error) {

	// 添加 节点模式任务 的 流程
	designWorkFlow := uc.workFlowDomainService.NewDesignWorkFlow(ctx.Ctx, &wf_service.GenerateWorkFlowReq{
		SpaceId:            ctx.SpaceId,
		WorkItemTypeId:     ctx.WorkItemTypeId,
		Ranking:            800,
		WorkItemStatusInfo: work_flow_facade.BuildWorkItemStatusInfo(ctx.WorkItemStatusInfo),
		WorkItemRoleInfo:   work_flow_facade.BuildWorkItemRoleInfo(ctx.Roles),
	})
	xuqiuWorkFlow := uc.workFlowDomainService.NewXuQiuWorkFlow(ctx.Ctx, &wf_service.GenerateWorkFlowReq{
		SpaceId:            ctx.SpaceId,
		WorkItemTypeId:     ctx.WorkItemTypeId,
		Ranking:            700,
		WorkItemStatusInfo: work_flow_facade.BuildWorkItemStatusInfo(ctx.WorkItemStatusInfo),
		WorkItemRoleInfo:   work_flow_facade.BuildWorkItemRoleInfo(ctx.Roles),
	})

	issueStateFlow := uc.workFlowDomainService.NewIssueStateFlow(ctx.Ctx, &wf_service.GenerateWorkFlowReq{
		SpaceId:            ctx.SpaceId,
		WorkItemTypeId:     ctx.StateWorkItemTypeId,
		Ranking:            600,
		WorkItemStatusInfo: work_flow_facade.BuildWorkItemStatusInfo(ctx.WorkItemStatusInfo),
		WorkItemRoleInfo:   work_flow_facade.BuildWorkItemRoleInfo(ctx.Roles),
	})

	zouchaWorkFlow := uc.workFlowDomainService.NewZouChaWorkFlow(ctx.Ctx, &wf_service.GenerateWorkFlowReq{
		SpaceId:            ctx.SpaceId,
		WorkItemTypeId:     ctx.WorkItemTypeId,
		Ranking:            500,
		WorkItemStatusInfo: work_flow_facade.BuildWorkItemStatusInfo(ctx.WorkItemStatusInfo),
		WorkItemRoleInfo:   work_flow_facade.BuildWorkItemRoleInfo(ctx.Roles),
	})

	subTaskFlow := uc.workFlowDomainService.NewSubTaskFlow(ctx.Ctx, &wf_service.GenerateWorkFlowReq{
		SpaceId:            ctx.SpaceId,
		WorkItemTypeId:     ctx.SubTaskWorkItemTypeId,
		Ranking:            400,
		WorkItemStatusInfo: work_flow_facade.BuildWorkItemStatusInfo(ctx.WorkItemStatusInfo),
		WorkItemRoleInfo:   work_flow_facade.BuildWorkItemRoleInfo(ctx.Roles),
	})

	bugWorkFlow := uc.workFlowDomainService.NewBugWorkFlow(ctx.Ctx, &wf_service.GenerateWorkFlowReq{
		SpaceId:            ctx.SpaceId,
		WorkItemTypeId:     ctx.WorkItemTypeId,
		Ranking:            300,
		WorkItemStatusInfo: work_flow_facade.BuildWorkItemStatusInfo(ctx.WorkItemStatusInfo),
		WorkItemRoleInfo:   work_flow_facade.BuildWorkItemRoleInfo(ctx.Roles),
	})

	result := &createSpaceWorkFlowResult{
		WorkFlowTemplate: []*wf.WorkFlowTemplate{
			designWorkFlow.WorkFlowTemplate,
			xuqiuWorkFlow.WorkFlowTemplate,
			bugWorkFlow.WorkFlowTemplate,
			zouchaWorkFlow.WorkFlowTemplate,
			subTaskFlow.WorkFlowTemplate,
			issueStateFlow.WorkFlowTemplate,
		},
		WorkFlow: []*wf.WorkFlow{
			designWorkFlow.WorkFlow,
			xuqiuWorkFlow.WorkFlow,
			bugWorkFlow.WorkFlow,
			zouchaWorkFlow.WorkFlow,
			subTaskFlow.WorkFlow,
			issueStateFlow.WorkFlow,
		},
	}

	return result, nil
}

func (uc *SpaceUsecase) GetMySpace(ctx context.Context, userId int64, spaceId int64) (*rsp.SpaceDetail, error) {

	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, userId)
	if err != nil { //不是成员，无权查看
		return nil, errs.NoPerm(ctx)
	}

	space, err := uc.repo.GetSpaceDetail(ctx, spaceId)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_DB_QUERY_FAIL)
	}

	spaceConfig, err := uc.repo.GetSpaceConfig(ctx, spaceId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	if spaceConfig.Id == 0 {
		err := uc.repo.CreateSpaceConfig(ctx, spaceConfig)
		if err != nil {
			return nil, errs.Internal(ctx, err)
		}
	}

	return &rsp.SpaceDetail{
		Id:                           space.Id,
		UserId:                       space.UserId,
		SpaceName:                    space.SpaceName,
		Describe:                     space.Describe,
		Notify:                       space.Notify,
		WorkingDay:                   string(spaceConfig.WorkingDay),
		CommentDeletable:             spaceConfig.CommentDeletable,
		CommentDeletableWhenArchived: spaceConfig.CommentDeletableWhenArchived,
		CommentShowPos:               spaceConfig.CommentShowPos,
		CreatedAt:                    space.CreatedAt,
		RoleId:                       member.RoleId,
		IsMember:                     member.Id != 0,
	}, nil
}

func (uc *SpaceUsecase) SetMySpaceBaseInfo(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, spaceName string, describe string) (*db.Space, error) {

	//检查创建空间的用户是否存在, 不是创建者不能修改
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if err != nil {
		//创建空间失败-创建人的信息不存在
		err := errs.New(ctx, comm.ErrorCode_PERMISSION_INSUFFICIENT_DATA_PERMISSIONS)
		return nil, err
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE,
	})

	if err != nil {
		return nil, err
	}

	//检查创建空间的用户是否存在, 不是创建者不能修改
	space, err := uc.repo.GetSpaceDetail(ctx, spaceId)
	if err != nil {
		//创建空间失败-创建人的信息不存在
		err := errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
		return nil, err
	}

	if space.SpaceName != spaceName {
		//检查空间名称是否存在
		isExist, err := uc.repo.IsExistBySpaceName(ctx, space.UserId, spaceName)
		if err != nil {
			//创建空间失败-数据库查询异常
			err := errs.New(ctx, comm.ErrorCode_DB_QUERY_FAIL)
			return nil, err
		}

		if isExist {
			//创建空间失败-空间名称重复
			err := errs.New(ctx, comm.ErrorCode_SPACE_NAME_EXIST)
			return nil, err
		}
	}

	err = uc.spaceServie.UpdateSpaceInfo(ctx, space, spaceName, describe, oper)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	//填充数据实体
	err = uc.repo.SaveSpace(ctx, space)
	if err != nil {
		//创建空间失败
		err := errs.Business(ctx, "修改项目信息失败")
		return nil, err
	}

	uc.domainMessageProducer.Send(ctx, space.GetMessages())

	return convert.SpaceEntityToPo(space), nil
}

func (uc *SpaceUsecase) SetMySpaceDescribe(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, describe string) (*db.Space, error) {

	//检查创建空间的用户是否存在, 不是创建者不能修改
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if err != nil {
		//创建空间失败-创建人的信息不存在
		err := errs.New(ctx, comm.ErrorCode_PERMISSION_INSUFFICIENT_DATA_PERMISSIONS)
		return nil, err
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE,
	})

	if err != nil {
		return nil, err
	}

	//检查创建空间的用户是否存在, 不是创建者不能修改
	space, err := uc.repo.GetSpaceDetail(ctx, spaceId)
	if err != nil {
		//创建空间失败-创建人的信息不存在
		err := errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
		return nil, err
	}

	space.UpdateDescribe(describe, oper)

	//填充数据实体
	err = uc.repo.SaveSpace(ctx, space)
	if err != nil {
		//创建空间失败
		err := errs.Business(ctx, "修改项目信息失败")
		return nil, err
	}

	uc.domainMessageProducer.Send(ctx, space.GetMessages())

	return convert.SpaceEntityToPo(space), nil
}

func (uc *SpaceUsecase) SetName(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, spaceName string) (any, error) {

	//检查创建空间的用户是否存在, 不是创建者不能修改
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE,
	})

	if err != nil {
		return nil, err
	}

	//检查创建空间的用户是否存在, 不是创建者不能修改
	space, err := uc.repo.GetSpace(ctx, spaceId)
	if err != nil {
		return nil, errs.New(ctx, comm.ErrorCode_SPACE_INFO_WRONG)
	}

	err = uc.spaceServie.UpdateSpaceName(ctx, space, spaceName, oper)
	if err != nil {
		err := errs.Business(ctx, "修改项目信息失败")
		return nil, err
	}

	err = uc.repo.SaveSpace(ctx, space)
	if err != nil {
		err := errs.Business(ctx, "修改项目信息失败")
		return nil, err
	}

	uc.domainMessageProducer.Send(ctx, space.GetMessages())

	return nil, nil
}

func (uc *SpaceUsecase) GetMySpaceList(ctx context.Context, ownerId int64) (*v1.SpaceListReplyData, error) {
	list, err := uc.repo.GetUserSpaceList(ctx, ownerId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	spaceIds := stream.Map(list, func(v *space.Space) int64 {
		return v.Id
	})

	memberMap, err := uc.memberRepo.UserSpaceMemberMapFromDB(ctx, ownerId, spaceIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	slices.SortFunc(list, func(a, b *space.Space) int {
		rankCmp := cmp.Compare(memberMap[b.Id].Ranking, memberMap[a.Id].Ranking)
		if rankCmp != 0 {
			return rankCmp
		}
		return cmp.Compare(b.Id, a.Id)
	})

	configMap, err := uc.repo.SpaceConfigMap(ctx, spaceIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	memberNum, err := uc.memberRepo.SpaceMemberNumMapBySpaceIds(ctx, spaceIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var items []*rsp.Space
	for _, v := range list {
		items = append(items, &rsp.Space{
			Id:         v.Id,
			UserId:     v.UserId,
			SpaceName:  v.SpaceName,
			RoleId:     memberMap[v.Id].RoleId,
			Ranking:    memberMap[v.Id].Ranking,
			Notify:     memberMap[v.Id].Notify,
			WorkingDay: string(configMap[v.Id].WorkingDay),
			MemberNum:  memberNum[v.Id],
			CreatedAt:  v.CreatedAt,
			IsMember:   true,
		})
	}

	userInfo, err := uc.userRepo.GetUserByUserId(ctx, ownerId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}
	if userInfo.IsEnterpriseAdmin() {
		allSpaceList, err := uc.repo.GetAllSpace(ctx)
		if err != nil {
			return nil, errs.Internal(ctx, err)
		}

		slices.SortedFunc(slices.Values(allSpaceList), func(a, b *space.Space) int {
			return cmp.Compare(a.Id, b.Id)
		})

		for _, v := range allSpaceList {
			if slices.ContainsFunc(items, func(s *rsp.Space) bool {
				return v.Id == s.Id
			}) {
				continue
			}

			items = append(items, &rsp.Space{
				Id:         v.Id,
				UserId:     v.UserId,
				SpaceName:  v.SpaceName,
				RoleId:     consts.MEMBER_ROLE_WATCHER,
				Ranking:    0,
				Notify:     0,
				WorkingDay: "[]",
				MemberNum:  0,
				CreatedAt:  v.CreatedAt,
			})
		}
	}

	return &v1.SpaceListReplyData{List: items}, nil
}

func (uc *SpaceUsecase) DelSpace(ctx context.Context, oper shared.Oper, spaceId int64, scene string) (*db.Space, error) {
	uid := oper.GetId()

	//检查创建空间的用户是否存在, 不是创建者不能删除
	space, err := uc.repo.GetSpace(ctx, spaceId)
	if err != nil {
		return nil, nil // 删除项目时，项目不存在的情况，不返回错误
	}

	if uid == space.UserId { //当前用户是项目创建人
		member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
		if err != nil {
			return nil, errs.NoPerm(ctx)
		}

		// 验证权限
		err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
			SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
			Perm:              consts.PERM_DELETE_SPACE,
		})
		if err != nil {
			return nil, err
		}
	} else if scene == "admin" {
		// 检查是否是系统管理员 && 当前用户角色权限大于
		userMap, _ := uc.userRepo.UserMap(ctx, []int64{uid, space.UserId})

		curUser := userMap[uid]
		spaceOwner := userMap[space.UserId]

		if !curUser.IsSystemAdmin() || !curUser.RoleGreaterThan(spaceOwner.Role) {
			return nil, errs.NoPerm(ctx)
		}
	} else {
		return nil, errs.NoPerm(ctx)
	}

	memberIds, err := uc.memberRepo.GetSpaceAllMemberIds(ctx, spaceId)
	if err != nil {
		return nil, err
	}

	space.OnDelete(oper)

	err = uc.tm.InTx(ctx, func(ctx context.Context) error {

		err = uc.repo.DelSpace(ctx, spaceId)
		if err != nil {
			return err
		}

		//删除空间配置
		err = uc.repo.DelSpaceConfig(ctx, spaceId)
		if err != nil {
			return err
		}

		//删除空间工作项
		_, err = uc.workItemRepo.DelSpaceWorkItemBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		//删除空间工作项-流程信息
		_, err = uc.workItemRepo.DelSpaceWorkItemFlowBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		//删除空间工作项-流程角色
		_, err = uc.workItemRepo.DelSpaceWorkItemFlowRoleBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		//删除空间模块
		_, err = uc.workObjectRepo.DelWorkObjectBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		//删除空间成员
		_, err = uc.memberRepo.DelSpaceMemberBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		//软 删除文件
		err = uc.fileInfoRepo.SoftDelFileBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		//删除空间标签
		err = uc.tagRepo.DelSpaceTagBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		//删除版本
		_, err = uc.workVersionRepo.DelWorkVersionBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		//删除工作状态
		err = uc.workItemStatusRepo.DelWorkItemStatusBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		//删除工作项类型
		err = uc.workItemTypeRepo.DelBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		//删除工作项类型-角色
		err = uc.workItemRoleRepo.DelWorkItemRoleBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		//删除工作项类型-工作流
		err = uc.workFlowRepo.DelWorkFlowBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		//删除工作类型-工作流模版
		err = uc.workFlowRepo.DelWorkFlowTemplateBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		// 清理视图
		err = uc.viewRepo.DeleteGlobalViewBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}
		err = uc.viewRepo.DeleteUserViewBySpaceId(ctx, spaceId)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, errs.Business(ctx, "删除空间失败")
	}

	uc.domainMessageProducer.Send(ctx, space.GetMessages())

	//这个事件只能额外加一个了 特殊处理了, 专门给推送消息用的
	bus.Emit(notify.Event_DeleteSpace, &event.DeleteSpace{
		Event:     notify.Event_DeleteSpace,
		Space:     space,
		Operator:  oper.GetId(),
		MemberIds: memberIds,
	})

	return convert.SpaceEntityToPo(space), nil
}

func (uc *SpaceUsecase) TransferSpaceOwnership(ctx context.Context, uid int64, spaceId, srcUserId, dstUserId int64) error {
	user, err := uc.userRepo.GetUserByUserId(ctx, uid)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	space, err := uc.repo.GetSpace(ctx, spaceId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	// 权限检查
	if !user.IsSystemAdmin() && // 不是系统管理员
		(uid != srcUserId || !space.IsCreator(uid)) { //不是项目创建者
		return errs.NoPerm(ctx)
	}

	srcMember, err := uc.memberRepo.GetSpaceMember(ctx, space.Id, srcUserId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	dstMember, err := uc.memberRepo.GetSpaceMember(ctx, space.Id, dstUserId)
	if err != nil {
		return errs.Business(ctx, comm.ErrorCode_SPACE_MEMBER_WRONG)
	}

	logOper := shared.UserOper(srcUserId)

	//变更创建人
	space.TransferSpace(dstMember.UserId, logOper)
	//把对方标记为项目创建者
	dstMember.ChangeRoleId(consts.MEMBER_ROLE_SPACE_CREATOR, logOper)
	//把自己标记为可管理
	srcMember.ChangeRoleId(consts.MEMBER_ROLE_MANAGER, logOper)

	err = uc.tm.InTx(ctx, func(ctx context.Context) error {
		err = uc.repo.SaveSpace(ctx, space)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = uc.memberRepo.SaveSpaceMember(ctx, srcMember)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = uc.memberRepo.SaveSpaceMember(ctx, dstMember)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	//派发领域日志
	uc.domainMessageProducer.Send(ctx, space.GetMessages())

	return nil
}

func (uc *SpaceUsecase) SetNotify(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, n int64) error {
	//检查创建空间的用户是否存在, 不是创建者不能修改
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	space, err := uc.repo.GetSpace(ctx, spaceId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE,
	})
	if err != nil {
		return err
	}

	space.UpdateNotify(n, oper)

	//填充数据实体
	err = uc.repo.SaveSpace(ctx, space)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	//派发领域日志
	uc.domainMessageProducer.Send(ctx, space.GetMessages())

	return nil
}

func (uc *SpaceUsecase) SetWorkingDay(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, weekDays []int64) error {
	//检查创建空间的用户是否存在, 不是创建者不能修改
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	//TODO 这里的设计有问题，需要调整
	space, err := uc.repo.GetSpace(ctx, spaceId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	spaceConfig, err := uc.repo.GetSpaceConfig(ctx, space.Id)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE,
	})

	if err != nil {
		return err
	}

	spaceConfig.UpdateWorkingDay(utils.ToStrArray(weekDays), oper)

	//填充数据实体
	err = uc.repo.SaveSpaceConfig(ctx, spaceConfig)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	//派发领域日志
	uc.domainMessageProducer.Send(ctx, spaceConfig.GetMessages())

	return nil
}

func (uc *SpaceUsecase) SetCommentDeletable(ctx context.Context, uid int64, spaceId int64, deletable int64) error {
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	space, err := uc.repo.GetSpace(ctx, spaceId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	config, err := uc.repo.GetSpaceConfig(ctx, space.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	config.UpdateCommentDeletable(deletable, utils.GetLoginUser(ctx))

	//填充数据实体
	err = uc.repo.SaveSpaceConfig(ctx, config)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	uc.domainMessageProducer.Send(ctx, config.GetMessages())

	return nil
}

func (uc *SpaceUsecase) SearchWorkItem(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, keyword string) (*v1.SearchWorkItemReplyData, error) {
	//检查创建空间的用户是否存在, 不是创建者不能修改
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if err != nil || member.Id == 0 {
		return nil, errs.NoPerm(ctx)
	}

	list, err := uc.searchRepo.SearchByName(ctx, spaceId, keyword)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	var items []*v1.SearchWorkItemReplyData_Item
	for _, v := range list {
		items = append(items, &v1.SearchWorkItemReplyData_Item{
			Id:   v.Id,
			Name: v.WorkItemName,
		})
	}

	return &v1.SearchWorkItemReplyData{Items: items}, nil
}

func (uc *SpaceUsecase) GetWorkItemTypes(ctx context.Context, oper *utils.LoginUserInfo, req *pb.GetWorkItemTypesRequest) ([]*pb.GetWorkItemTypesReplyData_Item, error) {
	//检查创建空间的用户是否存在, 不是创建者不能修改
	member, err := uc.memberRepo.GetSpaceMember(ctx, req.SpaceId, oper.UserId)
	if err != nil || member.Id == 0 {
		return nil, errs.NoPerm(ctx)
	}

	info, err := uc.workItemTypeRepo.QWorkItemTypeInfo(ctx, query.WorkItemTypeInfoQuery{
		SpaceId: req.SpaceId,
	})
	if err != nil || member.Id == 0 {
		return nil, errs.NoPerm(ctx)
	}

	var items []*pb.GetWorkItemTypesReplyData_Item
	for _, v := range info.Types() {
		items = append(items, &pb.GetWorkItemTypesReplyData_Item{
			Id:       v.Id,
			Name:     v.Name,
			Key:      v.Key,
			FlowMode: string(v.FlowMode),
			SpaceId:  int64(v.SpaceId),
		})
	}

	return items, nil
}

func (uc *SpaceUsecase) GetTempConfig(ctx context.Context, uid, spaceId int64, keys []string) (map[string]string, error) {
	//检查创建空间的用户是否存在, 不是创建者不能修改
	_, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	config := uc.repo.GetTempConfig(ctx, spaceId, keys...)

	key := consts.SpaceTempConfigKey_spaceWorkbenchCountConditions
	if slices.Contains(keys, key) {
		oldConfVal := config[key]
		confVal, needUpdate := uc.handleOldConditionGroupData(ctx, spaceId, oldConfVal)
		config[key] = confVal

		if needUpdate || oldConfVal == "" {
			uc.repo.SetTempConfig(ctx, spaceId, config)
		}
	}

	return config, nil
}

func (uc *SpaceUsecase) handleOldConditionGroupData(ctx context.Context, spaceId int64, confVal string) (string, bool) {

	statusList, err := uc.workItemStatusRepo.GetWorkItemStatusItemsBySpaceIds(ctx, []int64{spaceId})
	if err != nil {
		return "", false
	}

	if confVal == "" {
		return space_temp_config.GetDefaultByKey(consts.SpaceTempConfigKey_spaceWorkbenchCountConditions, statusList), false
	}

	var needUpdate bool

	var cond []space_temp_config.CondItem
	json.Unmarshal([]byte(confVal), &cond)

	for _, v := range cond {
		var newConditions []*search_es.Condition
		for _, condition := range v.Value.Conditions {
			switch condition.Field {
			case search_es.WorkItemStatusField:
				needUpdate = true
				condition.Field = search_es.WorkItemStatusIdField
				condition.SetAttr("flow_scope", string(consts.FlowScope_Workflow))
				var newValues []any
				for i := 0; i < len(condition.Values); i++ {
					item := statusList.GetStatusByVal(cast.ToString(condition.Values[i]))
					if item == nil {
						continue
					}

					newValues = append(newValues, cast.ToString(item.Id))
				}
				condition.Values = newValues
			}
			newConditions = append(newConditions, condition)
		}
		v.Value.Conditions = newConditions
	}

	return utils.ToJSON(cond), needUpdate
}

func (uc *SpaceUsecase) SetTempConfig(ctx context.Context, uid, spaceId int64, configs map[string]string) error {
	//检查创建空间的用户是否存在, 不是创建者不能修改
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	for key, _ := range configs {
		if !space_temp_config.IsValidKey(key) {
			return errs.Business(ctx, "wrong key")
		}
	}

	oldConfigs, _ := uc.GetTempConfig(ctx, uid, spaceId, stream.Keys(configs))

	err = uc.repo.SetTempConfig(ctx, spaceId, configs)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	msg := &domain_message.SetSpaceTempConfig{
		DomainMessageBase: shared.DomainMessageBase{
			Oper:     utils.GetLoginUser(ctx),
			OperTime: time.Now(),
		},
		SpaceId:   spaceId,
		OldValues: oldConfigs,
		NewValues: configs,
	}

	uc.domainMessageProducer.Send(ctx, []shared.DomainMessage{
		msg,
	})

	return nil
}

func (uc *SpaceUsecase) DelTempConfig(ctx context.Context, uid, spaceId int64, keys []string) error {
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	err = uc.repo.DelTempConfig(ctx, spaceId, keys...)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	return nil
}

func (uc *SpaceUsecase) Copy(ctx context.Context, oper shared.Oper, srcSpaceId int64, spaceName, spaceDescribe string) (*pb.CopyReply_Data, error) {
	uid := oper.GetId()

	srcSpace, err := uc.repo.GetSpaceDetail(ctx, srcSpaceId)
	if errs.IsDbRecordNotFoundErr(err) {
		return nil, errs.Custom(ctx, comm.ErrorCode_SPACE_NOT_EXIST, "项目不存在")
	}
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	//检查创建空间的用户是否存在, 不是创建者不能修改
	member, err := uc.memberRepo.GetSpaceMember(ctx, srcSpaceId, uid)
	if err != nil || member.Id == 0 {
		return nil, errs.NoPerm(ctx)
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_COPY_SPACE,
	})
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	spaceConfig, err := uc.repo.GetSpaceConfig(ctx, srcSpaceId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	spaceMembers, err := uc.memberRepo.GetSpaceMemberBySpaceId(ctx, srcSpaceId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	spaceRole, err := uc.workItemRoleRepo.GetWorkItemRoles(ctx, srcSpaceId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	spaceStatus, err := uc.workItemStatusRepo.GetWorkItemStatusItemsBySpace(ctx, srcSpaceId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	workItemType, err := uc.workItemTypeRepo.GetWorkItemTypeBySpaceId(ctx, srcSpaceId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	spaceWorkFlow, err := uc.workFlowRepo.GetWorkFlowBySpaceId(ctx, srcSpaceId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	flowAndWorkItemType := stream.InnerJoin2Map(spaceWorkFlow, workItemType, func(flow *wf.WorkFlow, workItemType *witem_type.WorkItemType) bool {
		return flow.WorkItemTypeId == workItemType.Id
	})

	templateIds := stream.Map(spaceWorkFlow, func(v *wf.WorkFlow) int64 {
		return v.LastTemplateId
	})

	templates, err := uc.workFlowRepo.GetFlowTemplateByIds(ctx, templateIds)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	flowAndTemplates := stream.InnerJoin2Map(spaceWorkFlow, templates, func(flow *wf.WorkFlow, template *wf.WorkFlowTemplate) bool {
		return flow.LastTemplateId == template.Id
	})

	workObjects, err := uc.workObjectRepo.GetSpaceWorkObjectBySpaceIds(ctx, []int64{srcSpaceId})
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	versions, err := uc.workVersionRepo.GetSpaceWorkVersionBySpaceId(ctx, srcSpaceId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	dstSpace, err := uc.spaceServie.CreateSpace(ctx, oper.GetId(), spaceName, spaceDescribe, oper)
	if err != nil {
		return nil, err
	}

	userIds := stream.Map(spaceMembers, func(v *member_domain.SpaceMember) int64 {
		return v.UserId
	})

	err = uc.tm.InTx(ctx, func(ctx context.Context) error {
		spaceId := dstSpace.Id
		dstSpace.Notify = srcSpace.Notify
		if spaceName != "" {
			dstSpace.SpaceName = spaceName
		}
		if spaceDescribe != "" {
			dstSpace.Describe = spaceDescribe
		}

		spaceConfig.Id = 0
		spaceConfig.SpaceId = spaceId

		for _, v := range spaceMembers {
			v.Id = 0
			v.SpaceId = spaceId
			v.HistoryRoleId = consts.MEMBER_ROLE_WATCHER
			v.Notify = 1
			v.Ranking = time.Now().UnixMilli()

			if v.RoleId == consts.MEMBER_ROLE_SPACE_CREATOR {
				v.RoleId = consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER
				v.HistoryRoleId = consts.MEMBER_ROLE_WATCHER
			}

			if v.UserId == uid {
				v.RoleId = consts.MEMBER_ROLE_SPACE_CREATOR
				v.HistoryRoleId = consts.MEMBER_ROLE_SPACE_SUPPER_MANAGER
			}

			v.CreatedAt = 0
			v.UpdatedAt = 0
		}

		roleMap := map[int64]int64{}
		for _, v := range spaceRole {
			oldId := v.Id
			newId := uc.idService.NewId(ctx, consts.BusinessId_Type_WorkItemRole).Id
			roleMap[oldId] = newId

			v.Id = newId
			v.SpaceId = spaceId
			v.Uuid = uuid.NewString()
			v.CreatedAt = 0
			v.UpdatedAt = 0
		}

		statusMap := map[int64]int64{}
		for _, v := range spaceStatus {
			oldId := v.Id
			newId := uc.idService.NewId(ctx, consts.BusinessId_Type_WorkItemStatus).Id
			statusMap[oldId] = newId

			v.Id = newId
			v.SpaceId = spaceId
			v.Uuid = uuid.NewString()
			v.CreatedAt = 0
			v.UpdatedAt = 0
		}

		for _, v := range workItemType {
			v.Id = uc.idService.NewId(ctx, consts.BusinessId_Type_WorkItemType).Id
			v.SpaceId = witem_type.SpaceId(spaceId)
			v.Uuid = uuid.NewString()
			v.CreatedAt = 0
			v.UpdatedAt = 0
		}

		for _, v := range templates {
			v.Id = uc.idService.NewId(ctx, consts.BusinessId_Type_WorkFlowTemplate).Id
			v.SpaceId = spaceId
			v.Uuid = uuid.NewString()
			v.CreatedAt = 0
			v.UpdatedAt = 0

			// 修改流程节点
			if v.WorkFLowConfig != nil {
				for _, node := range v.WorkFLowConfig.Nodes {
					// 调整角色
					for _, v := range node.Owner.OwnerRole {
						v.Id = cast.ToString(roleMap[cast.ToInt64(v.Id)])
					}

					// 调整状态
					for _, v := range node.OnPass {
						v.TargetSubState.Id = cast.ToString(statusMap[cast.ToInt64(v.TargetSubState.Id)])
					}
					for _, v := range node.OnReach {
						v.TargetSubState.Id = cast.ToString(statusMap[cast.ToInt64(v.TargetSubState.Id)])
					}
				}
			}

			// 修改状态流转
			if v.StateFlowConfig != nil {
				for _, node := range v.StateFlowConfig.StateFlowNodes {
					// 调整角色
					if node.Owner != nil {
						for _, v := range node.Owner.OwnerRole {
							v.Id = cast.ToString(roleMap[cast.ToInt64(v.Id)])
						}
					}

					// 调整状态
					node.SubStateId = cast.ToString(statusMap[cast.ToInt64(node.SubStateId)])
				}
			}
		}

		for _, v := range spaceWorkFlow {
			v.Id = uc.idService.NewId(ctx, consts.BusinessId_Type_WorkFlow).Id

			template := flowAndTemplates[v]
			template.WorkFlowId = v.Id
			template.WorkItemTypeId = flowAndWorkItemType[v].Id

			v.LastTemplateId = template.Id
			v.WorkItemTypeId = flowAndWorkItemType[v].Id
			v.SpaceId = spaceId
			v.Version = 1
			v.Uuid = uuid.NewString()
			v.CreatedAt = 0
			v.UpdatedAt = 0
		}

		for _, v := range workObjects {
			v.Id = uc.idService.NewId(ctx, consts.BusinessId_Type_SpaceWorkObject).Id
			v.SpaceId = spaceId
			v.WorkObjectGuid = uuid.NewString()
			v.CreatedAt = 0
			v.UpdatedAt = 0
		}

		for _, v := range versions {
			v.Id = uc.idService.NewId(ctx, consts.BusinessId_Type_SpaceWorkVersion).Id
			v.SpaceId = spaceId
			v.CreatedAt = 0
			v.UpdatedAt = 0
		}

		err = uc.repo.CreateSpace(ctx, dstSpace)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = uc.repo.CreateSpaceConfig(ctx, spaceConfig)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = uc.memberRepo.AddSpaceMembers(ctx, spaceMembers)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = uc.workItemRoleRepo.CreateWorkItemRoles(ctx, spaceRole)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = uc.workItemStatusRepo.CreateWorkItemStatusItems(ctx, spaceId, spaceStatus)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		for _, v := range workItemType {
			err = uc.workItemTypeRepo.CreateWorkItemType(ctx, v)
			if err != nil {
				return errs.Internal(ctx, err)
			}
		}

		err = uc.workFlowRepo.CreateWorkFlowTemplates(ctx, templates)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		err = uc.workFlowRepo.CreateWorkFlows(ctx, spaceWorkFlow)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		for _, v := range workObjects {
			err = uc.workObjectRepo.CreateSpaceWorkObject(ctx, v)
			if err != nil {
				return errs.Internal(ctx, err)
			}
		}

		for _, v := range versions {
			err = uc.workVersionRepo.CreateSpaceWorkVersion(ctx, v)
			if err != nil {
				return errs.Internal(ctx, err)
			}
		}

		// 创建公共视图
		globalViews, err := uc.viewService.InitSpacePublicView(ctx, dstSpace.Id, spaceStatus)
		if err != nil {
			return errs.Internal(ctx, err)
		}
		err = uc.viewRepo.CreateGlobalViews(ctx, globalViews)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		// 创建个人视图
		views, err := uc.viewService.InitUserGlobalView(ctx, dstSpace.Id, userIds)
		if err != nil {
			return errs.Internal(ctx, err)
		}
		err = uc.viewRepo.CreateUserViews(ctx, views)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		return nil
	})

	msg := &domain_message.CreateSpace{
		SrcSpaceId:   srcSpace.Id,
		SrcSpaceName: srcSpace.SpaceName,
		SpaceId:      dstSpace.Id,
		SpaceName:    dstSpace.SpaceName,
	}

	msg.SetOper(utils.GetLoginUser(ctx), time.Now())
	uc.domainMessageProducer.Send(ctx, shared.DomainMessages{
		msg,
	})

	return &pb.CopyReply_Data{
		SpaceId: dstSpace.Id,
	}, err
}

func (uc *SpaceUsecase) SetCommentDeletableWhenArchived(ctx context.Context, uid int64, spaceId int64, val int64) error {
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	space, err := uc.repo.GetSpace(ctx, spaceId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	config, err := uc.repo.GetSpaceConfig(ctx, space.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	config.UpdateCommentDeletableWhenArchived(val, utils.GetLoginUser(ctx))

	//填充数据实体
	err = uc.repo.SaveSpaceConfig(ctx, config)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	uc.domainMessageProducer.Send(ctx, config.GetMessages())

	return nil
}

func (uc *SpaceUsecase) SetCommentShowPos(ctx context.Context, uid int64, spaceId int64, val int64) error {
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	space, err := uc.repo.GetSpace(ctx, spaceId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	config, err := uc.repo.GetSpaceConfig(ctx, space.Id)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	config.UpdateCommentShowPos(val, utils.GetLoginUser(ctx))

	//填充数据实体
	err = uc.repo.SaveSpaceConfig(ctx, config)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	uc.domainMessageProducer.Send(ctx, config.GetMessages())

	return nil
}
