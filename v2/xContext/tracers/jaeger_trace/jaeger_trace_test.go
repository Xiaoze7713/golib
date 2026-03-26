package jaeger_trace

import (
	"testing"
	"time"

	"github.com/golib/v2/xContext"
	"github.com/golib/v2/xContext/loggers/xlog"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func newTestTracer(t *testing.T) (opentracing.Tracer, func()) {
	t.Helper()
	conf := &JaegerConfig{
		ServiceName: "test-service",
		Param:       1,
		Type:        APITypeByteAPM,
		EndPoint:    "http://localhost:14268/api/traces",
		//EndPoint: "",
	}

	// 每次创建新 tracer，绕过包级别 sync.Once
	tr, cl, err := NewJaegerTrace(conf)
	if err != nil {
		t.Fatalf("NewJaegerTrace error: %v", err)
	}
	return tr, func() {
		if cl != nil {
			_ = cl.Close()
		}
	}
}

func TestNewJaegerTrace_Default(t *testing.T) {
	conf := &JaegerConfig{
		ServiceName: "test-service",
		Param:       1,
		Type:        APITypeDefault,
		EndPoint:    "http://localhost:14268/api/traces",
		User:        "admin",
		Passwd:      "secret",
	}
	tr, cl, err := NewJaegerTrace(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("tracer should not be nil")
	}
	if cl != nil {
		_ = cl.Close()
	}
}

func TestNewJaegerTrace_AgentLocal(t *testing.T) {
	conf := &JaegerConfig{
		ServiceName:   "test-service-agent",
		Param:         1,
		Type:          APITypeAgentLocal,
		AgentHostPort: "127.0.0.1:6831",
	}
	// sync.Once 是局部变量，每次调用均独立
	tr, cl, err := NewJaegerTrace(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("tracer should not be nil")
	}
	if cl != nil {
		_ = cl.Close()
	}
}

func TestNewJaegerTrace_ByteAPM(t *testing.T) {
	conf := &JaegerConfig{
		ServiceName: "test-service-byteapm",
		Param:       1,
		Type:        APITypeByteAPM,
		EndPoint:    "",
		ApiKey:      "my-api-key",
	}
	tr, cl, err := NewJaegerTrace(conf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("tracer should not be nil")
	}
	if cl != nil {
		_ = cl.Close()
	}
}

func TestGetTracer(t *testing.T) {
	// GetTracer 返回包级别 tracer，初始为 nil
	_ = GetTracer()
	_ = GetCloser()
}

func TestInit_FileNotExist(t *testing.T) {
	_, _, err := Init("/nonexistent/path/jaeger.toml")
	if err == nil {
		t.Fatal("expected error for non-existent config file")
	}
}

func TestJaegerSpanID(t *testing.T) {
	tr, cleanup := newTestTracer(t)
	defer cleanup()

	span := tr.StartSpan("test-op")
	defer span.Finish()

	spanCtx := span.Context().(jaeger.SpanContext)
	spanID := JaegerSpanID{spanCtx}

	if spanID.String() == "" {
		t.Error("SpanID String() should not be empty")
	}
	_ = spanID.Value() // 验证不 panic
}

func TestJaegerTraceID(t *testing.T) {
	tr, cleanup := newTestTracer(t)
	defer cleanup()

	span := tr.StartSpan("test-op")
	defer span.Finish()

	spanCtx := span.Context().(jaeger.SpanContext)
	traceID := JaegerTraceID{spanCtx}

	if traceID.String() == "" {
		t.Error("TraceID String() should not be empty")
	}
	_ = traceID.Value()
}

func TestSpanIDFunc(t *testing.T) {
	tr, cleanup := newTestTracer(t)
	defer cleanup()

	xContext.Init(xlog.GetLogger(), tr, nil, SpanIDFunc, TraceIDFunc, DurFunc, nil, nil)

	x := xContext.NewXContext("test-span")
	defer x.Fin()

	id := SpanIDFunc(x)
	if id == nil {
		t.Error("SpanIDFunc should not return nil")
	}
	_ = id.String()
	_ = id.Value()
}

func TestTraceIDFunc(t *testing.T) {
	tr, cleanup := newTestTracer(t)
	defer cleanup()
	defer time.Sleep(3)
	xlog.SetupLogDefault()
	xContext.Init(xlog.GetLogger(), tr, nil, SpanIDFunc, TraceIDFunc, DurFunc, nil, nil)
	//xlog.SetLevel(xlog.DEBUG)
	x := xContext.NewXContext("test-span")
	defer x.Fin()

	id := TraceIDFunc(x)
	if id == nil {
		t.Error("TraceIDFunc should not return nil")
	}
	_ = id.String()
	_ = id.Value()
	xlog.Infof("id %s, value %d", id.String(), id.Value())
}

func TestDurFunc(t *testing.T) {
	tr, cleanup := newTestTracer(t)
	defer cleanup()

	xContext.Init(nil, tr, nil, SpanIDFunc, TraceIDFunc, DurFunc, nil, nil)

	x := xContext.NewXContext("test-span")
	time.Sleep(1 * time.Millisecond)
	x.Fin()

	dur := DurFunc(x)
	if dur < 0 {
		t.Errorf("DurFunc returned negative duration: %v", dur)
	}
}

func TestSerializeToStringAndExtract(t *testing.T) {
	tr, cleanup := newTestTracer(t)
	defer cleanup()

	parent := tr.StartSpan("parent-op")
	spanStr := SerializeToString(parent)
	if spanStr == "" {
		t.Fatal("SerializeToString should not return empty string")
	}

	child, err := ExtractSpanFromString("child-op", spanStr, tr)
	if err != nil {
		t.Fatalf("ExtractSpanFromString error: %v", err)
	}
	if child == nil {
		t.Fatal("child span should not be nil")
	}
	child.Finish()
	parent.Finish()
}

func TestTraceReport(t *testing.T) {
	tr, cleanup := newTestTracer(t)
	defer cleanup()
	defer time.Sleep(3)
	xlog.SetupLogDefault()
	xContext.Init(xlog.GetLogger(), tr, nil, SpanIDFunc, TraceIDFunc, DurFunc, nil, nil)
	xctx := xContext.NewXContext("test-span")
	defer xctx.Fin()
	child := xContext.NewChildXContext(xctx, "child-op")
	go func() {
		t1 := time.Now()
		child.SetTag("v1", "hi")
		child.LogFields("active", "start_child")
		time.Sleep(500 * time.Millisecond)
		child.LogFields("time1", time.Since(t1).Milliseconds())
		time.Sleep(4 * time.Second)
		child.SetTag("v2", "bye")
		child.LogFields("time2", time.Since(t1).Milliseconds())
		child.LogFields("active", "end_child")
		child.Fin()
	}()
	follow := xContext.NewFollowXContext(xctx.Span.Context(), "follow-op")
	go func() {
		t1 := time.Now()
		follow.SetTag("v1", "hi")
		follow.LogFields("active", "start_follow")
		time.Sleep(200 * time.Millisecond)
		follow.LogFields("time1", time.Since(t1))
		time.Sleep(2200 * time.Millisecond)
		follow.SetTag("v2", "hello?")
		follow.LogFields("time2", time.Since(t1))
		follow.LogFields("active", "end_follow")
		follow.Fin()
	}()
	xlog.Info("waiting for fin")
	time.Sleep(10 * time.Second)
	xlog.Info("all finished")
}
