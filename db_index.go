package main

import (
	"errors"
	"github.com/emirpasic/gods/sets/hashset"
	"github.com/kljensen/snowball"
	"gorm.io/gorm"
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

func newDBIndex(connString string, useSqlite bool) *DBIndex {
	db, err := connectToDB(connString, useSqlite)
	if err != nil {
		log.Fatalf("Error connecting to DB: %v\n", err)
	}
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
	// Create the URL object
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
		names = append(names, word)
	}

	// Identify all the ones that have the matching name in the db
	var existingWords []*Word
	idx.db.Model(&Word{}).Select("name").Where("name IN ?", names).Find(&existingWords)

	// Now keep track of all the names that are already in the database
	seenNames := hashset.New()
	for _, word := range existingWords {
		seenNames.Add(word.Name)
	}

	// Finally go through the names again and if they are not in the database then create them
	var newWords []Record
	for word := range c.TermFrequency {
		if !seenNames.Contains(word) {
			newWords = append(newWords, &Word{Name: word})
		}
	}

	batchSize := 500 // Set the desired batch size
	if err := batchInsert(idx.db, newWords, &Word{}, batchSize); err != nil {
		log.Printf("Error inserting words: %v", err)
		return
	}

	// Now create the word frequency records array
	var wordFrequencyRecords []Record
	for _, record := range newWords {
		word := record.GetWord()
		if !seenNames.Contains(word.Name) {
			wordFrequencyRecords = append(wordFrequencyRecords, &WordFrequencyRecord{
				Url:    url,
				Word:   *word,
				WordID: word.ID,
				UrlID:  url.ID,
				Count:  c.TermFrequency[word.Name],
			})
		}
	}
	// TODO: issue here with passing generic type.
	if err = batchInsert(idx.db, wordFrequencyRecords, &WordFrequencyRecord{}, batchSize); err != nil {
		log.Printf("Error inserting word frequency records %v", err)
	}
}
