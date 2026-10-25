package main

import (
	"errors"
	"reflect"
	"testing"
)

type stubNotifier struct {
	messages []string
	err      error
}

func (n *stubNotifier) Send(msg string) error {
	n.messages = append(n.messages, msg)
	return n.err
}

func TestDispatchSendsToAll(t *testing.T) {
	first := &stubNotifier{}
	second := &stubNotifier{}
	third := &stubNotifier{}
	msg := "  Заказ #1 готов к выдаче  "

	if err := Dispatch([]Notifier{first, second, third}, msg); err != nil {
		t.Fatal(err)
	}
	for i, n := range []*stubNotifier{first, second, third} {
		if !reflect.DeepEqual(n.messages, []string{msg}) {
			t.Errorf("notifier %d received %q, want exactly one unchanged message %q", i, n.messages, msg)
		}
	}
}

func TestDispatchStopsOnError(t *testing.T) {
	wantErr := errors.New("delivery failed")
	first := &stubNotifier{}
	second := &stubNotifier{err: wantErr}
	third := &stubNotifier{}
	msg := "hello"

	err := Dispatch([]Notifier{first, second, third}, msg)
	if err != wantErr {
		t.Errorf("error = %v, want original error %v", err, wantErr)
	}
	for i, n := range []*stubNotifier{first, second} {
		if !reflect.DeepEqual(n.messages, []string{msg}) {
			t.Errorf("notifier %d received %q, want [%q]", i, n.messages, msg)
		}
	}
	if len(third.messages) != 0 {
		t.Error("notifier after the failed channel must not be called")
	}
}

func TestDispatchEmptyList(t *testing.T) {
	for _, all := range [][]Notifier{nil, {}} {
		if err := Dispatch(all, "hello"); err != nil {
			t.Errorf("Dispatch(%v) = %v, want nil", all, err)
		}
	}
}

func TestDispatchRejectsEmptyMessage(t *testing.T) {
	for _, msg := range []string{"", " \t\n", "\u00a0"} {
		t.Run(msg, func(t *testing.T) {
			n := &stubNotifier{}
			if err := Dispatch([]Notifier{n}, msg); !errors.Is(err, ErrEmptyMessage) {
				t.Errorf("error = %v, want %v", err, ErrEmptyMessage)
			}
			if len(n.messages) != 0 {
				t.Error("empty message must be rejected before calling a notifier")
			}
			if err := Dispatch(nil, msg); !errors.Is(err, ErrEmptyMessage) {
				t.Errorf("empty list error = %v, want %v", err, ErrEmptyMessage)
			}
		})
	}
}
