package main

import (
	"errors"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
)

type DBIndex struct {
	db       *gorm.DB
	rsClient *redis.Client
}

func (idx *DBIndex) containsUrl(url string) bool {
	urlObj := &Url{
		Name: url,
	}
	result := idx.db.Where(&urlObj).First(urlObj)
	return result.Error == nil
}

const batchSize = 1000

func (idx *DBIndex) fetchFromDB(searchTerm string) *SearchResult {
	wordObj := Word{Name: searchTerm}

	// Get total count of URLs in the database
	var totalURLs int64
	idx.db.Table("urls").Count(&totalURLs)

	// Initialize frequency map
	frequency := make(Frequency)
	urlMap := make(UrlMap)

	// Retrieve the word object from the database
	if err := getItem(idx.db, &wordObj); err != nil {
		return &SearchResult{
			urlMap,
			frequency,
			int(totalURLs),
			false,
		}
	}

	// Fetch word frequency records in batches
	var offset int
	for {
		var wordFrequencyRecords []WordFrequencyRecord
		result := idx.db.
			Where("word_id = ?", wordObj.ID).
			Preload("Url").
			Offset(offset).
			Limit(batchSize).
			Find(&wordFrequencyRecords)

		// Handle potential errors during fetching
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				log.Printf("No word frequency records found for word: %v\n", wordObj.Name)
			} else {
				log.Printf("Error fetching word frequency records: %v\n", result.Error)
			}
			break
		}

		// Break the loop if no more records are returned
		if len(wordFrequencyRecords) == 0 {
			break
		}

		// Populate frequency map
		for _, record := range wordFrequencyRecords {
			urlMap[record.Url.Name] = UrlEntry{
				record.Url.Count,
				record.Url.Title,
				record.Url.Description,
			}
			frequency[record.Url.Name] = record.Count
		}

		// Increment offset for the next batch
		offset += batchSize
	}

	// Return search result
	return &SearchResult{
		urlMap,
		frequency,
		int(totalURLs),
		true,
	}
}

func (idx *DBIndex) search(word string) *SearchResult {
	word, err := getStemmedWord(word)
	if err != nil {
		log.Printf("[WARNING] Could Not Stem Word %q\n", word)
	}
	result, err := fetchFromCache(idx.rsClient, word)
	if err != nil {
		log.Printf("Cache miss for word %q Fetching from DB now\n", word)
		result = idx.fetchFromDB(word)
	} else {
		log.Printf("Cache hit for term %q\n", word)
	}
	err = insertIntoCache(idx.rsClient, word, result)
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
	return &DBIndex{db, rsClient}
}

// To create a UUID for the index
func cantorPairing(wordID, urlID uint) uint {
	return (wordID+urlID)*(wordID+urlID+1)/2 + urlID
}

func (idx *DBIndex) insertCrawlResults(c *CrawlResult) {
	url := Url{
		Name:        c.Url,
		Count:       c.TotalWords,
		Title:       c.Title,
		Description: c.Description,
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
			Url:        url,
			Word:       *word,
			WordID:     word.ID,
			UrlID:      url.ID,
			Count:      c.TermFrequency[word.Name],
			IdxWordUrl: fmt.Sprintf("%d", cantorPairing(word.ID, url.ID)),
		})
	}

	if err = batchInsertWordFrequencyRecords(idx.db, wordFrequencyRecords, 250); err != nil {
		log.Printf("Error upserting word frequency records %v", err)
	}
}

func (idx *DBIndex) updateUserTable(word, email string) error {
	user := User{
		Email: email,
	}
	err := getItemOrCreate(idx.db, &user)
	if err != nil {
		log.Printf("[ERROR] Could Not Fetch User: %v\n", user)
		// TODO: Incorporate usage for cache here
	}

	stemmedWord, err := getStemmedWord(word)
	if err != nil {
		log.Printf("[ERROR] Could Not Get Stemmed Word: %v\n", err)
		return err
	}

	wordObj := Word{
		Name: stemmedWord,
	}

	err = getItem(idx.db, &wordObj)
	if err != nil {
		log.Printf("[ERROR] Could Not Fetch Word: %v\n", err)
		return err
	}

	return nil
}
