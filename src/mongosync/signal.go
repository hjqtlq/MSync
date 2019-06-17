package mongosync

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var s = &signalMonitor{}

type signalMonitor struct {
	ctx        context.Context
	cancel     context.CancelFunc
	signalChan chan os.Signal

	wg *sync.WaitGroup
}

func init() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	//创建监听退出chan
	s.signalChan = make(chan os.Signal)
	//监听指定信号 ctrl+c kill
	signal.Notify(s.signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		select {
		case <-s.signalChan:
			s.cancel()
			time.Sleep(time.Second * 1)
			os.Exit(0)
		}
	}()

	s.wg = &sync.WaitGroup{}
}

func (signalMonitor *signalMonitor) beforeClose(cb func()) {
	go func() {
		select {
		case <-s.ctx.Done():
			cb()
			time.Sleep(time.Second)
		}
	}()
}
