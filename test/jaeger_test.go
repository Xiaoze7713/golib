package test

import (
	"git.singularity-ai.com/backend/library/v2/xContext/tracers/jaeger_trace"
	"testing"
)

func TestJaeger(t *testing.T) {
	jaeger_trace.NewJaegerTrace(&jaeger_trace.JaegerConfig{
		ServiceName:   "test",
		Param:         1,
		AgentHostPort: "",
	})
}
