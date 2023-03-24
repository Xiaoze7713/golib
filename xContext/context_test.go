package xContext

import (
	"context"
	"git.singularity-ai.com/backend/library/xContext/loggers/xlog"
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
