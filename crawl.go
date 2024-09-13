package main

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"log"
	"net/url"
)

// frequency map[url]int
type frequency map[string]int

// index map[word]frequency
// index map[word][url]int
type index map[string]frequency

// parseURL takes in a rawURL or href and returns it as a *url.URL object
func parseURL(rawURL string) *url.URL {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("error parsing URL %s: %v", rawURL, err)
	}
	return parsedURL
}

// populateIndex takes in a slice of words and the sourceURL for those words and adds the data to the index
func populateIndex(indexMap index, words []string, sourceURL string) {
	for _, word := range words {
		stemmedWord := getStemmedWord(word) // Convert the word to its stemmed version
		_, found := indexMap[stemmedWord]
		if !found {
			indexMap[stemmedWord] = make(frequency)
		}
		indexMap[stemmedWord][sourceURL]++
	}
}

func validURL(visited *hashset.Set, rawURL, host string) bool {
	urlObj := parseURL(rawURL)
	if visited.Contains(rawURL) || urlObj.Host == host {
		return false
	}
	return true
}

// Crawl takes in a hostURL object and returns an index object
func crawl(hostURL *url.URL) index {
	queue := []string{hostURL.Path} // Queue for keeping track of hrefs to visit
	resultIndex := make(index)      // Index object to populate and return
	visited := hashset.New()        // Hashset to keep track of visited URLs

	for len(queue) > 0 {
		// Pop the last href off the queue
		href := queue[0]
		queue = queue[1:]

		// Clean the URL and verify it hasn't been visited before, or is outside the host domain,
		// otherwise skip the processing step
		cleanedURL := clean(hostURL, href)
		if !validURL(visited, cleanedURL, hostURL.Host) {
			continue
		}

		fmt.Printf("Download: url=%s result=ok\n", cleanedURL)
		body := download(cleanedURL)
		visited.Add(cleanedURL)

		// Populate the index and queue with the words and hrefs from the document
		words, hrefs := extract(body)
		populateIndex(resultIndex, words, cleanedURL)
		queue = append(queue, hrefs...)
	}

	return resultIndex
}
