package command

import (
	"context"
	"go-cs/api/comm"
	vo "go-cs/internal/bean/vo"
	"go-cs/internal/consts"
	file_repo "go-cs/internal/domain/file_info/repo"
	file_service "go-cs/internal/domain/file_info/service"
	perm_service "go-cs/internal/domain/perm/service"
	domain_message "go-cs/internal/domain/pkg/message"
	space_repo "go-cs/internal/domain/space/repo"
	space_file_service "go-cs/internal/domain/space_file_info/service"
	member_repo "go-cs/internal/domain/space_member/repo"
	member_service "go-cs/internal/domain/space_member/service"
	workObj_repo "go-cs/internal/domain/space_work_object/repo"
	wf_domain "go-cs/internal/domain/work_flow"
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

	workVersion_repo "go-cs/internal/domain/space_work_version/repo"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type CreateSpaceTaskCmd struct {
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
	fileInfoService    *file_service.FileInfoService

	domainMessageProducer *domain_message.DomainMessageProducer

	log *log.Helper
	tm  trans.Transaction
}

func NewCreateSpaceTaskCmd(
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
	fileInfoService *file_service.FileInfoService,

	domainMessageProducer *domain_message.DomainMessageProducer,

	logger log.Logger,
	tm trans.Transaction,
) *CreateSpaceTaskCmd {

	moduleName := "biz.CreateSpaceTaskCmd"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &CreateSpaceTaskCmd{
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
		fileInfoService:    fileInfoService,

		domainMessageProducer: domainMessageProducer,

		log: hlog,
		tm:  tm,
	}
}

func (cmd *CreateSpaceTaskCmd) Execute(ctx context.Context, oper *utils.LoginUserInfo, in vo.CreateSpaceWorkItemVoV2) (int64, error) {

	space, err := cmd.spaceRepo.GetSpace(ctx, in.SpaceId)
	if err != nil { //成员不存在 不允许操作
		return 0, errs.NoPerm(ctx)
	}

	// 检查所有负责人是否合法
	allUserIds := in.Owner.AllDirectors()
	isAllMember, err := cmd.spaceMemberService.CheckAllIsMember(ctx, in.SpaceId, allUserIds)
	if err != nil {
		return 0, errs.Internal(ctx, err)
	}

	if !isAllMember {
		return 0, errs.Business(ctx, "成员不存在")
	}

	// 判断是否存在对应的工作项信息
	workObject, err := cmd.spaceWorkObjectRepo.GetSpaceWorkObject(ctx, in.SpaceId, in.WorkObjectId)
	if err != nil {
		return 0, errs.NoPerm(ctx)
	}

	// 角色信息
	workItemRoles, err := cmd.workItemRoleRepo.GetWorkItemRoles(ctx, space.Id)
	if err != nil {
		return 0, errs.NoPerm(ctx)
	}

	// 检查Version
	workVersionInfo, err := cmd.workVersionRepo.GetSpaceWorkVersion(ctx, in.WorkVersionId)
	if err != nil || workVersionInfo.SpaceId != space.Id {
		return 0, err
	}

	// 任务状态
	statusInfo, err := cmd.workItemStatusRepo.GetWorkItemStatusInfo(ctx, space.Id)
	if err != nil {
		return 0, errs.NoPerm(ctx)
	}

	// 检查流程
	flowInfo, err := cmd.workFlowRepo.GetWorkFlow(ctx, in.WorkFlowId)
	if err != nil {
		return 0, errs.NoPerm(ctx)
	}
	if flowInfo.Status != wf_domain.WorkFlowStatus_Enable {
		return 0, errs.Business(ctx, comm.ErrorCode_WORK_FLOW_STATUS_NOT_ENABLE)
	}

	workItemTaskType, err := cmd.workItemTypeRepo.GetWorkItemType(ctx, flowInfo.WorkItemTypeId)
	if err != nil {
		return 0, errs.NoPerm(ctx)
	}
	// 检查任务类型
	if workItemTaskType.IsSubTaskType() {
		return 0, errs.Business(ctx, "创建子任务不能调用此接口")
	}

	// 检查流程模板
	flowTplt, err := cmd.workFlowRepo.GetFlowTemplate(ctx, flowInfo.LastTemplateId)
	if err != nil {
		return 0, errs.NoPerm(ctx)
	}
	switch flowInfo.FlowMode {
	case consts.FlowMode_WorkFlow:
		if flowTplt.WorkFlowConf() == nil {
			return 0, errs.NoPerm(ctx)
		}
	case consts.FlowMode_StateFlow:
		if flowTplt.StateFlowConf() == nil {
			return 0, errs.NoPerm(ctx)
		}
	}

	if flowTplt.WorkItemTypeId != workItemTaskType.Id {
		return 0, errs.NoPerm(ctx)
	}

	// 检查文件
	fileIds := make([]int64, 0)
	for _, v := range in.FileAdd {
		fileIds = append(fileIds, v.Id)
	}

	// 检查TAG
	tagIds, err := cmd.spaceTagService.FilterExistSpaceTagIds(ctx, space.Id, in.TagAdd)
	if err != nil {
		return 0, err
	}

	witemIconFlag := witem_domain.IconFlag(0)
	witemIconFlag.AddFlag(in.IconFlags...)

	var directors []witem_service.CreateWorkItemRequest_Directors
	for _, v := range in.Owner {
		workItemRole := workItemRoles.GetRoleById(cast.ToInt64(v.OwnerRole))
		if workItemRole == nil {
			return 0, errs.Business(ctx, "无效的角色")
		}

		directors = append(directors, witem_service.CreateWorkItemRequest_Directors{
			RoleId:    cast.ToString(workItemRole.Id),
			RoleKey:   workItemRole.Key,
			Directors: utils.ToStrArray(stream.Unique(v.Directors)),
		})
	}

	newWorkItem, err := cmd.workItemService.CreateSpaceTask(
		ctx,
		&witem_service.CreateWorkItemRequest{
			SpaceId:         space.Id,
			UserId:          oper.UserId,
			WorkObjectId:    workObject.Id,
			VersionId:       workVersionInfo.Id,
			WorkItemTypeId:  workItemTaskType.Id,
			WorkItemTypeKey: workItemTaskType.Key,
			Name:            in.WorkItemName,
			PlanTime: witem_domain.PlanTime{
				StartAt:    in.PlanStartAt,
				CompleteAt: in.PlanCompleteAt,
			},
			ProcessRate:            in.ProcessRate,
			Remark:                 in.Remark,
			Describe:               in.Describe,
			Priority:               in.Priority,
			IconFlag:               witemIconFlag,
			Tags:                   utils.ToStrArray(tagIds),
			Directors:              directors,
			Files:                  fileIds,
			Followers:              in.Followers,
			WorkFlowFacade:         witem_facade.BuildWorkFlowFacade(flowInfo),
			WorkFlowTemplateFacade: witem_facade.BuildWorkFlowTemplateFacade(flowTplt),
			WorkItemStatusFacade:   witem_facade.BuildWorkItemStatusFacade(statusInfo),
			FileInfoFacade:         witem_facade.BuildFileInfoFacade(cmd.fileInfoService),
		},
		oper,
	)
	if err != nil {
		return 0, err
	}

	//保存
	txErr := cmd.tm.InTx(ctx, func(ctx context.Context) error {

		err := cmd.repo.CreateWorkItem(ctx, newWorkItem)
		if err != nil {
			return err
		}

		err = cmd.repo.CreateWorkItemFlowNodes(ctx, newWorkItem.WorkItemFlowNodes...)
		if err != nil {
			return err
		}

		err = cmd.repo.CreateWorkItemFlowRoles(ctx, newWorkItem.WorkItemFlowRoles...)
		if err != nil {
			return err
		}

		err = cmd.repo.CreateWorkItemFiles(ctx, newWorkItem.WorkItemFiles)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return 0, errs.Internal(ctx, txErr)
	}

	// 同步一份到es，避免延迟
	_ = cmd.witemEsRepo.CreateWorkItemEs(ctx, newWorkItem)

	cmd.domainMessageProducer.Send(ctx, newWorkItem.GetMessages())

	return newWorkItem.Id, nil
}
