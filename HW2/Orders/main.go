package main

import "fmt"

func main() {
	order, err := NewOrder(1)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Empty order: total = %d cents, len = %d, cap = %d\n",
		Total(*order), len(order.Items), cap(order.Items))

	for _, item := range []Item{
		{Name: "Bread", Price: 5000},
		{Name: "Milk", Price: 9000},
	} {
		if err := AddItem(order, item); err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Added: %s; len = %d, cap = %d\n",
			item.Name, len(order.Items), cap(order.Items))
	}

	clone := CloneOrder(*order)
	clone.Items[0].Price = 1000
	fmt.Printf("Order #%d: %s, total = %d cents\n",
		order.ID, StatusNames[order.Status], Total(*order))
	fmt.Printf("Copy total after modification = %d cents\n", Total(clone))
}
