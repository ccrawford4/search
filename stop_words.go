package main

import (
	"encoding/json"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kljensen/snowball"
	"log"
	"strings"
)

// convertJSON converts a JSON file into a hashset
func convertJSON(filePath string, result *hashset.Set) {
	jsonContent, err := openAndReadFile(filePath)
	if err != nil {
		log.Printf("[WARNING] could not open file %q\n", filePath)
		return
	}
	err = json.Unmarshal(jsonContent, &result)
	if err != nil {
		log.Fatalf("Error Parsing JSON %v\n", err)
	}
}

// getStopWords() returns a hashset containing all the unique stop words (english)
func getStopWords() *hashset.Set {
	var stopWords hashset.Set
	convertJSON("./documents/stopwords-en.json", &stopWords)
	return &stopWords
}

func getStemmedWord(word string, stopWords *hashset.Set) string {
	word = strings.ToLower(word)
	if !stopWords.Contains(word) {
		_, err := snowball.Stem(word, "english", true)
		if err != nil {
			log.Fatalf("Error stemming word %q: %v\n", word, err.Error())
		}
	}
	return word
}
