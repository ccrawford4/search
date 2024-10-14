package main

import (
	"github.com/go-test/deep"
	"testing"
)

func TestDisallow(t *testing.T) {
	tests := []struct {
		name, host string
		expected   *SearchResult
	}{
		{
			"rabbit",
			"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/",
			&SearchResult{
				Frequency{
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap11.html": 8,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap01.html": 9,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap12.html": 8,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap08.html": 6,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap02.html": 4,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/":            2,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap10.html": 1,
				},
				12,
			},
		},
		{
			"jekyll",
			"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/",
			&SearchResult{
				Frequency{
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/":            8,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap06.html": 8,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap07.html": 6,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap08.html": 14,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap09.html": 12,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap10.html": 30,
				},
				6,
			},
		},
		{
			"Nicolo",
			"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/",
			&SearchResult{
				Frequency{
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/":            5,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap00.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap01.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap02.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap03.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap04.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap05.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap06.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap07.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap08.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap09.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap10.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap11.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap12.html": 2,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap13.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap14.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap15.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap16.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap17.html": 1,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap18.html": 2,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap19.html": 1,
				},
				21,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var idx Index
			url, err := parseURL(test.host)
			if err != nil {
				t.Fatal(err)
			}

			idx = newMemoryIndex()
			crawl(&idx, url, false)
			if diff := deep.Equal(idx.search(test.name), test.expected); diff != nil {
				t.Error(diff)
			}

		})
	}

}
