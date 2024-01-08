package test

import (
	"github.com/Xiaoze7713/golib/v3/xContext/tracers/jaeger_trace"
	"testing"
)

func TestJaeger(t *testing.T) {
	jaeger_trace.NewJaegerTrace(&jaeger_trace.JaegerConfig{
		ServiceName:   "test",
		Param:         1,
		AgentHostPort: "",
	})
}
