package xContext

import (
	"context"
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
	"sync"
	"testing"
	"time"
)

func TestCtx(t *testing.T) {
	xlog.SetupLogDefault()
	ctxF := context.Background()
	ctxFC, fc := context.WithCancel(ctxF)
	defer fc()
	ctxC, cc := context.WithCancel(ctxFC)
	defer cc()
	w := &sync.WaitGroup{}
	w.Add(1)
	go func() {
		select {
		case <-ctxC.Done():
			xlog.Info("ctxC done")
		}
		w.Done()
	}()
	w.Add(1)
	go func() {
		time.Sleep(time.Second * 3)
		fc()
		select {
		case <-ctxFC.Done():
			xlog.Info("ctxF done")
		}
		w.Done()
	}()
	w.Wait()
	xlog.Info("all fin")
	time.Sleep(time.Second * 3)
}

func TestLog(t *testing.T) {
	lc := xlog.LogConfig{
		Level: "debug",
		FW: xlog.ConfFileWriter{
			On:              true,
			LogPath:         "/Users/cangxiaoze/user/git/logxx/chat-api.log.info",
			RotateLogPath:   "/Users/cangxiaoze/user/git/logxx/chat-api.log.info-%Y%M%D",
			WfLogPath:       "/Users/cangxiaoze/user/git/logxx/chat-api.log.wf",
			RotateWfLogPath: "/Users/cangxiaoze/user/git/logxx/chat-api.log.wf-%Y%M%D",
		},
		CW: xlog.ConfConsoleWriter{
			On:    true,
			Color: true,
		},
	}
	xlog.SetupLogWithConf(lc)
	defer xlog.Close()
	InitByNullOptNoLog()
	xlog.GetLogger().SetFmtJson()
	ctx := NewXContext("test")
	ctx.Info("你好", "hello")
	//xlog.GetLogger().SetFmtRaw()
	ctx.Info("你好", "hello")
	time.Sleep(time.Second * 3)
}
