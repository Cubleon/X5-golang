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
	var amount int

	_, err := fmt.Scan(&amount)
	if err != nil {
		fmt.Println("Input error: ", err)
		return
	}

	category, err := orderCategory(amount)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(category)
}
