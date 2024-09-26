package main

import (
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kljensen/snowball"
	"log"
	"strings"
)

// getStemmedWord takes in a word and returns the stemmed version of said word
func getStemmedWord(word string, stopWords *hashset.Set) string {
	// Only get the stemmed version if it's not a stop word
	word = strings.ToLower(word)
	if !isStopWord(word, stopWords) {
		word, err := snowball.Stem(word, "english", true)
		if err != nil {
			log.Fatalf("Error stemming word %q: %v\n", word, err.Error())
		}
	}
	return word
}

// search takes in an index and a search term and returns the
// Frequency result and a bool indicating whether the term was found or not
func search(index *Index, searchTerm string, stopWords *hashset.Set) (Frequency[float64], bool) {
	// If the search term is not a stop word then get the stemmed version of the word
	searchTerm = getStemmedWord(searchTerm, stopWords)
	freq, found := (*index)[searchTerm]
	return freq, found
}
