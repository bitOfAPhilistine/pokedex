package main

import (
	"testing"
)


func TestCleanInput(t *testing.T) {
	cleanInputCases := []struct {
		input string
		expected []string
	}{
		{
			input: "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "  HELLO  WORLD  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "Hello World",
			expected: []string{"hello", "world"},
		},
		{
			input: "hello",
			expected: []string{"hello"},
		},
		{
			input: "",
			expected: []string{},
		},
		{
			input: "   ",
			expected: []string{},
		},
	}

	for _, c := range cleanInputCases {
		result := cleanInput(c.input)
		for i, v := range result {
			if v != c.expected[i] {
				t.Errorf("cleanInput(%q) = %v, want %v", c.input, result, c.expected)
				break
			}
		}
	}
}