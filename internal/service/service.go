package service

import (
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewUserService,
	NewConfigService,
	NewLoginService,
	NewSpaceService,
	NewSpaceTagService,
	NewSpaceWorkItemService,
	NewSpaceMemberService,
	NewSpaceWorkObjectService,
	NewSearchService,
	NewWorkbenchService,
	NewSpaceWorkItemFlowService,
	NewLogService,
	NewSpaceWorkVersionService,
	NewAdminService,
	NewWorkItemRoleService,
	NewWorkFlowService,
	NewWorkItemStatusService,
	NewRptService,
	NewSpaceViewService,
)
