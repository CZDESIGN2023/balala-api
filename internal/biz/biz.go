package biz

import (
	"go-cs/internal/biz/command"
	"go-cs/internal/biz/query"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	wire.NewSet(
		command.ProviderSet,
		query.ProviderSet,
	),

	NewMiddlearesecase,
	NewUserUsecase,
	NewUploadUsecase,
	NewConfigUsecase,
	NewLoginUsecase,
	NewSpaceUsecase,
	NewSpaceTagUsecase,
	NewSpaceWorkItemFlowUsecase,
	NewSpaceWorkObjectUsecase,
	NewSpaceMemberUsecase,
	NewSearchUsecase,
	NewFileInfoUsecase,
	NewStaticsUsecase,
	NewLogUsecase,
	NewSpaceWorkVersionUsecase,
	NewAdminUsecase,
	NewWorkFlowUsecase,
	NewWorkItemStatusUsecase,
	NewWorkItemRoleUsecase,
	NewSpaceFileInfoUsecase,
	NewSpaceWorkItemUsecase,
	NewRptUsecase,
	NewSpaceViewUsecase,
)
