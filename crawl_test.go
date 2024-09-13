package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCrawl(t *testing.T) {
	simpleDoc := []byte("<html><body>Hello CS 272, there are no links here.</body></html>")
	hrefDoc := []byte("<html><body>For a simple example, see <a href=\"/tests/project01/simple.html\">simple.html</a></body></html>")
	styleDoc := []byte("\n<html>\n<head>\n  <title>Style</title>\n  <style>\n    a.blue {\n      color: blue;\n    }\n    a.red {\n      color: red;\n    }\n  </style>\n<body>\n  <p>\n    Here is a blue link to <a class=\"blue\" href=\"/tests/project01/href.html\">href.html</a>\n  </p>\n  <p>\n    And a red link to <a class=\"red\" href=\"/tests/project01/simple.html\">simple.html</a>\n  </p>\n</body>\n</html>")
	repeatDoc := []byte("<html><body><a href=\"/repeat-href\"></a><a href=\"/repeat-href\"></a></body></html>")

	simpleWords, _ := extract(simpleDoc)
	hrefWords, _ := extract(hrefDoc)
	hrefWords = append(hrefWords, simpleWords...)
	styleWords, _ := extract(styleDoc)
	styleWords = append(styleWords, hrefWords...)

	tests := []struct {
		name                         string
		expectedHrefs, expectedWords []string
		serverContent                map[string][]byte
	}{
		{
			"simple",
			[]string{"/"},
			simpleWords,
			map[string][]byte{
				"/": simpleDoc,
			},
		},
		{
			"href",
			[]string{
				"/",
				"/tests/project01/simple.html",
			},
			hrefWords,
			map[string][]byte{
				"/":                            hrefDoc,
				"/tests/project01/simple.html": simpleDoc,
			},
		},
		{
			"style",
			[]string{
				"/",
				"/tests/project01/href.html",
				"/tests/project01/simple.html",
			},
			styleWords,
			map[string][]byte{
				"/":                            styleDoc,
				"/tests/project01/href.html":   hrefDoc,
				"/tests/project01/simple.html": simpleDoc,
			},
		},
		{
			"repeat-href",
			[]string{
				"/",
				"repeat-href",
				"repeat-href",
			},
			styleWords,
			map[string][]byte{
				"/":           repeatDoc,
				"repeat-href": repeatDoc,
			},
		},
		{
			"outside-domain",
			[]string{
				"/",
			},
			nil,
			map[string][]byte{
				"/": []byte("<html><body><a href=\"https://wikipedia.org\"></a></body></html>"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write(test.serverContent[r.URL.Path])
				if err != nil {
					log.Fatalf("Error writing response: %v\n", err.Error())
				}
			}))
			fmt.Printf("ts: %v\n", ts.URL)
			defer ts.Close()

			resultIndex := crawl(parseURL(ts.URL))
			expectedIndex := make(index)

			for path, doc := range test.serverContent {
				fullURL := clean(parseURL(ts.URL), path)
				words, _ := extract(doc)
				populateIndex(expectedIndex, words, fullURL)
			}

			if !reflect.DeepEqual(expectedIndex, resultIndex) {
				t.Errorf("expected: %v\n, got: %v\n", expectedIndex, resultIndex)
			}
		})
	}
}
