package main

import (
	"testing"
)

func TestCleanHref(t *testing.T) {
	tests := []struct {
		name, host, url, expectedOutput string
	}{
		{
			"home",
			"https://cs272-f24.github.io/",
			"/",
			"https://cs272-f24.github.io/",
		},
		{
			"help",
			"https://cs272-f24.github.io/",
			"/help/",
			"https://cs272-f24.github.io/help/",
		},
		{
			"syllabus",
			"https://cs272-f24.github.io/",
			"/syllabus/",
			"https://cs272-f24.github.io/syllabus/",
		},
		{
			"different-domain",
			"https://cs272-f24.github.io/",
			"https://gobyexample.com/",
			"https://gobyexample.com/",
		},
		{
			"documents",
			"https://cs272-f24.github.io/",
			"documents",
			"https://cs272-f24.github.io/documents",
		},
		{
			"no /",
			"https://cs272-f24.github.io",
			"test",
			"https://cs272-f24.github.io/test",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := clean(test.host, test.url)
			if err != nil {
				t.Fatalf("Error cleaning url: %v\n", err)
			}
			if got != test.expectedOutput {
				t.Errorf("got %q, want %q", got, test.expectedOutput)
			}
		})
	}
}
