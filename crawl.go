package main

import (
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

type CrawlResult struct {
	TermFrequency Frequency // frequency per word
	Url, Title    string    // the url crawled
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
func validUrl(crawlerPolicy *CrawlerPolicy, newUrl, host string) bool {
	if violatesPolicy(crawlerPolicy, newUrl) ||
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

	downloadChan := make(chan Download)
	var wg sync.WaitGroup
	wg.Add(1)
	go download(policyPath, &wg, downloadChan)

	go func() {
		wg.Wait()
		close(downloadChan)
	}()

	downloadObj, ok := <-downloadChan
	if !ok {
		return getDefaultCrawlerPolicy(seedUrl)
	}

	var err error
	if downloadObj.Err != nil {
		log.Printf("error downloading robots.txt: %v", downloadObj.Err)
		crawlerPolicy = getDefaultCrawlerPolicy(seedUrl)
	} else {
		if crawlerPolicy, err = getCrawlerPolicy(seedUrl, string(downloadObj.Body)); err != nil {
			log.Printf("error parsing robots.txt: %v", err)
		}
	}
	return crawlerPolicy
}

type Download struct {
	Body []byte
	Url  string
	Err  error
}

func crawl(index *Index, seedUrl string, testCrawl bool) {
	before := time.Now()
	urlObj, err := parseURL(seedUrl)
	if err != nil {
		log.Fatalf("error parsing URL %s: %v", seedUrl, err)
	}
	initialFullPath, err := clean(seedUrl, urlObj.Path) // keep track of the initial full path
	if err != nil {
		log.Fatalf("Could not clean host url %s: %v", seedUrl, err)
	}
	var crawlerPolicy *CrawlerPolicy
	if testCrawl {
		crawlerPolicy = getTestCrawlerPolicy(initialFullPath)
	} else {
		crawlerPolicy = getPolicy(initialFullPath)
	}

	visited := hashset.New()
	queue := []string{urlObj.Path} // Queue for keeping track of hrefs to visit
	for len(queue) > 0 {
		var cleanedUrls []string
		for _, href := range queue {
			cleanedUrl, err := clean(seedUrl, href)
			if err != nil {
				log.Printf("error crawling url %s: %v", seedUrl, err)
				continue
			}
			if !validUrl(crawlerPolicy, cleanedUrl, initialFullPath) {
				continue
			}
			if visited.Contains(href) {
				continue
			}
			visited.Add(href)
			cleanedUrls = append(cleanedUrls, cleanedUrl)
		}

		queue = queue[:0]

		// start wait group
		var wg sync.WaitGroup
		downloadChannel := make(chan Download, 10000)
		for i, cleanedUrl := range cleanedUrls {
			wg.Add(1)
			// For usfca.edu crawling

			if testCrawl {
				time.Sleep(crawlerPolicy.delay + time.Duration(i) + time.Millisecond)
			} else {
				time.Sleep(time.Duration(i) + (10 * time.Millisecond))
			}

			// For testing
			// time.Sleep(crawlerPolicy.delay + time.Duration(i) + time.Millisecond)
			go download(cleanedUrl, &wg, downloadChannel)
		}

		go func() {
			wg.Wait()
			close(downloadChannel)
		}()

		var allHrefs []string
		for downloadObj := range downloadChannel {
			// Extract content from the downloaded body
			words, hrefs, title := extract(downloadObj.Body)
			if words == nil && hrefs == nil {
				log.Printf("Could not parse content from url %q\n", downloadObj.Url)
				continue
			}
			allHrefs = append(allHrefs, hrefs...)

			// Process word frequencies
			wordFreq := make(Frequency)
			for _, word := range words {
				stemmedWord := getStemmedWord(word, (*index).getStopWords())
				wordFreq[stemmedWord]++
			}

			// Insert crawl results into the index
			(*index).insertCrawlResults(&CrawlResult{
				wordFreq,
				downloadObj.Url,
				title,
				len(words),
			})
		}
		queue = append(queue, allHrefs...)
		fmt.Println("Finished a round.")
	}
	duration := time.Now().Sub(before)
	log.Printf("Crawl Duration: %v\n", duration.String())
}
