package xRequest

import (
	"fmt"
	"git.singularity-ai.com/backend/library/v2/utils"
	"git.singularity-ai.com/backend/library/v2/xContext"
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
	"git.singularity-ai.com/backend/library/v2/xContext/metrics/xmetric"
	"git.singularity-ai.com/backend/library/v2/xContext/tracers/jaeger_trace"
	"net/http"
	"testing"
	"time"
)

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
		SamplingUrl:   "",
		User:          "",
		Passwd:        "",
	}, nil)
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
func TestHttpBase_Request(t *testing.T) {
	cc := HttpBase{
		Url: "http://39.99.233.6:7775/qa_rank/match",
		Cli: &http.Client{Timeout: time.Second * 4},
	}
	err := xInit()
	if err != nil {
		return
	}
	ctx := xContext.NewXContext("request_xx")
	resp := map[string]interface{}{}
	req := map[string]string{"hi": "hello"}
	err = cc.Post(ctx, req, &resp)
	fmt.Printf("%v, %s", err, utils.MustJson(resp))
}
