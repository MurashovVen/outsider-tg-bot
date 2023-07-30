package tg

import (
	"errors"
)

var (
	ErrUnsupportedUpdate       = errors.New("unsupported update")
	ErrUnsupportedCallbackData = errors.New("unsupported callback data")
	ErrMsgIsNotCommand         = errors.New("message is not command")
	ErrUnsupportedMessage      = errors.New("unsupported message")

	ErrInvalidInterval = errors.New("invalid interval")
)
