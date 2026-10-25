package main

func Dispatch(all []Notifier, msg string) error {
	if err := validateMessage(msg); err != nil {
		return err
	}
	for _, n := range all {
		if err := n.Send(msg); err != nil {
			return err
		}
	}
	return nil
}
