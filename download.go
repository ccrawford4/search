package main

import (
	"io"
	"log"
	"net/http"
)

// download takes in a URL, makes an HTTP GET request and returns the data as an array of bytes
func download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v\n", err)
		}
	}(resp.Body)

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}

	return bytes, nil
}
