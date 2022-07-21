/**
 * @Author: wenliangzhang
 * @Description:
 * @File: web_server
 * @Version: 1.0.0
 * @Date: 2022/4/12 3:35 PM
 */
package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
)

// WebServer 基于http协议的服务
// 这里的实现是基于gin框架，封装了gin的所有的方法
type WebServer struct {
	// 重写所有的路由相关的方法
	*WebRoute
	// 继承引擎本身的其他方法
	*gin.Engine
}

// NewWebServer 创建WebServer
func NewWebServer(mode string) *WebServer {
	gin.SetMode(mode)

	server := &WebServer{
		Engine: gin.New(),
		WebRoute: &WebRoute{
			root: true,
		},
	}
	server.WebRoute.server = server
	server.WebRoute.RouterGroup = &server.Engine.RouterGroup
	// visit http://127.0.0.1:port/debug/pprof/ and you'll see what you want
	if mode == "debug" {
		ginpprof.Wrap(server.Engine)
	}
	return server
}

// RouterRegister 路由注册者
type RouterRegister func(*WebServer)

// RegisterRouter 注册路由
func (webServer *WebServer) RegisterRouter(rr RouterRegister) {
	rr(webServer)
}

const (
	defaultCancelTimeout = time.Second * 5
)

// RunGrace 实现Server接口
func (webServer *WebServer) RunGrace(addr string, cancelFuncs []context.CancelFunc, timeouts ...time.Duration) (err error) {
	cancelTimeout := defaultCancelTimeout
	if len(timeouts) > 0 {
		cancelTimeout = timeouts[0]

	}
	srv := &http.Server{
		Addr:    addr,
		Handler: webServer,
	}
	// 由于服务是在另外的协程中启动，因此需要用channel传递服务异常信息
	fail := make(chan struct{})
	go func(fail chan<- struct{}) {
		log.Printf("[debug] Listening and serving on %s\n", addr)
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("Serve err: %s\n", err)
			fail <- struct{}{}
		}
	}(fail)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of some seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	select {
	case <-fail:
	case <-quit:
	}
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), cancelTimeout)
	defer func() {
		cancel()
		if len(cancelFuncs) > 0 {
			for _, cancelFunc := range cancelFuncs {
				cancelFunc()
			}
		}
	}()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Shutdown err:", err)
	}
	log.Println("Server exist")
	return
}

// Delims 设置模板的分解符
// 重写gin方法
func (webServer *WebServer) Delims(left, right string) *WebServer {
	webServer.Engine.Delims(left, right)
	return webServer
}

// SecureJsonPrefix sets the secureJsonPrefix used in Context.SecureJSON.
// 重写gin方法
func (webServer *WebServer) SecureJsonPrefix(prefix string) *WebServer {
	webServer.SecureJsonPrefix(prefix)
	return webServer
}

// HandleContext re-enter a context that has been rewritten.
// This can be done by setting c.Request.URL.Path to your new target.
// Disclaimer: You can loop yourself to death with this, use wisely.
/*func (webServer *WebServer) HandleContext(wc *WebContext) {
	webServer.Engine.HandleContext(wc.XContext)
}*/

// NoRoute adds handlers for NoRoute. It return a 404 code by default.
// 重写gin方法
func (webServer *WebServer) NoRoute(handlers ...WebHandlerFunc) {
	webServer.Engine.NoRoute(decorateWebHandlers(handlers)...)
}

// NoMethod sets the handlers called when... TODO.
// 重写gin方法
func (webServer *WebServer) NoMethod(handlers ...WebHandlerFunc) {
	webServer.Engine.NoMethod(decorateWebHandlers(handlers)...)
}

// Use adds middleware to the group, see example code in github.
func (webServer *WebServer) Use(middleware ...WebHandlerFunc) WebRouter {
	return webServer.WebRoute.Use(middleware...)
}
