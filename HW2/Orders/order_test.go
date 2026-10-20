package main

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func TestNewOrder(t *testing.T) {
	for _, id := range []int{-1, 0, 1, 42} {
		t.Run(strconv.Itoa(id), func(t *testing.T) {
			order, err := NewOrder(id)
			if id <= 0 {
				if !errors.Is(err, ErrInvalidID) || order != nil {
					t.Fatalf("NewOrder(%d) = (%v, %v), want (nil, %v)", id, order, err, ErrInvalidID)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if order.ID != id || order.Status != StatusNew || len(order.Items) != 0 {
				t.Fatalf("unexpected order: %+v", order)
			}
			if got := Total(*order); got != 0 {
				t.Errorf("empty order total = %d, want 0", got)
			}
		})
	}
}

func TestAddItem(t *testing.T) {
	order, err := NewOrder(1)
	if err != nil {
		t.Fatal(err)
	}
	items := []Item{{Name: "Bread", Price: 5000}, {Name: "Milk", Price: 9000}}
	for _, item := range items {
		if err := AddItem(order, item); err != nil {
			t.Fatal(err)
		}
	}
	if !reflect.DeepEqual(order.Items, items) {
		t.Errorf("items = %v, want %v", order.Items, items)
	}
	if got := Total(*order); got != 14000 {
		t.Errorf("total = %d, want 14000", got)
	}
}

func TestAddItemInvalidData(t *testing.T) {
	tests := []struct {
		name    string
		item    Item
		wantErr error
	}{
		{"empty name", Item{Price: 100}, ErrEmptyItemName},
		{"whitespace name", Item{Name: " \t\n", Price: 100}, ErrEmptyItemName},
		{"negative price", Item{Name: "Bread", Price: -1}, ErrNegativePrice},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := Order{ID: 1, Status: StatusNew, Items: []Item{{Name: "Milk", Price: 9000}}}
			before := CloneOrder(order)
			beforeCap := cap(order.Items)
			if err := AddItem(&order, tt.item); !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(order, before) || cap(order.Items) != beforeCap {
				t.Errorf("invalid item changed order: %+v", order)
			}
		})
	}
}

func TestAddItemNilOrder(t *testing.T) {
	if err := AddItem(nil, Item{Name: "Bread", Price: 5000}); !errors.Is(err, ErrNilOrder) {
		t.Errorf("error = %v, want %v", err, ErrNilOrder)
	}
}

func TestAddItemZeroPrice(t *testing.T) {
	order := Order{ID: 1}
	if err := AddItem(&order, Item{Name: "Gift", Price: 0}); err != nil {
		t.Fatal(err)
	}
	if len(order.Items) != 1 || Total(order) != 0 {
		t.Errorf("unexpected order: %+v", order)
	}
}

func TestTotal(t *testing.T) {
	tests := []struct {
		name  string
		items []Item
		want  int
	}{
		{"nil slice", nil, 0},
		{"empty slice", []Item{}, 0},
		{"one item", []Item{{Name: "Bread", Price: 5000}}, 5000},
		{"two items", []Item{{Name: "Bread", Price: 5000}, {Name: "Milk", Price: 9000}}, 14000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Total(Order{Items: tt.items}); got != tt.want {
				t.Errorf("Total() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCloneOrderIndependent(t *testing.T) {
	original := Order{ID: 1, Status: StatusNew, Items: make([]Item, 1, 4)}
	original.Items[0] = Item{Name: "Bread", Price: 5000}
	before := CloneOrder(original)
	clone := CloneOrder(original)
	if !reflect.DeepEqual(clone, original) {
		t.Fatalf("clone = %+v, want %+v", clone, original)
	}
	clone.ID = 2
	clone.Status = StatusCompleted
	clone.Items[0].Name = "Bun"
	clone.Items[0].Price = 1000
	if err := AddItem(&clone, Item{Name: "Milk", Price: 9000}); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(original, before) {
		t.Errorf("changing clone changed original: %+v", original)
	}
	original.Items[0].Price = 7000
	if clone.Items[0].Price != 1000 {
		t.Error("changing original changed clone")
	}
}

func TestCloneOrderEmpty(t *testing.T) {
	for _, items := range [][]Item{nil, {}} {
		original := Order{ID: 1, Status: StatusNew, Items: items}
		clone := CloneOrder(original)
		if !reflect.DeepEqual(clone, original) {
			t.Fatalf("clone = %+v, want %+v", clone, original)
		}
		if err := AddItem(&clone, Item{Name: "Bread", Price: 5000}); err != nil {
			t.Fatal(err)
		}
		if len(original.Items) != 0 {
			t.Error("adding to clone changed original")
		}
	}
}

func TestStatusNames(t *testing.T) {
	want := map[string]string{
		StatusNew: "New", StatusProcessing: "Processing",
		StatusCompleted: "Completed", StatusCancelled: "Cancelled",
	}
	if !reflect.DeepEqual(StatusNames, want) {
		t.Errorf("StatusNames = %v, want %v", StatusNames, want)
	}
	names := makeStatusNames()
	names["test"] = "Test"
	if names["test"] != "Test" {
		t.Error("map entry was not saved")
	}
}

func TestAddItemSliceGrowth(t *testing.T) {
	order := Order{ID: 1, Items: make([]Item, 0, 2)}
	for i := 0; i < 3; i++ {
		beforeLen, beforeCap := len(order.Items), cap(order.Items)
		if err := AddItem(&order, Item{Name: "Product", Price: 100}); err != nil {
			t.Fatal(err)
		}
		t.Logf("append %d: len %d -> %d, cap %d -> %d",
			i+1, beforeLen, len(order.Items), beforeCap, cap(order.Items))
		if len(order.Items) != beforeLen+1 {
			t.Fatal("append must increase len by one")
		}
		if beforeLen < beforeCap && cap(order.Items) != beforeCap {
			t.Fatal("capacity changed despite available space")
		}
		if beforeLen == beforeCap && cap(order.Items) <= beforeCap {
			t.Fatal("capacity did not grow when slice was full")
		}
	}
}
