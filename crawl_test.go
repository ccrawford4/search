package main

import (
	"fmt"
	"github.com/go-test/deep"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCrawl(t *testing.T) {
	simpleDoc := []byte("<html><body>Hello CS 272, there are no links here.</body></html>")
	hrefDoc := []byte("<html><body>For a simple example, see <a href=\"/tests/project01/simple.html\">simple.html</a></body></html>")
	styleDoc := []byte("\n<html>\n<head>\n  <title>Style</title>\n  <style>\n    a.blue {\n      color: blue;\n    }\n    a.red {\n      color: red;\n    }\n  </style>\n<body>\n  <p>\n    Here is a blue link to <a class=\"blue\" href=\"/tests/project01/href.html\">href.html</a>\n  </p>\n  <p>\n    And a red link to <a class=\"red\" href=\"/tests/project01/simple.html\">simple.html</a>\n  </p>\n</body>\n</html>")
	repeatDoc := []byte("<html><body><a href=\"/repeat-href\"></a><a href=\"/repeat-href\"></a></body></html>")

	simpleWords, _, _, _ := extract(simpleDoc)
	hrefWords, _, _, _ := extract(hrefDoc)
	hrefWords = append(hrefWords, simpleWords...)
	styleWords, _, _, _ := extract(styleDoc)
	styleWords = append(styleWords, hrefWords...)

	tests := []struct {
		name                         string
		expectedHrefs, expectedWords []string
		serverContent                map[string][]byte
		numDocs                      float64
		index                        Index
		expectedIndex                *MemoryIndex
	}{
		{
			"simple",
			[]string{"/"},
			simpleWords,
			map[string][]byte{
				"/": simpleDoc,
			},
			1,
			newMemoryIndex(),
			&MemoryIndex{
				UrlMap{
					"http://127.0.0.1:8080/": UrlEntry{8,
						"",
						"",
					},
				},
				IndexMap{
					"there": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"are": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"no": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"link": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"here": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"hello": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"cs": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"272": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
				},
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
			2,
			newMemoryIndex(),
			&MemoryIndex{
				UrlMap{
					"http://127.0.0.1:8080/": UrlEntry{
						7,
						"",
						"",
					},
					"http://127.0.0.1:8080/tests/project01/simple.html": UrlEntry{
						8,
						"",
						"",
					},
				},
				IndexMap{
					"there": Frequency{
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"are": Frequency{
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"no": Frequency{
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"link": Frequency{
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"here": Frequency{
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"hello": Frequency{
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"cs": Frequency{
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"272": Frequency{
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"a": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"see": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"simpl": Frequency{
						"http://127.0.0.1:8080/": 2,
					},
					"exampl": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"for": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"html": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
				},
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
			3,
			newMemoryIndex(),
			&MemoryIndex{
				UrlMap{
					"http://127.0.0.1:8080/": UrlEntry{
						16,
						"Style",
						"",
					},
					"http://127.0.0.1:8080/tests/project01/href.html": UrlEntry{
						7,
						"",
						"",
					},
					"http://127.0.0.1:8080/tests/project01/simple.html": UrlEntry{
						8,
						"",
						"",
					},
				},
				IndexMap{
					"there": Frequency{
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"are": Frequency{
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"no": Frequency{
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"link": Frequency{
						"http://127.0.0.1:8080/":                             2,
						"http://127.0.0.1:56038/tests/project01/simple.html": 1,
					},
					"here": Frequency{
						"http://127.0.0.1:8080/":                            1,
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"hello": Frequency{
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"cs": Frequency{
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"272": Frequency{
						"http://127.0.0.1:8080/tests/project01/simple.html": 1,
					},
					"a": Frequency{
						"http://127.0.0.1:8080/":                          2,
						"http://127.0.0.1:8080/tests/project01/href.html": 1,
					},
					"see": Frequency{
						"http://127.0.0.1:8080/tests/project01/href.html": 1,
					},
					"simpl": Frequency{
						"http://127.0.0.1:8080/":                          1,
						"http://127.0.0.1:8080/tests/project01/href.html": 2,
					},
					"exampl": Frequency{
						"http://127.0.0.1:8080/tests/project01/href.html": 1,
					},
					"for": Frequency{
						"http://127.0.0.1:8080/tests/project01/href.html": 1,
					},
					"html": Frequency{
						"http://127.0.0.1:8080/":                          2,
						"http://127.0.0.1:8080/tests/project01/href.html": 1,
					},
					"is": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"to": Frequency{
						"http://127.0.0.1:8080/": 2,
					},
					"and": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"red": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"blue": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"href": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
					"style": Frequency{
						"http://127.0.0.1:8080/": 1,
					},
				},
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
				"/":            repeatDoc,
				"/repeat-href": repeatDoc,
			},
			2,
			newMemoryIndex(),

			&MemoryIndex{
				UrlMap{
					"http://127.0.0.1:8080/": UrlEntry{
						0,
						"",
						"",
					},
					"http://127.0.0.1:8080/repeat-href": UrlEntry{
						0,
						"",
						"",
					},
				},
				IndexMap{},
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
			1,
			newMemoryIndex(),
			&MemoryIndex{
				UrlMap{
					"http://127.0.0.1:8080/": UrlEntry{
						0,
						"",
						"",
					},
				},
				IndexMap{},
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

			svrURL, err := parseURL(ts.URL)
			if err != nil {
				log.Fatalf("Error parsing URL: %v\n", err.Error())
			}
			crawl(&test.index, ts.URL, true)

			hostParts := strings.Split(svrURL.Host, ":")
			mockSVRPort := hostParts[len(hostParts)-1]

			// Iterate over the map and update the keys with the new port number
			for oldURL := range test.expectedIndex.WordCount {
				// Parse the URL
				u, err := url.Parse(oldURL)
				if err != nil {
					fmt.Printf("Error parsing URL: %v\n", err)
					continue
				}

				// Split the host to update the port
				hostParts = strings.Split(u.Host, ":")
				if len(hostParts) > 1 {
					// Continue if you already have the correct server port
					if hostParts[1] == mockSVRPort {
						continue
					}
					// Replace the old port with the new port number
					u.Host = hostParts[0] + ":" + mockSVRPort
				}

				// Insert the new URL into the map
				test.expectedIndex.WordCount[u.String()] = test.expectedIndex.WordCount[oldURL]

				// Delete the old key
				delete(test.expectedIndex.WordCount, oldURL)
			}

			for term := range test.expectedIndex.Index {
				for oldURL := range test.expectedIndex.Index[term] {
					// Parse the URL
					u, err := url.Parse(oldURL)
					if err != nil {
						fmt.Printf("Error parsing URL: %v\n", err)
						continue
					}

					// Split the host to update the port
					hostParts = strings.Split(u.Host, ":")
					if len(hostParts) > 1 {
						// Continue if you already have the correct server port
						if hostParts[1] == mockSVRPort {
							continue
						}
						// Replace the old port with the new port number
						u.Host = hostParts[0] + ":" + mockSVRPort
					}

					// Insert the new URL into the map
					test.expectedIndex.Index[term][u.String()] = test.expectedIndex.Index[term][oldURL]

					// Delete the old key
					delete(test.expectedIndex.Index[term], oldURL)
				}
			}

			if diff := deep.Equal(test.index, test.expectedIndex); diff != nil {
				t.Error(diff)
			}
		})
	}
}
