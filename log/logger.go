/**
 * @Author: wenliangzhang
 * @Description:
 * @File: logger
 * @Version: 1.0.0
 * @Date: 2022/4/21 11:12 AM
 */
package log

import (
	"fmt"

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

func (l Logger) WriteLogf(level logrus.Level, format string, depth int, args ...interface{}) {
	names := buildWriterNames(level)
	msg := FormatMsg(level, format, depth, false, args...)
	for _, name := range names {
		l[name].l.Log(level, msg)
	}
}

func sprintln(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}

func (l Logger) WriteLogln(level logrus.Level, depth int, args ...interface{}) {
	names := buildWriterNames(level)
	msg := FormatMsg(level, "", depth, true, args...)
	for _, name := range names {
		l[name].l.Log(level, msg)
	}
}

func (l Logger) WritePrintf(format string, depth int, args ...interface{}) {
	names := buildWriterNames(logrus.InfoLevel)
	msg := FormatMsg(logrus.InfoLevel, format, depth, false, args...)
	for _, name := range names {
		l[name].l.Print(msg)
	}
}

func (l Logger) WritePrintln(depth int, args ...interface{}) {
	names := buildWriterNames(logrus.InfoLevel)
	msg := FormatMsg(logrus.InfoLevel, "", depth, true, args...)
	for _, name := range names {
		l[name].l.Print(msg)
	}
}

func (l Logger) Tracef(format string, depth int, args ...interface{}) {
	l.WriteLogf(logrus.TraceLevel, format, depth, args...)
}

func (l Logger) Debugf(format string, depth int, args ...interface{}) {
	globalLogger.WriteLogf(logrus.DebugLevel, format, depth, args...)
}

func (l Logger) Infof(format string, depth int, args ...interface{}) {
	globalLogger.WriteLogf(logrus.InfoLevel, format, depth, args...)
}

func (l Logger) Printf(format string, depth int, args ...interface{}) {
	globalLogger.WritePrintf(format, depth, args...)
}

func (l Logger) Warnf(format string, depth int, args ...interface{}) {
	globalLogger.WriteLogf(logrus.WarnLevel, format, depth, args...)
}

func (l Logger) Warningf(format string, depth int, args ...interface{}) {
	l.Warnf(format, depth, args...)
}

func (l Logger) Errorf(format string, depth int, args ...interface{}) {
	globalLogger.WriteLogf(logrus.ErrorLevel, format, depth, args...)
}

func (l Logger) Fatalf(format string, depth int, args ...interface{}) {
	globalLogger.WriteLogf(logrus.FatalLevel, format, depth, args...)
	// logger.Exit(1)
}

func (l Logger) Panicf(format string, depth int, args ...interface{}) {
	globalLogger.WriteLogf(logrus.PanicLevel, format, depth, args...)
}

func (l Logger) Trace(depth int, args ...interface{}) {
	globalLogger.WriteLogf(logrus.TraceLevel, "", depth, args...)
}

func (l Logger) Debug(depth int, args ...interface{}) {
	globalLogger.WriteLogf(logrus.DebugLevel, "", depth, args...)
}

func (l Logger) Info(depth int, args ...interface{}) {
	globalLogger.WriteLogf(logrus.InfoLevel, "", depth, args...)
}

func (l Logger) Print(depth int, args ...interface{}) {
	globalLogger.WritePrintf("", depth, args...)
}

func (l Logger) Warn(depth int, args ...interface{}) {
	globalLogger.WriteLogf(logrus.WarnLevel, "", depth, args...)
}

func (l Logger) Warning(depth int, args ...interface{}) {
	l.Warn(depth, args...)
}

func (l Logger) Error(depth int, args ...interface{}) {
	globalLogger.WriteLogf(logrus.ErrorLevel, "", depth, args...)
}

func (l Logger) Fatal(depth int, args ...interface{}) {
	globalLogger.WriteLogf(logrus.FatalLevel, "", depth, args...)
	// logger.Exit(1)
}

func (l Logger) Panic(depth int, args ...interface{}) {
	globalLogger.WriteLogf(logrus.PanicLevel, "", depth, args...)
}

func (l Logger) Traceln(depth int, args ...interface{}) {
	globalLogger.WriteLogln(logrus.TraceLevel, depth, args...)
}

func (l Logger) Debugln(depth int, args ...interface{}) {
	globalLogger.WriteLogln(logrus.DebugLevel, depth, args...)
}

func (l Logger) Infoln(depth int, args ...interface{}) {
	globalLogger.WriteLogln(logrus.InfoLevel, depth, args...)
}

func (l Logger) Println(depth int, args ...interface{}) {
	globalLogger.WritePrintln(depth, args...)
}

func (l Logger) Warnln(depth int, args ...interface{}) {
	globalLogger.WriteLogln(logrus.WarnLevel, depth, args...)
}

func (l Logger) Warningln(depth int, args ...interface{}) {
	l.Warnln(depth, args...)
}

func (l Logger) Errorln(depth int, args ...interface{}) {
	globalLogger.WriteLogln(logrus.ErrorLevel, depth, args...)
}

func (l Logger) Fatalln(depth int, args ...interface{}) {
	globalLogger.WriteLogln(logrus.FatalLevel, depth, args...)
	// logger.Exit(1)
}

func (l Logger) Panicln(depth int, args ...interface{}) {
	globalLogger.WriteLogln(logrus.PanicLevel, depth, args...)
}

func GetLogger() Logger {
	return globalLogger
}
