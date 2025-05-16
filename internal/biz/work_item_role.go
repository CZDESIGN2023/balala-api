package biz

import (
	"context"
	pb "go-cs/api/work_item_role/v1"
	"go-cs/internal/bean/vo/query"
	"go-cs/internal/consts"
	perm_facade "go-cs/internal/domain/perm/facade"
	perm_service "go-cs/internal/domain/perm/service"
	domain_message "go-cs/internal/domain/pkg/message"
	member_repo "go-cs/internal/domain/space_member/repo"
	flow_repo "go-cs/internal/domain/work_flow/repo"
	witem_repo "go-cs/internal/domain/work_item/repo"
	"go-cs/internal/domain/work_item_role"
	witem_role_repo "go-cs/internal/domain/work_item_role/repo"
	witem_role_service "go-cs/internal/domain/work_item_role/service"
	witem_type_repo "go-cs/internal/domain/work_item_type/repo"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/pkg/trans"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/cast"
)

type WorkItemRoleUsecase struct {
	repo          witem_role_repo.WorkItemRoleRepo
	memberRepo    member_repo.SpaceMemberRepo
	witemTypeRepo witem_type_repo.WorkItemTypeRepo
	witemRepo     witem_repo.WorkItemRepo
	flowRepo      flow_repo.WorkFlowRepo

	roleService *witem_role_service.WorkItemRoleService
	permService *perm_service.PermService

	domainMessageProducer *domain_message.DomainMessageProducer

	log *log.Helper
	tm  trans.Transaction
}

func NewWorkItemRoleUsecase(
	repo witem_role_repo.WorkItemRoleRepo,
	memberRepo member_repo.SpaceMemberRepo,
	witemTypeRepo witem_type_repo.WorkItemTypeRepo,
	witemRepo witem_repo.WorkItemRepo,
	flowRepo flow_repo.WorkFlowRepo,

	roleService *witem_role_service.WorkItemRoleService,
	permService *perm_service.PermService,

	domainMessageProducer *domain_message.DomainMessageProducer,

	tm trans.Transaction,
	logger log.Logger,

) *WorkItemRoleUsecase {
	moduleName := "WorkItemRoleUsecase"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	return &WorkItemRoleUsecase{
		roleService: roleService,
		permService: permService,

		repo:          repo,
		memberRepo:    memberRepo,
		witemTypeRepo: witemTypeRepo,
		witemRepo:     witemRepo,
		flowRepo:      flowRepo,

		domainMessageProducer: domainMessageProducer,

		log: hlog,
		tm:  tm,
	}
}

func (uc *WorkItemRoleUsecase) SetSpaceWorkItemRoleRanking(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, rankingList []map[string]int64) error {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		err := errs.NoPerm(ctx)
		return err
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_ModifyWorkFlowRole,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	var wItemRoles work_item_role.WorkItemRoles
	for _, v := range rankingList {
		roleId := cast.ToInt64(v["id"])
		newRanking := cast.ToInt64(v["ranking"])

		//调整排序
		witemRole, err := uc.repo.GetWorkItemRole(ctx, roleId)
		if err != nil {
			return errs.Internal(ctx, err)
		}

		witemRole.ChangeRanking(newRanking, oper)

		wItemRoles = append(wItemRoles, witemRole)
	}

	txErr := uc.tm.InTx(ctx, func(ctx context.Context) error {

		for _, v := range wItemRoles {
			err = uc.repo.SaveWorkItemRole(ctx, v)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return errs.Internal(ctx, txErr)
	}

	msg := &domain_message.ChangeRoleOrder{
		DomainMessageBase: shared.DomainMessageBase{
			Oper:     oper,
			OperTime: time.Now(),
		},
		SpaceId:   spaceId,
		FlowScope: wItemRoles[0].FlowScope,
	}

	uc.domainMessageProducer.Send(ctx, shared.DomainMessages{msg})

	return nil
}

func (uc *WorkItemRoleUsecase) QSpaceWorkItemRoleList(ctx context.Context, uid int64, req *pb.SpaceRoleListQueryRequest) (*pb.SpaceRoleListQueryResult, error) {

	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, req.SpaceId, uid)
	if member == nil || err != nil {
		//不是改空间成员，不允许查看该空间的其它成员列表
		err := errs.NoPerm(ctx)
		return nil, err
	}

	qResult, err := uc.repo.QSpaceWorkItemRoleList(ctx, &query.SpaceWorkItemRoleQuery{
		SpaceId:   req.SpaceId,
		FlowScope: consts.FlowScope(req.FlowScope),
	})

	if err != nil {
		return nil, errs.Internal(ctx, err)
	}

	result := &pb.SpaceRoleListQueryResult{}
	result.Total = int32(qResult.Total)
	result.List = make([]*pb.SpaceRoleListQueryResult_ListItem, 0)

	for _, v := range qResult.List {
		result.List = append(result.List, &pb.SpaceRoleListQueryResult_ListItem{
			Id:             uint64(v.Id),
			Name:           v.Name,
			Ranking:        uint64(v.Ranking),
			SpaceId:        uint64(v.SpaceId),
			WorkItemTypeId: uint64(v.WorkItemTypeId),
			Key:            v.Key,
			IsSys:          uint32(v.IsSys),
			Status:         uint32(v.Status),
			CreatedAt:      uint64(v.CreatedAt),
			UpdatedAt:      uint64(v.UpdatedAt),
			FlowScope:      v.FlowScope,
		})
	}

	return result, nil
}

func (uc *WorkItemRoleUsecase) SeSpaceWorkItemRoleName(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, roleId int64, newName string) error {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		err := errs.NoPerm(ctx)
		return err
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_ModifyWorkFlowRole,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	workItemRole, err := uc.repo.GetWorkItemRole(ctx, roleId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	err = uc.roleService.CheckAndUpdateSpaceWorkItemRoleName(ctx, workItemRole, newName, oper)
	if err != nil {
		return err
	}

	err = uc.repo.SaveWorkItemRole(ctx, workItemRole)
	if err != nil {
		return err
	}

	uc.domainMessageProducer.Send(ctx, workItemRole.GetMessages())

	return nil
}

func (uc *WorkItemRoleUsecase) DelSpaceWorkItemRole(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, roleId int64) error {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		err := errs.NoPerm(ctx)
		return err
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_DeleteWorkFlowRole,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	workItemRole, err := uc.repo.GetWorkItemRole(ctx, roleId)
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if !workItemRole.IsSameSpace(spaceId) {
		return errs.NoPerm(ctx)
	}

	//检查是不是关联了流程配置
	tpltIds, err := uc.flowRepo.SearchTaskWorkFlowTemplateByOwnerRoleRule(ctx, spaceId, cast.ToString(workItemRole.Id))
	if err != nil {
		return errs.Internal(ctx, err)
	}

	if len(tpltIds) > 0 {
		return errs.Business(ctx, "该角色已被应用至任务流程设置中，无法删除")
	}

	err = workItemRole.OnDelete(oper)
	if err != nil {
		return errs.Business(ctx, err.Error())
	}

	err = uc.repo.DelWorkItemRole(ctx, roleId)
	if err != nil {
		return err
	}

	uc.domainMessageProducer.Send(ctx, workItemRole.GetMessages())

	return nil
}

func (uc *WorkItemRoleUsecase) CreateSpaceWorkItemRole(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, name string, scope consts.FlowScope) error {

	uid := oper.UserId
	// 判断当前用户是否在要查询的项目空间内
	member, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, uid)
	if member == nil || err != nil {
		err := errs.NoPerm(ctx)
		return err
	}

	// 验证权限
	err = uc.permService.CheckSpaceOperatePerm(ctx, &perm_service.CheckSpaceOperatePermRequest{
		SpaceMemberFacade: perm_facade.BuildSpaceMemberFacade(member),
		Perm:              consts.PERM_CreateWorkFlowRole,
	})
	if err != nil {
		return errs.NoPerm(ctx)
	}

	wItemTypeInfo, err := uc.witemTypeRepo.QWorkItemTypeInfo(ctx, query.WorkItemTypeInfoQuery{
		SpaceId: spaceId,
	})

	if err != nil {
		return errs.Internal(ctx, err)
	}

	newRole, err := uc.roleService.CreateSpaceRole(ctx, witem_role_service.CreateSpaceRoleRequest{
		SpaceId:        spaceId,
		WorkItemTypeId: wItemTypeInfo.GetWorkFlowTaskType().Id,
		Name:           name,
		Uid:            uid,
		FlowScope:      scope,
	}, oper)
	if err != nil {
		return err
	}

	err = uc.repo.CreateWorkItemRole(ctx, newRole)
	if err != nil {
		return err
	}

	uc.domainMessageProducer.Send(ctx, newRole.GetMessages())

	return nil
}

func (uc *WorkItemRoleUsecase) QSpaceWorkItemRelationCount(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, flowRoleId int64) (int64, error) {

	// 判断当前用户是否在要查询的项目空间内
	_, err := uc.memberRepo.GetSpaceMember(ctx, spaceId, oper.UserId)
	if err != nil {
		return 0, errs.NoPerm(ctx)
	}

	totalNum, err := uc.witemRepo.CountWorkFlowRoleRelatedSpaceWorkItem(ctx, spaceId, flowRoleId)
	return totalNum, err
}

func (uc *WorkItemRoleUsecase) QSpaceTemplateRelationCount(ctx context.Context, oper *utils.LoginUserInfo, spaceId int64, flowRoleId int64) (int64, error) {

	//检查是不是关联了流程配置
	tpltIds, err := uc.flowRepo.SearchTaskWorkFlowTemplateByOwnerRoleRule(ctx, spaceId, cast.ToString(flowRoleId))
	if err != nil {
		return 0, err
	}

	return int64(len(tpltIds)), nil
}
