package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func main() {
	msg := flag.String("message", "Заказ #1 готов к выдаче", "notification message")
	flag.Parse()

	all := []Notifier{
		EmailNotifier{Channel: "email"},
		PushNotifier{Channel: "push"},
	}
	if err := Dispatch(all, *msg); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		if errors.Is(err, ErrEmptyMessage) {
			fmt.Fprintln(os.Stderr, "Укажите непустое сообщение через -message.")
		}
		if sendErr, ok := errors.AsType[*SendError](err); ok {
			fmt.Fprintln(os.Stderr, "Failed channel:", sendErr.Channel)
		}
		os.Exit(1)
	}
}
