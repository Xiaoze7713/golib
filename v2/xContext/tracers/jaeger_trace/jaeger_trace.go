package jaeger_trace

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/golib/v2/xContext"
	"github.com/BurntSushi/toml"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

type JaegerConfig struct {
	ServiceName   string  `json:"service_name" toml:"service_name"`
	Param         float64 `json:"param" toml:"param"`
	AgentHostPort string  `json:"agent_host_port" toml:"agent_host_port"`
	EndPoint      string  `json:"end_point" toml:"end_point"`
	SamplingUrl   string  `json:"sampling_url" toml:"sampling_url"`
	User          string  `json:"user" toml:"user"`
	Passwd        string  `json:"passwd" toml:"passwd"`
}

var tracer opentracing.Tracer
var closer io.Closer

var defaultJaegerConfigPath = "conf/service/jaeger.toml"

func Init(filePath string) (opentracing.Tracer, io.Closer, error) {
	var config JaegerConfig
	if filePath == "" {
		filePath = defaultJaegerConfigPath
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("jaeger conf file IsNotExist")
		return nil, nil, err
	}
	if _, err := toml.DecodeFile(filePath, &config); err != nil {
		panic(fmt.Sprintf("Can't load config file, %s", err.Error()))
	}
	return NewJaegerTrace(&config, nil)
}

func GetTracer() opentracing.Tracer {
	return tracer
}

func GetCloser() io.Closer {
	return closer
}

func NewJaegerTrace(jConf *JaegerConfig, options ...jaegercfg.Option) (opentracing.Tracer, io.Closer, error) {
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
			},
			RPCMetrics: false,
		}
		if jConf.AgentHostPort != "" {
			cfg.Reporter.LocalAgentHostPort = jConf.AgentHostPort
		} else if jConf.EndPoint != "" {
			cfg.Reporter.CollectorEndpoint = jConf.EndPoint
			cfg.Reporter.User = jConf.User
			cfg.Reporter.Password = jConf.Passwd
		}

		// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
		// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
		// frameworks.

		// Initialize tracer with a logger and a metrics factory
		tracer, closer, err = cfg.NewTracer(options...)
	})
	return tracer, closer, err
}

type JaegerSpanID struct {
	jaeger.SpanContext
}

func (m JaegerSpanID) String() string {
	return m.SpanContext.SpanID().String()
}

func (m JaegerSpanID) Value() int64 {
	return int64(m.SpanContext.SpanID())
}

func SpanIDFunc(x *xContext.XContext) xContext.IDType {
	spanCtx := x.Span.Context().(jaeger.SpanContext)
	return JaegerSpanID{spanCtx}
}

type JaegerTraceID struct {
	jaeger.SpanContext
}

func (m JaegerTraceID) String() string {
	return m.SpanContext.TraceID().String()
}

func (m JaegerTraceID) Value() int64 {
	return int64(m.SpanContext.TraceID().Low)
}
func TraceIDFunc(x *xContext.XContext) xContext.IDType {
	spanCtx := x.Span.Context().(jaeger.SpanContext)
	return JaegerTraceID{spanCtx}
}

func DurFunc(x *xContext.XContext) time.Duration {
	spanCtx := x.Span.(*jaeger.Span)
	return spanCtx.Duration()
}

func ExtractSpanFromString(operationName, spanStr string, tracer opentracing.Tracer) (newSpan opentracing.Span, err error) {
	spanCtx, err := jaeger.ContextFromString(spanStr)
	newSpan = tracer.StartSpan(operationName, opentracing.ChildOf(spanCtx))
	return newSpan, err
}

func SerializeToString(span opentracing.Span) string {
	jSpan := span.(*jaeger.Span)
	return jSpan.String()
}
