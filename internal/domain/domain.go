package domain

import (
	notify_snapshot "go-cs/internal/domain/notify_snapshot/service"
	space_view "go-cs/internal/domain/space_view/service"

	file_info "go-cs/internal/domain/file_info/service"
	file_task "go-cs/internal/domain/file_info/task"
	perm "go-cs/internal/domain/perm/service"
	space "go-cs/internal/domain/space/service"
	space_task "go-cs/internal/domain/space/task"
	space_file_info "go-cs/internal/domain/space_file_info/service"
	space_member "go-cs/internal/domain/space_member/service"
	space_tag "go-cs/internal/domain/space_tag/service"
	space_work_item_comment "go-cs/internal/domain/space_work_item_comment/service"
	space_work_object "go-cs/internal/domain/space_work_object/service"
	space_work_version "go-cs/internal/domain/space_work_version/service"
	user "go-cs/internal/domain/user/service"
	user_login_log "go-cs/internal/domain/user_login_log/service"
	work_flow "go-cs/internal/domain/work_flow/service"
	work_item "go-cs/internal/domain/work_item/service"
	work_item_task "go-cs/internal/domain/work_item/task"
	work_item_role "go-cs/internal/domain/work_item_role/service"
	work_item_status "go-cs/internal/domain/work_item_status/service"
	work_item_type "go-cs/internal/domain/work_item_type/service"

	message "go-cs/internal/domain/pkg/message"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(

	wire.NewSet(
		space.NewSpaceService,
		space_task.NewSpaceTask,
	),

	wire.NewSet(
		file_info.NewFileInfoService,
		file_task.NewFileInfoTask,
	),

	wire.NewSet(
		work_item.NewWorkItemService,
		work_item_task.NewWorkItemTask,
	),

	space_tag.NewSpaceTagService,
	space_member.NewSpaceMemberService,
	user.NewUserService,
	user_login_log.NewUserLoginLogService,
	work_item_role.NewWorkItemRoleService,
	work_item_status.NewWorkItemStatusService,
	perm.NewPermService,
	work_flow.NewWorkFlowService,
	work_item_type.NewWorkItemTypeService,
	notify_snapshot.NewNotifySnapShotService,

	space_work_object.NewSpaceWorkObjectService,
	space_file_info.NewSpaceFileInfoService,
	space_work_version.NewSpaceWorkVersionService,
	space_work_item_comment.NewSpaceWorkItemCommentService,
	space_view.NewSpaceViewService,

	wire.NewSet(
		message.NewDomainMessageConsumer,
		message.NewDomainMessageProducer,
	),
)
