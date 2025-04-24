package main

import (
	"testing"

	"github.com/go-test/deep"
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
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap10.html",
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
						0.021404760590749916,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap09.html",
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
						0.01717976635463083,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap12.html",
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
						0.0013251625756327977,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/index.html",
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
						0.0004432634609959085,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Iliad of Homer/chap15.html",
						"The Project Gutenberg eBook of The Iliad of Homer",
						"",
						0.00019274614776964563,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Iliad of Homer/illus46.html",
						"The Project Gutenberg eBook of The Iliad of Homer",
						"",
						0.00019274614776964563,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Iliad of Homer/illus47.html",
						"The Project Gutenberg eBook of The Iliad of Homer",
						"",
						0.00019274614776964563,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Iliad of Homer/chap11.html",
						"The Project Gutenberg eBook of The Iliad of Homer",
						"",
						0.00017710101638674188,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Iliad of Homer/illus37.html",
						"The Project Gutenberg eBook of The Iliad of Homer",
						"",
						0.00017710101638674188,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Iliad of Homer/illus38.html",
						"The Project Gutenberg eBook of The Iliad of Homer",
						"",
						0.00017710101638674188,
					},
				},
				"turtle",
			},
			0,
		},
		{
			TemplateData{
				Hits{
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap04.html",
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
						0.008269737644028363,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap11.html",
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
						0.0058065462630735145,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap01.html",
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
						0.00579569843547606,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap12.html",
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
						0.005163424479964146,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap08.html",
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
						0.0033169782647825665,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap02.html",
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
						0.0025864083968533454,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/index.html",
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
						0.0008635760802733841,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll/chap10.html",
						"The Project Gutenberg eBook of Alice’s Adventures in Wonderland, by Lewis Carroll",
						"",
						0.0006726007645081016,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg EBook of A Tale of Two Cities, by Charles Dickens/link2H_4_0014.html",
						"The Project Gutenberg EBook of A Tale of Two Cities, by Charles Dickens",
						"",
						0.0003039499820799758,
					},
					{
						"http://127.0.0.1:8080/documents/top10/Dracula | Project Gutenberg/chap19.html",
						"Dracula | Project Gutenberg",
						"",
						0.0002468538222517581,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg EBook of A Tale of Two Cities, by Charles Dickens/link2H_4_0043.html",
						"The Project Gutenberg EBook of A Tale of Two Cities, by Charles Dickens",
						"",
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
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg EBook of A Tale of Two Cities, by Charles Dickens/link2H_4_0009.html",
						"The Project Gutenberg EBook of A Tale of Two Cities, by Charles Dickens",
						"",
						0.0007550692926003138,
					},
					{
						"http://127.0.0.1:8080/documents/top10/The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson/chap08.html",
						"The Project Gutenberg eBook of The Strange Case Of Dr. Jekyll And Mr. Hyde, by Robert Louis Stevenson",
						"",
						0.0004254983802128141,
					},
					{
						"http://127.0.0.1:8080/documents/top10/Dracula | Project Gutenberg/chap11.html",
						"Dracula | Project Gutenberg",
						"",
						0.00035796404729308,
					},
				},
				"monkey",
			},
			0,
		},
	}

	var dbIdx Index
	dbIdx = newDBIndex("test.db", true, nil)
	crawl(&dbIdx, "http://127.0.0.1:8080/documents/top10/", true)

	var memIdx Index
	memIdx = newMemoryIndex()
	crawl(&memIdx, "http://127.0.0.1:8080/documents/top10/", true)

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
