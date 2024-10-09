package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

type CrawlResult struct {
	TermFrequency Frequency // frequency per word
	Url           string    // the url crawled
	TotalWords    int       // the total # of words in the document
}

// parseURL takes in a rawURL or href and returns a new pointer to an url.URL object
func parseURL(rawURL string) (*url.URL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("error parsing URL %s: %v", rawURL, err)
	}
	return parsedURL, err
}

// validUrl returns true if url has not been crawled already and is within the host domain
func validUrl(index *Index, newUrl, host string) bool {
	if (*index).containsUrl(newUrl) || !strings.HasPrefix(newUrl, host) {
		return false
	}
	return true
}

func crawl(index *Index, seedUrl *url.URL) {
	before := time.Now()
	initialFullPath, err := clean(seedUrl, seedUrl.Path) // keep track of the initial full path
	if err != nil {
		log.Printf("error crawling seedUrl %s: %v", seedUrl, err)
		return
	}
	queue := []string{seedUrl.Path} // Queue for keeping track of hrefs to visit

	for len(queue) > 0 {
		// Pop the last href off the queue
		href := queue[0]
		queue = queue[1:]

		// Clean the URL and check to see if it's a valid URL
		cleanedUrl, err := clean(seedUrl, href)
		if err != nil {
			log.Printf("error crawling url %s: %v", seedUrl, err)
			continue
		}
		if !validUrl(index, cleanedUrl, initialFullPath) {
			continue
		}

		/// Download the contents of the page and add the cleanedURL to the visited set
		body := download(cleanedUrl)
		words, hrefs := extract(body)
		fmt.Printf("Download: url=%s result=ok\n", cleanedUrl)

		// Create the wordFrequency map
		wordFreq := make(Frequency)
		for _, word := range words {
			stemmedWord := (*index).getStemmedWord(word)
			wordFreq[stemmedWord] += 1
		}

		(*index).insertCrawlResults(&CrawlResult{
			wordFreq,
			cleanedUrl,
			len(words),
		})

		queue = append(queue, hrefs...)
	}
	duration := time.Now().Sub(before)
	log.Printf("Crawl Duration: %v\n", duration.String())
}
