package main

import "fmt"

func main() {
	all := []Notifier{
		EmailNotifier{Channel: "email"},
		PushNotifier{Channel: "push"},
	}

	if err := Dispatch(all, "Заказ #1 готов к выдаче"); err != nil {
		fmt.Println("Error:", err)
		return
	}

	for _, n := range all {
		if email, ok := n.(EmailNotifier); ok {
			fmt.Println("Email channel:", email.Channel)
		}
	}
}
