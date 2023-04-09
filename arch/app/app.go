/**
 * @Author: wenliangzhang
 * @Description:
 * @File: app
 * @Version: 1.0.0
 * @Date: 2022/4/12 12:19 PM
 */
package app

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"

	"git.singularity-ai.com/backend/library/v2/arch/rpc"
	"git.singularity-ai.com/backend/library/v2/arch/store/mysql"
	"git.singularity-ai.com/backend/library/v2/arch/store/redis"
	"git.singularity-ai.com/backend/library/v2/arch/web"
	"git.singularity-ai.com/backend/library/v2/env"
	"git.singularity-ai.com/backend/library/v2/kafka"
	logger "git.singularity-ai.com/backend/library/v2/log"
	"git.singularity-ai.com/backend/library/v2/service"
)

// AppConfig struct
type AppConfig struct {
	// AppName  应用名称
	AppName string `toml:"app_name"`

	// RunMode 运行模式，可选 debug, release
	RunMode string `toml:"run_mode"`

	// HttpListen  Web服务监听的地址, eg: 0.0.0.0:8080
	HTTPListen string `toml:"http_listen"`

	// PRCListen RPC服务监听的地址，eg: 0.0.0.0:8082
	RPCListen string `toml:"rpc_listen"`

	// IDC 当前应用所部署的机房，eg：jx,tx,test
	IDC string `toml:"idc"`

	// AutoSetMaxProcs 是否自动通话PaaS分配的CPU QUOTE调整线程数
	// 值为空 或者 auto 为生效
	AutoSetMaxProcs string `toml:"AutoSetMaxProcs"`
}

// App App由一个或多个Engine组成，每个Engine对应一个Web或RPC的Server。
type App struct {
	config *AppConfig

	// webEngine 目前web引擎使用gin
	webServer *web.WebServer

	// RPCServer
	rpcServer *rpc.RPCServer
}

// DefaultApp 默认的App。这个默认的APP做了很多定制，如需要自定义，可以自己创建。
var DefaultApp = &App{}

// NewDefaultApp 生成一个按照默认环境生成的应用，默认使用conf目录为配置目录
func NewDefaultApp() *App {
	return new(App).Init()
}

// Init 返回一个按照默认环境生成的应用，默认会使用conf目录的父目录作为应用的根
// 目录，等同于NewDefaultApp()。若想使用非默认的方式生成App，请使用NewAppWithFile。
func (app *App) Init() *App {
	return app.InitWithConfigName("app.toml")
}

// InitWithConfigName 使用指定的配置文件进行初始化
// fn string必须在配置文件根目录下，默认为conf，如果需要设置为其它名称如etc,config等，需要先调用
// env.SetConfDirName()进行调整
func (app *App) InitWithConfigName(fn string) *App {
	config := &AppConfig{}
	absPath := filepath.Join(env.ConfRootPath(), fn)
	if _, err := toml.DecodeFile(absPath, config); err != nil {
		panic(fmt.Sprintf("Can't load config file %s: %s", fn, err.Error()))
	}
	return app.InitWithConfig(config)
}

// InitWithConfig 使用配置结构体进行初始化
func (app *App) InitWithConfig(config *AppConfig) *App {
	if config == nil {
		panic("Can't initial App with nil config")
	}

	app.config = config

	if app.config.IDC != "" {
		env.SetIDC(app.config.IDC)
	}

	if app.config.AppName != "" {
		env.SetAppName(app.config.AppName)
	}

	if app.config.RunMode != "" {
		env.SetRunMode(app.config.RunMode)
	}

	if app.config.HTTPListen != "" {
		arr := strings.Split(app.config.HTTPListen, ":")
		if len(arr) >= 2 {
			env.SetRunMode(app.config.RunMode)
			env.SetHttpPort(arr[1])
		}
	}

	// 自动依据PaaS 分配的cpu quota 调整线程数
	// 开启后，能提升25%左右的性能
	if app.config.AutoSetMaxProcs == "" || app.config.AutoSetMaxProcs == "auto" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	if app.config.HTTPListen != "" {
		app.webServer = web.NewWebServer(app.config.RunMode)
		if len(DefaultWebServerMiddlewares) > 0 {
			app.webServer.Use(DefaultWebServerMiddlewares...)
		}
	}

	if app.config.RPCListen != "" {
		app.rpcServer = rpc.NewRPCServer()
	}

	app.InitLog()
	app.InitService()

	log.Printf("[app.Init] inited with: root_path= %s, config_dir= %s, idc= %s, app_name= %s, run_mode= %s",
		env.RootPath(), env.ConfRootPath(), env.IDC(), env.AppName(), env.RunMode())

	return app
}

// DefaultWebServerMiddlewares 默认的Http Server中间件
// todo:多加一个recovery来保证业务日志崩溃后依旧有访问日志
var DefaultWebServerMiddlewares = []web.WebHandlerFunc{
	// gin.Logger(),
	web.GinHandler2WebHandler(gin.Recovery()),
}

// WebServer 获取WebServer的指针
func (app *App) WebServer() *web.WebServer {
	return app.webServer
}

// RPCServer 获取RPCServer的指针
func (app *App) RPCServer() *rpc.RPCServer {
	return app.rpcServer
}

// RunWebServer 运行web应用
func (app *App) RunWebServer() {
	err := app.webServer.Run(app.config.HTTPListen)
	if err != nil {
		log.Fatalf("Can't RunWebServer: %s", err.Error())
	}
}

// RunWebGraceServer 运行web应用
func (app *App) RunWebGraceServer(cancelFuncs []context.CancelFunc) {
	err := app.webServer.RunGrace(app.config.HTTPListen, cancelFuncs)
	if err != nil {
		log.Fatalf("Can't RunWebGraceEngine, HTTPListen=%s, err=%s", app.config.HTTPListen, err.Error())
	}
}

// InitService 初始化Service模块
func (app *App) InitService() {

	confDirName := "services"
	confPath := filepath.Join(env.ConfRootPath(), confDirName)

	if err := loadServicesFromFolder(confPath, env.IDC()); err != nil {
		log.Printf("LOAD Service FAIL: %s\n", err.Error())
	} else {
		log.Println("LOAD Service SUCC")
	}

	redis.Init("")
	mysql.Init("")
	kafka.Init("")
	service.Init("")
	// todo 配置文件热加载
}

func loadServicesFromFolder(confPath, idc string) error {
	return nil
}

// InitLog 初始化Log模块
func (app *App) InitLog() {
	logger.Init("")
}
