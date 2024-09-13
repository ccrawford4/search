package main

import (
	"github.com/kljensen/snowball"
	"log"
)

func getStemmedWord(word string) string {
	stemmed, err := snowball.Stem(word, "english", true)
	if err != nil {
		log.Fatalf("Error stemming word %q: %v\n", word, err.Error())
	}
	return stemmed
}

func search(host, term string) frequency {
	indexMap := crawl(parseURL(host))
	return indexMap[getStemmedWord(term)]
}
