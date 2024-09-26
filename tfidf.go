package main

import (
	"github.com/emirpasic/gods/sets/hashset"
	"math"
	"sort"
)

// createWordFrequency creates a Frequency map
// with words as the keys, and their counts as values
func createWordFrequency(words []string, stopWords *hashset.Set) Frequency[int] {
	freq := make(Frequency[int])
	for _, word := range words {
		stemmedWord := getStemmedWord(word, stopWords)
		freq[stemmedWord] += 1
	}
	return freq
}

func calculateIDF(index *Index, numDocs float64, word string) float64 {
	docsContainingWord := (float64)(len((*index)[word]))
	return math.Log10(numDocs / (docsContainingWord + 1))
}

// calculateTF returns the Term Frequency of a word given
// the termCount and the total number of words in the document
func calculateTF(termCount, totalWords float64) float64 {
	return termCount / totalWords
}

func populateTFIDFValues(index *Index, numDocs float64, wordsInDoc Frequency[int]) {
	for word, frequency := range *index {
		idf := calculateIDF(index, numDocs, word)
		for URL, termCount := range frequency {
			totalWords := wordsInDoc[URL]
			// Compute the TF-IDF by multiplying the TF * IDF
			tf := calculateTF(termCount, (float64)(totalWords))
			(*index)[word][URL] = tf * idf
		}
	}
}

// getTemplateData takes in a Frequency object and
// a searchTerm and returns the formated TemplateData response
func getTemplateData(freq Frequency[float64], searchTerm string) *TemplateData {
	// Iterate through the frequency map and populate the hits array
	var hits Hits
	for url, tf := range freq {
		hits = append(hits, Hit{url, tf})
	}
	// Sort the hits array based on TF-IDF score
	sort.Sort(hits)
	return &TemplateData{
		HITS: hits,
		TERM: searchTerm,
	}
}
