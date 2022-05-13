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
	loggerDef.l.Logf(logrus.TraceLevel, format, args...)
}

func Debugf(format string, args ...interface{}) {
	loggerDef.l.Logf(logrus.DebugLevel, format, args...)
}

func Infof(format string, args ...interface{}) {
	loggerDef.l.Logf(logrus.InfoLevel, format, args...)
}

func Printf(format string, args ...interface{}) {
	loggerDef.l.Printf(format, args)
}

func Warnf(format string, args ...interface{}) {
	loggerWf.l.Logf(logrus.WarnLevel, format, args...)
}

func Warningf(format string, args ...interface{}) {
	Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	loggerWf.l.Logf(logrus.ErrorLevel, format, args...)
}

func Fatalf(format string, args ...interface{}) {
	loggerWf.l.Logf(logrus.FatalLevel, format, args...)
	//logger.Exit(1)
}

func Panicf(format string, args ...interface{}) {
	loggerWf.l.Logf(logrus.PanicLevel, format, args...)
}

func Trace(args ...interface{}) {
	loggerDef.l.Log(logrus.TraceLevel, args...)
}

func Debug(args ...interface{}) {
	loggerDef.l.Log(logrus.DebugLevel, args...)
}

func Info(args ...interface{}) {
	loggerDef.l.Log(logrus.InfoLevel, args...)
}

func Print(args ...interface{}) {
	loggerDef.l.Print(args)
}

func Warn(args ...interface{}) {
	loggerWf.l.Log(logrus.WarnLevel, args...)
}

func Warning(args ...interface{}) {
	Warn(args...)
}

func Error(args ...interface{}) {
	loggerWf.l.Log(logrus.ErrorLevel, args...)
}

func Fatal(args ...interface{}) {
	loggerWf.l.Log(logrus.FatalLevel, args...)
	//logger.Exit(1)
}

func Panic(args ...interface{}) {
	loggerWf.l.Log(logrus.PanicLevel, args...)
}

func Traceln(args ...interface{}) {
	loggerDef.l.Logln(logrus.TraceLevel, args...)
}

func Debugln(args ...interface{}) {
	loggerDef.l.Logln(logrus.DebugLevel, args...)
}

func Infoln(args ...interface{}) {
	loggerDef.l.Logln(logrus.InfoLevel, args...)
}

func Println(args ...interface{}) {
	loggerDef.l.Println(args)
}

func Warnln(args ...interface{}) {
	loggerWf.l.Logln(logrus.WarnLevel, args...)
}

func Warningln(args ...interface{}) {
	Warnln(args...)
}

func Errorln(args ...interface{}) {
	loggerWf.l.Logln(logrus.ErrorLevel, args...)
}

func Fatalln(args ...interface{}) {
	loggerWf.l.Logln(logrus.FatalLevel, args...)
	//logger.Exit(1)
}

func Panicln(args ...interface{}) {
	loggerWf.l.Logln(logrus.PanicLevel, args...)
}
