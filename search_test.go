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
		pathFrequency Frequency[float64]
	}{
		{
			"Verona",
			"/tests/rnj/sceneI_30.0.html",
			Frequency[float64]{
				"/tests/rnj/sceneI_30.0.html": -0.00316873679646296,
			},
		},
		{
			"Benvolio",
			"/tests/rnj/sceneI_30.1.html",
			Frequency[float64]{
				"/tests/rnj/sceneI_30.1.html": -0.011392692703440337,
			},
		},
		{
			"Romeo",
			"/tests/rnj/",
			Frequency[float64]{
				"/tests/rnj/sceneI_30.0.html":  -0.0007955486503031538,
				"/tests/rnj/sceneI_30.1.html":  -0.0012101140313927157,
				"/tests/rnj/sceneI_30.3.html":  -0.00019478639633711237,
				"/tests/rnj/sceneI_30.4.html":  -0.0013581512370397391,
				"/tests/rnj/sceneI_30.5.html":  -0.0011544366870488737,
				"/tests/rnj/sceneII_30.2.html": -0.002859674878116742,
				"/tests/rnj/":                  -0.0024192420543789886,
				"/tests/rnj/sceneI_30.2.html":  -0.0014065221174714565,
				"/tests/rnj/sceneII_30.0.html": -0.0012060179007255256,
				"/tests/rnj/sceneII_30.1.html": -0.001749470411546287,
				"/tests/rnj/sceneII_30.3.html": -0.0012312062445167854,
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

			index := make(Index)
			stopWords := getStopWords()
			crawl(&index, parseURL(testURL), stopWords)
			got, _ := search(&index, test.name, stopWords)

			expected := make(Frequency[float64])
			for path, freq := range test.pathFrequency {
				expected[clean(hostURL, path)] += freq
			}

			if diff := deep.Equal(got, expected); diff != nil {
				t.Error(diff)
			}
		})
	}
}
