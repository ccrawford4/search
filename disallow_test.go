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
				UrlMap{
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap11.html": UrlEntry{
						1959,
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap01.html": UrlEntry{
						2208,
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap12.html": UrlEntry{
						2203,
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap08.html": UrlEntry{
						2572,
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap02.html": UrlEntry{
						2199,
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/": UrlEntry{
						3293,
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap10.html": UrlEntry{
						2114,
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap09.html": UrlEntry{
						2379,
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap07.html": UrlEntry{
						2400,
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap05.html": UrlEntry{
						2256,
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap03.html": UrlEntry{
						1752,
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap06.html": UrlEntry{
						2682,
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
					},
				},
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
				true,
			},
		},
		{
			"jekyll",
			"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/",
			&SearchResult{
				UrlMap{
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/": UrlEntry{
						3281,
						"The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap06.html": UrlEntry{
						1527,
						"The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap07.html": UrlEntry{
						580,
						"The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap08.html": UrlEntry{
						4463,
						"The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap09.html": UrlEntry{
						2845,
						"The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap10.html": UrlEntry{
						6995,
						"The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson",
						"",
					},
				},
				Frequency{
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/":            8,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap06.html": 8,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap07.html": 6,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap08.html": 14,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap09.html": 12,
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap10.html": 30,
				},
				6,
				true,
			},
		},
		{
			"Nicolo",
			"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/",
			&SearchResult{
				UrlMap{
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/": UrlEntry{
						3585,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap00.html": UrlEntry{
						12,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap01.html": UrlEntry{
						144,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap02.html": UrlEntry{
						266,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap03.html": UrlEntry{
						3119,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap04.html": UrlEntry{
						973,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap05.html": UrlEntry{
						467,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap06.html": UrlEntry{
						1225,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap07.html": UrlEntry{
						2763,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap08.html": UrlEntry{
						1569,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap09.html": UrlEntry{
						1324,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap10.html": UrlEntry{
						730,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap11.html": UrlEntry{
						899,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap12.html": UrlEntry{
						1993,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap13.html": UrlEntry{
						1407,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap14.html": UrlEntry{
						902,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap15.html": UrlEntry{
						530,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap16.html": UrlEntry{
						822,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap17.html": UrlEntry{
						1097,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap18.html": UrlEntry{
						1289,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
					"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli/chap19.html": UrlEntry{
						3604,
						"The Project Gutenberg eBook of The Prince, by Nicolo Machiavelli",
						"",
					},
				},
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
				true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var idx Index
			idx = newMemoryIndex()
			crawl(&idx, test.host, false)
			if diff := deep.Equal(idx.search(test.name), test.expected); diff != nil {
				t.Error(diff)
			}

		})
	}

}
