package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Notifier interface {
	Send(msg string) error
}

type BaseNotifier struct {
	Channel string
	Output  io.Writer
}

func (b BaseNotifier) Log(msg string) error {
	output := b.Output
	if output == nil {
		output = os.Stdout
	}
	_, err := fmt.Fprintln(output, b.Channel, msg)
	return err
}

type EmailNotifier struct {
	BaseNotifier
}

func (n EmailNotifier) Send(msg string) error {
	if err := validateMessage(msg); err != nil {
		return &SendError{Channel: n.Channel, Err: err}
	}
	if err := n.Log(msg); err != nil {
		return &SendError{Channel: n.Channel, Err: err}
	}
	return nil
}

type PushNotifier struct {
	BaseNotifier
}

func (n PushNotifier) Send(msg string) error {
	if err := validateMessage(msg); err != nil {
		return &SendError{Channel: n.Channel, Err: err}
	}
	if err := n.Log(msg); err != nil {
		return &SendError{Channel: n.Channel, Err: err}
	}
	return nil
}

var (
	_ Notifier = EmailNotifier{}
	_ Notifier = PushNotifier{}
)

func validateMessage(msg string) error {
	if strings.TrimSpace(msg) == "" {
		return ErrEmptyMessage
	}
	return nil
}
