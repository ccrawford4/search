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

func validUrl(index *Index, crawlerPolicy *CrawlerPolicy, newUrl, host string) bool {
	if (*index).containsUrl(newUrl) ||
		violatesPolicy(crawlerPolicy, newUrl) ||
		!strings.HasPrefix(strings.ToLower(newUrl), strings.ToLower(host)) {
		return false
	}
	return true
}

func getPolicy(seedUrl string) *CrawlerPolicy {
	var crawlerPolicy *CrawlerPolicy
	var subPath string
	if strings.HasSuffix(seedUrl, "/") {
		subPath = "robots.txt"
	} else {
		subPath = "/robots.txt"
	}

	policyPath := seedUrl + subPath
	policyContent, err := download(policyPath)
	if err != nil {
		log.Printf("error downloading robots.txt: %v", err)
		crawlerPolicy = getDefaultCrawlerPolicy(seedUrl)
	} else {
		if crawlerPolicy, err = getCrawlerPolicy(seedUrl, string(policyContent)); err != nil {
			log.Printf("error parsing robots.txt: %v", err)
		}
	}
	return crawlerPolicy
}

func crawl(index *Index, seedUrl *url.URL, testCrawl bool) {
	before := time.Now()
	initialFullPath, err := clean(seedUrl, seedUrl.Path) // keep track of the initial full path
	if err != nil {
		log.Fatalf("Could not clean host url %s: %v", seedUrl, err)
	}
	var crawlerPolicy *CrawlerPolicy
	if testCrawl {
		crawlerPolicy = getTestCrawlerPolicy(initialFullPath)
	} else {
		crawlerPolicy = getPolicy(initialFullPath)
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

		if !validUrl(index, crawlerPolicy, cleanedUrl, initialFullPath) {
			continue
		}

		// Sleep to avoid network bans
		time.Sleep(crawlerPolicy.delay)

		/// Download the contents of the page and add the cleanedURL to the visited set
		body, err := download(cleanedUrl)
		if err != nil {
			log.Printf("Error downloading content for %v: %v", cleanedUrl, err)
			continue
		}

		words, hrefs := extract(body)
		fmt.Printf("Download: url=%s result=ok\n", cleanedUrl)

		// Create the wordFrequency map
		wordFreq := make(Frequency)
		stopWords := getStopWords()
		for _, word := range words {
			stemmedWord := getStemmedWord(word, stopWords)
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
