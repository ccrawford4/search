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
type index map[string]frequency

// parseURL takes in a rawURL or href and returns a new pointer to an url.URL object
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

// validURL takes in a visited *hashset.Set containing all the cleaned URLs that have been visited,
// newURL which is the newly constructed URL from the host and href, and host which is the host domain as a string
// validURL returns a true if the newURL has not been visited already and is within the host's domain, and false if
// either of these requirements are not met
func validURL(visited *hashset.Set, newURL, host string) bool {
	urlObj := parseURL(newURL)
	if visited.Contains(newURL) || urlObj.Host != host {
		return false
	}
	return true
}

// crawl takes in a hostURL object, crawls all the pages within the subdomain, and returns an index object
// containing all the stemmed words and their corresponding counts organized by URL
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

		/// Download the contents of the page and add the cleanedURL to the visited set
		fmt.Printf("Download: url=%s result=ok\n", cleanedURL)
		body := download(cleanedURL)
		visited.Add(cleanedURL)

		// Populate the index and queue with the relevant data from the document
		words, hrefs := extract(body)
		populateIndex(resultIndex, words, cleanedURL)
		queue = append(queue, hrefs...)
	}

	return resultIndex
}
