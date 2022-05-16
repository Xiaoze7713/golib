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
	GetDefaultLogger().l.Logf(logrus.TraceLevel, format, args...)
}

func Debugf(format string, args ...interface{}) {
	GetDefaultLogger().l.Logf(logrus.DebugLevel, format, args...)
}

func Infof(format string, args ...interface{}) {
	GetDefaultLogger().l.Logf(logrus.InfoLevel, format, args...)
}

func Printf(format string, args ...interface{}) {
	GetDefaultLogger().l.Printf(format, args)
}

func Warnf(format string, args ...interface{}) {
	GetWfLogger().l.Logf(logrus.WarnLevel, format, args...)
}

func Warningf(format string, args ...interface{}) {
	Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	GetWfLogger().l.Logf(logrus.ErrorLevel, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	GetWfLogger().l.Logf(logrus.FatalLevel, format, args...)
	//logger.Exit(1)
}

func Panicf(format string, args ...interface{}) {
	GetWfLogger().l.Logf(logrus.PanicLevel, format, args...)
}

func Trace(args ...interface{}) {
	GetDefaultLogger().l.Log(logrus.TraceLevel, args...)
}

func Debug(args ...interface{}) {
	GetDefaultLogger().l.Log(logrus.DebugLevel, args...)
}

func Info(args ...interface{}) {
	GetDefaultLogger().l.Log(logrus.InfoLevel, args...)
}

func Print(args ...interface{}) {
	GetDefaultLogger().l.Print(args)
}

func Warn(args ...interface{}) {
	GetWfLogger().l.Log(logrus.WarnLevel, args...)
}

func Warning(args ...interface{}) {
	Warn(args...)
}

func Error(args ...interface{}) {
	GetWfLogger().l.Log(logrus.ErrorLevel, args...)
}

func Fatal(args ...interface{}) {
	GetWfLogger().l.Log(logrus.FatalLevel, args...)
	//logger.Exit(1)
}

func Panic(args ...interface{}) {
	GetWfLogger().l.Log(logrus.PanicLevel, args...)
}

func Traceln(args ...interface{}) {
	GetDefaultLogger().l.Logln(logrus.TraceLevel, args...)
}

func Debugln(args ...interface{}) {
	GetDefaultLogger().l.Logln(logrus.DebugLevel, args...)
}

func Infoln(args ...interface{}) {
	GetDefaultLogger().l.Logln(logrus.InfoLevel, args...)
}

func Println(args ...interface{}) {
	GetDefaultLogger().l.Println(args)
}

func Warnln(args ...interface{}) {
	GetWfLogger().l.Logln(logrus.WarnLevel, args...)
}

func Warningln(args ...interface{}) {
	Warnln(args...)
}

func Errorln(args ...interface{}) {
	GetWfLogger().l.Logln(logrus.ErrorLevel, args...)
}

func Fatalln(args ...interface{}) {
	GetWfLogger().l.Logln(logrus.FatalLevel, args...)
	//logger.Exit(1)
}

func Panicln(args ...interface{}) {
	GetWfLogger().l.Logln(logrus.PanicLevel, args...)
}
