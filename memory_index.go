package main

import (
	"log"
)

// IndexMap can be used to map word to their frequencies
type IndexMap map[string]Frequency

type MemoryIndex struct {
	WordCount UrlMap   // map[url]{total words, title}
	Index     IndexMap // map[word][url]count
}

func (memoryIndex *MemoryIndex) containsUrl(s string) bool {
	return memoryIndex.Index[s] != nil
}

func newMemoryIndex() *MemoryIndex {
	return &MemoryIndex{make(UrlMap), make(IndexMap)}
}

func (memoryIndex *MemoryIndex) search(word string) *SearchResult {
	word, err := getStemmedWord(word)
	if err != nil {
		log.Printf("[WARNING] Could not stem word %q: %v\n", word, err)
	}
	freq, found := memoryIndex.Index[word]
	return &SearchResult{
		memoryIndex.WordCount,
		freq,
		len(memoryIndex.WordCount),
		found,
	}
}

func (memoryIndex *MemoryIndex) getTotalWords(url string) int {
	return memoryIndex.WordCount[url].TotalWords
}

func (memoryIndex *MemoryIndex) insertCrawlResults(c *CrawlResult) {
	memoryIndex.WordCount[c.Url] = UrlEntry{
		c.TotalWords,
		c.Title,
		c.Description,
	}

	for term, frequency := range c.TermFrequency {
		_, found := memoryIndex.Index[term]
		if !found {
			memoryIndex.Index[term] = make(Frequency)
		}
		memoryIndex.Index[term][c.Url] = frequency
	}
}
