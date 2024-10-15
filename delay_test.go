package main

import (
	"fmt"
	"testing"
	"time"
)

func TestCrawlDelay(t *testing.T) {
	tests := []struct {
		name, host    string
		expectedDelay time.Duration
		numDocs       int
	}{
		{
			"one-second",
			"http://localhost:8080/documents/top10/The Project Gutenberg EBook of A Tale of Two Cities, by Charles Dickens/",
			time.Second,
			49,
		},
		{
			"two-second",
			"http://localhost:8080/documents/top10/Dracula%20%7C%20Project%20Gutenberg/",
			2 * time.Second,
			28,
		},
		{
			"default 0.5 seconds",
			"http://localhost:8080/documents/top10/The Project Gutenberg eBook of Frankenstein, by Mary Wollstonecraft Shelley/",
			500 * time.Millisecond,
			29,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var idx Index
			url, err := parseURL(test.host)
			if err != nil {
				t.Fatal(err)
			}

			t1 := time.Now()
			idx = newMemoryIndex()
			fmt.Printf("%v %v\n", idx, url)
			crawl(&idx, url, false)
			t2 := time.Now()

			if t2.Sub(t1) < (test.expectedDelay * time.Duration(test.numDocs)) {
				t.Errorf("TestDelay was too fast\n")
			}
		})
	}
}
