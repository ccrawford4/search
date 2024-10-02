package main

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"log"
	"net/url"
	"strings"
)

// parseURL takes in a rawURL or href and returns a new pointer to an url.URL object
func parseURL(rawURL string) *url.URL {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("error parsing URL %s: %v", rawURL, err)
	}
	return parsedURL
}

func validURL(visited *hashset.Set, newURL, host string) bool {
	if visited.Contains(newURL) || !strings.HasPrefix(newURL, host) {
		return false
	}
	return true
}

func populateIndexValues(index *Index, URL string, frequency *Frequency) {
	for word, value := range *frequency {
		_, found := (*index)[word]
		if !found {
			(*index)[word] = make(Frequency)
		}
		(*index)[word][URL] = value
	}
}

// crawl takes in a pointer
func crawl(index *Index, wordsInDoc *Frequency, hostURL *url.URL, stopWords *hashset.Set) {
	initialFullPath := clean(hostURL, hostURL.Path)
	queue := []string{hostURL.Path} // Queue for keeping track of hrefs to visit
	visited := hashset.New()        // Visited hashset to keep track of URLs crawled

	for len(queue) > 0 {
		// Pop the last href off the queue
		href := queue[0]
		queue = queue[1:]

		// Clean the URL and verify it hasn't been visited before, or is outside the host domain,
		// otherwise skip the processing step
		cleanedURL := clean(hostURL, href)
		if !validURL(visited, cleanedURL, initialFullPath) {
			continue
		}

		/// Download the contents of the page and add the cleanedURL to the visited set
		fmt.Printf("Download: url=%s result=ok\n", cleanedURL)
		body := download(cleanedURL)
		visited.Add(cleanedURL)

		// Populate the index and queue with the relevant data from the document
		words, hrefs := extract(body)
		(*wordsInDoc)[cleanedURL] += len(words)
		wordFrequency := createWordFrequency(words, stopWords)
		populateIndexValues(index, cleanedURL, &wordFrequency)
		queue = append(queue, hrefs...)
	}
}
