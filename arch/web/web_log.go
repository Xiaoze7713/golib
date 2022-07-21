/**
 * @Author: wenliangzhang
 * @Description:
 * @File: log
 * @Version: 1.0.0
 * @Date: 2022/7/18 6:28 PM
 */
package web

import (
	"fmt"
)

func (ctx *WebContext) ContextPrefix() string {
	return fmt.Sprintf("trace_id[%v] span_id[%v] ", ctx.TraceID().String(), ctx.SpanID().String())
}

func (ctx *WebContext) Debugf(format string, args ...interface{}) {
	format = ctx.ContextPrefix() + format
	ctx.XContext.LoggerIF.Debugf(format, args...)
}

func (ctx *WebContext) Infof(format string, args ...interface{}) {
	format = ctx.ContextPrefix() + format
	ctx.XContext.LoggerIF.Infof(format, args...)
}

func (ctx *WebContext) Warnf(format string, args ...interface{}) {
	format = ctx.ContextPrefix() + format
	ctx.XContext.LoggerIF.Warnf(format, args...)
}

func (ctx *WebContext) Errorf(format string, args ...interface{}) {
	format = ctx.ContextPrefix() + format
	ctx.XContext.LoggerIF.Errorf(format, args...)
}

func (ctx *WebContext) Fatalf(format string, args ...interface{}) {
	format = ctx.ContextPrefix() + format
	ctx.XContext.LoggerIF.Fatalf(format, args...)
}

func (ctx *WebContext) Debug(args ...interface{}) {
	ctx.XContext.LoggerIF.Debug(ctx.ContextPrefix(), ctx.ArgsFormats(args...))
}

func (ctx *WebContext) Info(args ...interface{}) {
	ctx.XContext.LoggerIF.Info(ctx.ContextPrefix(), ctx.ArgsFormats(args...))
}

func (ctx *WebContext) Warn(args ...interface{}) {
	ctx.XContext.LoggerIF.Warn(ctx.ContextPrefix(), ctx.ArgsFormats(args...))
}

func (ctx *WebContext) Error(args ...interface{}) {
	ctx.XContext.LoggerIF.Error(ctx.ContextPrefix(), ctx.ArgsFormats(args...))
}

func (ctx *WebContext) Fatal(args ...interface{}) {
	ctx.XContext.LoggerIF.Fatal(ctx.ContextPrefix(), ctx.ArgsFormats(args...))

}
