package logger

import (
	"context"
	"errors"
	"log"

	"go.uber.org/zap"

	"tg-bot/pkg/app"
)

type Logger struct {
	*zap.Logger
}

func MustCreateLogger(environment app.Environment) *Logger {
	l, err := CreateLogger(environment)
	if err != nil {
		panic("creating logger: " + err.Error())
	}

	return l
}

func CreateLogger(environment app.Environment) (*Logger, error) {
	switch environment {
	case app.DevEnv:
		l, err := zap.NewDevelopment()
		return &Logger{Logger: l}, err

	default:
		return nil, errors.New("unknown environment")
	}
}

func NewNop() *Logger {
	return &Logger{
		Logger: zap.NewNop(),
	}
}
func (l Logger) SyncWaiter(ctx context.Context) func() error {
	return func() error {
		select {
		case <-ctx.Done():
			// todo fix
			if err := l.Sync(); err != nil {
				log.Printf("sync logger error: %v", err)
			}
		}

		return nil
	}
}

func (l Logger) Named(name string) *Logger {
	return &Logger{
		Logger: l.Logger.Named(name),
	}
}
