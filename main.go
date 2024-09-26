package main

import (
	"time"
)

func main() {
	go startServer("8080") // For serving the contents of the web page

	index := make(Index)
	stopWords := getStopWords()
	initHandlers(&index, stopWords)
	crawl(&index, parseURL("http://127.0.0.1:8080/documents/top11/"), stopWords)
	// Loop indefinitely
	for {
		time.Sleep(100 * time.Millisecond)
	}
}
