package main

import (
	"io"
	"log"
	"net/http"
	"sync"
)

// download takes in a URL, makes an HTTP GET request and returns the data as an array of bytes
func download(url string, wg *sync.WaitGroup, ch chan Download) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		ch <- Download{nil, "", err}
		log.Printf("ERROR Downloading %q: %v\n", url, err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v\n", err)
		}
	}(resp.Body)

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- Download{nil, "", err}
		return
	}
	ch <- Download{
		bytes,
		url,
		nil,
	}
}
