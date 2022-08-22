package sign

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Wait 重要的wait事件， 任意一个done全局，close && collection
func Wait(ctx context.Context, waitFunction []func() <-chan struct{}, signals []os.Signal, collectFunctions []func()) (cancel context.CancelFunc) {
	newCtx, cancel := signal.NotifyContext(ctx, signals...)
	for _, wf := range waitFunction {
		go func(wf func() <-chan struct{}) {
			for {
				select {
				case <-wf():
					cancel()
				}
			}
		}(wf)
	}
	for {
		select {
		case <-newCtx.Done():
			for _, f := range collectFunctions {
				f()
			}
		default:
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func CommonlyUsedSign() []os.Signal {
	return []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGKILL,
		syscall.SIGTERM,
	}
}

func CC() {
	ctx := context.Background()
	tc := make(chan interface{}, 1)
	select {
	case tc <- ctx.Done():

	}
}
