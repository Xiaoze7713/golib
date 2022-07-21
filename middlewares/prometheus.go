/**
 * @Author: wenliangzhang
 * @Description:
 * @File: prometheus
 * @Version: 1.0.0
 * @Date: 2022/6/6 3:47 PM
 */
package middlewares

import (
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"git.singularity-ai.com/backend/library/env"
)

const (
	metricsPath = "/metrics"
	faviconPath = "/favicon.ico"
)

var (
	httpHistogramCostGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "service",
		Subsystem:   "",
		Name:        "costVer",
		Help:        "Histogram of response latency (Milliseconds) of http handlers.",
		ConstLabels: nil,
	}, []string{"idc", "code", "uri"})
	httpHistogramCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   "service",
		Subsystem:   "",
		Name:        "count",
		Help:        "Histogram of response count of http handlers.",
		ConstLabels: nil,
	}, []string{"idc", "code", "errno", "uri"})
)

// init 初始化prometheus模型
func init() {
	prometheus.MustRegister(httpHistogramCount, httpHistogramCostGauge)
}

// handlerPath 定义采样路由struct
type handlerPath struct {
	sync.Map
}

// get 获取path
func (hp *handlerPath) get(handler string) string {
	v, ok := hp.Load(handler)
	if !ok {
		return ""
	}
	return v.(string)
}

// set 保存path到sync.Map
func (hp *handlerPath) set(ri gin.RouteInfo) {
	hp.Store(ri.Handler, ri.Path)
}

// GinPrometheus gin调用Prometheus的struct
type GinPrometheus struct {
	engine  *gin.Engine
	ignored map[string]bool
	pathMap *handlerPath
	updated bool
}

type Option func(*GinPrometheus)

// Ignore 添加忽略的路径
func Ignore(path ...string) Option {
	return func(gp *GinPrometheus) {
		for _, p := range path {
			gp.ignored[p] = true
		}
	}
}

// New new gin prometheus
func NewPrometheus(e *gin.Engine, options ...Option) *GinPrometheus {
	if e == nil {
		return nil
	}

	gp := &GinPrometheus{
		engine: e,
		ignored: map[string]bool{
			metricsPath: true,
			faviconPath: true,
		},
		pathMap: &handlerPath{},
	}

	for _, o := range options {
		o(gp)
	}
	return gp
}

// updatePath 更新path
func (gp *GinPrometheus) updatePath() {
	gp.updated = true
	for _, ri := range gp.engine.Routes() {
		gp.pathMap.set(ri)
	}
}

// Middleware set gin middleware
func (gp *GinPrometheus) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !gp.updated {
			gp.updatePath()
		}
		// 过滤请求
		if gp.ignored[c.Request.URL.String()] {
			c.Next()
			return
		}

		start := time.Now()
		c.Next()

		end := time.Now()
		// 执行时间
		cost := end.Sub(start)
		httpHistogramCostGauge.WithLabelValues(
			env.IDC(),
			strconv.Itoa(c.Writer.Status()),
			c.Request.URL.Path,
		).Set(float64(cost.Milliseconds()))

		errno, isExit := c.Get("errno")
		if !isExit {
			errno = 0
		}
		httpHistogramCount.WithLabelValues(
			env.IDC(),
			strconv.Itoa(c.Writer.Status()),
			strconv.Itoa(errno.(int)),
			c.Request.URL.Path,
		).Inc()
	}
}
