package main

import (
	"errors"
	"fmt"
)

var ErrEmptyMessage = errors.New("empty message")

type SendError struct {
	Channel string
	Err     error
}

func (e *SendError) Error() string {
	return fmt.Sprintf("send failed (%s): %v", e.Channel, e.Err)
}

func (e *SendError) Unwrap() error {
	return e.Err
}

var _ error = (*SendError)(nil)
