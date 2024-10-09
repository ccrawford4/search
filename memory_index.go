package main

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kljensen/snowball"
	"log"
	"strings"
)

// IndexMap can be used to map word to their frequencies
type IndexMap map[string]Frequency

type MemoryIndex struct {
	StopWords *hashset.Set // Keeps track of all the StopWords
	WordCount Frequency    // map[url]total word count
	Index     IndexMap     // map[word][url]count
}

func (memoryIndex *MemoryIndex) isStopWord(s string) bool {
	return memoryIndex.StopWords.Contains(strings.ToLower(s))
}

func (memoryIndex *MemoryIndex) containsUrl(s string) bool {
	return memoryIndex.Index[s] != nil
}

func newMemoryIndex() *MemoryIndex {
	return &MemoryIndex{getStopWords(), make(Frequency), make(IndexMap)}
}

func (memoryIndex *MemoryIndex) getStemmedWord(word string) string {
	// Only get the stemmed version if it's not a stop word
	word = strings.ToLower(word)
	if !memoryIndex.isStopWord(word) {
		word, err := snowball.Stem(word, "english", true)
		if err != nil {
			log.Fatalf("Error stemming word %q: %v\n", word, err.Error())
		}
	}
	return word
}

func (memoryIndex *MemoryIndex) search(word string) *SearchResult {
	searchTerm := memoryIndex.getStemmedWord(word)
	freq, found := memoryIndex.Index[searchTerm]
	return &SearchResult{
		freq,
		len(memoryIndex.WordCount),
		found,
	}
}

func (memoryIndex *MemoryIndex) getTotalWords(url string) int {
	return memoryIndex.WordCount[url]
}

func (memoryIndex *MemoryIndex) insertCrawlResults(c *CrawlResult) {
	memoryIndex.WordCount[c.Url] = c.TotalWords
	for term, frequency := range c.TermFrequency {
		_, found := memoryIndex.Index[term]
		if !found {
			memoryIndex.Index[term] = make(Frequency)
		}
		memoryIndex.Index[term][c.Url] = frequency
	}
}
