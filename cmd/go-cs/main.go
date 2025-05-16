package main

import (
	"context"
	"flag"
	"go-cs/api/comm"
	"go-cs/internal/dwh"
	"go-cs/internal/utils/third_platform"
	"go-cs/internal/utils/third_platform/adapter"
	"go-cs/pkg/i18n"
	"go-cs/pkg/qqwry"

	"github.com/go-kratos/kratos/v2/config/env"

	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"go.uber.org/zap"

	"go-cs/internal/conf"
	"go-cs/pkg/log"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/transport/http"

	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	kconfig "github.com/go-kratos/kratos/v2/config"
	klog "github.com/go-kratos/kratos/v2/log"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	hostname, _ = os.Hostname()
)

var (
	logger  klog.Logger    // kratos log interface
	hLogger *klog.Helper   // kratos log helper
	zLogger *zap.Logger    // zap logger
	config  conf.Bootstrap // kratos config interface
)

const WaitSignalTimeout = 4

const VERSION = "1.4.14.1"

func init() {
	Name = "balala-service"
	flag.StringVar(&flagconf, "c", "./configs/config.yaml", "config path, eg: -c config.yaml")
	//这几个入参做标记而已
	flag.String("tag", "", "eg: -tag test")
	flag.String("ver", "", "eg: -ver 1.0")
}

func setup() {
	var err error
	var c = config.Log
	var writer io.Writer

	var logFilePath = c.LogPath + c.LogFilename
	var logLevel = c.LogLevel
	var logFormat = "json"
	var logCallerSkip = 3 // 從zap往上跳三層才能確實定位到寫log時所在的檔案及行數

	if c.LogFormat != "" {
		logFormat = c.LogFormat
	}

	writer = os.Stdout
	if c.LogFilename != "" { // 日志文件
		writer, err = rotatelogs.New(
			logFilePath+".%Y%m%d",
			rotatelogs.WithClock(rotatelogs.Local),
			rotatelogs.WithLinkName(logFilePath),
			rotatelogs.WithMaxAge(time.Hour*24*7),     // 保留最近7天日志
			rotatelogs.WithRotationTime(time.Hour*24), // 按天分割
		)
		if err != nil {
			panic(err)
		}
	}

	zLogger = log.New(
		log.WithLevel(log.Level(logLevel)),
		log.WithFormat(log.Format(logFormat)),
		log.WithWriter(writer),
		log.WithCallerSkip(logCallerSkip),
	)

	logger = klog.With(
		kzap.NewLogger(zLogger),
		//"service.id", hostname,
		//"service.name", Name,
	)

	// 设置默认 logger
	klog.SetLogger(logger)
	klog.DefaultMessageKey = "message" // 更改預設key避免和zap內部輸出的msg欄位重複, 不利log收集分析

	hLogger = klog.NewHelper(logger)

	hLogger.Info("initializing resource ...")
	hLogger.Infof("the log output directory: %s", filepath.Dir(logFilePath))
}

func newApp(logger klog.Logger, zLogger *zap.Logger, hs *http.Server) *kratos.App {
	// TODO:使用端口加主机名作为服务ID
	return kratos.New(
		kratos.ID(hostname),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
		),
	)
}

func main() {
	flag.Parse()

	c := kconfig.New(
		kconfig.WithSource(
			env.NewSource(),
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap

	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	config = bc

	//加载语言包
	i18n.InitI18n(bc.I18N.I18NPath)

	//加载IP库
	qqwry.InitQQwry(bc.Data.Qqwry.DatPath)

	setup()
	// setTracerProvider()

	//初始化 thirdPlatform 推送
	tpClient := initThirdPlatformClient(&bc)

	dwhApp, cleanup, err := dwh.NewApp(bc.Dwh, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	//初始化biz事件处理handler
	app, cleanup, err := wireApp(&bc, bc.Server, bc.Data, bc.Jwt, bc.FileConfig, bc.Dwh, dwhApp, dwhApp.Service(), tpClient, logger, zLogger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	app.SetVersion(VERSION)

	bgCtx := context.Background()
	appStop, err := app.Start(bgCtx)
	if err != nil {
		panic(err)
	}

	// monitor signal
	signalCtx, signalStop := signal.NotifyContext(bgCtx, syscall.SIGINT, syscall.SIGTERM)
	defer signalStop()

	// waiting exit signal ...
	select {
	case err = <-appStop:
		hLogger.Error(err)
	case <-signalCtx.Done():
	}
	signalStop()

	hLogger.Info("the app is shutting down ...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(WaitSignalTimeout)*time.Second)
	defer cancel()

	if err = app.Stop(ctx); err != nil {
		hLogger.Error(err)
	}
}

// 初始化第三方平台推送客户端
func initThirdPlatformClient(config *conf.Bootstrap) *third_platform.Client {
	tpClient := third_platform.NewClient()

	if conf := config.TeaIm; conf != nil {
		cli := adapter.NewIMClient(conf.PlatformCode, conf.PrivateKey, conf.Domain)
		tpClient.Add(comm.ThirdPlatformCode_pf_IM, cli)
	}

	if conf := config.Ql; conf != nil {
		cli := adapter.NewIMClient(conf.PlatformCode, conf.PrivateKey, conf.Domain)
		tpClient.Add(comm.ThirdPlatformCode_pf_QL, cli)
	}

	if conf := config.Halala; conf != nil {
		cli := adapter.NewIMClient(conf.PlatformCode, conf.PrivateKey, conf.Domain)
		tpClient.Add(comm.ThirdPlatformCode_pf_Halala, cli)
	}

	return tpClient
}
