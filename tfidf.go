package main

import (
	"math"
	"sort"
)

// calculateIDF returns the IDF score based on the number of
// docs containing a word and the total number of docs searched
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

// getTemplateData takes in a Frequency object and
// a searchTerm and returns the formated TemplateData response
func getTemplateData(index *Index, searchTerm string) *TemplateData {
	searchResults := (*index).search(searchTerm)
	if !searchResults.Found {
		return nil
	}

	var hits Hits
	docsContainingWord := (float64)(len(searchResults.TermFrequency))
	numDocs := float64(searchResults.TotalDocsSearched)

	// Iterate through the frequency map and populate the hits array
	for url, count := range searchResults.TermFrequency {
		if url == "" {
			continue
		}

		// Calculate TFIDF
		totalWords := searchResults.UrlMap[url].TotalWords
		tfidf := calculateTFIDF((float64)(count), (float64)(totalWords), docsContainingWord, numDocs)

		// Pull out data from the url map and clean up titles that are not found
		urlEntry := searchResults.UrlMap[url]
		title := urlEntry.Title

		if urlEntry.Title == "" {
			title = url
		}

		// Update the hits array
		hits = append(hits, Hit{url, title, urlEntry.Description, tfidf})
	}
	// Sort the hits array based on TF-IDF score
	sort.Sort(hits)
	return &TemplateData{
		HITS: hits,
		TERM: searchTerm,
	}
}
