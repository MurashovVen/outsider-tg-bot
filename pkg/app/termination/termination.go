package termination

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"tg-bot/pkg/app/logger"
)

type (
	Waiter struct {
		logger *logger.Logger
	}

	WaiterOption func(*Waiter)
)

func NewWaiter(options ...WaiterOption) *Waiter {
	w := &Waiter{
		logger: logger.NewNop(),
	}
	for _, opt := range options {
		opt(w)
	}

	return w
}

func WaiterWithLogger(logger *logger.Logger) WaiterOption {
	return func(waiter *Waiter) {
		waiter.logger = logger.Named("TerminationWaiter")
	}
}

var ErrStopped = errors.New("application is stopped")

func (w *Waiter) WaitFunc(ctx context.Context) func() error {
	return func() error {
		signalsChan := make(chan os.Signal, 1)
		signal.Notify(signalsChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

		select {
		case sig := <-signalsChan:
			w.logger.Debug("got termination signal", zap.String("signal", sig.String()))

			return ErrStopped

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
