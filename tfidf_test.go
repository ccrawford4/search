package main

import (
	"fmt"
	"github.com/go-test/deep"
	"log"
	"net/http"
	"testing"
)

func TestTfIdf(t *testing.T) {
	tests := []struct {
		expectedTemplateData TemplateData
	}{
		{
			TemplateData{
				Hits{
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap10.html",
						0.10530369316483629,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap09.html",
						0.08700977324729901,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap12.html",
						0.006118026177740771,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/index.html",
						0.0026074832313376555,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of The Iliad of Homer/chap15.html",
						0.0009804136949829584,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of The Iliad of Homer/illus46.html",
						0.0009804136949829584,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of The Iliad of Homer/illus47.html",
						0.0009804136949829584,
					},
				},
				"turtle",
			},
		},
		{
			TemplateData{
				Hits{
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap08.html",
						0.002080903332378551,
					},
					{
						"http://127.0.0.1:8080/documents/top11/Dracula | Project Gutenberg/chap11.html",
						0.0018178188288181252,
					},
				},
				"monkey",
			},
		},
		{
			TemplateData{},
		},
	}

	go func() {
		port := 8080
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			log.Fatalf("Error starting server: %v", err.Error())
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		if urlPath == "/" {
			urlPath = "/documents/top11/index.html"
		}
		filePath := "." + urlPath
		fileContent, err := openAndReadFile(filePath)
		if err != nil {
			_, err = w.Write([]byte("404 No Page found!"))
			return
		}
		_, err = w.Write(fileContent)
		if err != nil {
			log.Fatalf("Error writing response: %v", err.Error())
		}
	})

	for _, test := range tests {
		t.Run(test.expectedTemplateData.TERM, func(t *testing.T) {
			index := make(Index)
			stopWords := getStopWords()
			crawl(&index, parseURL("http://127.0.0.1:8080"), stopWords)

			freq, _ := search(&index, test.expectedTemplateData.TERM, stopWords)
			got := getTemplateData(freq, test.expectedTemplateData.TERM)

			if diff := deep.Equal(*got, test.expectedTemplateData); diff != nil {
				t.Error(diff)
			}
		})
	}
}