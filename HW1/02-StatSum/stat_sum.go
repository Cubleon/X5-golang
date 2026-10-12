package main

import (
	"errors"
	"fmt"
)

var errEmptyOrders = errors.New("No orders: cannot calculate statistics")

func calculateStats(amounts []int) (min, max, sum int, average float64, err error) {
	if len(amounts) == 0 {
		return 0, 0, 0, 0, errEmptyOrders
	}

	min, max = amounts[0], amounts[0]
	for _, amount := range amounts {
		if amount < min {
			min = amount
		}
		if amount > max {
			max = amount
		}
		sum += amount
	}

	average = float64(sum) / float64(len(amounts))
	return min, max, sum, average, nil
}

func main() {
	amounts := []int{12000, 7500, 31000, 4500}
	min, max, sum, average, err := calculateStats(amounts)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Minimum:", min)
	fmt.Println("Maximum:", max)
	fmt.Println("Sum:", sum)
	fmt.Println("Average:", average)
}
