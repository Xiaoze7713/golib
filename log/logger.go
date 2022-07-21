/**
 * @Author: wenliangzhang
 * @Description:
 * @File: logger
 * @Version: 1.0.0
 * @Date: 2022/4/21 11:12 AM
 */
package log

import (
	"github.com/sirupsen/logrus"
)

// 文件后缀分类 用于区分普通日志、wf错误日志、标准输出
const (
	TxtLogNormal    = "normal"
	TxtLogWarnfatal = "warnfatal"
	TxtLogStdout    = "stdout"
)

type Logger map[string]*LogWriter

var globalLogger Logger

func initGlobalLogger(conf LogConfig) {
	if globalLogger == nil {
		globalLogger = make(Logger)
	}
	globalLogger[TxtLogStdout] = NewStdWriter(conf)
	globalLogger[TxtLogNormal] = NewWriter(conf, "")
	globalLogger[TxtLogWarnfatal] = NewWriter(conf, ".wf")
}

func filterWriter(level logrus.Level) string {
	if level <= logrus.WarnLevel {
		return TxtLogWarnfatal
	}
	return TxtLogNormal
}

func buildWriterNames(level logrus.Level) []string {
	names := []string{
		filterWriter(level),
	}
	if globalConfig.Stdout {
		names = append(names, TxtLogStdout)
	}
	return names
}

func (l Logger) WriteLogf(level logrus.Level, format string, args ...interface{}) {
	names := buildWriterNames(level)
	for _, name := range names {
		if format != "" {
			l[name].l.Logf(level, format, args...)
		} else {
			l[name].l.Log(level, args...)
		}
	}
}

func (l Logger) WriteLogln(level logrus.Level, args ...interface{}) {
	names := buildWriterNames(level)
	for _, name := range names {
		l[name].l.Logln(level, args...)
	}
}

func (l Logger) WritePrintf(format string, args ...interface{}) {
	names := buildWriterNames(logrus.InfoLevel)
	for _, name := range names {
		if format != "" {
			l[name].l.Printf(format, args...)
		} else {
			l[name].l.Print(args...)
		}
	}
}

func (l Logger) WritePrintln(args ...interface{}) {
	names := buildWriterNames(logrus.InfoLevel)
	for _, name := range names {
		l[name].l.Println(args...)
	}
}

func (l Logger) Tracef(format string, args ...interface{}) {
	l.WriteLogf(logrus.TraceLevel, format, args...)
}

func (l Logger) Debugf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.DebugLevel, format, args...)
}

func (l Logger) Infof(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.InfoLevel, format, args...)
}

func (l Logger) Printf(format string, args ...interface{}) {
	globalLogger.WritePrintf(format, args...)
}

func (l Logger) Warnf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.WarnLevel, format, args...)
}

func (l Logger) Warningf(format string, args ...interface{}) {
	Warnf(format, args...)
}

func (l Logger) Errorf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.ErrorLevel, format, args...)
}

func (l Logger) Fatalf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.FatalLevel, format, args...)
	// logger.Exit(1)
}

func (l Logger) Panicf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.PanicLevel, format, args...)
}

func (l Logger) Trace(args ...interface{}) {
	globalLogger.WriteLogf(logrus.TraceLevel, "", args...)
}

func (l Logger) Debug(args ...interface{}) {
	globalLogger.WriteLogf(logrus.DebugLevel, "", args...)
}

func (l Logger) Info(args ...interface{}) {
	globalLogger.WriteLogf(logrus.InfoLevel, "", args...)
}

func (l Logger) Print(args ...interface{}) {
	globalLogger.WritePrintf("", args...)
}

func (l Logger) Warn(args ...interface{}) {
	globalLogger.WriteLogf(logrus.WarnLevel, "", args...)
}

func (l Logger) Warning(args ...interface{}) {
	Warn(args...)
}

func (l Logger) Error(args ...interface{}) {
	globalLogger.WriteLogf(logrus.ErrorLevel, "", args...)
}

func (l Logger) Fatal(args ...interface{}) {
	globalLogger.WriteLogf(logrus.FatalLevel, "", args...)
	// logger.Exit(1)
}

func (l Logger) Panic(args ...interface{}) {
	globalLogger.WriteLogf(logrus.PanicLevel, "", args...)
}

func (l Logger) Traceln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.TraceLevel, args...)
}

func (l Logger) Debugln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.DebugLevel, args...)
}

func (l Logger) Infoln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.InfoLevel, args...)
}

func (l Logger) Println(args ...interface{}) {
	globalLogger.WritePrintln(args...)
}

func (l Logger) Warnln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.WarnLevel, args...)
}

func (l Logger) Warningln(args ...interface{}) {
	Warnln(args...)
}

func (l Logger) Errorln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.ErrorLevel, args...)
}

func (l Logger) Fatalln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.FatalLevel, args...)
	// logger.Exit(1)
}

func (l Logger) Panicln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.PanicLevel, args...)
}

func GetLogger() Logger {
	return globalLogger
}
