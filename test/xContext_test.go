package test

import (
	"context"
	"fmt"
	"github.com/Xiaoze7713/golib/v3/xContext"
	"github.com/Xiaoze7713/golib/v3/xContext/loggers/xlog"
	"github.com/Xiaoze7713/golib/v3/xContext/metrics/xmetric"
	"github.com/Xiaoze7713/golib/v3/xContext/tracers/jaeger_trace"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"testing"
	"time"
)

const name = iota

func RunFunc(ctx *xContext.XContext) {
	for i := 0; i < 1; i++ {
		time.Sleep(time.Second)
		ctx.Info(ctx.OperationName(), i)
		xlog.Info(ctx.OperationName(), i)
		ctx.Info(runtime.Caller(1))
		ctx.SetTag("val", fmt.Sprintf("%v", i))
		ctx.LogFields("hi", "i`am xiaoai")
	}
	ctx.Info(ctx.SerializeSpanContext())
}

//日志自定义格式
func TestContext(t *testing.T) {
	err := xInit()
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := xContext.NewXContext("start")
	defer ctx.Fin()
	ctx1 := xContext.NewChildXContext(ctx, "count1")
	ctx2 := xContext.NewChildXContext(ctx, "count2")
	ctx3 := xContext.NewChildXContext(ctx, "count3")
	ctx4 := xContext.NewChildXContext(ctx, "count4")
	var ctxList = []*xContext.XContext{ctx1, ctx2, ctx3, ctx4}
	w := sync.WaitGroup{}
	for _, c := range ctxList {
		w.Add(1)
		go func(xc *xContext.XContext) {
			RunFunc(xc)
			defer xc.Fin()
			w.Done()
		}(c)
	}
	ctx.Info("lalala ")
	w.Wait()
	time.Sleep(time.Second * 4)
}

func GinServer() (ginInstance *gin.Engine) {
	//if err := g.InitTrans("zh"); err != nil {
	//	panic(fmt.Errorf("init trans failed, err:%v", err))
	//}
	router := gin.New()
	//gin.Logger()
	router.Use(
		// logger.Logger(utils.SetLogger("room.log")),
		gin.Recovery(),
		xContext.DoRequest(),
	)
	//router.Use(requests.DoRequest())
	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "url Not Found",
		})
	})
	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "function not found",
		})
	})
	return router
}

func Sign() (exit chan os.Signal) {
	exit = make(chan os.Signal)
	stopSigs := []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGKILL,
		syscall.SIGTERM,
	}
	signal.Notify(exit, stopSigs...)
	return
}

func TestPointer(t *testing.T) {
	xc := xContext.NewXContext("xx")
	ctx := context.Context(xc)
	ctx.Done()
}

func xInit() error {
	xlog.SetupLogDefault()
	serverName := "x_content_test"
	xmetric.InitXMetric(&xmetric.Config{
		ServerName: serverName,
		//SubSystem:   "metric",
		Environment: "dev",
		IDC:         "beijing",
		IP:          "127.0.0.1",
	}, func(err error) {
		xlog.Error(err)
	})
	trace, _, err := jaeger_trace.NewJaegerTrace(&jaeger_trace.JaegerConfig{
		ServiceName:   serverName,
		Param:         0,
		AgentHostPort: "127.0.0.1:6832",
	})
	if err != nil {
		return err
	}
	xContext.Init(xlog.GetLogger(), trace, xmetric.Handler,
		jaeger_trace.TraceIDFunc,
		jaeger_trace.SpanIDFunc,
		jaeger_trace.DurFunc,
		jaeger_trace.ExtractSpanFromString,
		jaeger_trace.SerializeToString)
	return err
}

func TestContext2(t *testing.T) {
	err := xInit()
	if err != nil {
		fmt.Println(err)
		return
	}
	g := GinServer()
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))
	gg := g.Group("api")
	gg.POST("hello", func(c *gin.Context) {
		ctxI, ok := c.Get("xContext")
		if !ok {
			return
		}
		ctx := ctxI.(*xContext.XContext)
		ctx.Info("pupupu", c.GetHeader("trace_id"))
		ctx1 := xContext.NewChildXContext(ctx, "count1")
		ctx2 := xContext.NewChildXContext(ctx, "count2")
		ctx3 := xContext.NewChildXContext(ctx, "count3")
		ctx4 := xContext.NewChildXContext(ctx, "count4")
		var ctxList = []*xContext.XContext{ctx1, ctx2, ctx3, ctx4}
		w := sync.WaitGroup{}
		for _, c := range ctxList {
			w.Add(1)
			go func(xc *xContext.XContext) {
				RunFunc(xc)
				defer xc.Fin()
				w.Done()
			}(c)
		}
		ctx.Info("lalala ")
		w.Wait()

		c.Set("resp", map[string]interface{}{"data": map[string]interface{}{"reply": "hello"}})
	})

	_ = g.Run(":9033")
	c := <-Sign()
	fmt.Printf("%s", c.String())
}

func TestContext3(t *testing.T) {
	err := xInit()
	if err != nil {
		fmt.Println(err)
		return
	}
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/api/hello", func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context().(*xContext.XContext)
		defer ctx.Fin()
		ctx1 := xContext.NewChildXContext(ctx, "count1")
		ctx2 := xContext.NewChildXContext(ctx, "count2")
		ctx3 := xContext.NewChildXContext(ctx, "count3")
		ctx4 := xContext.NewChildXContext(ctx, "count4")
		ctx5 := xContext.NewXContextWithContextString(ctx, "request_in", "6bb39019d4295881:6187ccb70fd3667f:21aafde365567e92:0")
		var ctxList = []*xContext.XContext{ctx1, ctx2, ctx3, ctx4, ctx5}
		w := sync.WaitGroup{}
		for _, c := range ctxList {
			w.Add(1)
			go func(xc *xContext.XContext) {
				RunFunc(xc)
				defer xc.Fin()
				w.Done()
			}(c)
		}
		ctx.Info("lalala ")
		w.Wait()

		data, _ := jsoniter.Marshal(map[string]interface{}{"data": map[string]interface{}{"reply": "hello"}})
		_, err := writer.Write(data)
		if err != nil {
			return
		}
	})
	s := http.Server{
		Addr:    ":9033",
		Handler: xContext.HttpIntercept(mux),
	}
	err = s.ListenAndServe()

	c := <-Sign()
	fmt.Printf("%s", c.String())
}
func TestName(t *testing.T) {
	a := map[string]string{
		"a": "a",
		"b": "b",
	}
	s, _ := jsoniter.MarshalIndent(a, "\t", string("4"))
	fmt.Println(s)
}

func Context2(ctx xContext.XContext) {

}

func DoubleWrite() {
	ctx := context.WithValue(context.Background(), "span", 1)
	ctx = context.WithValue(ctx, "traceID", "tid")
	fmt.Printf("ctx %v\n", ctx.Value("span"))
	ctx2 := context.WithValue(ctx, "span", 2)
	fmt.Printf("ctx %v trace %v\n", ctx2.Value("span"), ctx2.Value("traceID"))
	fmt.Printf("ctx %v trace %v\n", ctx.Value("span"), ctx2.Value("traceID"))
}

func TestCtx(t *testing.T) {
	DoubleWrite()
}
