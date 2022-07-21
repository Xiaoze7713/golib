/**
 * @Author: wenliangzhang
 * @Description:
 * @File: init
 * @Version: 1.0.0
 * @Date: 2022/5/13 4:39 PM
 */
package service

import (
	"git.singularity-ai.com/backend/library/log"
	"git.singularity-ai.com/backend/library/service/token"
	"git.singularity-ai.com/backend/library/xContext"
	"git.singularity-ai.com/backend/library/xContext/tracers/jaeger_trace"
)

// todo 待完善
func Init(configPath string) {
	token.InitToken("")
	jaeger_trace.Init("")
	xContext.Init(log.GetLogger(), jaeger_trace.GetTracer(), nil,
		jaeger_trace.TraceIDFunc,
		jaeger_trace.SpanIDFunc,
		jaeger_trace.DurFunc,
		jaeger_trace.ExtractSpanFromString,
		jaeger_trace.SerializeToString)
}
