package main

import (
	"github.com/kljensen/snowball"
	"log"
)

// getStemmedWord takes in a word and returns the stemmed version of said word
func getStemmedWord(word string) string {
	stemmed, err := snowball.Stem(word, "english", true)
	if err != nil {
		log.Fatalf("Error stemming word %q: %v\n", word, err.Error())
	}
	return stemmed
}

// search takes in a host and a search term and returns the frequency object for the
// stemmed version of the search term
func search(host, term string) frequency {
	indexMap := crawl(parseURL(host))
	return indexMap[getStemmedWord(term)]
}
