package gorm_xlog

import (
	"context"
	"github.com/golib/v2/xContext"
	"gorm.io/gorm/logger"
	"time"
)

type Interface interface {
	LogMode(logger.LogLevel) Interface
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

const TraceLevel = logger.LogLevel(5)

type GormXLog struct {
	Ctx *xContext.XContext
	//xLogger *xlog.Logger
	level logger.LogLevel
}

func NewGormXLog(ctx *xContext.XContext) (loggerIF logger.Interface) {
	return &GormXLog{Ctx: ctx}
}

func (m *GormXLog) LogMode(level logger.LogLevel) logger.Interface {
	m.level = level
	return m
}

func (m *GormXLog) Info(ctx context.Context, s string, args ...interface{}) {
	if m.level >= logger.Info {
		if xCtx, ok := ctx.(*xContext.XContext); ok {
			xCtx.Info(s, args)
		} else {
			m.Ctx.Info(m.Ctx.OperationName(), s, args)
		}
	}
}

func (m *GormXLog) Warn(ctx context.Context, s string, args ...interface{}) {
	if m.level >= logger.Warn {
		if xCtx, ok := ctx.(*xContext.XContext); ok {
			xCtx.Warn(s, args)
		} else {
			m.Ctx.Warn(m.Ctx.OperationName(), s, args)
		}
	}
}

func (m *GormXLog) Error(ctx context.Context, s string, args ...interface{}) {
	if m.level >= logger.Error {
		if xCtx, ok := ctx.(*xContext.XContext); ok {
			xCtx.Error(s, args)
		} else {
			m.Ctx.Error(m.Ctx.OperationName(), s, args)
		}
	}
}

func (m *GormXLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if m.level >= TraceLevel {
		sql, rows := fc()
		if xCtx, ok := ctx.(*xContext.XContext); ok {
			xCtx.Debugf("[%s][T:%s][%s]ROWS %d, RES %v", xCtx.OperationName(), begin.Format("2006-01-02 15:04:05.999"), sql, rows, err)
		} else {
			m.Ctx.Debugf("[%s][T:%s][%s]ROWS %d, RES %v", m.Ctx.OperationName(), begin.Format("2006-01-02 15:04:05.999"), sql, rows, err)
		}
	}
}
