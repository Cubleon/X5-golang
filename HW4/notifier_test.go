package main

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type errorWriter struct {
	err error
}

func (w errorWriter) Write([]byte) (int, error) {
	return 0, w.err
}

func TestNotifiers(t *testing.T) {
	tests := []struct {
		channel string
		new     func(BaseNotifier) Notifier
	}{
		{"email", func(base BaseNotifier) Notifier { return EmailNotifier{BaseNotifier: base} }},
		{"push", func(base BaseNotifier) Notifier { return PushNotifier{BaseNotifier: base} }},
	}
	for _, tt := range tests {
		t.Run(tt.channel, func(t *testing.T) {
			t.Run("success", func(t *testing.T) {
				var output bytes.Buffer
				n := tt.new(BaseNotifier{Channel: tt.channel, Output: &output})
				msg := "  Заказ готов  "
				if err := n.Send(msg); err != nil {
					t.Fatal(err)
				}
				if got, want := output.String(), tt.channel+" "+msg+"\n"; got != want {
					t.Errorf("output = %q, want %q", got, want)
				}
			})

			t.Run("empty message", func(t *testing.T) {
				var output bytes.Buffer
				n := tt.new(BaseNotifier{Channel: tt.channel, Output: &output})
				for _, msg := range []string{"", " \t\n", "\u00a0"} {
					err := n.Send(msg)
					if !errors.Is(err, ErrEmptyMessage) {
						t.Errorf("Send(%q) = %v, want ErrEmptyMessage in chain", msg, err)
					}
					var sendErr *SendError
					if !errors.As(err, &sendErr) || sendErr.Channel != tt.channel {
						t.Errorf("Send(%q) = %v, want SendError for %s", msg, err, tt.channel)
					}
				}
				if output.Len() != 0 {
					t.Errorf("invalid message produced output: %q", output.String())
				}
			})

			t.Run("output error", func(t *testing.T) {
				n := tt.new(BaseNotifier{Channel: tt.channel, Output: errorWriter{err: io.ErrClosedPipe}})
				err := Dispatch([]Notifier{n}, "hello")
				if !errors.Is(err, io.ErrClosedPipe) {
					t.Errorf("error = %v, want output failure preserved", err)
				}
				var sendErr *SendError
				if !errors.As(err, &sendErr) || sendErr.Channel != tt.channel {
					t.Errorf("error = %v, want wrapped SendError for %s", err, tt.channel)
				}
			})
		})
	}
}

func ExampleDispatch() {
	all := []Notifier{
		EmailNotifier{BaseNotifier: BaseNotifier{Channel: "email"}},
		PushNotifier{BaseNotifier: BaseNotifier{Channel: "push"}},
	}
	if err := Dispatch(all, "Заказ #1 готов к выдаче"); err != nil {
		panic(err)
	}

	// Output:
	// email Заказ #1 готов к выдаче
	// push Заказ #1 готов к выдаче
}
