package main

import (
	"github.com/kljensen/snowball"
	"log"
)

func getStemmedWord(word string) (string, error) {
	result, err := snowball.Stem(word, "english", false)
	if err != nil {
		log.Printf("[WARNING] could not stem word %q\n", word)
		return word, err
	}
	return result, nil
}
