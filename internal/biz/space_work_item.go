package biz

import (
	"context"
	"fmt"
	"go-cs/api/comm"
	"go-cs/api/notify"
	vo "go-cs/internal/bean/vo"
	"go-cs/internal/bean/vo/event"
	"go-cs/internal/bean/vo/rsp"
	"go-cs/internal/biz/command"
	"go-cs/internal/biz/query"
	"go-cs/internal/consts"
	file_repo "go-cs/internal/domain/file_info/repo"
	file_service "go-cs/internal/domain/file_info/service"
	perm_facade "go-cs/internal/domain/perm/facade"
	perm_service "go-cs/internal/domain/perm/service"
	domain_message "go-cs/internal/domain/pkg/message"
	space_repo "go-cs/internal/domain/space/repo"
	space_file_repo "go-cs/internal/domain/space_file_info/repo"
	space_file_service "go-cs/internal/domain/space_file_info/service"
	member_repo "go-cs/internal/domain/space_member/repo"
	member_service "go-cs/internal/domain/space_member/service"
	workObj_repo "go-cs/internal/domain/space_work_object/repo"
	user_repo "go-cs/internal/domain/user/repo"
	wf_repo "go-cs/internal/domain/work_flow/repo"
	"go-cs/pkg/bus"
	"go-cs/pkg/stream"
	"time"

	witem "go-cs/internal/domain/work_item"
	witem_facade "go-cs/internal/domain/work_item/facade"
	witem_repo "go-cs/internal/domain/work_item/repo"
	witem_service "go-cs/internal/domain/work_item/service"

	witem_role_repo "go-cs/internal/domain/work_item_role/repo"
	witem_status_repo "go-cs/internal/domain/work_item_status/repo"
	witem_status_service "go-cs/internal/domain/work_item_status/service"
	witem_type_repo "go-cs/internal/domain/work_item_type/repo"

	tag_repo "go-cs/internal/domain/space_tag/repo"
	tag_service "go-cs/internal/domain/space_tag/service"
	comment_repo "go-cs/internal/domain/space_work_item_comment/repo"
	comment_service "go-cs/internal/domain/space_work_item_comment/service"

	workVersion_repo "go-cs/internal/domain/space_work_version/repo"

	statics_repo "go-cs/internal/domain/statics/repo"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"

	cache "github.com/Code-Hex/go-generics-cache"
	goCache "github.com/Code-Hex/go-generics-cache"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/panjf2000/ants/v2"
	"github.com/spf13/cast"
)

type SpaceWorkItemUsecase struct {
	repo              witem_repo.WorkItemRepo
	witemEsRepo       witem_repo.WorkItemEsRepo
	userRepo          user_repo.UserRepo
	spaceRepo         space_repo.SpaceRepo
	spaceMemberRepo   member_repo.SpaceMemberRepo
	spaceFileInfoRepo space_file_repo.SpaceFileInfoRepo
	fileInfoRepo      file_repo.FileInfoRepo
	tagRepo           tag_repo.SpaceTagRepo
	staticsRepo       statics_repo.StaticsRepo

	spaceWorkObjectRepo workObj_repo.SpaceWorkObjectRepo
	workVersionRepo     workVersion_repo.SpaceWorkVersionRepo
	commentRepo         comment_repo.SpaceWorkItemCommentRepo

	commentPool *ants.Pool
	log         *log.Helper
	tm          trans.Transaction
	cacheLocker *goCache.Cache[string, interface{}]

	workItemStatusRepo witem_status_repo.WorkItemStatusRepo
	workItemTypeRepo   witem_type_repo.WorkItemTypeRepo
	workFlowRepo       wf_repo.WorkFlowRepo
	workItemRoleRepo   witem_role_repo.WorkItemRoleRepo

	permService                 *perm_service.PermService
	workItemService             *witem_service.WorkItemService
	spaceMemberService          *member_service.SpaceMemberService
	spaceTagService             *tag_service.SpaceTagService
	spaceFileInfoService        *space_file_service.SpaceFileInfoService
	spaceWorkItemCommentService *comment_service.SpaceWorkItemCommentService
	fileService                 *file_service.FileInfoService
	witemStatusService          *witem_status_service.WorkItemStatusService

	createSpaceTaskCommand    *command.CreateSpaceTaskCmd
	createSpaceSubTaskCommand *command.CreateSpaceSubTaskCmd
	addWorkItemCommentCommand *command.AddWorkItemCommentCmd
	workItemDetailQuery       *query.WorkItemDetailQuery

	domainMessageProducer *domain_message.DomainMessageProducer
}

func NewSpaceWorkItemUsecase(
	repo witem_repo.WorkItemRepo,
	witemEsRepo witem_repo.WorkItemEsRepo,
	workVersionRepo workVersion_repo.SpaceWorkVersionRepo,
	staticsRepo statics_repo.StaticsRepo,

	tm trans.Transaction,
	spaceFileInfoRepo space_file_repo.SpaceFileInfoRepo,
	fileInfoRepo file_repo.FileInfoRepo,
	spaceWorkObjectRepo workObj_repo.SpaceWorkObjectRepo,
	tagRepo tag_repo.SpaceTagRepo,
	workItemStatusRepo witem_status_repo.WorkItemStatusRepo,
	workItemTypeRepo witem_type_repo.WorkItemTypeRepo,
	workFlowRepo wf_repo.WorkFlowRepo,
	userRepo user_repo.UserRepo,
	spaceRepo space_repo.SpaceRepo,
	spaceMemberRepo member_repo.SpaceMemberRepo,
	commentRepo comment_repo.SpaceWorkItemCommentRepo,
	workItemRoleRepo witem_role_repo.WorkItemRoleRepo,

	permService *perm_service.PermService,
	workItemService *witem_service.WorkItemService,
	spaceMemberService *member_service.SpaceMemberService,
	spaceTagService *tag_service.SpaceTagService,
	spaceFileInfoService *space_file_service.SpaceFileInfoService,
	spaceWorkItemCommentService *comment_service.SpaceWorkItemCommentService,
	fileService *file_service.FileInfoService,
	witemStatusService *witem_status_service.WorkItemStatusService,

	createSpaceTaskCommand *command.CreateSpaceTaskCmd,
	createSpaceSubTaskCommand *command.CreateSpaceSubTaskCmd,
	addWorkItemCommentCommand *command.AddWorkItemCommentCmd,
	workItemDetailQuery *query.WorkItemDetailQuery,

	domainMessageProducer *domain_message.DomainMessageProducer,

	logger log.Logger,

) *SpaceWorkItemUsecase {

	moduleName := "SpaceWorkItemUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &SpaceWorkItemUsecase{
		log:                 hlog,
		repo:                repo,
		witemEsRepo:         witemEsRepo,
		userRepo:            userRepo,
		spaceRepo:           spaceRepo,
		spaceMemberRepo:     spaceMemberRepo,
		spaceFileInfoRepo:   spaceFileInfoRepo,
		spaceWorkObjectRepo: spaceWorkObjectRepo,
		staticsRepo:         staticsRepo,
		commentRepo:         commentRepo,
		workItemRoleRepo:    workItemRoleRepo,
		workVersionRepo:     workVersionRepo,

		tagRepo:                     tagRepo,
		fileInfoRepo:                fileInfoRepo,
		tm:                          tm,
		cacheLocker:                 goCache.New(goCache.AsFIFO[string, interface{}]()),
		workItemStatusRepo:          workItemStatusRepo,
		workItemTypeRepo:            workItemTypeRepo,
		workFlowRepo:                workFlowRepo,
		permService:                 permService,
		workItemService:             workItemService,
		spaceMemberService:          spaceMemberService,
		spaceTagService:             spaceTagService,
		spaceFileInfoService:        spaceFileInfoService,
		spaceWorkItemCommentService: spaceWorkItemCommentService,
		fileService:                 fileService,
		witemStatusService:          witemStatusService,

		createSpaceTaskCommand:    createSpaceTaskCommand,
		createSpaceSubTaskCommand: createSpaceSubTaskCommand,
		addWorkItemCommentCommand: addWorkItemCommentCommand,
		workItemDetailQuery:       workItemDetailQuery,

		domainMessageProducer: domainMessageProducer,
	}
}

func (s *SpaceWorkItemUsecase) CreateTask(ctx context.Context, oper *utils.LoginUserInfo, in vo.CreateSpaceWorkItemVoV2) (int64, error) {

	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, in.SpaceId, oper.UserId)
	if member == nil || err != nil {
		return 0, errs.NoPerm(ctx)
	}

	//检查操作权限
	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CREATE_SPACE_WORK_ITEM,
	})
	if err != nil {
		return 0, err
	}

	return s.createSpaceTaskCommand.Execute(ctx, oper, in)
}

func (s *SpaceWorkItemUsecase) CreateSubTask(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, in vo.CreateSpaceWorkItemTaskVoV2) (int64, error) {

	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return 0, errs.NoPerm(ctx)
	}

	//检查操作权限
	err = s.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CREATE_SPACE_WORK_ITEM,
	})

	if err != nil {
		return 0, err
	}

	return s.createSpaceSubTaskCommand.Excute(ctx, oper, spaceId, workItemId, in)
}

func (s *SpaceWorkItemUsecase) GetWorkItemDetail(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64) (*rsp.SpaceWorkItemDetail, error) {
	return s.workItemDetailQuery.Execute(ctx, oper, spaceId, workItemId)
}

func (s *SpaceWorkItemUsecase) ConfirmWorkFlowMain(ctx context.Context, oper *utils.LoginUserInfo, workItemId int64, flowNodeCode string) ([]int64, error) {
	workItem, err := s.repo.GetWorkItem(ctx, workItemId,
		&witem_repo.WithDocOption{Directors: true},
		&witem_repo.WithOption{FlowNodes: true},
	)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	space, err := s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	//判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, workItem.SpaceId, oper.UserId)
	if member == nil || err != nil {
		return nil, errs.NoPerm(ctx)
	}

	statusInfo, err := s.workItemStatusRepo.GetWorkItemStatusInfo(ctx, space.Id)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE,
		FlowNodeCode:      flowNodeCode,
	})
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	//子任务不处理
	if !workItem.IsSameSpace(workItem.SpaceId) || workItem.IsSubTask() {
		return nil, errs.NoPerm(ctx)
	}

	flowTplt, err := s.workFlowRepo.GetFlowTemplate(ctx, workItem.WorkFlowTemplateId)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	//------- 开始处理业务
	tasks, err := s.workItemService.ConfirmWorkFlowMain(
		ctx,
		workItem,
		&witem_service.ConfirmSpaceTaskNodeState{
			NodeCode:               flowNodeCode,
			WorkFlowTemplateFacade: witem_facade.BuildWorkFlowTemplateFacade(flowTplt),
			WorkItemStatusFacade:   witem_facade.BuildWorkItemStatusFacade(statusInfo),
		},
		oper,
	)

	if err != nil {
		return nil, err
	}

	//保存
	err = s.tm.InTx(ctx, func(ctx context.Context) error {

		for _, v := range tasks {
			err = s.repo.SaveWorkItem(ctx, v)
			if err != nil {
				return err
			}
		}

		for _, v := range workItem.WorkItemFlowNodes {
			err = s.repo.SaveWorkItemFlowNode(ctx, v)
			if err != nil {
				return err
			}
		}

		for _, v := range workItem.WorkItemFlowRoles {
			err = s.repo.SaveWorkItemFlowRole(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil
	})

	s.domainMessageProducer.Send(ctx, tasks.GetMessages())

	return tasks.GetIds(), nil
}

func (s *SpaceWorkItemUsecase) SetWorkItemTag(ctx context.Context, oper *utils.LoginUserInfo, in vo.SetSpaceWorkItemTagVoV2) error {

	//判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, in.SpaceId, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	_, err = s.spaceRepo.GetSpace(ctx, in.SpaceId)
	if err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	workItem, err := s.repo.GetWorkItem(ctx, in.WorkItemId, &witem_repo.WithDocOption{Tags: true, Directors: true}, nil)
	if err != nil || !workItem.IsSameSpace(in.SpaceId) {
		return errs.NoPerm(ctx)
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE_WORK_ITEM,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	//修改标签
	s.workItemService.ChangeWorkItemTag(ctx, workItem, &witem_service.ChangeWorkItemTagRequest{
		TagAdd:                      utils.ToStrArray(in.TagAdd),
		TagRemove:                   utils.ToStrArray(in.TagRemove),
		WorkItemStatusServiceFacade: witem_facade.BuildWorkItemStatusServiceFacade(s.witemStatusService),
	}, oper)

	err = s.repo.SaveWorkItem(ctx, workItem)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	s.domainMessageProducer.Send(ctx, workItem.GetMessages())

	return nil
}

func (s *SpaceWorkItemUsecase) SetWorkItemFileInfo(ctx context.Context, oper *utils.LoginUserInfo, in vo.SetSpaceWorkItemFileInfoVoV2) error {
	workItem, err := s.repo.GetWorkItem(ctx, in.WorkItemId, &witem_repo.WithDocOption{Directors: true, Describe: true, Remark: true}, &witem_repo.WithOption{FileInfos: true})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, workItem.SpaceId, oper.UserId)
	if member == nil || err != nil {
		return errs.NoPerm(ctx)
	}

	_, err = s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE_WORK_ITEM,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	//
	addWorkItemFiles, removeWorkItemFiles, err := s.workItemService.SetWorkItemFileInfo(ctx, workItem, &witem_service.SetWorkItemFileInfoRequest{
		AddFileInfoIds:              in.FileInfoAdd,
		RemoveFileInfoIds:           in.FileInfoRemove,
		FileInfoFacade:              witem_facade.BuildFileInfoFacade(s.fileService),
		WorkItemStatusServiceFacade: witem_facade.BuildWorkItemStatusServiceFacade(s.witemStatusService),
	}, oper)
	if err != nil {
		return err
	}

	err = s.tm.InTx(ctx, func(ctx context.Context) error {

		err = s.repo.SaveWorkItem(ctx, workItem)
		if err != nil {
			return err
		}

		if len(addWorkItemFiles) > 0 {
			err = s.repo.CreateWorkItemFiles(ctx, addWorkItemFiles)
			if err != nil {
				return err
			}
		}

		for _, v := range removeWorkItemFiles {
			err = s.spaceFileInfoRepo.SoftDelSpaceWorkItemFileInfo(ctx, v.FileInfo.FileInfoId, workItem.SpaceId, in.WorkItemId)
			if err != nil {
				s.log.Error(err)
				return err
			}
		}

		return nil
	})

	if err != nil {
		return errs.Internal(ctx, err)
	}

	s.domainMessageProducer.Send(ctx, workItem.GetMessages())

	return nil

}

func (s *SpaceWorkItemUsecase) ModifyWorkItemName(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, newVal string) error {

	return s.modifyWorkItemField(ctx, oper, spaceId, workItemId, shared.PropDiffSet{
		witem.Diff_Name: newVal,
	})
}

func (s *SpaceWorkItemUsecase) ModifyWorkItemPlanTime(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, startAt, completeAt int64) error {
	nodes, err := s.repo.GetWorkItemFlowNodes(ctx, workItemId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	for _, node := range nodes {
		if node.PlanTimeHasSet() {
			if completeAt < node.PlanTime.CompleteAt {
				return errs.Business(ctx, "不可小于已选排期")
			}
			if completeAt > node.PlanTime.CompleteAt && startAt > node.PlanTime.StartAt {
				return errs.Business(ctx, "不可大于已选排期")
			}
		}
	}

	return s.modifyWorkItemField(ctx, oper, spaceId, workItemId, shared.PropDiffSet{
		witem.Diff_PlanTime: []int64{startAt, completeAt},
	})
}

func (s *SpaceWorkItemUsecase) ModifyWorkItemProcessRate(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, newVal int64) error {
	return s.modifyWorkItemField(ctx, oper, spaceId, workItemId, shared.PropDiffSet{
		witem.Diff_ProcessRate: newVal,
	})
}

func (s *SpaceWorkItemUsecase) ModifyWorkItemDescribe(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, newVal string, iconFlags []uint32) error {
	return s.modifyWorkItemField(ctx, oper, spaceId, workItemId, shared.PropDiffSet{
		witem.Diff_Describe: newVal,
	})
}

func (s *SpaceWorkItemUsecase) ModifyWorkItemPriority(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, newVal string) error {
	return s.modifyWorkItemField(ctx, oper, spaceId, workItemId, shared.PropDiffSet{
		witem.Diff_Priority: newVal,
	})
}

func (s *SpaceWorkItemUsecase) ModifyWorkItemRemark(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, newVal string) error {

	return s.modifyWorkItemField(ctx, oper, spaceId, workItemId, shared.PropDiffSet{
		witem.Diff_Remark: newVal,
	})

}

func (s *SpaceWorkItemUsecase) modifyWorkItemField(ctx context.Context, oper *utils.LoginUserInfo, spaceId, workItemId int64, propDiffs shared.PropDiffSet) error {

	_, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	// 判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if member == nil || err != nil {
		return errs.NoPerm(ctx)
	}

	//查看一下对应的工作用例是否存在
	docOpt := &witem_repo.WithDocOption{All: true}

	var opt *witem_repo.WithOption
	// 这两个字段需要重新计算iconFlag
	if propDiffs.HasProp(witem.Diff_Describe) || propDiffs.HasProp(witem.Diff_Remark) {
		opt = &witem_repo.WithOption{
			FileInfos: true,
		}
	}

	workItem, err := s.repo.GetWorkItem(ctx, workItemId, docOpt, opt)
	if err != nil {
		return errs.RepoErr(ctx, err)
	}

	if !workItem.IsSameSpace(spaceId) {
		return errs.NoPerm(ctx)
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE_WORK_ITEM,
	})
	if err != nil {
		return err
	}

	err = s.workItemService.ModifyWorkItemField(ctx, workItem, &witem_service.ModifyWorkItemFieldRequest{
		PropDiffs:                   propDiffs,
		WorkItemStatusServiceFacade: witem_facade.BuildWorkItemStatusServiceFacade(s.witemStatusService),
	}, oper)
	if err != nil { //成员不存在 不允许操作
		return err
	}

	if propDiffs.HasProp(witem.Diff_Describe) || propDiffs.HasProp(witem.Diff_Remark) {
		workItem.CalIconFlag()
	}

	err = s.tm.InTx(ctx, func(ctx context.Context) error {
		err = s.repo.SaveWorkItem(ctx, workItem)
		if err != nil {
			return errs.Internal(ctx, err)
		}
		return nil
	})

	s.domainMessageProducer.Send(ctx, workItem.GetMessages())

	if err != nil {
		return err
	}

	return err
}

func (s *SpaceWorkItemUsecase) DelWorkItem(ctx context.Context, oper *utils.LoginUserInfo, workItemId int64) error {

	//查看一下对应的工作用例是否存在
	workItem, err := s.repo.GetWorkItem(ctx, workItemId,
		&witem_repo.WithDocOption{Directors: true, Followers: true, Participators: true, PlanTime: true, ProcessRate: true},
		&witem_repo.WithOption{FlowRoles: true, FlowNodes: true})

	if err != nil {
		return errs.RepoErr(ctx, err)
	}

	// 判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, workItem.SpaceId, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		//成员不存在 不允许操作
		errInfo := errs.New(ctx, comm.ErrorCode_PERMISSION_INSUFFICIENT_DATA_PERMISSIONS)
		return errInfo
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_DELETE_SPACE_WORK_ITEM,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	space, err := s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil { //成员不存在 不允许操作
		return errs.New(ctx, comm.ErrorCode_PERMISSION_INSUFFICIENT_DATA_PERMISSIONS)
	}

	tasks, err := s.workItemService.DelSpaceTask(ctx, workItem, oper)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	allWorkItemIds := tasks.GetIds()

	_, err = s.spaceFileInfoService.DeleteSpaceFileInfoByWorkItemIds(ctx, space.Id, allWorkItemIds)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	err = s.tm.InTx(ctx, func(ctx context.Context) error {

		//清理附件
		err = s.spaceFileInfoRepo.SoftDelSpaceWorkItemsAllFile(ctx, allWorkItemIds)
		if err != nil {
			return err
		}

		//清理节点
		_, err = s.repo.DelWorkItemFlowNodeByWorkItemIds(ctx, allWorkItemIds...)
		if err != nil {
			return err
		}

		//清楚角色负责人
		_, err = s.repo.DelWorkItemFlowRoleByWorkItemIds(ctx, allWorkItemIds...)
		if err != nil {
			return err
		}

		//清理当前任务+子任务
		_, err = s.repo.DelSpaceWorkItem(ctx, allWorkItemIds...)
		if err != nil {
			return err
		}

		if workItem.IsSubTask() {
			err = s.repo.ResetChildTaskNum(ctx, workItem.Pid)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return errs.Internal(ctx, err)
	}

	//日志
	s.domainMessageProducer.Send(ctx, tasks.GetMessages())

	//特殊处理的通知
	var subWorkItems []*witem.WorkItem
	for _, v := range tasks {
		if v.IsSubTask() {
			subWorkItems = append(subWorkItems, v)
		}
	}

	bus.Emit(notify.Event_DeleteWorkItem, &event.DeleteWorkItem{
		Event:        notify.Event_DeleteWorkItem,
		Operator:     oper.UserId,
		Space:        space,
		WorkItem:     workItem,
		SubWorkItems: subWorkItems,
	})

	return nil
}

func (s *SpaceWorkItemUsecase) SetSubDirector(ctx context.Context, oper *utils.LoginUserInfo, workItemId int64, addUserIds, removeUserIds []int64) error {

	item, err := s.repo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{Directors: true}, &witem_repo.WithOption{FlowNodes: true})
	if err != nil {
		return errs.RepoErr(ctx, err)
	}

	if !item.IsSubTask() {
		return errs.Business(ctx, "仅子任务调用此接口")
	}

	// 判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, item.SpaceId, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(item),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE_WORK_ITEM,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	_, err = s.spaceRepo.GetSpace(ctx, item.SpaceId)
	if err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	ok, err := s.spaceMemberRepo.AllIsMember(ctx, item.SpaceId, append(addUserIds, oper.UserId)...)
	if err != nil { //成员不存在 不允许操作
		return errs.Internal(ctx, err)
	}
	if !ok {
		return errs.NoPerm(ctx)
	}

	err = s.workItemService.SetDirectorsForSubTask(ctx, item, &witem_service.SetDirectorsForSubTaskRequest{
		AddDirectors:                utils.ToStrArray(addUserIds),
		RemoveDirectors:             utils.ToStrArray(removeUserIds),
		WorkItemStatusServiceFacade: witem_facade.BuildWorkItemStatusServiceFacade(s.witemStatusService),
	}, oper)
	if err != nil {
		return err
	}

	err = s.tm.InTx(ctx, func(ctx context.Context) error {

		err = s.repo.SaveWorkItem(ctx, item)
		if err != nil {
			return err
		}

		for _, v := range item.WorkItemFlowNodes {
			err = s.repo.SaveWorkItemFlowNode(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	s.domainMessageProducer.Send(ctx, item.GetMessages())

	return nil
}

func (s *SpaceWorkItemUsecase) ConfirmSub(ctx context.Context, oper *utils.LoginUserInfo, workItemId int64, nextStatusKey string, reason string) error {

	workItem, err := s.repo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{Directors: true, ProcessRate: true}, nil)
	if err != nil {
		return errs.RepoErr(ctx, err)
	}

	if !workItem.IsSubTask() {
		return errs.Business(ctx, "非子任务不适用此接口")
	}

	_, err = s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	//判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, workItem.SpaceId, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CHANGE_STATE_SPACE_WORK_ITEM,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	statusInfo, err := s.workItemStatusRepo.GetWorkItemStatusInfo(ctx, workItem.SpaceId)
	if err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	flowTplt, err := s.workFlowRepo.GetFlowTemplate(ctx, workItem.WorkFlowTemplateId)
	if err != nil || flowTplt.StateFlowConf() == nil {
		return errs.Internal(ctx, err)
	}

	err = s.workItemService.ConfirmSpaceSubTaskState(ctx,
		workItem,
		&witem_service.ConfirmSpaceSubTaskState{
			NextStatusKey:          cast.ToString(nextStatusKey),
			WorkFlowTemplateFacade: witem_facade.BuildWorkFlowTemplateFacade(flowTplt),
			WorkItemStatusFacade:   witem_facade.BuildWorkItemStatusFacade(statusInfo),
			Reason:                 reason,
		},
		oper,
	)

	if err != nil {
		return err
	}

	txErr := s.tm.InTx(ctx, func(ctx context.Context) error {

		err = s.repo.SaveWorkItem(ctx, workItem)
		if err != nil {
			return err
		}

		for _, v := range workItem.WorkItemFlowNodes {
			err = s.repo.SaveWorkItemFlowNode(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	s.domainMessageProducer.Send(ctx, workItem.GetMessages())

	return nil
}

func (s *SpaceWorkItemUsecase) SetFlowMainDirectorByRoleKey(ctx context.Context, oper *utils.LoginUserInfo, workItemId int64, roleKey string, addUserIds, removeUserIds []int64) error {

	addUserIds = stream.Unique(addUserIds)
	removeUserIds = stream.Unique(removeUserIds)

	workItem, err := s.repo.GetWorkItem(
		ctx,
		workItemId,
		&witem_repo.WithDocOption{Directors: true},
		&witem_repo.WithOption{FlowRoles: true, FlowNodes: true},
	)
	if err != nil {
		return errs.RepoErr(ctx, err)
	}

	space, err := s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	//判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, space.Id, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	// 权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE_WORK_ITEM,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	ok, err := s.spaceMemberRepo.AllIsMember(ctx, workItem.SpaceId, append(addUserIds, oper.UserId)...)
	if err != nil {
		return errs.Internal(ctx, err)
	}
	if !ok {
		return errs.NoPerm(ctx)
	}

	// 获取模版
	template, err := s.workFlowRepo.GetFlowTemplate(ctx, workItem.WorkFlowTemplateId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	switch {
	case workItem.IsWorkFlowMainTask():
		err = s.workItemService.SetDirectorsForWorkFlowMainTaskByRoleKey(ctx, workItem, &witem_service.SetDirectorsForWorkFlowMainRequest{
			RoleKey:                     roleKey,
			AddDirectors:                utils.ToStrArray(addUserIds),
			RemoveDirectors:             utils.ToStrArray(removeUserIds),
			WorkItemStatusServiceFacade: witem_facade.BuildWorkItemStatusServiceFacade(s.witemStatusService),
		}, oper)
	case workItem.IsStateFlowMainTask():
		stateKeys := template.StateFlowConfig.GetNodeKeysByRoleKey(roleKey)
		if len(stateKeys) == 0 {
			return errs.Business(ctx, "找不到状态")
		}

		flowNodeDirectors := workItem.WorkItemFlowRoles.GetByRoleKey(roleKey).Directors

		directors := stream.Unique(stream.Diff(append(utils.ToInt64Array(flowNodeDirectors), addUserIds...), removeUserIds))

		err = s.workItemService.SetDirectorsForStateFlowMainTask(ctx, workItem, &witem_service.SetDirectorsForStateFlowMainTaskRequest{
			Directors:                   utils.ToStrArray(directors),
			WorkItemStatusServiceFacade: witem_facade.BuildWorkItemStatusServiceFacade(s.witemStatusService),
			StateKeys:                   stateKeys,
			RoleKeys:                    []string{roleKey},
		}, oper)
	}

	if err != nil {
		return err
	}

	err = s.tm.InTx(ctx, func(ctx context.Context) error {

		err = s.repo.SaveWorkItem(ctx, workItem)
		if err != nil {
			return err
		}

		for _, v := range workItem.WorkItemFlowNodes {
			err = s.repo.SaveWorkItemFlowNode(ctx, v)
			if err != nil {
				return err
			}
		}

		for _, v := range workItem.WorkItemFlowRoles {
			err = s.repo.SaveWorkItemFlowRole(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil

	})

	if err != nil {
		s.log.Error(err)
		return err
	}

	s.domainMessageProducer.Send(ctx, workItem.GetMessages())

	return nil
}

func (s *SpaceWorkItemUsecase) TerminateWorkItem(ctx context.Context, oper *utils.LoginUserInfo, workItemId int64, reason string) ([]int64, error) {

	workItem, err := s.repo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{Directors: true}, nil)
	if err != nil {
		return nil, errs.RepoErr(ctx, err)
	}

	space, err := s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil { //成员不存在 不允许操作
		return nil, errs.NoPerm(ctx)
	}

	//判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, space.Id, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return nil, errs.NoPerm(ctx)
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CHANGE_STATE_SPACE_WORK_ITEM,
	})
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	statusInfo, err := s.workItemStatusRepo.GetWorkItemStatusInfo(ctx, workItem.SpaceId)
	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	tasks, err := s.workItemService.TerminateWorkItem(ctx, workItem, &witem_service.TerminateWorkItemRequest{
		Reason:               reason,
		WorkItemStatusFacade: witem_facade.BuildWorkItemStatusFacade(statusInfo),
	}, oper)

	if err != nil {
		return nil, err
	}

	txErr := s.tm.InTx(ctx, func(ctx context.Context) error {
		for _, v := range tasks {
			err = s.repo.SaveWorkItem(ctx, v)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	s.domainMessageProducer.Send(ctx, tasks.GetMessages())

	return tasks.GetIds(), nil
}

func (s *SpaceWorkItemUsecase) RestartTask(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, flowNodeCode string, reason string) ([]int64, error) {

	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil { //成员不存在 不允许操作
		return nil, errs.RepoErr(ctx, err)
	}

	//判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return nil, errs.RepoErr(ctx, err)
	}

	workItem, err := s.repo.GetWorkItem(
		ctx,
		workItemId,
		&witem_repo.WithDocOption{
			Directors:     true,
			Participators: true,
			ProcessRate:   true,
		},
		&witem_repo.WithOption{
			FlowNodes: true,
		},
	)

	if err != nil {
		return nil, errs.RepoErr(ctx, err)
	}

	if !workItem.IsSameSpace(space.Id) {
		return nil, errs.NoPerm(ctx)
	}

	statusInfo, err := s.workItemStatusRepo.GetWorkItemStatusInfo(ctx, spaceId)
	if err != nil {
		return nil, errs.RepoErr(ctx, err)
	}

	curState := statusInfo.GetItemByKey(workItem.WorkItemStatus.Key)
	//权限 如果是关闭到恢复, 检查节点状态变更权限
	var hasPerm bool
	if curState.IsClose() {

		err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
			Oper:              oper,
			WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
			SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
			Perm:              consts.PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE,
		})

		hasPerm = err == nil
	} else {

		err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
			Oper:              oper,
			WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
			SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
			Perm:              consts.PERM_CHANGE_STATE_SPACE_WORK_ITEM,
		})

		hasPerm = err == nil
	}

	if !hasPerm {
		return nil, errs.NoPerm(ctx)
	}

	flowTplt, err := s.workFlowRepo.GetFlowTemplate(ctx, workItem.WorkFlowTemplateId)
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	tasks, err := s.workItemService.RestartTask(ctx, workItem, &witem_service.RestartTaskRequest{
		FlowNodeCode:           flowNodeCode,
		Reason:                 reason,
		WorkFlowTemplateFacade: witem_facade.BuildWorkFlowTemplateFacade(flowTplt),
		WorkItemStatusFacade:   witem_facade.BuildWorkItemStatusFacade(statusInfo),
	}, oper)

	if err != nil {
		return nil, err
	}

	txErr := s.tm.InTx(ctx, func(ctx context.Context) error {

		for _, v := range tasks {
			err = s.repo.SaveWorkItem(ctx, v)
			if err != nil {
				return err
			}

			for _, vv := range v.WorkItemFlowNodes {
				err = s.repo.SaveWorkItemFlowNode(ctx, vv)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	s.domainMessageProducer.Send(ctx, tasks.GetMessages())

	return tasks.GetIds(), nil
}

func (s *SpaceWorkItemUsecase) CloseTask(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, flowNodeCode string, reason string) ([]int64, error) {

	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil { //成员不存在 不允许操作
		return nil, errs.RepoErr(ctx, err)
	}

	statusInfo, err := s.workItemStatusRepo.GetWorkItemStatusInfo(ctx, spaceId)
	if err != nil { //成员不存在 不允许操作
		return nil, errs.RepoErr(ctx, err)
	}

	//判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return nil, errs.RepoErr(ctx, err)
	}

	workItem, err := s.repo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{Directors: true}, &witem_repo.WithOption{FlowNodes: true})
	if err != nil {
		return nil, errs.RepoErr(ctx, err)
	}

	if !workItem.IsSameSpace(space.Id) || workItem.IsSubTask() {
		return nil, errs.NoPerm(ctx)
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE,
	})
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	tasks, err := s.workItemService.CloseTask(ctx, workItem, &witem_service.CloseTaskRequest{
		Reason:               reason,
		FlowNodeCode:         flowNodeCode,
		WorkItemStatusFacade: witem_facade.BuildWorkItemStatusFacade(statusInfo),
	}, oper)

	if err != nil {
		return nil, err
	}

	txErr := s.tm.InTx(ctx, func(ctx context.Context) error {

		for _, v := range tasks {
			err = s.repo.SaveWorkItem(ctx, v)
			if err != nil {
				return err
			}

			for _, node := range v.WorkItemFlowNodes {
				err = s.repo.SaveWorkItemFlowNode(ctx, node)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	s.domainMessageProducer.Send(ctx, tasks.GetMessages())

	return tasks.GetIds(), nil
}

func (s *SpaceWorkItemUsecase) RollbackTask(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, flowNodeCode string, reason string) ([]int64, error) {

	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil { //成员不存在 不允许操作
		return nil, errs.RepoErr(ctx, err)
	}

	//判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return nil, errs.RepoErr(ctx, err)
	}

	workItem, err := s.repo.GetWorkItem(
		ctx,
		workItemId,
		&witem_repo.WithDocOption{
			Directors:     true,
			Participators: true,
		},
		&witem_repo.WithOption{
			FlowNodes: true,
			FlowRoles: true,
		},
	)

	if err != nil {
		return nil, errs.RepoErr(ctx, err)
	}

	if !workItem.IsSameSpace(space.Id) || workItem.IsSubTask() {
		return nil, errs.NoPerm(ctx)
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE,
	})
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	statusInfo, err := s.workItemStatusRepo.GetWorkItemStatusInfo(ctx, spaceId)
	if err != nil { //成员不存在 不允许操作
		return nil, errs.RepoErr(ctx, err)
	}

	//流程模版不存在
	flowTplt, err := s.workFlowRepo.GetFlowTemplate(ctx, workItem.WorkFlowTemplateId)
	if err != nil || flowTplt.WorkFlowConf() == nil {
		//不是节点模式，不走这里
		return nil, errs.NoPerm(ctx)
	}

	tasks, err := s.workItemService.RollbackTask(ctx, workItem, &witem_service.RollbackTaskRequest{
		FlowNodeCode:           flowNodeCode,
		Reason:                 reason,
		WorkFlowTemplateFacade: witem_facade.BuildWorkFlowTemplateFacade(flowTplt),
		WorkItemStatusFacade:   witem_facade.BuildWorkItemStatusFacade(statusInfo),
	}, oper)

	if err != nil {
		return nil, err
	}

	txErr := s.tm.InTx(ctx, func(ctx context.Context) error {

		for _, v := range tasks {
			err = s.repo.SaveWorkItem(ctx, v)
			if err != nil {
				return err
			}
		}

		for _, v := range workItem.WorkItemFlowNodes {
			err = s.repo.SaveWorkItemFlowNode(ctx, v)
			if err != nil {
				return err
			}
		}

		for _, v := range workItem.WorkItemFlowRoles {
			err = s.repo.SaveWorkItemFlowRole(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	s.domainMessageProducer.Send(ctx, tasks.GetMessages())

	return tasks.GetIds(), nil
}

func (s *SpaceWorkItemUsecase) ResumeTask(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, reason string) ([]int64, error) {

	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil { //成员不存在 不允许操作
		return nil, errs.RepoErr(ctx, err)
	}

	statusInfo, err := s.workItemStatusRepo.GetWorkItemStatusInfo(ctx, spaceId)
	if err != nil {
		return nil, errs.RepoErr(ctx, err)
	}

	//判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return nil, errs.RepoErr(ctx, err)
	}

	workItem, err := s.repo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{Directors: true}, nil)
	if err != nil {
		return nil, errs.RepoErr(ctx, err)
	}

	if !workItem.IsSameSpace(space.Id) {
		return nil, errs.NoPerm(ctx)
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE,
	})
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	tasks, err := s.workItemService.ResumeTask(ctx, workItem, &witem_service.ResumeTaskRequest{
		Reason:               reason,
		WorkItemStatusFacade: witem_facade.BuildWorkItemStatusFacade(statusInfo),
	}, oper)

	if err != nil {
		return nil, err
	}

	txErr := s.tm.InTx(ctx, func(ctx context.Context) error {

		for _, v := range tasks {
			err = s.repo.SaveWorkItem(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return nil, errs.Business(ctx, "恢复任务失败")
	}

	s.domainMessageProducer.Send(ctx, tasks.GetMessages())

	return tasks.GetIds(), nil
}

func (s *SpaceWorkItemUsecase) FollowWorkItem(ctx context.Context, oper *utils.LoginUserInfo, workItemId int64, unfollow bool) error {

	uid := oper.UserId

	workItem, err := s.repo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{Followers: true}, nil)
	if err != nil {
		return errs.RepoErr(ctx, err)
	}

	space, err := s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil { //成员不存在 不允许操作
		return errs.RepoErr(ctx, err)
	}

	//判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, space.Id, uid)
	if member == nil || err != nil { //成员不存在 不允许操作
		return errs.RepoErr(ctx, err)
	}

	if unfollow {
		workItem.UnFollow([]string{cast.ToString(uid)}, oper)
	} else {
		workItem.Follow([]string{cast.ToString(uid)}, oper)
	}

	err = s.repo.SaveWorkItem(ctx, workItem)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	s.domainMessageProducer.Send(ctx, workItem.GetMessages())

	return nil
}

func (s *SpaceWorkItemUsecase) ChangeTaskVersion(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, workVersionId int64) (effectIds []int64, err error) {

	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil { //成员不存在 不允许操作
		return effectIds, errs.RepoErr(ctx, err)
	}

	// 判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, space.Id, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return effectIds, errs.RepoErr(ctx, err)
	}

	workVersionInfo, err := s.workVersionRepo.GetSpaceWorkVersion(ctx, workVersionId)
	if err != nil {
		return effectIds, errs.RepoErr(ctx, err)
	}

	if workVersionInfo.SpaceId != space.Id {
		return effectIds, errs.NoPerm(ctx)
	}

	workItem, err := s.repo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{Directors: true}, nil)
	if err != nil {
		return effectIds, errs.RepoErr(ctx, err)
	}

	if !workItem.IsSameSpace(space.Id) || workItem.IsSubTask() {
		return effectIds, errs.NoPerm(ctx)
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE_WORK_ITEM,
	})
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	tasks, err := s.workItemService.ChangeWorkItemVersion(ctx, workItem, &witem_service.ChangeVersionRequest{
		WorkVersionId:               workVersionId,
		WorkItemStatusServiceFacade: witem_facade.BuildWorkItemStatusServiceFacade(s.witemStatusService),
	}, oper)

	if err != nil {
		return nil, err
	}

	err = s.tm.InTx(ctx, func(ctx context.Context) (txErr error) {
		for _, v := range tasks {
			err = s.repo.SaveWorkItem(ctx, v)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	s.domainMessageProducer.Send(ctx, tasks.GetMessages())

	return tasks.GetIds(), nil
}

func (s *SpaceWorkItemUsecase) RemindWork(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, nodeCode string) error {

	space, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		return errs.RepoErr(ctx, err)
	}

	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if member == nil || err != nil {
		return errs.RepoErr(ctx, err)
	}

	//查看一下对应的工作用例是否存在
	workItem, err := s.repo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{Directors: true}, &witem_repo.WithOption{FlowNodes: true})
	if err != nil {
		return errs.RepoErr(ctx, err)
	}

	if !workItem.IsSameSpace(space.Id) {
		return errs.Internal(ctx, err)
	}

	statusInfo, err := s.workItemStatusRepo.GetWorkItemStatusInfo(ctx, spaceId)
	if err != nil {
		return errs.RepoErr(ctx, err)
	}

	//归档类的状态，不允许操作了
	if statusInfo.HasArchivedItem(workItem.WorkItemStatus.Key) {
		return errs.Business(ctx, "任务已归档，不允许操作")
	}

	lockKey := fmt.Sprintf("remind:%v:%v", oper.UserId, workItem.Id)
	_, isOk := s.cacheLocker.Get(lockKey)
	if isOk {
		return errs.Business(ctx, "操作过于频繁，请稍后再试")
	}
	//5秒通知一次
	s.cacheLocker.Set(lockKey, 1, cache.WithExpiration(time.Second*5))

	flowNode := workItem.WorkItemFlowNodes.GetNodeByCode(nodeCode)
	if flowNode == nil {
		return errs.Business(ctx, "节点不存在")
	}

	//通知自己以外的负责人
	directorIds := stream.Map(flowNode.Directors, func(v string) int64 {
		return cast.ToInt64(v)
	})

	bus.Emit(notify.Event_RemindWork, &event.RemindWork{
		Event:     notify.Event_RemindWork,
		Space:     space,
		WorkItem:  workItem,
		Operator:  oper.UserId,
		TargetIds: stream.Diff(directorIds, []int64{oper.UserId}),
	})

	return err
}

func (s *SpaceWorkItemUsecase) OperationPermissions(ctx context.Context, userId int64, spaceId int64, workItemId int64, scene string) (map[string]interface{}, error) {

	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, userId)
	if member == nil || err != nil {
		return nil, errs.RepoErr(ctx, err)
	}

	_, err = s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		return nil, errs.RepoErr(ctx, err)
	}

	// 任务表单权限角色
	var permSubFuncRole map[string]interface{}
	switch scene {
	case "work_item_create":

		permSubFuncRole = s.permService.GetPermissionWithSceneForCreateScence(ctx, member.GetRole())

	case "work_item_edit":

		workItem, err := s.repo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{Directors: true}, nil)
		if err != nil {
			return nil, errs.RepoErr(ctx, err)
		}

		permSubFuncRole = s.permService.GetPermissionWithSceneForWorkItemScence(ctx, &perm_service.GetPermissionWithSceneForWorkItemScenceRequest{
			Uid:               userId,
			WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
			SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		})
	}

	if permSubFuncRole == nil {
		return nil, errs.NoPerm(ctx)
	}

	return permSubFuncRole, nil
}

func (s *SpaceWorkItemUsecase) ChangeWorkItemObject(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, newWorkObjectId int64) ([]int64, error) {

	_, err := s.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil {
		return nil, errs.RepoErr(ctx, err)
	}

	// 判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return nil, errs.RepoErr(ctx, err)
	}

	workItem, err := s.repo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{Directors: true}, nil)
	if err != nil || workItem.SpaceId != spaceId {
		return nil, errs.RepoErr(ctx, err)
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_MODIFY_SPACE_WORK_ITEM,
	})
	if err != nil {
		return nil, errs.NoPerm(ctx)
	}

	//检查模块是否存在
	_, err = s.spaceWorkObjectRepo.GetSpaceWorkObject(ctx, spaceId, newWorkObjectId)
	if err != nil {
		return nil, errs.RepoErr(ctx, err)
	}

	tasks, err := s.workItemService.ChangeWorkItemObject(ctx, workItem, &witem_service.ChangeWorkObjectRequest{
		WorkObjectId:                newWorkObjectId,
		WorkItemStatusServiceFacade: witem_facade.BuildWorkItemStatusServiceFacade(s.witemStatusService),
	}, oper)

	err = s.tm.InTx(ctx, func(ctx context.Context) (txErr error) {
		for _, v := range tasks {
			err = s.repo.SaveWorkItem(ctx, v)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	s.domainMessageProducer.Send(ctx, tasks.GetMessages())

	return tasks.GetIds(), nil
}

func (s *SpaceWorkItemUsecase) SetFollowers(ctx context.Context, oper *utils.LoginUserInfo, workItemId int64, userIds []int64) error {

	item, err := s.repo.GetWorkItem(ctx, workItemId,
		&witem_repo.WithDocOption{Directors: true, Followers: true},
		&witem_repo.WithOption{FlowNodes: true})
	if err != nil {
		return errs.RepoErr(ctx, err)
	}

	// 判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, item.SpaceId, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	_, err = s.spaceRepo.GetSpace(ctx, item.SpaceId)
	if err != nil {
		return errs.NoPerm(ctx)
	}

	ok, err := s.spaceMemberRepo.AllIsMember(ctx, item.SpaceId, userIds...)
	if err != nil {
		return errs.Internal(ctx, err)
	}
	if !ok {
		return errs.NoPerm(ctx)
	}

	item.UpdateFollower(utils.ToStrArray(userIds), oper)

	err = s.tm.InTx(ctx, func(ctx context.Context) error {

		err = s.repo.SaveWorkItem(ctx, item)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	s.domainMessageProducer.Send(ctx, item.GetMessages())
	return nil
}

func (s *SpaceWorkItemUsecase) ConfirmStateFlowMain(ctx context.Context, oper *utils.LoginUserInfo, workItemId int64, nextStatusKey string, reason, remark string) error {

	workItem, err := s.repo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{Directors: true, ProcessRate: true}, &witem_repo.WithOption{
		FlowRoles: true,
		FlowNodes: true,
	})
	if err != nil {
		return errs.RepoErr(ctx, err)
	}

	if !workItem.IsStateFlowMainTask() {
		return errs.Business(ctx, "仅状态流主任务使用1")
	}

	_, err = s.spaceRepo.GetSpace(ctx, workItem.SpaceId)
	if err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	//判断是不是这个空间的成员，并且是否有相关的基本操作权限
	member, err := s.spaceMemberRepo.GetSpaceMember(ctx, workItem.SpaceId, oper.UserId)
	if member == nil || err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	//权限
	err = s.permService.CheckWorkItemEditPerm(ctx, &perm_service.CheckWorkItemEditPermRequest{
		Oper:              oper,
		WorkItemFacade:    perm_facade.BuildWorkItemFacade(workItem),
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CONFIRM_SPACE_WORK_ITEM_FLOW_NODE_STATE,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	statusInfo, err := s.workItemStatusRepo.GetWorkItemStatusInfo(ctx, workItem.SpaceId)
	if err != nil { //成员不存在 不允许操作
		return errs.NoPerm(ctx)
	}

	flowTplt, err := s.workFlowRepo.GetFlowTemplate(ctx, workItem.WorkFlowTemplateId)
	if err != nil || flowTplt.StateFlowConf() == nil {
		return errs.Internal(ctx, err)
	}

	err = s.workItemService.ConfirmStateFlowMainTaskState(ctx,
		workItem,
		&witem_service.ConfirmStateFlowMainTaskStateByStateKey{
			NextStatusKey:          cast.ToString(nextStatusKey),
			WorkFlowTemplateFacade: witem_facade.BuildWorkFlowTemplateFacade(flowTplt),
			WorkItemStatusFacade:   witem_facade.BuildWorkItemStatusFacade(statusInfo),
			Reason:                 reason,
			Remark:                 remark,
		},
		oper,
	)

	if err != nil {
		return err
	}

	txErr := s.tm.InTx(ctx, func(ctx context.Context) error {

		err = s.repo.SaveWorkItem(ctx, workItem)
		if err != nil {
			return err
		}

		for _, v := range workItem.WorkItemFlowNodes {
			err = s.repo.SaveWorkItemFlowNode(ctx, v)
			if err != nil {
				return err
			}
		}

		for _, v := range workItem.WorkItemFlowRoles {
			err = s.repo.SaveWorkItemFlowRole(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	s.domainMessageProducer.Send(ctx, workItem.GetMessages())

	return nil
}
