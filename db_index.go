package main

import (
	"errors"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-redis/redis/v8"
	"github.com/kljensen/snowball"
	"gorm.io/gorm"
	"log"
	"strings"
)

type DBIndex struct {
	StopWords *hashset.Set
	db        *gorm.DB
	rsClient  *redis.Client
}

func (idx *DBIndex) isStopWord(s string) bool {
	return idx.StopWords.Contains(strings.ToLower(s))
}

func (idx *DBIndex) containsUrl(url string) bool {
	result := idx.db.Where("name = ?", url).First(
		&Url{
			Name: url,
		})
	return result.Error == nil
}

func (idx *DBIndex) fetchFromDB(searchTerm string) *SearchResult {
	wordObj := Word{Name: searchTerm}

	// Get total count of URLs in the database
	var totalURLs int64
	idx.db.Table("urls").Count(&totalURLs)

	// Initialize frequency map
	frequency := make(Frequency)

	// Retrieve the word object from the database
	if err := getItem(idx.db, &wordObj); err != nil {
		return &SearchResult{
			frequency,
			int(totalURLs),
		}
	}

	// Fetch word frequency records associated with the word object
	var wordFrequencyRecords []WordFrequencyRecord
	result := idx.db.
		Where("word_id = ?", wordObj.ID).
		Preload("Url").
		Find(&wordFrequencyRecords)

	// Handle potential errors during fetching
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Printf("No word frequency records found for word: %v\n", wordObj.Name)
		} else {
			log.Printf("Error fetching word frequency records: %v\n", result.Error)
		}
	}

	// Populate frequency map
	for _, record := range wordFrequencyRecords {
		frequency[record.Url.Name] = record.Count
	}

	// Return search result
	return &SearchResult{
		frequency,
		int(totalURLs),
	}
}

func (idx *DBIndex) search(word string) *SearchResult {
	searchTerm := idx.getStemmedWord(word)
	result, err := fetchFromCache(idx.rsClient, searchTerm)
	if err != nil {
		log.Printf("Cache miss for word %q Fetching from DB now\n", word)
		result = idx.fetchFromDB(word)
	} else {
		log.Printf("Cache hit for term %q\n", word)
	}
	err = insertIntoCache(idx.rsClient, searchTerm, result)
	if err != nil {
		log.Printf("Failed to insert %q into cache: %v\n", word, err)
	}
	return result
}

func (idx *DBIndex) getTotalWords(url string) int {
	urlObj := Url{
		Name: url,
	}
	err := idx.db.Where(&urlObj).First(&urlObj).Error
	if err != nil {
		log.Printf("URL %s Not found in the DB: %v\n", url, err)
		return 0
	}
	return urlObj.Count
}

func newDBIndex(connString string, useSqlite bool, rsClient *redis.Client) *DBIndex {
	db, err := connectToDB(connString, useSqlite)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v\n", err)
	}
	return &DBIndex{StopWords: getStopWords(), db: db, rsClient: rsClient}
}

func (idx *DBIndex) getStemmedWord(word string) string {
	// Only get the stemmed version if it's not a stop word
	word = strings.ToLower(word)
	if !idx.isStopWord(word) {
		_, err := snowball.Stem(word, "english", true)
		if err != nil {
			log.Fatalf("Error stemming word %q: %v\n", word, err.Error())
		}
	}
	return word
}

func (idx *DBIndex) insertCrawlResults(c *CrawlResult) {
	url := Url{
		Name:  c.Url,
		Count: c.TotalWords,
	}
	err := getItemOrCreate(idx.db, &url)
	if err != nil {
		log.Printf("Error fetching or Creating URL %v: %v\n", url, err)
		return
	}

	// Extract all the names from the termFrequency
	names := make([]string, 0, len(c.TermFrequency))
	for word := range c.TermFrequency {
		if word == "the" {
			fmt.Printf("found the!: %v\n", word)
		}
		names = append(names, word)
	}

	for _, word := range names {
		if word == "the" {
			fmt.Printf("found the: %v\n", word)
		}
	}

	// Identify all the ones that have the matching name in the db
	var existingWords []*Word
	idx.db.Model(&Word{}).Select("name", "id").Where("name IN ?", names).Find(&existingWords)

	for _, word := range existingWords {
		if word.Name == "the" {
			fmt.Printf("Found the!: %v\n", word)
		}
	}

	// Now keep track of all the names that are already in the database
	seenNames := hashset.New()
	for _, word := range existingWords {
		seenNames.Add(word.Name)
	}

	// Finally go through the names again and if they are not in the database then create them
	var newWords []*Word
	for word := range c.TermFrequency {
		if !seenNames.Contains(word) {
			newWords = append(newWords, &Word{Name: word})
		}
	}
	// Set the desired batch size
	if err := batchInsertWords(idx.db, newWords, 500); err != nil {
		log.Printf("Error inserting words: %v", err)
		return
	}

	var allWords []*Word
	allWords = append(newWords, existingWords...)

	// Now create the word frequency records array
	var wordFrequencyRecords []*WordFrequencyRecord
	for _, word := range allWords {
		if word.Name == "the" {
			fmt.Printf("found the!: %v\n", word)
		}
		wordFrequencyRecords = append(wordFrequencyRecords, &WordFrequencyRecord{
			Url:    url,
			Word:   *word,
			WordID: word.ID,
			UrlID:  url.ID,
			Count:  c.TermFrequency[word.Name],
		})
	}

	if err = batchInsertWordFrequencyRecords(idx.db, wordFrequencyRecords, 250); err != nil {
		log.Printf("Error upserting word frequency records %v", err)
	}
}
