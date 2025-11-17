package xspan_base

import (
	"github.com/opentracing/opentracing-go"
)

type SpanIF interface {
	opentracing.Span
	//SpanID() string
}
