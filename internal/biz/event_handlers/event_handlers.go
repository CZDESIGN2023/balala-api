package event_handlers

import (
	comment_handler "go-cs/internal/biz/event_handlers/comment"
	notify_handler "go-cs/internal/biz/event_handlers/notify"
	"go-cs/internal/biz/event_handlers/ops_log"
	"go-cs/internal/utils"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewAppEventHandlers,
	ops_log.NewOpsLogEventHandlers,
	notify_handler.NewNotify,
	comment_handler.NewCommentHandler,
)

type AppEventHandlers struct {
	opsLogEventHandler *ops_log.OpsLogEventHandlers
	notifyHandler      *notify_handler.Notify
	commentHandler     *comment_handler.CommentHandler
	log                *log.Helper
}

func NewAppEventHandlers(
	logger log.Logger,

	opsLogEventHandler *ops_log.OpsLogEventHandlers,
	notifyHandler *notify_handler.Notify,
	commentHandler *comment_handler.CommentHandler,

) *AppEventHandlers {

	moduleName := "AppEventHandlers"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	evtHandlers := &AppEventHandlers{
		log:                hlog,
		opsLogEventHandler: opsLogEventHandler,
		notifyHandler:      notifyHandler,
		commentHandler:     commentHandler,
	}

	evtHandlers.init()
	return evtHandlers
}

func (s *AppEventHandlers) init() {
	s.opsLogEventHandler.Init()
	s.notifyHandler.Init()
	s.commentHandler.Init()
}
