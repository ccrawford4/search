package main

import (
	"testing"
)

func TestStop(t *testing.T) {
	var tests = []struct {
		name     string
		expected bool
	}{
		{
			"yourselves",
			true,
		},
		{
			"yourselves's",
			false,
		},
		{
			"calum",
			false,
		},
		{
			"i",
			true,
		},
		{
			"i's",
			false,
		},
		{
			"ourselves",
			true,
		},
		{
			"ourselves's",
			false,
		},
		{
			"they",
			true,
		},
		{
			"they're",
			true,
		},
		{
			"myself",
			true,
		},
		{
			"you",
			true,
		},
		{
			"herself",
			true,
		},
		{
			"them",
			true,
		},
		{
			"himself",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stopWords := getStopWords()
			result := isStopWord(test.name, stopWords)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}

}
