package main

type SearchResult struct {
	TermFrequency     Frequency `json:"termFrequency,omitempty"`
	TotalDocsSearched int       `json:"totalDocsSearched,omitempty"`
}

type Index interface {
	isStopWord(word string) bool
	containsUrl(url string) bool
	search(word string) *SearchResult
	getTotalWords(url string) int
	getStemmedWord(word string) string
	insertCrawlResults(c *CrawlResult)
}
