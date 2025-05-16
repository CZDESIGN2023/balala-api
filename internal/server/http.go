package server

import (
	adminV1 "go-cs/api/admin/v1"
	configV1 "go-cs/api/config/v1"
	logV1 "go-cs/api/log/v1"
	loginV1 "go-cs/api/login/v1"
	searchV1 "go-cs/api/search/v1"
	spaceV1 "go-cs/api/space/v1"
	spaceMemberV1 "go-cs/api/space_member/v1"
	spaceTagV1 "go-cs/api/space_tag/v1"
	spaceWorkItemV1 "go-cs/api/space_work_item/v1"
	spaceWorkItemFlowV1 "go-cs/api/space_work_item_flow/v1"
	spaceWorkObjectV1 "go-cs/api/space_work_object/v1"
	spaceWorkVersionV1 "go-cs/api/space_work_version/v1"
	userv1 "go-cs/api/user/v1"
	workFlowV1 "go-cs/api/work_flow/v1"
	workbenchV1 "go-cs/api/workbench/v1"

	rptV1 "go-cs/api/rpt/v1"
	viewV1 "go-cs/api/space_view/v1"
	workItemRoleV1 "go-cs/api/work_item_role/v1"
	workItemStatusV1 "go-cs/api/work_item_status/v1"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/handlers"

	"go-cs/internal/biz"
	"go-cs/internal/conf"
	"go-cs/internal/server/auth/server3auth"
	"go-cs/internal/server/file"
	"go-cs/internal/server/middleware"
	ws "go-cs/internal/server/websock"
	"go-cs/internal/service"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
	_ *server3auth.AuthServer,
	c *conf.Server,
	wss *ws.Server,
	fileServer *file.Server,
	user *service.UserService,
	login *service.LoginService,
	config *service.ConfigService,
	space *service.SpaceService,
	spaceMember *service.SpaceMemberService,
	spaceTag *service.SpaceTagService,
	spaceWorkItem *service.SpaceWorkItemService,
	spaceWorkObject *service.SpaceWorkObjectService,
	search *service.SearchService,
	admin *service.AdminService,
	workFlow *service.WorkFlowService,

	rdb *redis.Client,
	workbench *service.WorkbenchService,
	logger log.Logger,
	spaceWorkItemFlow *service.SpaceWorkItemFlowService,
	logService *service.LogService,
	middlewareBiz *biz.MiddleareUsecase,
	spaceWorkVersion *service.SpaceWorkVersionService,

	workItemRole *service.WorkItemRoleService,
	workItemStatus *service.WorkItemStatusService,
	rpt *service.RptService,
	view *service.SpaceViewService,

) *http.Server {

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			metadata.Server(),
			selector.Server(
				server3auth.Server(rdb),
				middleware.HttpAccessLog(logger),
			).Match(NewWhiteListMatcher()).Build(),
			tracing.Server(),
			middlewareBiz.OpLoggerMiddleware(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	corsOrigins := []string{"*"}
	if len(c.Http.CorsOrigins) > 0 {
		corsOrigins = c.Http.CorsOrigins
	}

	//配置跨域
	opts = append(opts, http.Filter(
		handlers.CORS(
			handlers.AllowCredentials(),
			handlers.AllowedHeaders([]string{"Content-Type", "Access-Control-Allow-Origin", "X-Requested-With", "Access-Control-Allow-Credentials", "User-Agent", "Content-Length", "Authorization"}),
			handlers.AllowedOrigins(corsOrigins),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}),
		),
	))

	srv := http.NewServer(opts...)

	// 消息推送
	wss.SetHttpServer(srv)
	srv.Route("/").GET("/ws/chat", wss.UpgradeHandler)

	userv1.RegisterUserHTTPServer(srv, user)
	configV1.RegisterConfigHTTPServer(srv, config)
	loginV1.RegisterLoginHTTPServer(srv, login)
	spaceV1.RegisterSpaceHTTPServer(srv, space)
	spaceTagV1.RegisterSpaceTagHTTPServer(srv, spaceTag)
	spaceWorkItemV1.RegisterSpaceWorkItemHTTPServer(srv, spaceWorkItem)
	spaceWorkObjectV1.RegisterSpaceWorkObjectHTTPServer(srv, spaceWorkObject)
	spaceMemberV1.RegisterSpaceMemberHTTPServer(srv, spaceMember)
	searchV1.RegisterSearchHTTPServer(srv, search)
	workbenchV1.RegisterWorkbenchHTTPServer(srv, workbench)
	spaceWorkItemFlowV1.RegisterSpaceWorkItemFlowHTTPServer(srv, spaceWorkItemFlow)
	logV1.RegisterLogHTTPServer(srv, logService)
	spaceWorkVersionV1.RegisterSpaceWorkVersionHTTPServer(srv, spaceWorkVersion)
	adminV1.RegisterAdminHTTPServer(srv, admin)
	workFlowV1.RegisterWorkFlowHTTPServer(srv, workFlow)
	workItemRoleV1.RegisterWorkItemRoleHTTPServer(srv, workItemRole)
	workItemStatusV1.RegisterWorkItemStatusHTTPServer(srv, workItemStatus)
	rptV1.RegisterRptHTTPServer(srv, rpt)
	viewV1.RegisterSpaceViewHTTPServer(srv, view)

	// 文件相关
	fileServer.RegisterHTTPServer(srv)

	return srv
}
