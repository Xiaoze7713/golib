/**
 * @Author: wenliangzhang
 * @Description:
 * @File: log
 * @Version: 1.0.0
 * @Date: 2022/4/21 7:55 PM
 */
package log

import (
	"github.com/sirupsen/logrus"
)

func Tracef(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.TraceLevel, format, args...)
}

func Debugf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.DebugLevel, format, args...)
}

func Infof(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.InfoLevel, format, args...)
}

func Printf(format string, args ...interface{}) {
	globalLogger.WritePrintf(format, args)
}

func Warnf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.WarnLevel, format, args...)
}

func Warningf(format string, args ...interface{}) {
	Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.ErrorLevel, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.FatalLevel, format, args...)
	// logger.Exit(1)
}

func Panicf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.PanicLevel, format, args...)
}

func Trace(args ...interface{}) {
	globalLogger.WriteLogf(logrus.TraceLevel, "", args...)
}

func Debug(args ...interface{}) {
	globalLogger.WriteLogf(logrus.DebugLevel, "", args...)
}

func Info(args ...interface{}) {
	globalLogger.WriteLogf(logrus.InfoLevel, "", args...)
}

func Print(args ...interface{}) {
	globalLogger.WritePrintf("", args)
}

func Warn(args ...interface{}) {
	globalLogger.WriteLogf(logrus.WarnLevel, "", args...)
}

func Warning(args ...interface{}) {
	Warn(args...)
}

func Error(args ...interface{}) {
	globalLogger.WriteLogf(logrus.ErrorLevel, "", args...)
}

func Fatal(args ...interface{}) {
	globalLogger.WriteLogf(logrus.FatalLevel, "", args...)
	// logger.Exit(1)
}

func Panic(args ...interface{}) {
	globalLogger.WriteLogf(logrus.PanicLevel, "", args...)
}

func Traceln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.TraceLevel, args...)
}

func Debugln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.DebugLevel, args...)
}

func Infoln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.InfoLevel, args...)
}

func Println(args ...interface{}) {
	globalLogger.WritePrintln(args)
}

func Warnln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.WarnLevel, args...)
}

func Warningln(args ...interface{}) {
	Warnln(args...)
}

func Errorln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.ErrorLevel, args...)
}

func Fatalln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.FatalLevel, args...)
	// logger.Exit(1)
}

func Panicln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.PanicLevel, args...)
}
