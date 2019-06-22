package mongosync

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var SignalMonitor = &signalMonitor{}

type signalMonitor struct {
	Context    context.Context
	cancel     context.CancelFunc
	signalChan chan os.Signal

	wg *sync.WaitGroup
}

func init() {
	SignalMonitor.Context, SignalMonitor.cancel = context.WithCancel(context.Background())
	//创建监听退出chan
	SignalMonitor.signalChan = make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(SignalMonitor.signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		select {
		case <-SignalMonitor.signalChan:
			SignalMonitor.cancel()
			time.Sleep(time.Second * 1)
			os.Exit(0)
		}
	}()

	SignalMonitor.wg = &sync.WaitGroup{}
}

func (signalMonitor *signalMonitor) BeforeClose(cb func()) {
	go func() {
		select {
		case <-SignalMonitor.Context.Done():
			cb()
			time.Sleep(time.Second)
		}
	}()
}
