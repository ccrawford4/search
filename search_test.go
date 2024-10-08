package main

import (
	"github.com/go-test/deep"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name, path    string
		pathFrequency Frequency
		indexType     IndexType
		expected      SearchResult
	}{
		{
			"Verona",
			"/tests/rnj/sceneI_30.0.html",
			Frequency{
				"/tests/rnj/sceneI_30.0.html": 1,
			},
			0,
			SearchResult{
				TotalDocsSearched: 1,
				Found:             true,
			},
		},
		{
			"Benvolio",
			"/tests/rnj/sceneI_30.1.html",
			Frequency{
				"/tests/rnj/sceneI_30.1.html": 26,
			},
			0,
			SearchResult{
				TotalDocsSearched: 1,
				Found:             true,
			},
		},
		{
			"Romeo",
			"/tests/rnj/",
			Frequency{
				"/tests/rnj/sceneI_30.0.html":  2,
				"/tests/rnj/sceneI_30.1.html":  22,
				"/tests/rnj/sceneI_30.3.html":  2,
				"/tests/rnj/sceneI_30.4.html":  17,
				"/tests/rnj/sceneI_30.5.html":  15,
				"/tests/rnj/sceneII_30.2.html": 42,
				"/tests/rnj/":                  200,
				"/tests/rnj/sceneI_30.2.html":  15,
				"/tests/rnj/sceneII_30.0.html": 3,
				"/tests/rnj/sceneII_30.1.html": 10,
				"/tests/rnj/sceneII_30.3.html": 13,
			},
			0,
			SearchResult{
				TotalDocsSearched: 11,
				Found:             true,
			},
		},
		{
			"Verona",
			"/tests/rnj/sceneI_30.0.html",
			Frequency{
				"/tests/rnj/sceneI_30.0.html": 1,
			},
			1,
			SearchResult{
				TotalDocsSearched: 1,
				Found:             true,
			},
		},
		{
			"Benvolio",
			"/tests/rnj/sceneI_30.1.html",
			Frequency{
				"/tests/rnj/sceneI_30.1.html": 26,
			},
			1,
			SearchResult{
				TotalDocsSearched: 1,
				Found:             true,
			},
		},
		{
			"Romeo",
			"/tests/rnj/",
			Frequency{
				"/tests/rnj/sceneI_30.0.html":  2,
				"/tests/rnj/sceneI_30.1.html":  22,
				"/tests/rnj/sceneI_30.3.html":  2,
				"/tests/rnj/sceneI_30.4.html":  17,
				"/tests/rnj/sceneI_30.5.html":  15,
				"/tests/rnj/sceneII_30.2.html": 42,
				"/tests/rnj/":                  200,
				"/tests/rnj/sceneI_30.2.html":  15,
				"/tests/rnj/sceneII_30.0.html": 3,
				"/tests/rnj/sceneII_30.1.html": 10,
				"/tests/rnj/sceneII_30.3.html": 13,
			},
			1,
			SearchResult{
				TotalDocsSearched: 11,
				Found:             true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				urlPath := r.URL.Path
				if urlPath == "/tests/rnj/" {
					urlPath += "/index.html"
				}
				filePath := "./documents" + urlPath
				reader, err := os.Open(filePath)
				if err != nil {
					t.Fatalf("Could not open file %q\n", filePath)
				}

				bytes, err := io.ReadAll(reader)
				_, err = w.Write(bytes)
				if err != nil {
					log.Fatalf("Error writing response: %v", err.Error())
				}
			}))
			defer ts.Close()

			hostURL, err := parseURL(ts.URL)
			if err != nil {
				t.Fatalf("Could not parse URL: %v\n", ts.URL)
			}
			testURL, err := clean(hostURL, test.path)
			if err != nil {
				t.Fatalf("Could not clean URL: %v\n", test.path)
			}
			var index Index
			if test.indexType == Memory {
				index = newMemoryIndex()
			} else {
				index = newDBIndex("test.db", true)
			}
			url, err := parseURL(testURL)
			if err != nil {
				t.Fatalf("Could not parse URL: %v\n", testURL)
			}
			crawl(&index, url)
			got := index.search(test.name)

			expectedTermFrequency := make(Frequency)
			for path, freq := range test.pathFrequency {
				cleanedUrl, err := clean(hostURL, path)
				if err != nil {
					t.Fatalf("Could not clean URL: %v\n", path)
				}
				expectedTermFrequency[cleanedUrl] += freq
			}
			test.expected.TermFrequency = expectedTermFrequency
			dropDatabase("test.db")

			if diff := deep.Equal(got, &test.expected); diff != nil {
				t.Error(diff)
			}

		})
	}
}
