package main

import (
	"encoding/json"
	"github.com/emirpasic/gods/sets/hashset"
	"log"
)

func convertJSON[V *[]string | *hashset.Set | *map[string]interface{}](filePath string, result V) {
	jsonContent, err := openAndReadFile(filePath)
	if err != nil {
		log.Printf("Error opening file %q\n", filePath)
		return
	}
	err = json.Unmarshal(jsonContent, &result)
	if err != nil {
		log.Fatalf("Error Parsing JSON %v\n", err)
	}
}

func getStopWords() *hashset.Set {
	var stopWords hashset.Set
	convertJSON("./documents/stopwords-en.json", &stopWords)
	return &stopWords
}

func isStopWord(word string, stopWords *hashset.Set) bool {
	return stopWords.Contains(word)
}
