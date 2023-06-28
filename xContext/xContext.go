package xContext

import (
	"context"
	"errors"
	"fmt"
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
	"github.com/opentracing/opentracing-go"
	"strings"
	"time"

	"git.singularity-ai.com/backend/library/v2/utils"
	"git.singularity-ai.com/backend/library/v2/xContext/base_if/xlog_base"
	"git.singularity-ai.com/backend/library/v2/xContext/base_if/xmetric_base"
	"git.singularity-ai.com/backend/library/v2/xContext/base_if/xspan_base"
	"git.singularity-ai.com/backend/library/v2/xContext/base_if/xtrace_base"
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/null_log"
	"git.singularity-ai.com/backend/library/v2/xContext/metrics/null_metric"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	HttpAction     = ext.StringTagName("action")
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
	Alert          = ext.StringTagName("alert")
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

// 监控 链路跟踪 log
var mMetric xmetric_base.MetricsIF
var mTrace xtrace_base.TracerIF
var mLogger xlog_base.LoggerIF

type GetIDType func(xContext *XContext) IDType
type GetDurType func(xContext *XContext) time.Duration
type GetErrorCodeByErr func(err error) int64
type ExtractSpanFunc func(operationName, spanStr string, tracer opentracing.Tracer) (opentracing.Span, error)
type SerializeToString func(span opentracing.Span) string

type IDType interface {
	String() string
	Value() int64
}

type Trace struct {
	TraceID string `json:"trace_id"`
}

type NullIDType int64

func (m NullIDType) String() string {
	return fmt.Sprintf("%x", m.Value())
}

func (m NullIDType) Value() int64 {
	return int64(m)
}

// 几个id的获取
var spanIDFunc GetIDType
var traceIDFunc GetIDType
var durFunc GetDurType
var toErrCode GetErrorCodeByErr
var extractFunc ExtractSpanFunc
var serializeFunc SerializeToString

func emptyIDFunc(x *XContext) IDType {
	return NullIDType(0)
}
func emptyDurFunc(x *XContext) time.Duration {
	return 0
}

func (x *XContext) ErrorAlert(e error) {
	x.SetTag(string(ExecStatus), "error")
	x.LogFields(string(Error), e.Error())
	x.MetricsIF.CounterBy("err", x.TagNames2PLabels(Error)).Add(1)
}

func (x *XContext) Alert(i interface{}) {
	x.LogFields(string(Alert), utils.MustJson(i))
	x.MetricsIF.CounterBy("alert", x.TagNames2PLabels(Alert)).Add(1)
}

// type ContextOptions interface {
//	Apply() error
// }
//
// type NullTraceID struct {
//	ContextOptions
// }
//
// func (NullTraceID) Apply(strFunc GetIDType) {
//	if strFunc != nil {
//		traceIDFunc = strFunc
//		return
//	}
//	traceIDFunc = func(xContext *XContext) string {
//		return ""
//	}
// }
//
// type NullSpanID struct {
//	ContextOptions
// }
//
// func (NullSpanID) Apply() {
//	spanIDFunc = func(xContext *XContext) string {
//		return ""
//	}
// }
//
// type NullDur struct {
//	ContextOptions
// }
//
// func (NullDur) Apply() {
//	durFunc = func(xContext *XContext) time.Duration {
//		return -1
//	}
// }
//

func InitByNullOpt() (err error) {
	xlog.SetupLogDefault()
	//serverName := "x_content_test"
	metrics := null_metric.MetricsNull{}
	trace := xtrace_base.XTraceNoop{}
	//trace, _, err := jaeger_trace.NewJaegerTrace(&jaeger_trace.JaegerConfig{
	//	ServiceName:   serverName,
	//	Param:         0,
	//	AgentHostPort: "127.0.0.1:6832",
	//})
	//if err != nil {
	//	return err
	//}
	Init(xlog.GetLogger(), trace, metrics,
		nil, nil, nil, nil, nil)
	return err
}

func InitByNullOptNoLog() (err error) {
	//serverName := "x_content_test"
	metrics := null_metric.MetricsNull{}
	trace := xtrace_base.XTraceNoop{}
	//trace, _, err := jaeger_trace.NewJaegerTrace(&jaeger_trace.JaegerConfig{
	//	ServiceName:   serverName,
	//	Param:         0,
	//	AgentHostPort: "127.0.0.1:6832",
	//})
	//if err != nil {
	//	return err
	//}
	Init(xlog.GetLogger(), trace, metrics,
		nil, nil, nil, nil, nil)
	return err
}

func Init(logger xlog_base.LoggerIF,
	trace xtrace_base.TracerIF,
	metrics xmetric_base.MetricsIF,
	traceF, spanF GetIDType, durF GetDurType, extractF ExtractSpanFunc, serializeF SerializeToString) {
	if logger != nil {
		mLogger = logger
	} else {
		fmt.Println("warning! not set logger")
		mLogger = null_log.LoggerNull{}
	}
	if trace != nil {
		mTrace = trace
	} else {
		fmt.Println("warning! not set trace")
		mTrace = xtrace_base.XTraceNoop{}
	}
	if metrics != nil {
		mMetric = metrics
	} else {
		fmt.Println("warning! not set metric")
		mMetric = null_metric.MetricsNull{}
	}
	traceIDFunc = traceF
	spanIDFunc = spanF
	durFunc = durF
	extractFunc = extractF
	serializeFunc = serializeF
}

type KVMType map[any]interface{}

func (m KVMType) ToAnyArgs() []interface{} {
	var args []interface{}
	for k, v := range m {
		args = append(args, k, v)
	}
	return args
}

func (m KVMType) ToAnyArgsWithValues(args ...interface{}) []interface{} {
	for k, v := range m {
		args = append(args, k, v)
	}
	return args
}

func (x *XContext) ToAnyArgs(keys ...any) []interface{} {
	var args []interface{}
	for _, k := range keys {
		args = append(args, k, x.Value(k))
	}
	return args
}

type XContext struct {
	context.Context
	xlog_base.LoggerIF
	xmetric_base.MetricsIF
	Span          xspan_base.SpanIF
	CancelList    []context.CancelFunc
	operationName string
}

func (x *XContext) OperationName() string {
	return x.operationName
}

//
// type GXContext struct {
//	gin.Context
//	xlog_base.LoggerIF
//	xmetric_base.MetricsIF
//	opentracing.Span
//	CancelList []context.CancelFunc
// }

var notGrandFather = errors.New("not grand father")

func GetXContextFromGrandFather(ctx context.Context, operationName string) (*XContext, error) {
	gCtx := ctx.Value("xContext")
	xCtx, ok := gCtx.(*XContext)
	if !ok {
		mLogger.Error(nil, notGrandFather)
		xCtx = NewXContextWithContext(ctx, operationName)
		return xCtx, nil
	}
	return xCtx.CopyWithContext(ctx), nil
}

func (x *XContext) CopyWithContext(ctx context.Context) *XContext {
	return &XContext{
		Context:       ctx,
		LoggerIF:      x.LoggerIF,
		MetricsIF:     x.MetricsIF,
		Span:          x.Span,
		CancelList:    nil,
		operationName: x.operationName,
	}
}

func NewXContextWithContextString(ctx context.Context, operationName, spanCtxString string) *XContext {
	xCtx := &XContext{
		Context:       ctx,
		LoggerIF:      mLogger,
		MetricsIF:     mMetric,
		operationName: operationName,
		// lock:      &sync.RWMutex{},
	}
	xCtx.Span = newSpanWithString(spanCtxString, operationName)
	return xCtx
}

func NewXContextWithContext(ctx context.Context, operationName string) *XContext {
	xCtx := &XContext{
		Context:       ctx,
		LoggerIF:      mLogger,
		MetricsIF:     mMetric,
		operationName: operationName,
		// lock:      &sync.RWMutex{},
	}
	xCtx.Span = mTrace.StartSpan(operationName,
		opentracing.StartTime{},
	)
	return xCtx
}

func NewXContext(operationName string) *XContext {
	xCtx := &XContext{
		Context:       context.Background(),
		LoggerIF:      mLogger,
		MetricsIF:     mMetric,
		operationName: operationName,
		// lock:      &sync.RWMutex{},
	}
	xCtx.Span = mTrace.StartSpan(operationName,
		opentracing.StartTime{},
	)
	return xCtx
}

func NewChildXContext(parents *XContext, operationName string) *XContext {
	childCtx, cancel := context.WithCancel(parents.Context)
	parents.CancelList = append(parents.CancelList, cancel)
	child := &XContext{
		Context:    childCtx,
		LoggerIF:   mLogger,
		MetricsIF:  mMetric,
		CancelList: []context.CancelFunc{},
		// lock:       parents.lock,
	}
	child.Span = mTrace.StartSpan(operationName,
		opentracing.ChildOf(parents.Span.Context()),
		opentracing.StartTime{},
	)
	return child
}

func NewFollowXContext(origin *XContext, operationName string) *XContext {
	brother := &XContext{
		Context:    context.Background(),
		LoggerIF:   mLogger,
		MetricsIF:  mMetric,
		CancelList: []context.CancelFunc{},
		// lock:       &sync.RWMutex{},
	}
	brother.Span = mTrace.StartSpan(operationName,
		opentracing.FollowsFrom(origin.Span.Context()),
		opentracing.StartTime{},
	)
	return brother
}

func (x *XContext) SpanID() IDType {
	if spanIDFunc == nil {
		return emptyIDFunc(x)
	} else {
		return spanIDFunc(x)
	}
}

func (x *XContext) TraceID() IDType {
	// return x.Load(TraceID)
	if traceIDFunc == nil {
		return emptyIDFunc(x)
	}
	return traceIDFunc(x)
}

func (x *XContext) Duration() time.Duration {
	if durFunc == nil {
		return emptyDurFunc(x)
	}
	return durFunc(x)
}

func (x *XContext) ErrCode(err error) int64 {
	if toErrCode != nil {
		return toErrCode(err)
	}
	return -1
}

// 获取kv

func GrpcGWIntercept() {

}

func (x *XContext) SetKVs(kvs KVMType) {
	for k, v := range kvs {
		x.SetKV(k, v)
	}
}

type ContextKey interface {
	string | ext.StringTagName | any
}

func (x *XContext) SetKV(key any, value interface{}) {
	x.Context = context.WithValue(x.Context, key, value)
}

// extTag 常量转 metric label

func (x *XContext) TagNames2PLabels(tags ...ext.StringTagName) prometheus.Labels {
	pl := prometheus.Labels{}
	for _, tag := range tags {
		pl[string(tag)] = x.LoadString(tag)
		// x.Debug(string(tag), x.LoadString(tag))
	}
	return pl
}

// set span tag
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

func (x *XContext) Keys2KVMap(keys ...any) KVMType {
	kvMap := KVMType{}
	for _, key := range keys {
		val := x.LoadString(key)
		kvMap[key] = val
	}
	return kvMap
}

func (x *XContext) LogTags(kvm KVMType) {
	for k, v := range kvm {
		extStr, ok := k.(ext.StringTagName)
		if !ok {
			s := fmt.Sprintf("%v", k)
			x.SetTag(s, v)
			continue
		}
		x.SpanSTagSet(extStr)
	}
}

func (x *XContext) SpanSTagSet(tags ...ext.StringTagName) {
	for _, tag := range tags {
		tag.Set(x.Span, x.LoadString(tag))
	}
}

func (x *XContext) SetTag(key string, value interface{}) {
	x.Span.SetTag(key, value)
}

func (x *XContext) LogFields(kvs ...interface{}) { // use k1, v1 , k2, v2,
	fields, err := log.InterleavedKVToFields(kvs...)
	if err == nil {
		x.Span.LogFields(fields...)
	}
}

func (x *XContext) XContextPrefix() string {
	return fmt.Sprintf("[trace_id:%v][span_id:%v] ", x.TraceID().String(), x.SpanID().String())
}

func (x *XContext) ArgsFormats(args ...interface{}) string {
	var outs []string
	for _, arg := range args {
		outs = append(outs, fmt.Sprintf("%v", arg))
	}
	return strings.Join(outs, " ")
}

func (x *XContext) KVsFormats(kvs KVMType) string {
	var outs []string
	for k, v := range kvs {
		outs = append(outs, fmt.Sprintf("%v=%v", k, v))
	}
	return strings.Join(outs, "||")
}
func (x *XContext) Fin() {
	x.Span.FinishWithOptions(opentracing.FinishOptions{
		FinishTime: time.Time{},
		// LogRecords:  nil,
		// BulkLogData: nil,
		//LogRecords:  nil,
		//BulkLogData: nil,
	})
	for _, cancel := range x.CancelList {
		cancel()
	}
}

func newSpanWithString(ctxStr, operationName string) opentracing.Span {
	if extractFunc == nil {
		return mTrace.StartSpan(operationName)
	}
	span, err := extractFunc(operationName, ctxStr, mTrace)
	if err != nil {
		return mTrace.StartSpan(operationName)
	}
	return span
}

func (m *XContext) SerializeSpanContext() string {
	if serializeFunc == nil {
		return ""
	}
	return serializeFunc(m.Span)
}
