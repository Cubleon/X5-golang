package main

import "testing"

func TestFizzBuzz(t *testing.T) {
	want := []string{
		"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz",
		"11", "Fizz", "13", "14", "FizzBuzz", "16", "17", "Fizz", "19", "Buzz",
		"Fizz", "22", "23", "Fizz", "Buzz", "26", "Fizz", "28", "29", "FizzBuzz",
	}

	for i, expected := range want {
		n := i + 1
		if got := fizzBuzz(n); got != expected {
			t.Errorf("fizzBuzz(%d) = %q, want %q", n, got, expected)
		}
	}
}
