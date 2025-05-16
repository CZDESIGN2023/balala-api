package file

import (
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	kratosHttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cast"
	"go-cs/internal/biz"
	"go-cs/internal/conf"
	file_repo "go-cs/internal/domain/file_info/repo"
	member_repo "go-cs/internal/domain/space_member/repo"
	user_repo "go-cs/internal/domain/user/repo"
	witem_repo "go-cs/internal/domain/work_item/repo"
	"go-cs/internal/server/file/result"
	"go-cs/internal/server/file/result/errs"
	"go-cs/internal/utils"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	ServerTypeDFS   = 1
	ServerTypeS3    = 2
	ServerTypeLocal = 0
	MaxUploadSize   = 5 * 1024 * 1024 // 5MB
	MaxFileSize     = 2048 << 20
)

type Server struct {
	log *log.Helper

	confJwt           *conf.Jwt
	confFile          *conf.FileConfig
	userRepo          user_repo.UserRepo
	fileRepo          file_repo.FileInfoRepo
	spaceMemberRepo   member_repo.SpaceMemberRepo
	spaceWorkItemRepo witem_repo.WorkItemRepo
	configRepo        biz.ConfigRepo
	rdb               *redis.Client
	uploadUsecase     *biz.UploadUsecase
	spaceFileInfoUc   *biz.SpaceFileInfoUsecase
	engine            *gin.Engine
}

func NewServer(logger log.Logger, uploadUsecase *biz.UploadUsecase, spaceFileInfoUc *biz.SpaceFileInfoUsecase,
	userRepo user_repo.UserRepo, spaceMemberRepo member_repo.SpaceMemberRepo, spaceWorkItemRepo witem_repo.WorkItemRepo, configRepo biz.ConfigRepo,
	fileRepo file_repo.FileInfoRepo,
	confJwt *conf.Jwt, confFile *conf.FileConfig, serverFile *conf.Server, rdb *redis.Client) *Server {
	moduleName := "FileServer"
	_, hlog := utils.InitModuleLogger(logger, moduleName)

	gin.SetMode(gin.ReleaseMode)
	g := gin.New()
	g.Use(gin.Recovery(), Cors(serverFile))

	s := &Server{
		log:               hlog,
		userRepo:          userRepo,
		fileRepo:          fileRepo,
		confJwt:           confJwt,
		confFile:          confFile,
		spaceMemberRepo:   spaceMemberRepo,
		spaceFileInfoUc:   spaceFileInfoUc,
		spaceWorkItemRepo: spaceWorkItemRepo,
		configRepo:        configRepo,
		rdb:               rdb,
		uploadUsecase:     uploadUsecase,
		engine:            g,
	}

	return s
}

func (s *Server) RegisterHTTPServer(srv *kratosHttp.Server) {
	srv.HandlePrefix("/", s.engine)

	auth := authMiddleware2Gin(s.confJwt, s.rdb)

	file := s.engine.Group("/file").Use(gin.Logger())
	{
		file.POST("/upload", auth, s.Upload)
		file.POST("/rapid_verify", auth, s.RapidVerify)
		file.POST("/rapid_chunk", auth, s.RapidChunk)
		file.POST("/rapid_merge", auth, s.RapidMerge)

		file.POST("/get_download_token", auth, s.GetDownloadToken)
		file.POST("/download", auth, s.Download)
		file.GET("/download/*filename", s.Download)

		//file.Any("/proxy", s.Proxy)
	}

	//s.engine.GET("/space_file/*filepath", auth, s.spaceFile)
	s.engine.StaticFS("/avatar", DirFs(s.confFile.LocalPath+"/avatar"))
	s.engine.StaticFS("/space_file", DirFs(s.confFile.LocalPath+"/space_file"))
	s.engine.StaticFS("/user_default", DirFs(s.confFile.LocalPath+"/user_default"))
}

func (s *Server) Proxy(ctx *gin.Context) {
	var args ProxyArgs

	if err := ctx.Bind(&args); err != nil {
		result.Fail(ctx, errs.Param(err))
		return
	}

	target, _ := url.Parse(args.Url)
	proxy := httputil.ReverseProxy{Director: func(request *http.Request) {
		request.Host = target.Host
		request.RequestURI = target.RequestURI()
		request.URL = target
		request.Header.Del("Referer")
	}}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}

func (s *Server) RapidVerify(ctx *gin.Context) {
	var args VerifyArgs

	if err := ctx.Bind(&args); err != nil {
		result.Fail(ctx, errs.Param(err))
		return
	}

	uid := GetUserIdFromCtx(ctx)

	data, err := s.verify(ctx, uid, args)
	if err != nil {
		result.Fail(ctx, err)
		return
	}

	result.Ok(ctx, data)
}

func (s *Server) RapidMerge(ctx *gin.Context) {
	var args MergeArgs

	if err := ctx.Bind(&args); err != nil {
		result.Fail(ctx, errs.Param(err))
		return
	}

	uid := GetUserIdFromCtx(ctx)

	data, err := s.merge(ctx, uid, args)
	if err != nil {
		result.Fail(ctx, err)
		return
	}

	result.Ok(ctx, data)
}

func (s *Server) Upload(ctx *gin.Context) {
	var args UploadArgs

	if err := ctx.Bind(&args); err != nil {
		result.Fail(ctx, errs.Param(err))
		return
	}

	uid := GetUserIdFromCtx(ctx)

	data, err := s.uploadHandler(ctx, uid, args)
	if err != nil {
		result.Fail(ctx, err)
		return
	}

	result.Ok(ctx, data)
}

func (s *Server) RapidChunk(ctx *gin.Context) {
	var args Upload2Args

	if err := ctx.Bind(&args); err != nil {
		result.Fail(ctx, errs.Param(err))
		return
	}

	uid := GetUserIdFromCtx(ctx)

	data, err := s.rapidChunk(ctx, uid, args)
	if err != nil {
		result.Fail(ctx, err)
		return
	}

	result.Ok(ctx, data)
}

func (s *Server) GetDownloadToken(ctx *gin.Context) {
	var args GetDownloadTokenArgs

	if err := ctx.Bind(&args); err != nil {
		result.Fail(ctx, errs.Param(err))
		return
	}

	uid := GetUserIdFromCtx(ctx)

	data, err := s.getDownloadToken(ctx, uid, args)
	if err != nil {
		result.Fail(ctx, err)
		return
	}

	result.Ok(ctx, data)
}

func (s *Server) Download(ctx *gin.Context) {
	var args DownloadArgs

	if err := ctx.Bind(&args); err != nil {
		result.Fail(ctx, errs.Param(err))
		return
	}

	_, err := s.download(ctx, args)
	if err != nil {
		result.Fail(ctx, err)
		return
	}
}

func (s *Server) spaceFile(ctx *gin.Context) {
	uid := GetUserIdFromCtx(ctx)

	filePath := ctx.Param("filepath")

	// 获取空间id
	spaceId := extractSpaceIdFromPath(filePath)
	if spaceId <= 0 {
		result.Forbidden(ctx)
		return
	}

	// 判断是否是空间成员
	_, err := s.spaceMemberRepo.GetSpaceMember(ctx, cast.ToInt64(spaceId), uid)
	if err != nil {
		result.Forbidden(ctx)
		return
	}

	ctx.FileFromFS(filePath, DirFs(s.confFile.LocalPath+"/space_file"))
}
