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
	globalLogger.WriteLogf(logrus.TraceLevel, format, 0, args...)
}

func Debugf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.DebugLevel, format, 0, args...)
}

func Infof(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.InfoLevel, format, 0, args...)
}

func Printf(format string, args ...interface{}) {
	globalLogger.WritePrintf(format, 0, args)
}

func Warnf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.WarnLevel, format, 0, args...)
}

func Warningf(format string, args ...interface{}) {
	Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.ErrorLevel, format, 0, args...)
}

func Fatalf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.FatalLevel, format, 0, args...)
	// logger.Exit(1)
}

func Panicf(format string, args ...interface{}) {
	globalLogger.WriteLogf(logrus.PanicLevel, format, 0, args...)
}

func Trace(args ...interface{}) {
	globalLogger.WriteLogf(logrus.TraceLevel, "", 0, args...)
}

func Debug(args ...interface{}) {
	globalLogger.WriteLogf(logrus.DebugLevel, "", 0, args...)
}

func Info(args ...interface{}) {
	globalLogger.WriteLogf(logrus.InfoLevel, "", 0, args...)
}

func Print(args ...interface{}) {
	globalLogger.WritePrintf("", 0, args)
}

func Warn(args ...interface{}) {
	globalLogger.WriteLogf(logrus.WarnLevel, "", 0, args...)
}

func Warning(args ...interface{}) {
	Warn(args...)
}

func Error(args ...interface{}) {
	globalLogger.WriteLogf(logrus.ErrorLevel, "", 0, args...)
}

func Fatal(args ...interface{}) {
	globalLogger.WriteLogf(logrus.FatalLevel, "", 0, args...)
	// logger.Exit(1)
}

func Panic(args ...interface{}) {
	globalLogger.WriteLogf(logrus.PanicLevel, "", 0, args...)
}

func Traceln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.TraceLevel, 0, args...)
}

func Debugln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.DebugLevel, 0, args...)
}

func Infoln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.InfoLevel, 0, args...)
}

func Println(args ...interface{}) {
	globalLogger.WritePrintln(0, args)
}

func Warnln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.WarnLevel, 0, args...)
}

func Warningln(args ...interface{}) {
	Warnln(args...)
}

func Errorln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.ErrorLevel, 0, args...)
}

func Fatalln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.FatalLevel, 0, args...)
	// logger.Exit(1)
}

func Panicln(args ...interface{}) {
	globalLogger.WriteLogln(logrus.PanicLevel, 0, args...)
}
