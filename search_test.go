package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name, path    string
		pathFrequency frequency
	}{
		{
			"Verona",
			"/tests/rnj/sceneI_30.0.html",
			frequency{
				"/tests/rnj/sceneI_30.0.html": 1,
			},
		},
		{
			"Benvolio",
			"/tests/rnj/sceneI_30.1.html",
			frequency{
				"/tests/rnj/sceneI_30.1.html": 26,
			},
		},
		{
			"Romeo",
			"/tests/rnj/",
			frequency{
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

			hostURL := parseURL(ts.URL)
			testURL := clean(hostURL, test.path)
			got := search(testURL, test.name)

			expected := make(frequency)
			for path, freq := range test.pathFrequency {
				expected[clean(hostURL, path)] += freq
			}

			if !reflect.DeepEqual(got, expected) {
				t.Errorf("got %v\n, want %v\n", got, expected)
			}
		})
	}
}
