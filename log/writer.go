/**
 * @Author: wenliangzhang
 * @Description:
 * @File: writer
 * @Version: 1.0.0
 * @Date: 2022/7/18 4:20 PM
 */
package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"

	"git.singularity-ai.com/backend/library/env"
)

type LogWriter struct {
	l *logrus.Logger
}

// 日志自定义格式
type LogFormatter struct{}

// 格式详情
func (s *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05")
	_, file, len, err := runtime.Caller(8) // 固定层数为7
	if err != true {
		file = filepath.Base(entry.Caller.File)
		len = entry.Caller.Line
	} else {
		file = env.AppName() + strings.Replace(file, env.RootPath(), "", 1)
	}
	msg := fmt.Sprintf("%s: %s %s:%d %s\n", strings.ToUpper(entry.Level.String()), timestamp, file, len, entry.Message)
	return []byte(msg), nil
}

func NewWriter(config LogConfig, suffix string) *LogWriter {

	if env.AppName() == "unknown" {
		env.SetAppName(config.AppName)
	}
	// 日志目录
	path := env.LogRootPath()
	// 日志文件
	fileName := filepath.Join(path, env.AppName()+suffix+".log")

	w := NewRotatelog(fileName, time.Duration(config.RotateUnit)*time.Hour, config.RotateCount)

	// 实例化
	loggerWrite := logrus.New()

	// 设置输出
	if config.Stdout == true {
		writers := []io.Writer{
			w,
			os.Stdout}
		// 同时写文件和屏幕
		fileAndStdoutWriter := io.MultiWriter(writers...)
		loggerWrite.Out = fileAndStdoutWriter
	} else {
		loggerWrite.Out = w
	}

	// 设置日志级别
	level, _ := logrus.ParseLevel(config.Level)
	loggerWrite.SetLevel(level)

	// 设置日志格式
	loggerWrite.SetReportCaller(true)
	loggerWrite.SetFormatter(new(LogFormatter))

	return &LogWriter{loggerWrite}
}

func NewStdWriter(config LogConfig) *LogWriter {

	// 实例化
	loggerWrite := logrus.New()
	// 设置输出
	loggerWrite.Out = os.Stdout
	// 设置日志级别
	level, _ := logrus.ParseLevel(config.Level)
	loggerWrite.SetLevel(level)
	// 设置日志格式
	loggerWrite.SetReportCaller(true)
	loggerWrite.SetFormatter(new(LogFormatter))

	return &LogWriter{loggerWrite}
}

func NewRotatelog(path string, rotationTime time.Duration, maxNums int) *rotatelogs.RotateLogs {
	writer, _ := rotatelogs.New(
		path+".%Y%m%d%H",
		rotatelogs.WithLinkName(path),
		rotatelogs.WithMaxAge(time.Duration(maxNums)*rotationTime),
		rotatelogs.WithRotationTime(rotationTime),
	)
	return writer
}
