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

	"github.com/golib/v2/log"
)

type WebLogger struct {
	Logger log.Logger
	Depth  int
}

func (l *WebLogger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, l.Depth, args...)
}

func (l *WebLogger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, l.Depth, args...)
}

func (l *WebLogger) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, l.Depth, args...)
}

func (l *WebLogger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, l.Depth, args...)
}

func (l *WebLogger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, l.Depth, args...)
}

func (l *WebLogger) Debug(args ...interface{}) {
	l.Logger.Debug(l.Depth, args...)
}

func (l *WebLogger) Info(args ...interface{}) {
	l.Logger.Info(l.Depth, args...)
}

func (l *WebLogger) Warn(args ...interface{}) {
	l.Logger.Warn(l.Depth, args...)
}

func (l *WebLogger) Error(args ...interface{}) {
	l.Logger.Error(l.Depth, args...)
}

func (l *WebLogger) Fatal(args ...interface{}) {
	l.Logger.Fatal(l.Depth, args...)
}

func (ctx *WebContext) ContextPrefix() string {
	return fmt.Sprintf("trace_id[%v] span_id[%v] ", ctx.TraceID().String(), ctx.SpanID().String())
}

func (ctx *WebContext) Debugf(format string, args ...interface{}) {
	format = ctx.ContextPrefix() + format
	ctx.XContext.LoggerIF.Debugf(nil, format, args...)
}

func (ctx *WebContext) Infof(format string, args ...interface{}) {
	format = ctx.ContextPrefix() + format
	ctx.XContext.LoggerIF.Infof(nil, format, args...)
}

func (ctx *WebContext) Warnf(format string, args ...interface{}) {
	format = ctx.ContextPrefix() + format
	ctx.XContext.LoggerIF.Warnf(nil, format, args...)
}

func (ctx *WebContext) Errorf(format string, args ...interface{}) {
	format = ctx.ContextPrefix() + format
	ctx.XContext.LoggerIF.Errorf(nil, format, args...)
}

func (ctx *WebContext) Fatalf(format string, args ...interface{}) {
	format = ctx.ContextPrefix() + format
	ctx.XContext.LoggerIF.Fatalf(nil, format, args...)
}

func (ctx *WebContext) Debug(args ...interface{}) {
	ctx.XContext.LoggerIF.Debug(nil, ctx.ContextPrefix(), ctx.ArgsFormats(args...))
}

func (ctx *WebContext) Info(args ...interface{}) {
	ctx.XContext.LoggerIF.Info(nil, ctx.ContextPrefix(), ctx.ArgsFormats(args...))
}

func (ctx *WebContext) Warn(args ...interface{}) {
	ctx.XContext.LoggerIF.Warn(nil, ctx.ContextPrefix(), ctx.ArgsFormats(args...))
}

func (ctx *WebContext) Error(args ...interface{}) {
	ctx.XContext.LoggerIF.Error(nil, ctx.ContextPrefix(), ctx.ArgsFormats(args...))
}

func (ctx *WebContext) Fatal(args ...interface{}) {
	ctx.XContext.LoggerIF.Fatal(nil, ctx.ContextPrefix(), ctx.ArgsFormats(args...))

}
