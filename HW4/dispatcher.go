package main

import "fmt"

func Dispatch(all []Notifier, msg string) error {
	if err := validateMessage(msg); err != nil {
		return fmt.Errorf("dispatch: %w", err)
	}
	for i, n := range all {
		if err := safeSend(n, msg); err != nil {
			return fmt.Errorf("dispatch notifier %d: %w", i, err)
		}
	}
	return nil
}

func safeSend(n Notifier, msg string) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if cause, ok := recovered.(error); ok {
				err = fmt.Errorf("notifier panic: %w", cause)
			} else {
				err = fmt.Errorf("notifier panic: %v", recovered)
			}
		}
	}()
	return n.Send(msg)
}
