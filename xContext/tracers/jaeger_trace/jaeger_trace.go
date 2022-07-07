package jaeger_trace

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
	"io"
	"sync"
	"time"
)

type JaegerConfig struct {
	ServiceName   string  `json:"service_name" toml:"service_name"`
	Param         float64 `json:"param" toml:"param"`
	AgentHostPort string  `json:"agent_host_port" toml:"agent_host_port"`
}

var tracer opentracing.Tracer
var closer io.Closer

func NewJaegerTrace(jConf *JaegerConfig) (opentracing.Tracer, io.Closer, error) {
	once := sync.Once{}
	var err error
	once.Do(func() {
		cfg := jaegercfg.Configuration{
			ServiceName: jConf.ServiceName, // TODO(qingwen): move to config
			Sampler: &jaegercfg.SamplerConfig{
				Type:  jaeger.SamplerTypeConst,
				Param: jConf.Param,
			},
			Reporter: &jaegercfg.ReporterConfig{
				QueueSize:           100,
				BufferFlushInterval: time.Second,
				LogSpans:            true,
				LocalAgentHostPort:  jConf.AgentHostPort, // TODO(qingwen): move to config
			},
			RPCMetrics: false,
		}

		// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
		// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
		// frameworks.
		jLogger := jaeger.NullLogger
		jMetricsFactory := metrics.NullFactory

		// Initialize tracer with a logger and a metrics factory
		tracer, closer, err = cfg.NewTracer(
			jaegercfg.Logger(jLogger),
			jaegercfg.Metrics(jMetricsFactory),
		)
	})
	return tracer, closer, err
}
