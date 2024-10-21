package main

import (
	"github.com/emirpasic/gods/sets/hashset"
)

// IndexMap can be used to map word to their frequencies
type IndexMap map[string]Frequency

type MemoryIndex struct {
	StopWords *hashset.Set // Keeps track of all the StopWords
	WordCount UrlMap       // map[url]{total words, title}
	Index     IndexMap     // map[word][url]count
}

func (memoryIndex *MemoryIndex) containsUrl(s string) bool {
	return memoryIndex.Index[s] != nil
}

func newMemoryIndex() *MemoryIndex {
	return &MemoryIndex{getStopWords(), make(UrlMap), make(IndexMap)}
}

func (memoryIndex *MemoryIndex) search(word string) *SearchResult {
	searchTerm := getStemmedWord(word, getStopWords())
	freq, found := memoryIndex.Index[searchTerm]
	return &SearchResult{
		memoryIndex.WordCount,
		freq,
		len(memoryIndex.WordCount),
		found,
	}
}

func (memoryIndex *MemoryIndex) getStopWords() *hashset.Set {
	return memoryIndex.StopWords
}

func (memoryIndex *MemoryIndex) getTotalWords(url string) int {
	return memoryIndex.WordCount[url].TotalWords
}

func (memoryIndex *MemoryIndex) insertCrawlResults(c *CrawlResult) {
	memoryIndex.WordCount[c.Url] = UrlEntry{
		c.TotalWords,
		c.Title,
	}

	for term, frequency := range c.TermFrequency {
		_, found := memoryIndex.Index[term]
		if !found {
			memoryIndex.Index[term] = make(Frequency)
		}
		memoryIndex.Index[term][c.Url] = frequency
	}
}
