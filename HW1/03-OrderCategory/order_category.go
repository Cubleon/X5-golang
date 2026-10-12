package main

import (
	"errors"
	"fmt"
)

var errNegativeAmount = errors.New("Order amount cannot be negative")

func orderCategory(amount int) (category string, err error) {
	if amount < 0 {
		return "", errNegativeAmount
	}

	switch {
	case amount < 10_000:
		return "small", nil
	case amount < 50_000:
		return "regular", nil
	default:
		return "large", nil
	}
}

func main() {
	amount := 50_000
	category, err := orderCategory(amount)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(category)
}
