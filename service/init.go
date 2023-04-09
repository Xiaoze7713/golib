/**
 * @Author: wenliangzhang
 * @Description:
 * @File: init
 * @Version: 1.0.0
 * @Date: 2022/5/13 4:39 PM
 */
package service

import (
	"git.singularity-ai.com/backend/library/v2/apollo"
	"git.singularity-ai.com/backend/library/v2/arch/web"
	"git.singularity-ai.com/backend/library/v2/log"
	"git.singularity-ai.com/backend/library/v2/service/token"
	"git.singularity-ai.com/backend/library/v2/xContext"
	"git.singularity-ai.com/backend/library/v2/xContext/tracers/jaeger_trace"
)

// todo 待完善
func Init(configPath string) {
	token.InitToken("")
	jaeger_trace.Init("")
	apollo.Init("")
	logger := log.GetLogger()
	weblogger := web.WebLogger{
		logger,
		2,
	}
	xContext.Init(&weblogger, jaeger_trace.GetTracer(), nil,
		jaeger_trace.TraceIDFunc,
		jaeger_trace.SpanIDFunc,
		jaeger_trace.DurFunc,
		jaeger_trace.ExtractSpanFromString,
		jaeger_trace.SerializeToString)
}
