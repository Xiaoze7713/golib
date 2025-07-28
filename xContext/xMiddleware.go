package xContext

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golib/v2/utils"
	"github.com/golib/v2/xContext/metrics"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
)

var reqBodySkip = map[string]bool{}
var respBodySkip = map[string]bool{}

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
	case SkipResp:
		respBodySkip[path] = true
	case SkipAll:
		reqBodySkip[path] = true
		respBodySkip[path] = true
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
	return ok
}

func RespBodySkip(path string) bool {
	lock.RLock()
	defer lock.RUnlock()
	_, ok := respBodySkip[path]
	return ok
}

type GinCtx2XCtx func(gc *gin.Context, ctx *XContext)

const (
	HeaderDevice  = "Device"
	HeaderUA      = "User-Agent"
	HeaderTraceID = "trace_id"
)

func NormalHeader(headers ...string) (gc2xc GinCtx2XCtx) {
	return func(gc *gin.Context, ctx *XContext) {
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

func DoWsConn(gc2xcList ...GinCtx2XCtx) gin.HandlerFunc {
	return func(gc *gin.Context) {
		traceID := ParseTraceID(gc.GetHeader("trace_id"))
		ctx := NewXContextWithContextString(gc, "request", traceID)
		for _, gc2xc := range gc2xcList {
			gc2xc(gc, ctx)
		}
		defer func() {
			ctx.Fin()
			ctx.SetKV(CostMs, ctx.Duration().Milliseconds())
			ctx.LogTags(ctx.Keys2KVMap(ReqStatus, CostMs))
			ctx.LogFields("resp", "resp_out")
			if RespBodySkip(gc.Request.URL.Path) {
				ctx.Info(append([]interface{}{metrics.MetricResponseOut}, ctx.ToAnyArgs(HttpMethod, HttpPath, ReqBodyLen, ReqStatus, CostMs)...))
			} else {
				ctx.Info(append([]interface{}{metrics.MetricResponseOut}, ctx.ToAnyArgs(HttpMethod, HttpPath, ReqBodyLen, ReqStatus, CostMs, RespBody)...))
			}
			//ctx.Info(ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus))
			ctx.CounterBy(metrics.MetricDoRequest, ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus)).Add(1)
			// todo code
			ctx.SummaryBy(metrics.MetricDoRequestCost, ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus)).Observe(float64(ctx.Duration().Milliseconds()))
		}()
		// gc.Set("RequestID", traceID)
		body, _ := io.ReadAll(gc.Request.Body)
		bodyStr := string(body)
		// 写回
		gc.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		requestIn := KVMType{
			ClientIP:   gc.ClientIP(),
			HttpMethod: gc.Request.Method,
			HttpPath:   gc.Request.URL.Path,
			ReqBodyLen: len(body),
			ReqBody:    bodyStr,
			UserAgent:  gc.Request.UserAgent(),
		}
		ctx.SetKVs(requestIn)
		ctx.LogTags(requestIn)
		if ReqBodySkip(gc.Request.URL.Path) {
			ctx.Info(append([]interface{}{metrics.MetricRequestIn}, ctx.ToAnyArgs(HttpMethod, HttpPath, ClientIP, ReqBodyLen)...))
		} else {
			ctx.Info(append([]interface{}{metrics.MetricRequestIn}, ctx.ToAnyArgs(HttpMethod, HttpPath, ClientIP, ReqBodyLen, ReqBody)...))
		}
		ctx.LogFields("req", "req_in")
		defer func() {
			if e := recover(); e != any(nil) {
				ctx.Fatalf("err=%v||panic=Info[\n%s]", e, debug.Stack())
				ctx.SetKV(ReqStatus, "panic")
			} else {
				ctx.SetKV(ReqStatus, "success")
			}
		}()
		gc.Set("xContext", ctx)
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

func DoRequest(gc2xcList ...GinCtx2XCtx) gin.HandlerFunc {
	return func(gc *gin.Context) {
		traceID := ParseTraceID(gc.GetHeader("trace_id"))
		ctx := NewXContextWithContextString(gc, "request", traceID)
		for _, gc2xc := range gc2xcList {
			gc2xc(gc, ctx)
		}
		defer func() {
			ctx.Fin()
			ctx.SetKV(CostMs, ctx.Duration().Milliseconds())
			ctx.LogTags(ctx.Keys2KVMap(ReqStatus, CostMs))
			ctx.LogFields("resp", "resp_out")
			if RespBodySkip(gc.Request.URL.Path) {
				ctx.Info(append([]interface{}{metrics.MetricResponseOut}, ctx.ToAnyArgs(HttpMethod, HttpPath, ReqBodyLen, ReqStatus, CostMs)...))

			} else {
				ctx.Info(append([]interface{}{metrics.MetricResponseOut}, ctx.ToAnyArgs(HttpMethod, HttpPath, ReqBodyLen, ReqStatus, CostMs, RespBody)...))
			}
			//ctx.Info(ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus))
			ctx.CounterBy(metrics.MetricDoRequest, ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus)).Add(1)
			// todo code
			ctx.SummaryBy(metrics.MetricDoRequestCost, ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus)).Observe(float64(ctx.Duration().Milliseconds()))
		}()
		// gc.Set("RequestID", traceID)
		body, _ := io.ReadAll(gc.Request.Body)
		bodyStr := string(body)
		// 写回
		gc.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		requestIn := KVMType{
			ClientIP:   gc.ClientIP(),
			HttpMethod: gc.Request.Method,
			HttpPath:   gc.Request.URL.Path,
			ReqBodyLen: len(body),
			ReqBody:    bodyStr,
			UserAgent:  gc.Request.UserAgent(),
		}
		ctx.SetKVs(requestIn)
		ctx.LogTags(requestIn)
		if ReqBodySkip(gc.Request.URL.Path) { // 20k
			ctx.Info(append([]interface{}{metrics.MetricRequestIn}, ctx.ToAnyArgs(HttpMethod, HttpPath, ClientIP, ReqBodyLen)...))
		} else {
			ctx.Info(append([]interface{}{metrics.MetricRequestIn}, ctx.ToAnyArgs(HttpMethod, HttpPath, ClientIP, ReqBodyLen, ReqBody)...))
		}
		ctx.LogFields("req", "req_in")
		defer func() {
			if e := recover(); e != any(nil) {
				ctx.Fatalf("err=%v||panic=Info[\n%s]", e, debug.Stack())
				ctx.SetKV(ReqStatus, "panic")
			} else {
				ctx.SetKV(ReqStatus, "success")
			}
		}()
		gc.Set("xContext", ctx)
		gc.Next()
		// ctxI, _ := gc.Get("xContext")
		// ctx = ctxI.(*XContext)
		gc.Writer.Header().Set("X-Request-Id", fmt.Sprintf("%v", ctx.TraceID()))
		resp, ok := gc.Get("resp")
		respContentJson := utils.MustJson(resp)
		if ok {
			ctx.SetKV(RespBody, respContentJson)
			gc.JSON(http.StatusOK, resp)
		} else {
			respBytes, ok := gc.Get("resp_bytes")
			if !ok {
				resp = map[string]interface{}{"code": -1, "code_msg": ""}
				ctx.SetKV(RespBody, respContentJson)
				gc.JSON(http.StatusOK, resp)
			} else {
				gc.Writer.Write(respBytes.([]byte))
				gc.Status(http.StatusOK)
			}
		}
	}
}

func HttpIntercept(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := ParseTraceID(r.Header.Get("trace_id"))
		ctx := NewXContextWithContextString(r.Context(), "request", traceID)
		//for _, gc2xc := range gc2xcList {
		//	gc2xc(gc, ctx)
		//}
		defer func() {
			ctx.Fin()
			ctx.SetKV(CostMs, ctx.Duration().Milliseconds())
			ctx.LogTags(ctx.Keys2KVMap(ReqStatus, CostMs))
			ctx.LogFields("resp", "resp_out")
			if RespBodySkip(r.URL.Path) {
				ctx.Info(append([]interface{}{"response_out"}, ctx.ToAnyArgs(HttpMethod, HttpPath, ReqBodyLen, ReqStatus, CostMs)...))
			} else {
				ctx.Info(append([]interface{}{"response_out"}, ctx.ToAnyArgs(HttpMethod, HttpPath, ReqBodyLen, ReqStatus, CostMs, RespBody)...))
			}
			ctx.Info(append([]interface{}{"response_out"}, ctx.ToAnyArgs(HttpMethod, HttpPath, ReqBodyLen, ReqStatus, CostMs)...))
			// ctx.Info(ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus))
			ctx.CounterBy("do_request", ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus)).Add(1)
			// todo code
			ctx.SummaryBy("do_request_cost", ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus)).Observe(float64(ctx.Duration().Milliseconds()))
		}()
		// 读body
		body, _ := ioutil.ReadAll(r.Body)
		bodyStr := string(body)
		// 写回
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		requestIn := KVMType{
			ClientIP:   r.RemoteAddr,
			HttpMethod: r.Method,
			HttpPath:   r.URL.Path,
			ReqBodyLen: len(body),
			ReqBody:    bodyStr,
		}
		ctx.SetKVs(requestIn)
		ctx.LogTags(requestIn)
		if ReqBodySkip(r.URL.Path) { // 20k
			ctx.Info(append([]interface{}{"request_in"}, ctx.ToAnyArgs(HttpMethod, HttpPath, ClientIP, ReqBodyLen)...))
		} else {
			ctx.Info(append([]interface{}{"request_in"}, ctx.ToAnyArgs(HttpMethod, HttpPath, ClientIP, ReqBodyLen, ReqBody)...))
		}
		ctx.LogFields("req", "req_in")

		// r = r.WithContext(lc)
		// request in

		defer func() {
			if e := recover(); e != any(nil) {
				ctx.Fatalf("err=%v||panic=Info[\n%s]", e, debug.Stack())
				ctx.SetKV(ReqStatus, "panic")
			} else {
				ctx.SetKV(ReqStatus, "success")
			}
		}()
		serverStartCtx := NewFollowXContext(ctx, ctx.LoadString(HttpPath))
		serverStartCtx.SetKV("xContext", serverStartCtx)
		r = r.WithContext(serverStartCtx)
		h.ServeHTTP(w, r)
		serverStartCtx.Fin()
	})
	// x.SpanIF.SetTag(, x.Load(ext.HTTPMethod))
}
