package main

import (
	"errors"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kljensen/snowball"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"strings"
)

type DBIndex struct {
	StopWords *hashset.Set
	db        *gorm.DB
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

func (idx *DBIndex) search(word string) *SearchResult {
	// Get stemmed version of the input word
	searchTerm := idx.getStemmedWord(word)
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
			false,
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
		true,
	}
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

func newDBIndex(dbName string, db *gorm.DB) *DBIndex {
	return &DBIndex{StopWords: getStopWords(), db: db}
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
				IdxWordUrl: fmt.Sprintf("%d%d", word.ID, url.ID),
			})
		}
	}

	result := idx.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "idx_word_url"}},
		DoUpdates: clause.AssignmentColumns([]string{"count"}),
	}).Create(&wordFrequencyRecords)

	if result.Error != nil {
		log.Printf("Error inserting CrawlResults: %v\n", result.Error)
	}
}
