package main

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

type notifierFunc func(string) error

func (f notifierFunc) Send(msg string) error {
	return f(msg)
}

func TestDispatchOrderAndMessage(t *testing.T) {
	var calls []int
	msg := "  hello  "
	var all []Notifier
	for i := 0; i < 3; i++ {
		all = append(all, notifierFunc(func(got string) error {
			calls = append(calls, i)
			if got != msg {
				t.Errorf("message = %q, want %q", got, msg)
			}
			return nil
		}))
	}
	if err := Dispatch(all, msg); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(calls, []int{0, 1, 2}) {
		t.Errorf("calls = %v, want [0 1 2]", calls)
	}
}

func TestDispatchWrapsErrorAndStops(t *testing.T) {
	cause := errors.New("delivery unavailable")
	wantErr := &SendError{Channel: "push", Err: cause}
	var calls []int
	all := []Notifier{
		notifierFunc(func(string) error { calls = append(calls, 0); return nil }),
		notifierFunc(func(string) error { calls = append(calls, 1); return wantErr }),
		notifierFunc(func(string) error { calls = append(calls, 2); return nil }),
	}
	err := Dispatch(all, "hello")
	if !errors.Is(err, cause) {
		t.Errorf("error = %v, want original cause in chain", err)
	}
	var sendErr *SendError
	if !errors.As(err, &sendErr) || sendErr != wantErr {
		t.Errorf("error = %v, want original SendError in chain", err)
	}
	if err == nil || !strings.Contains(err.Error(), "dispatch notifier 1:") {
		t.Errorf("error = %v, want failed notifier index", err)
	}
	if !reflect.DeepEqual(calls, []int{0, 1}) {
		t.Errorf("calls = %v, want [0 1]", calls)
	}
}

func TestDispatchRejectsEmptyMessage(t *testing.T) {
	for _, msg := range []string{"", " \t\n", "\u00a0"} {
		t.Run(msg, func(t *testing.T) {
			calls := 0
			n := notifierFunc(func(string) error { calls++; return nil })
			for _, all := range [][]Notifier{nil, {}, {n}} {
				err := Dispatch(all, msg)
				if !errors.Is(err, ErrEmptyMessage) {
					t.Errorf("error = %v, want ErrEmptyMessage in chain", err)
				}
				if err == ErrEmptyMessage {
					t.Error("error must include dispatch context")
				}
			}
			if calls != 0 {
				t.Errorf("notifier called %d times for an empty message", calls)
			}
		})
	}
}

func TestDispatchEmptyList(t *testing.T) {
	for _, all := range [][]Notifier{nil, {}} {
		if err := Dispatch(all, "hello"); err != nil {
			t.Errorf("error = %v, want nil", err)
		}
	}
}

func TestDispatchRecoversPanic(t *testing.T) {
	cause := errors.New("broken transport")
	tests := []struct {
		name  string
		value any
	}{
		{"string", "broken transport"},
		{"error", cause},
		{"nil", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls []int
			all := []Notifier{
				notifierFunc(func(string) error { calls = append(calls, 0); return nil }),
				notifierFunc(func(string) error { calls = append(calls, 1); panic(tt.value) }),
				notifierFunc(func(string) error { calls = append(calls, 2); return nil }),
			}
			err := Dispatch(all, "hello")
			if err == nil {
				t.Fatal("panic must be converted to an error")
			}
			if !strings.Contains(err.Error(), "dispatch notifier 1: notifier panic:") {
				t.Errorf("error = %v, want panic and channel index context", err)
			}
			if panicErr, ok := tt.value.(error); ok && !errors.Is(err, panicErr) {
				t.Errorf("error = %v, want panic error preserved in chain", err)
			}
			if !reflect.DeepEqual(calls, []int{0, 1}) {
				t.Errorf("calls = %v, want [0 1]", calls)
			}
		})
	}
}
