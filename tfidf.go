package main

import (
	"github.com/emirpasic/gods/sets/hashset"
	"math"
	"sort"
)

// createWordFrequency creates a Frequency map
// with words as the keys, and their counts as values
func createWordFrequency(words []string, stopWords *hashset.Set) Frequency {
	freq := make(Frequency)
	for _, word := range words {
		stemmedWord := getStemmedWord(word, stopWords)
		freq[stemmedWord] += 1
	}
	return freq
}

func calculateIDF(docsContainingWord float64, numDocs float64) float64 {
	return math.Log10(numDocs / (docsContainingWord + 1))
}

// calculateTF returns the Term Frequency of a word given
// the termCount and the total number of words in the document
func calculateTF(termCount, totalWords float64) float64 {
	return termCount / totalWords
}

// calculateTFIDF calculates the TFIDF score
func calculateTFIDF(termCount, totalWords, docsContainingWord, numDocs float64) float64 {
	return calculateTF(termCount, totalWords) * calculateIDF(docsContainingWord, numDocs)
}

//func populateTFIDFValues(index *Index, numDocs float64, wordsInDoc Frequency) {
//	for word, frequency := range *index {
//		idf := calculateIDF(index, numDocs, word)
//		for URL, termCount := range frequency {
//			totalWords := wordsInDoc[URL]
//			// Compute the TF-IDF by multiplying the TF * IDF
//			tf := calculateTF((float64)(termCount), (float64)(totalWords))
//			(*index)[word][URL] = tf * idf
//		}
//	}
//}

// getTemplateData takes in a Frequency object and
// a searchTerm and returns the formated TemplateData response
func getTemplateData(freq Frequency, searchTerm string, numDocs float64, urlWordTotals *Frequency) *TemplateData {
	// Iterate through the frequency map and populate the hits array
	var hits Hits
	docsContainingWord := (float64)(len(freq))
	for url, count := range freq {
		totalWords := (float64)((*urlWordTotals)[url])
		tfidf := calculateTFIDF((float64)(count), totalWords, docsContainingWord, numDocs)
		hits = append(hits, Hit{url, tfidf})
	}
	// Sort the hits array based on TF-IDF score
	sort.Sort(hits)
	return &TemplateData{
		HITS: hits,
		TERM: searchTerm,
	}
}
