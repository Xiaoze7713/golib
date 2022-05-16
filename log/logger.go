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
	"git.singularity-ai.com/backend/library/env"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Logger struct {
	l *logrus.Logger
}

//日志自定义格式
type LogFormatter struct{}

//格式详情
func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	_, file, len, err := runtime.Caller(7) // 固定层数为7
	if err != true {
		file = filepath.Base(entry.Caller.File)
		len = entry.Caller.Line
	} else {
		file = env.AppName() + strings.Replace(file, env.RootPath(), "", 1)
	}
	msg := fmt.Sprintf("%s: %s %s:%d %s\n", strings.ToUpper(entry.Level.String()), timestamp, file, len, entry.Message)
	return []byte(msg), nil
}

func NewLogger(config LogConfig, suffix string) *logrus.Logger {

	if env.AppName() == "unknown" {
		env.SetAppName(config.AppName)
	}
	//日志目录
	path := env.LogRootPath()
	//日志文件
	fileName := filepath.Join(path, env.AppName()+suffix+".log")

	w := NewWriter(fileName, time.Duration(config.RotateUnit)*time.Hour, config.RotateCount)

	//实例化
	logger := logrus.New()

	//设置输出
	if config.Stdout == true {
		writers := []io.Writer{
			w,
			os.Stdout}
		//同时写文件和屏幕
		fileAndStdoutWriter := io.MultiWriter(writers...)
		logger.Out = fileAndStdoutWriter
	} else {
		logger.Out = w
	}

	//设置日志级别
	level, _ := logrus.ParseLevel(config.Level)
	logger.SetLevel(level)

	//设置日志格式
	logger.SetReportCaller(true)
	logger.SetFormatter(new(LogFormatter))

	return logger
}

func NewWriter(path string, rotationTime time.Duration, maxNums int) *rotatelogs.RotateLogs {
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(maxNums)*rotationTime),
		rotatelogs.WithRotationTime(rotationTime),
	)
	return writer
}

var defaultLoggerConfig = LogConfig{
	"singularity",
	1,
	48,
	"warn",
	false,
}

func GetDefaultLogger() *Logger {
	if loggerDef == nil {
		// 初始化日志目录
		initLogDir(env.LogRootPath())
		loggerDef = &Logger{
			NewLogger(defaultLoggerConfig, ""),
		}
	}
	return loggerDef
}

func GetWfLogger() *Logger {
	if loggerWf == nil {
		// 初始化日志目录
		initLogDir(env.LogRootPath())
		loggerWf = &Logger{
			NewLogger(defaultLoggerConfig, "wf"),
		}
	}
	return loggerWf
}
