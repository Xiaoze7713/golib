package gorm_xlog

import (
	"context"
	"git.singularity-ai.com/backend/library/xContext"
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

type GormXLog struct {
	Ctx *xContext.XContext
	//xLogger *xlog.Logger
	level logger.LogLevel
	trace bool
}

func NewGormXLog(ctx *xContext.XContext) (loggerIF logger.Interface) {
	return &GormXLog{Ctx: ctx}
}

func (m *GormXLog) LogMode(level logger.LogLevel) logger.Interface {
	m.level = level
	return m
}

func (m *GormXLog) SetTrace(trace bool) {
	m.trace = trace
}

func (m *GormXLog) Info(ctx context.Context, s string, args ...interface{}) {
	if m.level >= logger.Info {
		m.Ctx.Info(s, args)
	}
}

func (m *GormXLog) Warn(ctx context.Context, s string, args ...interface{}) {
	if m.level >= logger.Warn {
		m.Ctx.Info(s, args)
	}
}

func (m *GormXLog) Error(ctx context.Context, s string, args ...interface{}) {
	if m.level >= logger.Error {
		m.Ctx.Error(s, args)
	}
}

func (m *GormXLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if m.trace {
		sql, rows := fc()
		m.Ctx.Debugf("[EXEC:%s] SQL  %s, ROW %d, err %v", begin.Format("2006-01-02 15:04:05.999999999"), sql, rows, err.Error())
	}
}
