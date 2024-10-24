package main

import (
	"github.com/emirpasic/gods/sets/hashset"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
)

type CrawlResult struct {
	TermFrequency           Frequency // frequency per word
	Url, Title, Description string    // the url crawled
	TotalWords              int       // the total # of words in the document
}

// parseURL takes in a rawURL or href and returns a new pointer to an url.URL object
func parseURL(rawURL string) (*url.URL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Printf("error parsing URL %s: %v", rawURL, err)
	}
	return parsedURL, err
}

// normalizeHost removes the "www." if present to handle both cases
func normalizeHost(url string) string {
	return strings.ReplaceAll(strings.ToLower(url), "www.", "")
}

// validUrl returns false if the url violates the crawl policy or is not within the same host
func validUrl(crawlerPolicy *CrawlerPolicy, newUrl, host string) bool {
	normalizedHost := normalizeHost(host)
	normalizedNewUrl := normalizeHost(newUrl)

	if violatesPolicy(crawlerPolicy, newUrl) ||
		!strings.HasPrefix(normalizedNewUrl, normalizedHost) {
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
				log.Printf("[WARNING] %q Could Not Be Cleaned: %v\n", href, err)
				continue
			}

			if !validUrl(crawlerPolicy, cleanedUrl, initialFullPath) {
				log.Printf("[WARNING] %q Is Not A Valid Url\n", cleanedUrl)
				continue
			}

			if visited.Contains(href) {
				log.Printf("[WARNING] %q Has Already Been Processed\n", href)
				continue
			}

			visited.Add(href)
			cleanedUrls = append(cleanedUrls, cleanedUrl)
		}

		// Clear the queue
		queue = queue[:0]

		// start wait group
		var wg sync.WaitGroup
		downloadChannel := make(chan Download, 10000)

		// Go go go
		for i, cleanedUrl := range cleanedUrls {
			wg.Add(1)
			time.Sleep(time.Duration(i) + crawlerPolicy.delay)
			go download(cleanedUrl, &wg, downloadChannel)
		}

		// Wait for synchronization
		go func() {
			wg.Wait()
			close(downloadChannel)
		}()

		// Process the download results as they come in
		var allHrefs []string
		for downloadObj := range downloadChannel {
			if downloadObj.Err != nil {
				log.Printf("[WARNING] Could Not Process Url %q: %v\n", downloadObj.Url, downloadObj.Err)
				continue
			}

			// Extract content from the downloaded body
			words, hrefs, title, description := extract(downloadObj.Body)
			if words == nil && hrefs == nil {
				log.Printf("[WARNING] Could not parse content from url %q\n", downloadObj.Url)
				continue
			}

			// Only add hrefs we haven't seen already
			for _, href := range hrefs {
				if !visited.Contains(href) {
					allHrefs = append(allHrefs, href)
				}
			}

			// Process word frequencies
			wordFreq := make(Frequency)
			for _, word := range words {
				// Ensure you use the stemmed version of the word
				stemmedWord, err := getStemmedWord(word)
				if err != nil {
					log.Printf("[WARNING] Could not get stemming word %q: %v\n", word, err)
					continue
				}
				wordFreq[stemmedWord]++
			}

			// Insert crawl results into the index
			(*index).insertCrawlResults(&CrawlResult{
				wordFreq,
				downloadObj.Url,
				title,
				description,
				len(words),
			})
		}
		queue = append(queue, allHrefs...)
		log.Println("Finished Queue Round.")
	}
	duration := time.Now().Sub(before)
	log.Printf("Crawl Duration: %v\n", duration.String())
}
