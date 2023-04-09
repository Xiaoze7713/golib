package loggers

import (
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
	"testing"
	"time"
)

func TestLogColor(t *testing.T) {
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
