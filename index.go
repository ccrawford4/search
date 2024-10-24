package main

type UrlEntry struct {
	TotalWords         int
	Title, Description string
}

type UrlMap map[string]UrlEntry

type SearchResult struct {
	UrlMap            UrlMap
	TermFrequency     Frequency
	TotalDocsSearched int
	Found             bool
}

type Index interface {
	containsUrl(url string) bool
	search(word string) *SearchResult
	getTotalWords(url string) int
	insertCrawlResults(c *CrawlResult)
	updateUserTable(word, email string) error
}
