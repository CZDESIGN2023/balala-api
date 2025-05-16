package command

import (
	"context"
	vo "go-cs/internal/bean/vo"
	"go-cs/internal/bean/vo/query"
	file_repo "go-cs/internal/domain/file_info/repo"
	space_repo "go-cs/internal/domain/space/repo"
	space_file_service "go-cs/internal/domain/space_file_info/service"
	member_repo "go-cs/internal/domain/space_member/repo"
	member_service "go-cs/internal/domain/space_member/service"
	workObj_repo "go-cs/internal/domain/space_work_object/repo"
	wf_repo "go-cs/internal/domain/work_flow/repo"
	"go-cs/pkg/stream"

	witem_domain "go-cs/internal/domain/work_item"
	witem_facade "go-cs/internal/domain/work_item/facade"
	witem_repo "go-cs/internal/domain/work_item/repo"
	witem_service "go-cs/internal/domain/work_item/service"

	witem_role_repo "go-cs/internal/domain/work_item_role/repo"

	witem_status_repo "go-cs/internal/domain/work_item_status/repo"
	witem_type_repo "go-cs/internal/domain/work_item_type/repo"

	tag_service "go-cs/internal/domain/space_tag/service"

	perm_service "go-cs/internal/domain/perm/service"
	domain_message "go-cs/internal/domain/pkg/message"
	workVersion_repo "go-cs/internal/domain/space_work_version/repo"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"

	"github.com/go-kratos/kratos/v2/log"
)

type CreateSpaceSubTaskCmd struct {
	repo                witem_repo.WorkItemRepo
	witemEsRepo         witem_repo.WorkItemEsRepo
	spaceRepo           space_repo.SpaceRepo
	spaceMemberRepo     member_repo.SpaceMemberRepo
	spaceWorkObjectRepo workObj_repo.SpaceWorkObjectRepo
	workItemRoleRepo    witem_role_repo.WorkItemRoleRepo
	workVersionRepo     workVersion_repo.SpaceWorkVersionRepo
	workItemStatusRepo  witem_status_repo.WorkItemStatusRepo
	workItemTypeRepo    witem_type_repo.WorkItemTypeRepo
	workFlowRepo        wf_repo.WorkFlowRepo
	fileInfoRepo        file_repo.FileInfoRepo

	permService        *perm_service.PermService
	spaceMemberService *member_service.SpaceMemberService
	spaceTagService    *tag_service.SpaceTagService
	workItemService    *witem_service.WorkItemService

	domainMessageProducer *domain_message.DomainMessageProducer

	log *log.Helper
	tm  trans.Transaction
}

func NewCreateSpaceSubTaskCmd(
	repo witem_repo.WorkItemRepo,
	witemEsRepo witem_repo.WorkItemEsRepo,
	spaceRepo space_repo.SpaceRepo,
	spaceMemberRepo member_repo.SpaceMemberRepo,
	spaceWorkObjectRepo workObj_repo.SpaceWorkObjectRepo,
	workItemRoleRepo witem_role_repo.WorkItemRoleRepo,
	workVersionRepo workVersion_repo.SpaceWorkVersionRepo,
	workItemStatusRepo witem_status_repo.WorkItemStatusRepo,
	workItemTypeRepo witem_type_repo.WorkItemTypeRepo,
	workFlowRepo wf_repo.WorkFlowRepo,
	fileInfoRepo file_repo.FileInfoRepo,

	permService *perm_service.PermService,
	spaceMemberService *member_service.SpaceMemberService,
	spaceTagService *tag_service.SpaceTagService,
	workItemService *witem_service.WorkItemService,
	spaceFileInfoService *space_file_service.SpaceFileInfoService,

	domainMessageProducer *domain_message.DomainMessageProducer,

	logger log.Logger,
	tm trans.Transaction,
) *CreateSpaceSubTaskCmd {

	moduleName := "biz.CreateSpaceSubTaskCmd"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &CreateSpaceSubTaskCmd{
		repo:                repo,
		witemEsRepo:         witemEsRepo,
		spaceRepo:           spaceRepo,
		spaceMemberRepo:     spaceMemberRepo,
		spaceWorkObjectRepo: spaceWorkObjectRepo,
		workItemRoleRepo:    workItemRoleRepo,
		workVersionRepo:     workVersionRepo,

		fileInfoRepo: fileInfoRepo,

		workItemStatusRepo: workItemStatusRepo,
		workItemTypeRepo:   workItemTypeRepo,
		workFlowRepo:       workFlowRepo,
		permService:        permService,
		workItemService:    workItemService,
		spaceMemberService: spaceMemberService,
		spaceTagService:    spaceTagService,

		log: hlog,
		tm:  tm,
	}
}

func (cmd *CreateSpaceSubTaskCmd) Excute(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, workItemId int64, in vo.CreateSpaceWorkItemTaskVoV2) (int64, error) {

	space, err := cmd.spaceRepo.GetSpace(ctx, spaceId)
	if err != nil { //成员不存在 不允许操作
		return 0, errs.NoPerm(ctx)
	}

	//判断是否存在对应的工作项信息
	workItem, err := cmd.repo.GetWorkItem(ctx, workItemId, &witem_repo.WithDocOption{
		Priority: true,
	}, nil)
	if err != nil {
		return 0, errs.NoPerm(ctx)
	}

	if workItem.IsSubTask() {
		return 0, errs.Business(ctx, "子任务不允许创建子任务")
	}

	if workItem.IsStateFlowMainTask() {
		return 0, errs.Business(ctx, "状态流转主任务不允许创建子任务")
	}

	//检查所有负责人是否合法
	isAllMember, err := cmd.spaceMemberService.CheckAllIsMember(ctx, spaceId, in.DirectorAdd)
	if err != nil {
		return 0, errs.Internal(ctx, err)
	}

	if !isAllMember {
		return 0, errs.Business(ctx, "成员不存在")
	}

	//任务状态
	statusInfo, err := cmd.workItemStatusRepo.GetWorkItemStatusInfo(ctx, space.Id)
	if err != nil {
		return 0, errs.NoPerm(ctx)
	}

	//判断工作项类型, 不是主任务
	workItemTypeInfo, err := cmd.workItemTypeRepo.QWorkItemTypeInfo(ctx, query.WorkItemTypeInfoQuery{
		SpaceId: space.Id,
	})
	if err != nil || workItemTypeInfo.GetWorkFlowTaskType() == nil {
		return 0, errs.NoPerm(ctx)
	}

	subTaskFlow, err := cmd.workFlowRepo.GetWorkFlowBySpaceWorkItemTypeId(ctx, space.Id, workItemTypeInfo.GetSubTaskType().Id)
	if err != nil {
		return 0, errs.NoPerm(ctx)
	}

	subTaskTpltId := subTaskFlow[0].LastTemplateId
	subTaskTplt, err := cmd.workFlowRepo.GetFlowTemplate(ctx, subTaskTpltId)
	if err != nil || subTaskTplt.StateFlowConf() == nil {
		return 0, errs.NoPerm(ctx)
	}

	newWorkItem, err := cmd.workItemService.CreateSpaceSubTask(
		ctx,
		workItem,
		&witem_service.CreateSubTaskRequest{
			Name: in.WorkItemName,
			PlanTime: witem_domain.PlanTime{
				StartAt:    in.PlanStartAt,
				CompleteAt: in.PlanCompleteAt,
			},
			ProcessRate:            in.ProcessRate,
			Directors:              utils.ToStrArray(stream.Unique(in.DirectorAdd)),
			WorkFlowTemplateFacade: witem_facade.BuildWorkFlowTemplateFacade(subTaskTplt),
			WorkItemStatusFacade:   witem_facade.BuildWorkItemStatusFacade(statusInfo),
			WorkFlowFacade:         witem_facade.BuildWorkFlowFacade(subTaskFlow[0]),
		},
		oper,
	)

	if err != nil {
		return 0, err
	}

	//保存
	err = cmd.tm.InTx(ctx, func(ctx context.Context) error {
		err = cmd.repo.CreateWorkItem(ctx, newWorkItem)
		if err != nil {
			return err
		}

		err = cmd.repo.CreateWorkItemFlowNodes(ctx, newWorkItem.WorkItemFlowNodes...)
		if err != nil {
			return err
		}

		err = cmd.repo.ResetChildTaskNum(ctx, workItemId)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, errs.Internal(ctx, err)
	}

	//同步一份到es，避免延迟
	_ = cmd.witemEsRepo.CreateWorkItemEs(ctx, newWorkItem)

	cmd.domainMessageProducer.Send(ctx, newWorkItem.GetMessages())

	return newWorkItem.Id, nil
}
