package shutdown

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
}

var signals = make(chan os.Signal, 1)

func Chan() <-chan struct{} {
	shutdown := make(chan struct{})

	go func() {
		<-signals
		close(shutdown)
	}()
	return shutdown
}

func Context() context.Context {
	return WithContext(context.Background())
}

func WithContext(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		<-signals
		cancel()
	}()

	return ctx
}
