/**
 * @Author: wenliangzhang
 * @Description:
 * @File: init
 * @Version: 1.0.0
 * @Date: 2022/5/13 4:39 PM
 */
package service

import (
	"github.com/golib/v2/apollo"
	"github.com/golib/v2/service/token"
	"github.com/golib/v2/xContext"
	"github.com/golib/v2/xContext/loggers/xlog"
	"github.com/golib/v2/xContext/metrics/null_metric"
	"github.com/golib/v2/xContext/tracers/jaeger_trace"
)

// todo 待完善
func Init(configPath string) {
	token.InitToken("")
	jaeger_trace.Init("")
	apollo.Init("")
	xlog.SetupLogDefault()
	xContext.Init(xlog.GetLogger(), jaeger_trace.GetTracer(), null_metric.MetricsNull{},
		jaeger_trace.TraceIDFunc,
		jaeger_trace.SpanIDFunc,
		jaeger_trace.DurFunc,
		jaeger_trace.ExtractSpanFromString,
		jaeger_trace.SerializeToString)
}
