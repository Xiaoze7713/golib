package test

import (
	"fmt"
	"git.singularity-ai.com/backend/library/v2/xContext"
	"git.singularity-ai.com/backend/library/v2/xContext/loggers/xlog"
	"runtime"
	"sync"
	"testing"
	"time"
)

func RunFunc2(ctx *xContext.XContext) {
	for i := 0; i < 1; i++ {
		time.Sleep(time.Second)
		ctx.Info(ctx.OperationName(), i)
		xlog.Info(ctx.OperationName(), i)
		ctx.Info(runtime.Caller(1))
		ctx.SetTag("val", fmt.Sprintf("%v", i))
		ctx.LogFields("hi", "i`am xiaoai")
	}
	ctx.Info(ctx.SerializeSpanContext())
}

//日志自定义格式
func TestLog(t *testing.T) {
	err := xInit()
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := xContext.NewXContext("start")
	defer ctx.Fin()
	ctx1 := xContext.NewChildXContext(ctx, "count1")
	ctx2 := xContext.NewChildXContext(ctx, "count2")
	ctx3 := xContext.NewChildXContext(ctx, "count3")
	ctx4 := xContext.NewChildXContext(ctx, "count4")
	var ctxList = []*xContext.XContext{ctx1, ctx2, ctx3, ctx4}
	w := sync.WaitGroup{}
	for _, c := range ctxList {
		w.Add(1)
		go func(xc *xContext.XContext) {
			RunFunc2(xc)
			defer xc.Fin()
			w.Done()
		}(c)
	}
	ctx.Info("lalala ")
	w.Wait()
	time.Sleep(time.Second * 4)
}
