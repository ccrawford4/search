package main

type SearchResult struct {
	TermFrequency     Frequency
	TotalDocsSearched int
	Found             bool
}

type Index interface {
	isStopWord(word string) bool
	containsUrl(url string) bool
	search(word string) *SearchResult
	getTotalWords(url string) int
	getStemmedWord(word string) string
	insertCrawlResults(c *CrawlResult)
}
