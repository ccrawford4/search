package main

import (
	"io"
	"log"
	"net/http"
)

// download takes in a URL, makes an HTTP GET request and returns the data as an array of bytes
func download(url string) []byte {
	resp, err := http.Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v\n", err)
		}
	}(resp.Body)

	if err != nil {
		log.Fatalf("Error downloading file: %v\n", err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}
	return bytes
}
