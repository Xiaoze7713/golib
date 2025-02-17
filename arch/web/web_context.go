/**
 * @Author: wenliangzhang
 * @Description:
 * @File: web_server
 * @Version: 1.0.0
 * @Date: 2022/4/12 3:35 PM
 */
package web

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/golib/v3/xContext"
)

// WebHandlerFunc http请求的处理者
type WebHandlerFunc func(*WebContext)

// WebContext http 的context
// WebContext 继承了 xContext,包含gin.Context， 并且扩展了日志、监控、trace功能
type WebContext struct {
	*gin.Context
	*xContext.XContext
}

// NewWebContext 创建 http context
func NewWebContextWithGinCtx(ginContext *gin.Context) *WebContext {
	var xctx *xContext.XContext
	traceID := getTraceIDFromRequest(ginContext)
	if traceID != "" {
		xctx = xContext.NewXContextWithContextString(ginContext, ginContext.Request.URL.Path, traceID)
	} else {
		xctx = xContext.NewXContextWithContext(ginContext, ginContext.Request.URL.Path)
		ginContext.Set("trace_id", xctx.SerializeSpanContext())
	}
	return &WebContext{
		ginContext,
		xctx,
	}
}

func NewWebContextWithTraceID(name, traceID string) *WebContext {
	xctx := xContext.NewXContextWithContextString(context.Background(), name, traceID)
	return &WebContext{
		&gin.Context{},
		xctx,
	}
}

func NewWebContext(name string) *WebContext {
	xctx := xContext.NewXContext(name)
	return &WebContext{
		&gin.Context{},
		xctx,
	}
}

// getTraceIDFromRequest 获取traceID
// 优先级 header中X_TRACE_ID > header中trace_id > 表单中trace_id
func getTraceIDFromRequest(c *gin.Context) (traceID string) {
	if trace, exists := c.Get("trace_id"); exists {
		if traceIDStr, ok := trace.(string); ok {
			return traceIDStr
		}
	}
	request := c.Request
	// 上下游header透传时
	if traceID = strings.TrimSpace(request.Header.Get("X_TRACE_ID")); traceID != "" {
		return
	}
	if traceID = strings.TrimSpace(request.Header.Get("Trace_id")); traceID != "" {
		return
	}
	// 上下游querystring透传时
	form := request.URL.Query()
	if traceID = strings.TrimSpace(form.Get("trace_id")); traceID != "" {
		return
	}
	return
}
