package xContext

import (
	"bytes"
	"fmt"
	"git.singularity-ai.com/backend/library/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"runtime/debug"
)

func DoWsConn() gin.HandlerFunc {
	return func(gc *gin.Context) {
		traceID := gc.GetHeader("trace_id")
		ctx := NewXContextWithContextString(gc, "request", traceID)
		defer func() {
			ctx.Fin()
			ctx.SetKV(CostMs, ctx.Duration().Milliseconds())
			ctx.LogTags(ctx.Keys2KVMap(ReqStatus, CostMs))
			ctx.LogFields("resp", "resp_out")
			ctx.Info(append([]interface{}{"response_out"}, ctx.ToAnyArgs(HttpMethod, HttpPath, ReqBodyLen, ReqStatus, CostMs, RespBody)...))
			//ctx.Info(ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus))
			ctx.CounterBy("do_request", ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus)).Add(1)
			// todo code
			ctx.SummaryBy("do_request_cost", ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus)).Observe(float64(ctx.Duration().Milliseconds()))
		}()
		// gc.Set("RequestID", traceID)
		body, _ := ioutil.ReadAll(gc.Request.Body)
		bodyStr := string(body)
		// 写回
		gc.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
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
		ctx.Info(append([]interface{}{"request_in"}, ctx.ToAnyArgs(HttpMethod, HttpPath, ClientIP, ReqBodyLen, ReqBody)...))
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
func DoRequest() gin.HandlerFunc {
	return func(gc *gin.Context) {
		traceID := gc.GetHeader("trace_id")
		ctx := NewXContextWithContextString(gc, "request", traceID)
		defer func() {
			ctx.Fin()
			ctx.SetKV(CostMs, ctx.Duration().Milliseconds())
			ctx.LogTags(ctx.Keys2KVMap(ReqStatus, CostMs))
			ctx.LogFields("resp", "resp_out")
			ctx.Info(append([]interface{}{"response_out"}, ctx.ToAnyArgs(HttpMethod, HttpPath, ReqBodyLen, ReqStatus, CostMs, RespBody)...))
			//ctx.Info(ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus))
			ctx.CounterBy("do_request", ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus)).Add(1)
			// todo code
			ctx.SummaryBy("do_request_cost", ctx.TagNames2PLabels(HttpMethod, HttpPath, ReqStatus)).Observe(float64(ctx.Duration().Milliseconds()))
		}()
		// gc.Set("RequestID", traceID)
		body, _ := ioutil.ReadAll(gc.Request.Body)
		bodyStr := string(body)
		// 写回
		gc.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
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
		ctx.Info(append([]interface{}{"request_in"}, ctx.ToAnyArgs(HttpMethod, HttpPath, ClientIP, ReqBodyLen, ReqBody)...))
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
		resp, ok := gc.Get("resp")
		if ok {
			gc.Writer.Header().Set("X-Request-Id", fmt.Sprintf("%v", ctx.TraceID()))
			ctx.SetKV(RespBody, utils.MustJson(resp))
			gc.JSON(http.StatusOK, resp)
		} else {
			resp = map[string]interface{}{"code": -1, "code_msg": ""}
			gc.JSON(http.StatusOK, resp)
			ctx.SetKV(RespBody, utils.MustJson(resp))
		}
	}
}

func HttpIntercept(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("trace_id")
		ctx := NewXContextWithContextString(r.Context(), "request", traceID)
		defer func() {
			ctx.Fin()
			ctx.SetKV(CostMs, ctx.Duration().Milliseconds())
			ctx.LogTags(ctx.Keys2KVMap(ReqStatus, CostMs))
			ctx.LogFields("resp", "resp_out")
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
		ctx.Info(append([]interface{}{"request_in"}, ctx.ToAnyArgs(HttpMethod, HttpPath, ClientIP, ReqBodyLen, ReqBody)...))
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
