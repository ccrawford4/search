package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

func TestDownload(t *testing.T) {
	var tests = []struct {
		name         string
		url          string
		expectedBody []byte
	}{
		{
			"simple",
			"https://cs272-f24.github.io/tests/project01/simple.html",
			[]byte("<html><body>Hello CS 272, there are no links here.</body></html>"),
		},
		{
			"href",
			"https://cs272-f24.github.io/tests/project01/href.html",
			[]byte("<html><body>For a simple example, see <a href=\"/tests/project01/simple.html\">simple.html</a></body></html>"),
		},
		{
			"style",
			"https://cs272-f24.github.io/tests/project01/style.html",
			[]byte("<html><head>  <title>Style</title>  <style>    a.blue {     color: blue;    }    a.red {      color: red;    }  </style><body>  <p>    Here is a blue link to <a class=\"blue\" href=\"/tests/project01/href.html\">href.html</a>  </p>  <p>    And a red link to <a class=\"red\" href=\"/tests/project01/simple.html\">simple.html</a>  </p></body></html>"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err := w.Write(test.expectedBody)
				if err != nil {
					log.Fatalf("Error writing response: %v", err.Error())
				}
			}))
			defer ts.Close()

			var wg sync.WaitGroup
			downloadChannel := make(chan Download)
			wg.Add(1)
			go download(ts.URL, &wg, downloadChannel)

			go func() {
				wg.Wait()
				close(downloadChannel)
			}()

			downloadObj, ok := <-downloadChannel
			if !ok {
				t.Errorf("Could not receive download Obj from channel\n")
			}
			if !reflect.DeepEqual(test.expectedBody, downloadObj.Body) {
				t.Errorf("Expected %s, got %s", test.expectedBody, downloadObj.Body)
			}
		})
	}

}
