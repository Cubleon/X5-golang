package main

import (
	"errors"
	"fmt"
	"strings"
)

type Notifier interface {
	Send(msg string) error
}

var ErrEmptyMessage = errors.New("message must not be empty")

type BaseNotifier struct {
	Channel string
}

func (b BaseNotifier) Log(msg string) {
	fmt.Println(b.Channel, msg)
}

type EmailNotifier struct {
	BaseNotifier
}

func (n EmailNotifier) Send(msg string) error {
	if err := validateMessage(msg); err != nil {
		return err
	}
	n.Log(msg)
	return nil
}

type PushNotifier struct {
	BaseNotifier
}

func (n PushNotifier) Send(msg string) error {
	if err := validateMessage(msg); err != nil {
		return err
	}
	n.Log(msg)
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
