package xmetric

import (
	"fmt"
	"github.com/golib/v3/xContext/base_if/xmetric_base"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func RunForever(metric xmetric_base.MetricsIF) {
	for i := 0; i < 100; i++ {
		//metric.Summary().Observe(1)
		mod10 := prometheus.Labels{"mod": "10"}
		mod3 := prometheus.Labels{"mod": "3"}
		if i%10 == 0 {
			metric.SummaryBy("mod_s", mod10).Observe(1)
			metric.HistogramBy("mod_h", mod10).Observe(1)
			metric.CounterBy("mod_c", mod10).Inc()
		}
		metric.GaugeBy("mod_g", mod3).Set(float64(i % 3))
		metric.GaugeBy("mod_g", mod10).Set(float64(i % 10))
		if i%3 == 0 {
			metric.SummaryBy("mod_s", mod3).Observe(1)
			metric.HistogramBy("mod_h", mod3).Observe(1)
			metric.CounterBy("mod_c", mod3).Inc()
		}
		time.Sleep(time.Second)
	}
}

func GinServer() (ginInstance *gin.Engine) {
	//if err := g.InitTrans("zh"); err != nil {
	//	panic(fmt.Errorf("init trans failed, err:%v", err))
	//}
	router := gin.Default()
	router.Use(
		// logger.Logger(utils.SetLogger("room.log")),
		gin.Recovery(),
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

func TestXMetrics(t *testing.T) {
	InitXMetric(&Config{
		ServerName:  "xmetrics",
		SubSystem:   "test",
		Environment: "dev",
		IDC:         "beijing",
		IP:          "127.0.0.1",
	}, func(err error) {
		fmt.Println(err)
	})
	go func() {
		RunForever(Handler)
	}()
	g := GinServer()
	g.GET("/metrics", gin.WrapH(promhttp.Handler()))
	_ = g.Run(":9033")
	c := <-Sign()
	fmt.Printf("%s", c.String())
}
