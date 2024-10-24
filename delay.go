package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const invalidAgent = "go-http-client/1.1"
const defaultCrawlDelay = time.Millisecond * 10

type CrawlerPolicy struct {
	root       string
	disallowed []string
	delay      time.Duration
}

func extractEntry(line, key string) string {
	var entry string
	cleanedLine := strings.ToLower(strings.TrimSpace(line))
	if strings.HasPrefix(cleanedLine, key) {
		data := strings.Split(line, ":")
		if len(data) > 1 {
			entry = strings.TrimSpace(data[1])
		}
	}
	entry = strings.ReplaceAll(entry, "*", ".*")
	return entry
}

func getDefaultCrawlerPolicy(root string) *CrawlerPolicy {
	return &CrawlerPolicy{
		root,
		nil,
		defaultCrawlDelay,
	}
}

func getTestCrawlerPolicy(root string) *CrawlerPolicy {
	return &CrawlerPolicy{
		root,
		nil,
		defaultCrawlDelay,
	}
}

func getCrawlerPolicy(root string, content string) (*CrawlerPolicy, error) {
	lines := strings.Split(content, "\n")
	var disallowedHosts []string
	crawlDelay := defaultCrawlDelay
	for _, line := range lines {
		if agent := extractEntry(line, "user-agent"); agent == invalidAgent {
			return nil, errors.New("error! Invalid agent go-http-client/1.1 is not allowed to crawl this site")
		}
		if host := extractEntry(line, "disallow"); host != "" {
			disallowedHosts = append(disallowedHosts, host)
		}
		if delay := extractEntry(line, "crawl-delay"); delay != "" {
			duration, err := strconv.ParseInt(delay, 10, 64)
			if err != nil {
				log.Printf("Could not parse %q as an int for Crawl-Delay: %v\n", delay, err)
			}
			crawlDelay = time.Duration(duration) * time.Second
		}
	}
	return &CrawlerPolicy{root, disallowedHosts, crawlDelay}, nil
}

func violatesPolicy(policy *CrawlerPolicy, host string) bool {
	for _, disallowed := range policy.disallowed {
		if disallowed == host || strings.HasPrefix(host, disallowed) {
			return true
		}
		pattern := fmt.Sprintf("%s", disallowed)
		re := regexp.MustCompile(pattern)
		if re.MatchString(host) {
			return true
		}
	}
	return false
}
