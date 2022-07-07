package xContext

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	xlog_base2 "git.singularity-ai.com/backend/library/xContext/base_if/xlog_base"
	"git.singularity-ai.com/backend/library/xContext/base_if/xmetric_base"
	"git.singularity-ai.com/backend/library/xContext/base_if/xspan_base"
	"git.singularity-ai.com/backend/library/xContext/base_if/xtrace_base"
	"git.singularity-ai.com/backend/library/xContext/loggers/null_log"
	"git.singularity-ai.com/backend/library/xContext/metrics/null_metric"
	"github.com/gin-gonic/gin"
	"github.com/go-stack/stack"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/segmentio/ksuid"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

// 监控 链路跟踪 log
var mMetric xmetric_base.MetricsIF
var mTrace xtrace_base.TracerIF
var mLogger xlog_base2.LoggerIF

type GetStrFunc func(xContext *XContext) string
type GetDurFunc func(xContext *XContext) time.Duration
type GetErrorCodeByErr func(err error) int64

// 几个id的获取
var spanIDFunc GetStrFunc
var traceIDFunc GetStrFunc
var durFunc GetDurFunc
var toErrCode GetErrorCodeByErr

func Init(logger xlog_base2.LoggerIF,
	trace xtrace_base.TracerIF,
	metrics xmetric_base.MetricsIF,
	traceF, spanF GetStrFunc, durF GetDurFunc) {
	if logger != nil {
		mLogger = logger
	} else {
		fmt.Printf("warning! not set logger")
		mLogger = null_log.LoggerNull{}
	}
	if trace != nil {
		mTrace = trace
	} else {
		fmt.Printf("warning! not set trace")
		mTrace = xtrace_base.XTraceNoop{}
	}
	if metrics != nil {
		mMetric = metrics
	} else {
		fmt.Printf("warning! not set metric")
		mMetric = null_metric.MetricsNull{}
	}
	emptyIDFunc := func(x *XContext) string {
		return ""
	}
	emptyDurFunc := func(x *XContext) time.Duration {
		return 0
	}
	if traceF == nil {
		traceIDFunc = emptyIDFunc
	} else {
		traceIDFunc = traceF
	}
	if spanF == nil {
		spanIDFunc = emptyIDFunc
	} else {
		spanIDFunc = spanF
	}
	if durF == nil {
		durFunc = emptyDurFunc
	} else {
		durFunc = durF
	}
}

type XContext struct {
	context.Context
	xlog_base2.LoggerIF
	xmetric_base.MetricsIF
	Span          xspan_base.SpanIF
	CancelList    []context.CancelFunc
	operationName string
}

//
//type GXContext struct {
//	gin.Context
//	xlog_base.LoggerIF
//	xmetric_base.MetricsIF
//	opentracing.Span
//	CancelList []context.CancelFunc
//}

var notGrandFather = errors.New("not grand father")

func GetXContextFromGrandFather(ctx context.Context, operationName string) (*XContext, error) {
	gCtx := ctx.Value("xContext")
	xCtx, ok := gCtx.(*XContext)
	if !ok {
		mLogger.Error(notGrandFather)
		xCtx = NewXContextWithContext(ctx, operationName)
		return xCtx, nil
	}
	return xCtx, nil
}

func NewXContextWithContext(ctx context.Context, operationName string) *XContext {
	xCtx := &XContext{
		Context:   ctx,
		LoggerIF:  mLogger,
		MetricsIF: mMetric,
		//lock:      &sync.RWMutex{},
	}
	xCtx.Span = mTrace.StartSpan(operationName,
		opentracing.StartTime{},
	)
	return xCtx
}

func NewXContext(operationName string) *XContext {
	xCtx := &XContext{
		Context:   context.Background(),
		LoggerIF:  mLogger,
		MetricsIF: mMetric,
		//lock:      &sync.RWMutex{},
	}
	xCtx.Span = mTrace.StartSpan(operationName,
		opentracing.StartTime{},
	)
	return xCtx
}

func NewXContextWithParent(origin *XContext, operationName string) *XContext {
	childCtx, cancel := context.WithCancel(origin.Context)
	origin.CancelList = append(origin.CancelList, cancel)
	child := &XContext{
		Context:    childCtx,
		LoggerIF:   mLogger,
		MetricsIF:  mMetric,
		CancelList: []context.CancelFunc{},
		//lock:       origin.lock,
	}
	child.Span = mTrace.StartSpan(operationName,
		opentracing.ChildOf(origin.Span.Context()),
		opentracing.StartTime{},
	)
	return child
}

func NewFollowContext(origin XContext, operationName string) *XContext {
	brother := &XContext{
		Context:    context.Background(),
		LoggerIF:   mLogger,
		MetricsIF:  mMetric,
		CancelList: []context.CancelFunc{},
		//lock:       &sync.RWMutex{},
	}
	brother.Span = mTrace.StartSpan(operationName,
		opentracing.FollowsFrom(origin.Span.Context()),
		opentracing.StartTime{},
	)
	return brother
}

const (
	HttpPath       = ext.StringTagName("path")
	HttpMethod     = ext.StringTagName("method")
	ReqBody        = ext.StringTagName("req_body")
	ReqBodyLen     = ext.StringTagName("req_body_len")
	HttpReqHeader  = ext.StringTagName("req_header")
	RespBody       = ext.StringTagName("resp_body")
	HttpRespHeader = ext.StringTagName("resp_header")
	RespCode       = ext.StringTagName("resp_code")
	RespCodeMsg    = ext.StringTagName("resp_code_msg")
	ReqStatus      = ext.StringTagName("req_status")
	CostMs         = ext.StringTagName("cost_ms")
	ClientIP       = ext.StringTagName("client_ip")
	UserAgent      = ext.StringTagName("user_agent")
	Error          = ext.StringTagName("error")
	ExecStatus     = ext.StringTagName("exec_status")
	TopicName      = ext.StringTagName("topic_name")
	ConsumerSource = ext.StringTagName("consumer_source")
	ProductTarget  = ext.StringTagName("product_target")
	MessageBody    = ext.StringTagName("msg_body")
	EventType      = ext.StringTagName("event_type")
	EventSystem    = ext.StringTagName("event_system")
	EventBody      = ext.StringTagName("event_body")
	SpanID         = "span_id"
	TraceID        = "trace_id"
)

func (x *XContext) SetKVs(kvs map[string]interface{}) {
	for k, v := range kvs {
		x.Store(k, v)
	}
}

func (x *XContext) SpanID() string {
	if spanIDFunc == nil {
		return fmt.Sprintf("%x", x.Span.Context())
	} else {
		return spanIDFunc(x)
	}
}

func (x *XContext) TraceID() string {
	//return x.Load(TraceID)
	return traceIDFunc(x)
}

func (x *XContext) Duration() time.Duration {
	return durFunc(x)
}

func (x *XContext) ErrCode(err error) int64 {
	if toErrCode != nil {
		return toErrCode(err)
	}
	return -1
}

type ContextKey interface {
	string | ext.StringTagName | any
}

func (x *XContext) Store(key ContextKey, value interface{}) {
	x.Context = context.WithValue(x.Context, key, value)
}

// 获取kv

func GrpcGWIntercept() {

}

func HttpIntercept(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := NewXContextWithContext(r.Context(), "request")
		defer ctx.Fin()
		ctx.Store(HttpMethod, r.RequestURI)
		ctx.Store(HttpPath, r.URL.Path)
		// 读body
		var body []byte
		n, _ := r.Body.Read(body)
		bodyStr := string(body)
		ctx.Store(ReqBodyLen, n)
		ctx.Store(ReqBody, bodyStr)
		// 写回
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		//r = r.WithContext(lc)
		// request in
		ctx.Info("request_in",
			HttpMethod, ctx.LoadString(HttpMethod),
			HttpPath, ctx.LoadString(HttpPath),
			ReqBodyLen, ctx.LoadString(ReqBodyLen),
			ReqBody, ctx.LoadString(ReqBody))
		defer func() {
			if e := recover(); e != any(nil) {
				ctx.Fatalf("err=%v||panic=Info[\n%s]", e, debug.Stack())
				ctx.Store(ReqStatus, "panic")
			} else {
				ctx.Store(ReqStatus, "success")
			}
			//ctx.Info(ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus))
			ctx.CounterBy("do_request", ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus)).Add(1)
			// todo code
			ctx.Store(CostMs, ctx.Duration().Milliseconds())
			ctx.SummaryBy("do_request_cost", ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus)).Observe(float64(ctx.Duration().Milliseconds()))
			ctx.Info("response_out",
				HttpMethod, ctx.LoadString(HttpMethod),
				HttpPath, ctx.LoadString(HttpPath),
				CostMs, ctx.LoadString(CostMs),
			)
		}()
		ctx.Store("xContext", ctx)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
	//x.SpanIF.SetTag(, x.Load(ext.HTTPMethod))
}

func DoRequest() gin.HandlerFunc {
	return func(gc *gin.Context) {
		xCtx := NewXContextWithContext(gc, "request")
		ctx := xCtx
		traceID := gc.Request.Header.Get("X-Request-Id")
		if traceID == "" {
			traceID = ksuid.New().String()
		}
		// todo ctx set traceID
		//gc.Set("RequestID", traceID)
		ctx.Store(ClientIP, gc.ClientIP())
		ctx.Store(HttpMethod, gc.Request.Method)
		ctx.Store(HttpPath, gc.Request.URL.Path)
		ctx.Store(UserAgent, gc.Request.UserAgent())
		var body []byte
		n, _ := gc.Request.Body.Read(body)
		bodyStr := string(body)
		ctx.Store(ReqBodyLen, n)
		ctx.Store(ReqBody, bodyStr)
		// 写回
		gc.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		ctx.Info("request_in",
			HttpMethod, ctx.LoadString(HttpMethod),
			HttpPath, ctx.LoadString(HttpPath),
			ReqBodyLen, ctx.LoadString(ReqBodyLen),
			ReqBody, ctx.LoadString(ReqBody))
		defer func() {
			if e := recover(); e != any(nil) {
				ctx.Fatalf("err=%v||panic=Info[\n%s]", e, debug.Stack())
				ctx.Store(ReqStatus, "panic")
			} else {
				ctx.Store(ReqStatus, "success")
			}
			ctx.CounterBy("do_request", ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus)).Add(1)
			// todo code
			ctx.Store(CostMs, ctx.Duration().Milliseconds())
			ctx.SummaryBy("do_request_cost", ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus)).Observe(float64(ctx.Duration().Milliseconds()))
			ctx.Info("response_out",
				HttpMethod, ctx.LoadString(HttpMethod),
				HttpPath, ctx.LoadString(HttpPath),
				CostMs, ctx.LoadString(CostMs),
			)
		}()
		gc.Set("xContext", ctx)
		gc.Next()
		ctxI, ok := gc.Get("xContext")
		ctx = ctxI.(*XContext)
		ctx.Fin()
		resp, ok := gc.Get("resp")
		if ok {
			gc.Writer.Header().Set("X-Request-Id", fmt.Sprintf("%v", ctx.TraceID()))
			gc.JSON(http.StatusOK, resp)
		} else {
			gc.JSON(http.StatusOK, map[string]interface{}{"code": -1, "code_msg": ""})
		}
	}
}

type KVMType map[string]interface{}

func (m KVMType) ToAnyArgs() []interface{} {
	var args []interface{}
	for k, v := range m {
		args = append(args, k, v)
	}
	return args
}

// extTag 常量转 metric label

func (x *XContext) TagNames2PLabels(tags ...ext.StringTagName) prometheus.Labels {
	pl := prometheus.Labels{}
	for _, tag := range tags {
		pl[string(tag)] = x.LoadString(tag)
		//x.Debug(string(tag), x.LoadString(tag))
	}
	return pl
}

// set span tag

func (x *XContext) SpanSTagSet(tags ...ext.StringTagName) {
	for _, tag := range tags {
		tag.Set(x.Span, x.LoadString(tag))
	}
}

func (x *XContext) DoRequestSpanTag() {
	x.SpanSTagSet(HttpPath, HttpMethod, ReqBody, HttpReqHeader, ReqBodyLen)
}

func (x *XContext) DoResponseSpanTag() {
	x.SpanSTagSet(RespBody, HttpRespHeader, RespCode, RespCodeMsg, ReqStatus)
}

func (x *XContext) Load(key any) any {
	return x.Context.Value(key)
}

func (x *XContext) LoadString(key any) string {
	return fmt.Sprintf("%v", x.Load(key))
}

func (x *XContext) Keys2KVMap(keys ...string) map[string]interface{} {
	kvMap := map[string]interface{}{}
	for _, key := range keys {
		val := x.LoadString(key)
		kvMap[key] = val
	}
	return kvMap
}

func (x *XContext) SetTag(key string, value interface{}) {
	x.Span.SetTag(key, value)
}

func (x *XContext) LogFields(kvs ...interface{}) { // use k1, v1 , k2, v2,
	fields, err := log.InterleavedKVToFields(kvs...)
	if err == nil {
		x.Span.LogFields(fields...)
	} else {
		x.LoggerIF.Error(err)
	}
}

//func (x XContext) LogField(key string, value interface{}) {
//	x.Span.LogFields(ToLogField(key, value))
//}

//func ToLogField(key string, value interface{}) log.Field {
//	switch reflect.TypeOf(value).Kind() {
//	case reflect.Int64:
//		return log.Int64(key, value.(int64))
//	case reflect.Int32:
//		return log.Int32(key, value.(int32))
//	case reflect.Int:
//		return log.Int(key, value.(int))
//	case reflect.String:
//		return log.String(key, value.(string))
//	case reflect.Float64:
//		return log.Float64(key, value.(float64))
//	case reflect.Float32:
//		return log.Float32(key, value.(float32))
//	case reflect.Uint64:
//		return log.Uint64(key, value.(uint64))
//	case reflect.Uint32:
//		return log.Uint32(key, value.(uint32))
//	case reflect.Bool:
//		return log.Bool(key, value.(bool))
//	default:
//		s, _ := jsoniter.MarshalToString(value)
//		return log.String(fmt.Sprintf("UnknownType[%v]", key), s)
//	}
//}

func (x *XContext) XContextPrefix() string {
	return fmt.Sprintf("[trace_id:%v][span_id:%v] ", x.TraceID(), x.SpanID())
}

func (x *XContext) ArgsFormats(args ...interface{}) string {
	var outs []string
	for _, arg := range args {
		outs = append(outs, fmt.Sprintf("%v", arg))
	}
	return strings.Join(outs, " ")
}

func (x *XContext) KVsFormats(kvs map[string]interface{}) string {
	var outs []string
	for k, v := range kvs {
		outs = append(outs, fmt.Sprintf("%v=%v", k, v))
	}
	return strings.Join(outs, "||")
}

func (x *XContext) Info(args ...interface{}) {
	x.LoggerIF.Info(append([]interface{}{x.XContextPrefix()}, x.ArgsFormats(args...))...)
}

func (x *XContext) InfoKV(kvm KVMType) {
	x.LogFields(kvm.ToAnyArgs()...)
	x.LoggerIF.Info(append([]interface{}{x.XContextPrefix()}, x.KVsFormats(kvm))...)
}

func (x *XContext) Debug(args ...interface{}) {
	x.LoggerIF.Debug(append([]interface{}{x.XContextPrefix()}, x.ArgsFormats(args...))...)
}

func (x *XContext) DebugKV(kvm KVMType) {
	x.LogFields(kvm.ToAnyArgs()...)
	x.LoggerIF.Debug(append([]interface{}{x.XContextPrefix()}, x.KVsFormats(kvm))...)
}

func (x *XContext) Warn(args ...interface{}) {
	x.LoggerIF.Warn(append([]interface{}{x.XContextPrefix()}, x.ArgsFormats(args...))...)
}

func (x *XContext) Error(args ...interface{}) {
	x.LoggerIF.Error(append([]interface{}{x.XContextPrefix()}, x.ArgsFormats(args...))...)
}

func (x *XContext) ErrorE(e error) {
	x.SetTag(string(ExecStatus), "error")
	x.LogFields(string(Error), e.Error())
	x.MetricsIF.CounterBy("err", x.TagNames2PLabels(Error)).Add(1)
	x.LoggerIF.Error(append([]interface{}{x.XContextPrefix()}, e.Error(), stack.Caller(7))...)
}

func (x *XContext) Fatal(args ...interface{}) {
	x.LoggerIF.Fatal(append([]interface{}{x.XContextPrefix()}, x.ArgsFormats(args...))...)
}

func (x *XContext) Infof(format string, args ...interface{}) {
	format = x.XContextPrefix() + format
	x.LoggerIF.Infof(format, args...)
}

func (x *XContext) Debugf(format string, args ...interface{}) {
	format = x.XContextPrefix() + format
	x.LoggerIF.Debugf(format, args...)
}

func (x *XContext) Warnf(format string, args ...interface{}) {
	format = x.XContextPrefix() + format
	x.LoggerIF.Warnf(format, args...)
}

func (x *XContext) Errorf(format string, args ...interface{}) {
	format = x.XContextPrefix() + format
	x.LoggerIF.Errorf(format, args...)
}

func (x *XContext) Fatalf(format string, args ...interface{}) {
	format = x.XContextPrefix() + format
	x.LoggerIF.Fatalf(format, args...)
}

func (x *XContext) Fin() {
	x.Span.FinishWithOptions(opentracing.FinishOptions{
		FinishTime: time.Time{},
		//LogRecords:  nil,
		//BulkLogData: nil,
	})
	for _, cancel := range x.CancelList {
		cancel()
	}
}
