package command

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewCreateSpaceTaskCmd,
	NewCreateSpaceSubTaskCmd,
	NewAddWorkItemCommentCommand,
	NewAddUserLoginLogCommand,
)
