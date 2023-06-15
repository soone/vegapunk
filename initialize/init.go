package initialize

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/exp/slog"
)

var wg *sync.WaitGroup
var notifyChan chan os.Signal
var cancel context.CancelFunc
var ctx context.Context

func init() {
	ctx, cancel = context.WithCancel(context.Background())
	wg = &sync.WaitGroup{}
	notifyChan = make(chan os.Signal, 3)
	signal.Notify(notifyChan, syscall.SIGINT, syscall.SIGTERM)
}

func WGAdd(delta int) {
	wg.Add(delta)
}

func WGWait() {
	wg.Wait()
}

func WGDone() {
	wg.Done()
}

func Go2Stop() {
	notifyChan <- syscall.SIGQUIT
}

func NotifyChan() chan os.Signal {
	return notifyChan
}

func Cancel() {
	cancel()
}

func Wait2Quit(tips string) {
	<-notifyChan
	if len(tips) > 0 {
		slog.Warn(tips)
	}

	cancel()
}

func WG2Exec(handler func(args ...any), params ...any) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		handler(params...)
	}()
}

func GetContext() context.Context {
	return ctx
}

func SetContext(c context.Context) {
	ctx = c
}
