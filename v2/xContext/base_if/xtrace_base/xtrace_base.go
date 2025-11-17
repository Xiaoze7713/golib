package xtrace_base

import "github.com/opentracing/opentracing-go"

type TracerIF interface {
	opentracing.Tracer
}

type XTraceNoop struct {
	opentracing.NoopTracer
}
