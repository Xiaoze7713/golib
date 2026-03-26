package xhttp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golib/v2/utils"
	"github.com/golib/v2/xContext"
	"github.com/golib/v2/xContext/metrics"
)

var reqBodySkip = map[string]bool{}
var respBodySkip = map[string]bool{}
var reqBodyPrefixSkip = map[string]bool{}
var respBodyPrefixSkip = map[string]bool{}

const (
	KeySkipReq  = "SkipReq"
	KeySkipResp = "SkipResp"
)

const (
	SkipReq = iota
	SkipResp
	SkipAll
)

var lock = &sync.RWMutex{}

func AddBodySkip(path string, t int) {
	lock.Lock()
	defer lock.Unlock()
	switch t {
	case SkipReq:
		reqBodySkip[path] = true
		reqBodyPrefixSkip[path] = true
	case SkipResp:
		respBodySkip[path] = true
		respBodyPrefixSkip[path] = true
	case SkipAll:
		reqBodySkip[path] = true
		reqBodyPrefixSkip[path] = true
		respBodySkip[path] = true
		respBodyPrefixSkip[path] = true
	}
}

func RemoveBodySkip(path string, t int) {
	lock.Lock()
	defer lock.Unlock()
	switch t {
	case SkipReq:
		delete(reqBodySkip, path)
	case SkipResp:
		delete(respBodySkip, path)
	case SkipAll:
		delete(reqBodySkip, path)
		delete(respBodySkip, path)
	}
}

func ReqBodySkip(path string) bool {
	lock.RLock()
	defer lock.RUnlock()
	_, ok := reqBodySkip[path]
	if ok {
		return true
	}
	for prefix, _ := range reqBodyPrefixSkip {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func RespBodySkip(path string) bool {
	lock.RLock()
	defer lock.RUnlock()
	_, ok := respBodySkip[path]
	if ok {
		return true
	}
	for prefix, _ := range respBodyPrefixSkip {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

func GcReqBodySkip(gc *gin.Context) {
	gc.Set(KeySkipReq, true)
	return
}

func GcRespBodySkip(gc *gin.Context) {
	gc.Set(KeySkipResp, true)
	return
}

type GinCtx2XCtx func(gc *gin.Context, ctx *xContext.XContext)

const (
	HeaderDevice  = "Device"
	HeaderUA      = "User-Agent"
	HeaderTraceID = "trace_id"
)

func NormalHeader(headers ...string) (gc2xc GinCtx2XCtx) {
	return func(gc *gin.Context, ctx *xContext.XContext) {
		for _, h := range headers {
			val := gc.GetHeader(h)
			ctx.SetKV(h, val)
		}
	}
}

func ParseTraceID(tid string) string {
	traceID := tid
	if len(strings.Split(traceID, ":")) == 1 {
		traceID = strings.Join([]string{traceID, "0", "0", "0"}, ":")
	}
	return traceID
}

func RecoverProcess(gc *gin.Context, ctx *xContext.XContext, e any, stack []byte) {
	if e != any(nil) {
		resp := Response{
			Code:    -2,
			Message: "no response",
			Reason:  "",
		}
		ctx.SetKV(xContext.RespBody, utils.MustJson(resp))
		gc.JSON(http.StatusOK, resp)
		ctx.Fatalf("err=%v||panic=Info[\n%s]", e, string(stack))
		ctx.SetKV(xContext.ReqStatus, "panic")
		gc.Abort()
	} else {
		ctx.SetKV(xContext.ReqStatus, "success")
	}
}

func ResponseOutProcess(gc *gin.Context, ctx *xContext.XContext) {
	ctx.Fin()
	ctx.SetKV(xContext.CostMs, ctx.Duration().Milliseconds())
	ctx.LogTags(ctx.Keys2KVMap(xContext.ReqStatus, xContext.CostMs))
	ctx.LogFields("resp", "resp_out")
	parameters := []interface{}{metrics.MetricResponseOut}
	anyArgs := []interface{}{xContext.HttpMethod, xContext.HttpPath, xContext.ReqBodyLen, xContext.ReqStatus, xContext.CostMs}
	_, ok := gc.Get(KeySkipResp)
	if !ok && !RespBodySkip(gc.Request.URL.Path) { // 20k
		anyArgs = append(anyArgs, xContext.RespBody)
	}
	parameters = append(parameters, ctx.ToAnyArgs(anyArgs...))
	ctx.Info(parameters...)
	//ctx.Info(ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus))
	ctx.CounterBy(metrics.MetricDoRequest, ctx.TagNames2PLabels(xContext.HttpMethod, xContext.HttpPath, xContext.ReqStatus)).Add(1)
	// todo code
	ctx.SummaryBy(metrics.MetricDoRequestCost, ctx.TagNames2PLabels(xContext.HttpMethod, xContext.HttpPath, xContext.ReqStatus)).Observe(float64(ctx.Duration().Milliseconds()))
}
func RequestInProcess(gc *gin.Context, ctx *xContext.XContext) {
	// gc.Set("RequestID", traceID)
	body, _ := io.ReadAll(gc.Request.Body)
	bodyStr := string(body)
	// 写回
	gc.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	requestIn := xContext.KVMType{
		xContext.ClientIP:   gc.ClientIP(),
		xContext.HttpMethod: gc.Request.Method,
		xContext.HttpPath:   gc.Request.URL.Path,
		xContext.ReqBodyLen: len(body),
		xContext.ReqBody:    bodyStr,
		xContext.UserAgent:  gc.Request.UserAgent(),
	}
	ctx.SetKVs(requestIn)
	ctx.LogTags(requestIn)
	parameters := []interface{}{metrics.MetricRequestIn}
	anyArgs := []interface{}{xContext.HttpMethod, xContext.HttpPath, xContext.ClientIP, xContext.ReqBodyLen}
	_, ok := gc.Get(KeySkipReq)
	if !ok && !ReqBodySkip(gc.Request.URL.Path) { // 20k
		anyArgs = append(anyArgs, xContext.ReqBody)
	}
	parameters = append(parameters, ctx.ToAnyArgs(anyArgs...))
	ctx.Info(parameters...)
	ctx.LogFields("req", "req_in")
}

func ProcessResponse(gc *gin.Context, ctx *xContext.XContext) {
	gc.Writer.Header().Set("X-Request-Id", fmt.Sprintf("%v", ctx.TraceID()))
	httpCode := http.StatusOK
	httpCodeAny, ok := gc.Get(HttpResponseCode)
	if ok {
		code, ok := httpCodeAny.(int)
		if ok {
			httpCode = code
		}
	}
	resp, ok := gc.Get(RespJson)
	if ok {
		respContentJson := utils.MustJson(resp)
		ctx.SetKV(xContext.RespBody, respContentJson)
		gc.JSON(httpCode, resp)
		return
	}
	respBytes, ok := gc.Get(RespBytes)
	if ok {
		gc.Writer.Write(respBytes.([]byte))
		ctx.SetKV(xContext.RespBody, string(respBytes.([]byte)))
		gc.Status(httpCode)
		return
	}
	respStream, ok := gc.Get(RespStream)
	if ok {
		respStreamList, ok1 := respStream.([]string)
		if ok1 {
			ctx.SetKV(xContext.RespBody, utils.MustJson(respStreamList))
		}
		gc.Status(httpCode)
		return
	}
	resp = Response{
		Code:    -1,
		Message: "no response",
		Reason:  "",
	}
	ctx.SetKV(xContext.RespBody, utils.MustJson(resp))
	gc.JSON(http.StatusOK, resp)
}

func DoWsConn(gc2xcList ...GinCtx2XCtx) gin.HandlerFunc {
	return func(gc *gin.Context) {
		traceID := ParseTraceID(gc.GetHeader("trace_id"))
		ctx := xContext.NewXContextWithContextString(gc, "request", traceID)
		for _, gc2xc := range gc2xcList {
			gc2xc(gc, ctx)
		}
		defer ResponseOutProcess(gc, ctx)
		// gc.Set("RequestID", traceID)
		RequestInProcess(gc, ctx)
		defer func() {
			e := recover()
			stack := debug.Stack()
			RecoverProcess(gc, ctx, e, stack)
		}()
		gc.Set(xContext.XContextKey.String(), ctx)
		gc.Next()
		// ctxI, _ := gc.Get("xContext")
		// ctx = ctxI.(*XContext)
		//resp, ok := gc.Get("resp")
		//if ok {
		//	gc.Writer.Header().Set("X-Request-Id", fmt.Sprintf("%v", ctx.TraceID()))
		//	ctx.SetKV(RespBody, utils.MustJson(resp))
		//	gc.JSON(http.StatusOK, resp)
		//} else {
		//	resp = map[string]interface{}{"code": -1, "code_msg": ""}
		//	gc.JSON(http.StatusOK, resp)
		//	ctx.SetKV(RespBody, utils.MustJson(resp))
		//}
	}
}

// DoRequest 包裹所有请求
func DoRequest(gc2xcList ...GinCtx2XCtx) gin.HandlerFunc {
	return func(gc *gin.Context) {
		traceID := ParseTraceID(gc.GetHeader("trace_id"))
		ctx := xContext.NewXContextWithContextString(gc, "request", traceID)
		for _, gc2xc := range gc2xcList {
			gc2xc(gc, ctx)
		}
		defer ResponseOutProcess(gc, ctx)
		RequestInProcess(gc, ctx)
		// recover
		defer func() {
			e := recover()
			stack := debug.Stack()
			RecoverProcess(gc, ctx, e, stack)
		}()
		gc.Set("xContext", ctx)
		gc.Next()
		// ctxI, _ := gc.Get("xContext")
		// ctx = ctxI.(*XContext)
		ProcessResponse(gc, ctx)
	}
}

func HttpIntercept(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := ParseTraceID(r.Header.Get("trace_id"))
		ctx := xContext.NewXContextWithContextString(r.Context(), "request", traceID)
		//for _, gc2xc := range gc2xcList {
		//	gc2xc(gc, ctx)
		//}
		defer func() {
			ctx.Fin()
			ctx.SetKV(xContext.CostMs, ctx.Duration().Milliseconds())
			ctx.LogTags(ctx.Keys2KVMap(xContext.ReqStatus, xContext.CostMs))
			ctx.LogFields("resp", "resp_out")
			if RespBodySkip(r.URL.Path) {
				ctx.Info(append([]interface{}{"response_out"}, ctx.ToAnyArgs(xContext.HttpMethod, xContext.HttpPath, xContext.ReqBodyLen, xContext.ReqStatus, xContext.CostMs)...))
			} else {
				ctx.Info(append([]interface{}{"response_out"}, ctx.ToAnyArgs(xContext.HttpMethod, xContext.HttpPath, xContext.ReqBodyLen, xContext.ReqStatus, xContext.CostMs, xContext.RespBody)...))
			}
			ctx.Info(append([]interface{}{"response_out"}, ctx.ToAnyArgs(xContext.HttpMethod, xContext.HttpPath, xContext.ReqBodyLen, xContext.ReqStatus, xContext.CostMs)...))
			// ctx.Info(ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus))
			ctx.CounterBy("do_request", ctx.TagNames2PLabels(xContext.HttpMethod, xContext.HttpPath, xContext.ReqStatus)).Add(1)
			// todo code
			ctx.SummaryBy("do_request_cost", ctx.TagNames2PLabels(xContext.HttpMethod, xContext.HttpPath, xContext.ReqStatus)).Observe(float64(ctx.Duration().Milliseconds()))
		}()
		// 读body
		body, _ := io.ReadAll(r.Body)
		bodyStr := string(body)
		// 写回
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		requestIn := xContext.KVMType{
			xContext.ClientIP:   r.RemoteAddr,
			xContext.HttpMethod: r.Method,
			xContext.HttpPath:   r.URL.Path,
			xContext.ReqBodyLen: len(body),
			xContext.ReqBody:    bodyStr,
		}
		ctx.SetKVs(requestIn)
		ctx.LogTags(requestIn)
		if ReqBodySkip(r.URL.Path) { // 20k
			ctx.Info(append([]interface{}{"request_in"}, ctx.ToAnyArgs(xContext.HttpMethod, xContext.HttpPath, xContext.ClientIP, xContext.ReqBodyLen)...))
		} else {
			ctx.Info(append([]interface{}{"request_in"}, ctx.ToAnyArgs(xContext.HttpMethod, xContext.HttpPath, xContext.ClientIP, xContext.ReqBodyLen, xContext.ReqBody)...))
		}
		ctx.LogFields("req", "req_in")

		// r = r.WithContext(lc)
		// request in

		defer func() {
			if e := recover(); e != any(nil) {
				ctx.Fatalf("err=%v||panic=Info[\n%s]", e, debug.Stack())
				ctx.SetKV(xContext.ReqStatus, "panic")
			} else {
				ctx.SetKV(xContext.ReqStatus, "success")
			}
		}()
		serverStartCtx := xContext.NewFollowXContext(ctx.Span.Context(), ctx.LoadString(xContext.HttpPath))
		serverStartCtx.SetKV("xContext", serverStartCtx)
		r = r.WithContext(serverStartCtx)
		h.ServeHTTP(w, r)
		serverStartCtx.Fin()
	})
	// x.SpanIF.SetTag(, x.Load(ext.HTTPMethod))
}
