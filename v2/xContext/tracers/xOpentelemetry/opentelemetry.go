package xOpentelemetry

import (
	"github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/otel"
)

func NewOpenTelemetryTracer(name string) opentracing.Tracer {
	_ = otel.Tracer(name)
	return nil
}
