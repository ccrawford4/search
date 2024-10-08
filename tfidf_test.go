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
		indexType            IndexType
	}{
		{
			TemplateData{
				Hits{
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap10.html",
						0.02192105439051487,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap09.html",
						0.018807540860160833,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap12.html",
						0.001450721085089249,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/index.html",
						0.00048526245831333377,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of The Iliad of Homer/chap15.html",
						0.00021100875151535821,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of The Iliad of Homer/illus46.html",
						0.00021100875151535821,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of The Iliad of Homer/illus47.html",
						0.00021100875151535821,
					},
				},
				"turtle",
			},
			0,
		}, {
			TemplateData{
				Hits{
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap04.html",
						0.007752879041276591,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap11.html",
						0.0058065462630735145,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap01.html",
						0.00579569843547606,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap12.html",
						0.005163424479964146,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap08.html",
						0.0033169782647825665,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap02.html",
						0.0025864083968533454,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/index.html",
						0.0008635760802733841,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap10.html",
						0.0006726007645081016,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg EBook of A Tale of Two Cities, by Charles Dickens/link2H_4_0014.html",
						0.0003039499820799758,
					},
					{
						"http://127.0.0.1:8080/documents/top11/Dracula | Project Gutenberg/chap19.html",
						0.0002468538222517581,
					},
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg EBook of A Tale of Two Cities, by Charles Dickens/link2H_4_0043.html",
						0.00024264129968773495,
					},
				},
				"rabbit",
			},
			0,
		},
		{
			TemplateData{
				Hits{
					{
						"http://127.0.0.1:8080/documents/top11/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap08.html",
						0.0004534927195828118,
					},
					{
						"http://127.0.0.1:8080/documents/top11/Dracula | Project Gutenberg/chap11.html",
						0.00038151517577720813,
					},
				},
				"monkey",
			},
			0,
		},
	}

	go func() {
		port := 8080
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			log.Fatalf("Error starting server: %v", err.Error())
		}
	}()

	http.HandleFunc("/documents/top11/", func(w http.ResponseWriter, r *http.Request) {
		corpusHandler(&w, r)
	})

	var dbIdx Index
	dbIdx = newDBIndex("test.db", true)
	url, err := parseURL("http://127.0.0.1:8080/documents/top11/")
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err.Error())
	}
	crawl(&dbIdx, url)

	var memIdx Index
	memIdx = newMemoryIndex()
	url, err = parseURL("http://127.0.0.1:8080/documents/top11/")
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err.Error())
	}
	crawl(&memIdx, url)

	for _, test := range tests {
		var idx Index
		if test.indexType == Memory {
			idx = memIdx
		} else {
			idx = dbIdx
		}

		t.Run(test.expectedTemplateData.TERM, func(t *testing.T) {
			got := getTemplateData(&idx, test.expectedTemplateData.TERM)
			if diff := deep.Equal(*got, test.expectedTemplateData); diff != nil {
				t.Error(diff)
			}
		})
	}
	dropDatabase("test.db")
}
