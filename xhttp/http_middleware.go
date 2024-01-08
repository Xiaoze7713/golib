package xhttp

import (
	"github.com/Xiaoze7713/golib/v3/xContext"
	"github.com/gin-gonic/gin"
)

func NoRoute(gc *gin.Context) {
	ctx, _ := xContext.GetXContextFromGrandFather(gc, gc.Request.RequestURI)
	// defer ctx.Fin()
	ctx.SummaryBy("noRoute", ctx.TagNames2PLabels(xContext.HttpMethod, xContext.HttpPath, xContext.ReqStatus))
	SetErrorResponse(gc, -1, "path no route")
	gc.Abort()
	gc.Next()
}

func NoMethod(gc *gin.Context) {
	ctx, _ := xContext.GetXContextFromGrandFather(gc, gc.Request.RequestURI)
	// defer ctx.Fin()
	ctx.SummaryBy("noMethod", ctx.TagNames2PLabels(xContext.HttpMethod, xContext.HttpPath, xContext.ReqStatus))
	SetErrorResponse(gc, -1, "function not found")
	gc.Abort()
	gc.Next()
}

func Health(gc *gin.Context) {
	ctx, _ := xContext.GetXContextFromGrandFather(gc, gc.Request.RequestURI)
	// defer ctx.Fin()
	ctx.SummaryBy("health_check", ctx.TagNames2PLabels(xContext.HttpMethod, xContext.HttpPath, xContext.ReqStatus))
	SetErrorResponse(gc, 200, "i am healthy~")
	gc.Abort()
	gc.Next()
}
