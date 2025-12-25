package xhttp

import (
	"github.com/gin-gonic/gin"
	"github.com/golib/v2/xContext"
)

func NoRoute(gc *gin.Context) {
	ctx, _ := xContext.GetXContextFromCtxValue(gc, gc.Request.RequestURI)
	// defer ctx.Fin()
	ctx.SummaryBy("noRoute", ctx.TagNames2PLabels(xContext.HttpMethod, xContext.HttpPath, xContext.ReqStatus))
	SetErrorResponse(gc, -1, "path no route")
	gc.Abort()
	gc.Next()
}

func NoMethod(gc *gin.Context) {
	ctx, _ := xContext.GetXContextFromCtxValue(gc, gc.Request.RequestURI)
	// defer ctx.Fin()
	ctx.SummaryBy("noMethod", ctx.TagNames2PLabels(xContext.HttpMethod, xContext.HttpPath, xContext.ReqStatus))
	SetErrorResponse(gc, -1, "function not found")
	gc.Abort()
	gc.Next()
}

func Health(gc *gin.Context) {
	ctx, _ := xContext.GetXContextFromCtxValue(gc, gc.Request.RequestURI)
	// defer ctx.Fin()
	ctx.SummaryBy("health_check", ctx.TagNames2PLabels(xContext.HttpMethod, xContext.HttpPath, xContext.ReqStatus))
	SetErrorResponse(gc, 200, "i am healthy~")
	gc.Abort()
	gc.Next()
}

const (
	RespJson         = "resp"
	RespBytes        = "resp_bytes"
	RespStream       = "resp_stream"
	HttpResponseCode = "http_response_code"
)
