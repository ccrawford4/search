package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func newDBIndex(connString string, rsClient *redis.Client) *DBIndex {
	db, err := connectToDB(connString)
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
		Title:       c.Title,
		Description: c.Description,
		Count:       c.TotalWords,
	}
	err := getItemOrCreate(idx.db, &url)
	if err != nil {
		log.Printf("Error fetching or Creating URL %v: %v\n", url, err)
	}

	var words []*Word
	var termCounts []int

	terms := make([]string, 0, len(c.TermFrequency))
	for term, frequency := range c.TermFrequency {
		terms = append(terms, term)
		termCounts = append(termCounts, frequency)
	}

	// Create Word objects
	for _, term := range terms {
		words = append(words, &Word{Name: term})
	}

	// Batch Create of Words with OnConflict handling
	idx.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).Create(&words)

	// Retrieve all existing words to ensure we have correct references
	var existingWords []Word
	idx.db.Where("name IN ?", terms).Find(&existingWords)

	// Replace the original words slice with the existing ones from the database
	words = make([]*Word, len(existingWords))
	for i, w := range existingWords {
		words[i] = &w
	}

	var wordFrequencyRecords []*WordFrequencyRecord

	// Create a map to associate terms with their corresponding existing words
	wordMap := make(map[string]*Word)
	for _, w := range existingWords {
		wordMap[w.Name] = &w
	}

	// Populate the WordFrequencyRecord
	for term, frequency := range c.TermFrequency {
		if word, found := wordMap[term]; found {
			wordFrequencyRecords = append(wordFrequencyRecords, &WordFrequencyRecord{
				Url:        url,
				Word:       *word,
				WordID:     word.ID,
				UrlID:      url.ID,
				Count:      frequency,
				IdxWordUrl: fmt.Sprintf("%d", cantorPairing(word.ID, url.ID)),
			})
		}
	}

	if err = batchInsertWordFrequencyRecords(idx.db, wordFrequencyRecords, 250); err != nil {
		log.Printf("[WARNING] Could Not Insert Word Frequency Records: %v", err)
	}
}
