package service

import (
	"context"
	"github.com/spf13/cast"
	"go-cs/internal/consts"
	perm_domain "go-cs/internal/domain/perm"
	"go-cs/internal/domain/perm/facade"
	"go-cs/internal/domain/user/repo"
	shared "go-cs/internal/pkg/domain"
	"go-cs/internal/utils"
	"go-cs/internal/utils/errs"
	"slices"
)

type PermService struct {
	userRepo repo.UserRepo
}

func NewPermService(userRepo repo.UserRepo) *PermService {
	return &PermService{
		userRepo: userRepo,
	}
}

type CheckSpaceOperatePermRequest struct {
	SpaceMemberFacade *facade.SpaceMemberFacade
	Perm              string
}

func (s *PermService) CheckSpaceOperatePerm(ctx context.Context, req *CheckSpaceOperatePermRequest) error {

	member := req.SpaceMemberFacade.GetSpaceMember()

	var memberRoleId int64

	if member != nil {
		memberRoleId = member.RoleId
		if memberRoleId == 0 {
			memberRoleId = consts.MEMBER_ROLE_WATCHER
		}
	} else {
		user, err := s.userRepo.GetUserByUserId(ctx, utils.GetLoginUserId(ctx))
		if err != nil {
			return errs.NoPerm(ctx)
		}
		if user.Role == consts.SystemRole_Enterprise {
			memberRoleId = consts.MEMBER_ROLE_WATCHER
		} else {
			return errs.NoPerm(ctx)
		}
	}

	hasPerm := perm_domain.Instance().Check(memberRoleId, req.Perm)
	if !hasPerm {
		return errs.NoPerm(ctx)
	}
	return nil
}

func (s *PermService) CheckRoleOperatePerm(ctx context.Context, memberRoleId int64, targetMemberRoleId int64) error {

	if memberRoleId == 0 {
		memberRoleId = consts.MEMBER_ROLE_WATCHER
	}

	if targetMemberRoleId == 0 {
		targetMemberRoleId = consts.MEMBER_ROLE_WATCHER
	}

	hasPerm := perm_domain.Instance().CheckLevel(memberRoleId, targetMemberRoleId)
	if !hasPerm {
		return errs.NoPerm(ctx)
	}
	return nil
}

type CheckWorkItemEditPermRequest struct {
	Oper              shared.Oper
	WorkItemFacade    *facade.WorkItemFacade
	SpaceMemberFacade *facade.SpaceMemberFacade
	Perm              string
	FlowNodeCode      string
}

func (s *PermService) CheckWorkItemEditPerm(ctx context.Context, req *CheckWorkItemEditPermRequest) error {

	workItem := req.WorkItemFacade.GetWorkItem()
	member := req.SpaceMemberFacade.GetSpaceMember()
	uid := req.Oper.GetId()

	//是否当前负责人
	var curWorkItemRole string
	if workItem.UserId == uid {
		curWorkItemRole = consts.WORK_ITEM_ROLE_CREATOR
	} else {
		if req.FlowNodeCode != "" {
			if workItem.WorkItemFlowNodes != nil && slices.Contains(workItem.WorkItemFlowNodes.GetNodeByCode(req.FlowNodeCode).Directors, cast.ToString(uid)) {
				curWorkItemRole = consts.WORK_ITEM_ROLE_NODE_OWNER
			}
		} else {
			if slices.Contains(workItem.Doc.Directors, cast.ToString(uid)) {
				curWorkItemRole = consts.WORK_ITEM_ROLE_NODE_OWNER
			}
		}
	}

	memberRoleId := member.RoleId
	if memberRoleId == 0 {
		memberRoleId = consts.MEMBER_ROLE_WATCHER
	}

	var hasPerm bool
	if workItem.IsSubTask() {
		hasPerm = perm_domain.Instance().WorkItemEditPerm.CheckTask(memberRoleId, curWorkItemRole, req.Perm)
	} else {
		hasPerm = perm_domain.Instance().WorkItemEditPerm.Check(memberRoleId, curWorkItemRole, req.Perm)
	}

	if !hasPerm {
		return errs.NoPerm(ctx)
	}

	return nil
}

// func (s *PermService) CheckWorkItemEditPermFor(ctx context.Context, workItemRole string, memberRoleId int64, perm string) error {
// 	hasPerm := perm_domain.Instance().WorkItemEditPerm.Check(memberRoleId, workItemRole, perm)
// 	if !hasPerm {
// 		return errs.NoPerm(ctx)
// 	}
// 	return nil
// }

func (s *PermService) GetPermissionWithSceneForCreateScence(ctx context.Context, memberRoleId int64) map[string]interface{} {
	return perm_domain.Instance().WorkItemEditPerm.GetPermissionWithScene(memberRoleId, "", perm_domain.CreateWorkItemScene)
}

type GetPermissionWithSceneForWorkItemScenceRequest struct {
	Uid               int64
	WorkItemFacade    *facade.WorkItemFacade
	SpaceMemberFacade *facade.SpaceMemberFacade
}

func (s *PermService) GetPermissionWithSceneForWorkItemScence(ctx context.Context, req *GetPermissionWithSceneForWorkItemScenceRequest) map[string]interface{} {

	workItem := req.WorkItemFacade.GetWorkItem()
	member := req.SpaceMemberFacade.GetSpaceMember()

	workItemRole := workItem.GetRole(req.Uid)

	return perm_domain.Instance().WorkItemEditPerm.GetPermissionWithScene(member.GetRole(), workItemRole, perm_domain.EditWorkItemScene)
}

func (s *PermService) GetPermissionWithScene(role int64, itemRole string, scene perm_domain.WorkItemEditPermissionScene) map[string]interface{} {
	return perm_domain.Instance().WorkItemEditPerm.GetPermissionWithScene(role, itemRole, scene)
}

func (s *PermService) GetTaskPermission(role int64, itemRole string) map[string]interface{} {
	return perm_domain.Instance().WorkItemEditPerm.GetTaskPermission(role, itemRole)
}
