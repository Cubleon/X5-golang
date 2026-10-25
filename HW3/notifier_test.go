package main

import (
	"errors"
	"testing"
)

func TestNotifiersRejectEmptyMessage(t *testing.T) {
	tests := []struct {
		name     string
		notifier Notifier
	}{
		{"email", EmailNotifier{Channel: "email"}},
		{"push", PushNotifier{Channel: "push"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range []string{"", " \t\n", "\u00a0"} {
				if err := tt.notifier.Send(msg); !errors.Is(err, ErrEmptyMessage) {
					t.Errorf("Send(%q) error = %v, want %v", msg, err, ErrEmptyMessage)
				}
			}
		})
	}
}

func ExampleDispatch() {
	all := []Notifier{
		EmailNotifier{Channel: "email"},
		PushNotifier{Channel: "push"},
	}
	if err := Dispatch(all, "Заказ #1 готов к выдаче"); err != nil {
		panic(err)
	}
}
