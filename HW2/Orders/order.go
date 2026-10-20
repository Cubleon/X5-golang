package main

import (
	"errors"
	"strings"
)

const (
	StatusNew        = "new"
	StatusProcessing = "processing"
	StatusCompleted  = "completed"
	StatusCancelled  = "cancelled"
)

var (
	ErrInvalidID     = errors.New("order ID must be positive")
	ErrNilOrder      = errors.New("order must not be nil")
	ErrEmptyItemName = errors.New("item name must not be empty")
	ErrNegativePrice = errors.New("item price must not be negative")
)

var StatusNames = makeStatusNames()

func makeStatusNames() map[string]string {
	names := make(map[string]string)
	names[StatusNew] = "New"
	names[StatusProcessing] = "Processing"
	names[StatusCompleted] = "Completed"
	names[StatusCancelled] = "Cancelled"
	return names
}

// Price is stored in cents.
type Item struct {
	Name  string
	Price int
}

type Order struct {
	ID     int
	Status string
	Items  []Item
}

func NewOrder(id int) (*Order, error) {
	if id <= 0 {
		return nil, ErrInvalidID
	}
	return &Order{ID: id, Status: StatusNew}, nil
}

func AddItem(order *Order, item Item) error {
	if order == nil {
		return ErrNilOrder
	}
	if strings.TrimSpace(item.Name) == "" {
		return ErrEmptyItemName
	}
	if item.Price < 0 {
		return ErrNegativePrice
	}

	order.Items = append(order.Items, item)
	return nil
}

func Total(order Order) int {
	total := 0
	for _, item := range order.Items {
		total += item.Price
	}
	return total
}

func CloneOrder(order Order) Order {
	clone := order
	if order.Items != nil {
		clone.Items = make([]Item, len(order.Items))
		copy(clone.Items, order.Items)
	}
	return clone
}
