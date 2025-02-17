package loggers

import (
	"github.com/golib/v3/xContext/loggers/xlog"
	"testing"
	"time"
)

func TestLogDefault(t *testing.T) {
	xlog.SetupLogDefault()
	xlog.SetLevel(xlog.TRACE)
	xlog.Info("info")
	xlog.Debug("debug")
	xlog.Error("error")
	xlog.Fatal("fatal")
	xlog.Warn("warn")
	xlog.Trace("trace")
	time.Sleep(time.Second * 3)
}
func TestLogColor(t *testing.T) {
	err := xlog.SetupLogWithConf(xlog.LogConfig{
		Level: "trace",
		FW: xlog.ConfFileWriter{
			On:              true,
			Fmt:             xlog.FmtTypeJson,
			LogPath:         "./log/xlog.log.info",
			RotateLogPath:   "./log/xlog.log.info-%Y%M%D",
			WfLogPath:       "./log/xlog.log.wf",
			RotateWfLogPath: "./log/xlog.log.wf-%Y%M%D",
			//PublicLogPath:       "./log/xlog.log.public",
			//RotatePublicLogPath: "./log/xlog.log.public-%Y%M%D",
		},
		CW: xlog.ConfConsoleWriter{
			On:  true,
			Fmt: xlog.FmtTypeColorRaw,
		},
	})
	if err != nil {
		return
	}
	xlog.SetLevel(xlog.TRACE)
	xlog.Info("info")
	xlog.Debug("debug")
	xlog.Error("error")
	xlog.Fatal("fatal")
	xlog.Warn("warn")
	xlog.Trace("trace")
	time.Sleep(time.Second * 3)
}
