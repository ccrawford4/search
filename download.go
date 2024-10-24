package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

// download takes in a URL, makes an HTTP GET request and returns the data as an array of bytes
func download(url string, wg *sync.WaitGroup, ch chan Download) {
	defer wg.Done()

	// Download the url
	resp, err := http.Get(url)

	// Handle any potential errors
	if err != nil {
		log.Printf("[WARNING] ERROR Downloading %q: %v\n", url, err)
		ch <- Download{nil, url, err}
		return
	}

	// Handle status code errors (forbidden, 404, internal server errors, etc)
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("[WARNING] Could not download %q, status code: %d", url, resp.StatusCode)
		log.Printf("%v\n", err)
		ch <- Download{nil, url, err}
		return
	}

	// Read in the bytes
	bytes, err := io.ReadAll(resp.Body)

	// Handle parsing errors
	if err != nil {
		log.Printf("[WARNING] ERROR Parsing Url %q: %v\n", url, err)
		ch <- Download{nil, url, err}
		return
	}

	// Put the bytes into the channel
	ch <- Download{
		bytes,
		url,
		nil,
	}
}
