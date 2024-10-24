package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"log"
	"os"
)

type Query struct {
	gorm.Model
	User   User
	UserID uint
	Word   Word
	WordID uint
}

type User struct {
	gorm.Model
	FirstName string
	LastName  string
	Email     string `gorm:"index:email_idx,unique"`
}

type Word struct {
	gorm.Model
	Name string `gorm:"index:word_name_idx,unique"`
}

type Url struct {
	gorm.Model
	Name               string `gorm:"index:url_name_idx,unique"`
	Title, Description string
	Count              int
}

type WordFrequencyRecord struct {
	gorm.Model
	Count      int
	WordID     uint
	Word       Word
	UrlID      uint
	Url        Url
	IdxWordUrl string `gorm:"index:idx_word_url,unique"`
}

type QueryResult[T any] struct {
	Result T
	Error  error
}

// migrateTables migrates the Word, Url, and WordFrequencyRecord tables using autoMigrate
func migrateTables(db *gorm.DB) {
	err := db.AutoMigrate(&Word{}, &Url{}, &WordFrequencyRecord{})
	if err != nil {
		log.Fatalf("Error creating tables: %v\n", err)
	}
}

// dropDatabase to drop the database file
func dropDatabase(dbName string) {
	if err := os.Remove(dbName); err != nil && !os.IsNotExist(err) {
		log.Fatalf("Failed to drop the database: %v", err)
	}
	log.Println("Database dropped and will be recreated.")
}

// connectToDB connects to a sqlite DB given its name, migrates the tables, and then
// returns a pointer to the gorm.DB struct
func connectToDB(connString string, useSqlite bool) (*gorm.DB, error) {
	if useSqlite {
		db, err := gorm.Open(sqlite.Open(connString), &gorm.Config{})
		migrateTables(db)
		return db, err
	}
	db, err := gorm.Open(sqlserver.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	// Check the connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Error getting SQL DB from GORM: ", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	fmt.Println("Connected to the database using GORM!")
	migrateTables(db)
	return db, err
}

// getItem takes in a pointer to a struct, and fills the
// struct with data from the first entry of the respective table that matches the filter
func getItem[K *Word | *WordFrequencyRecord | *Url | *User](db *gorm.DB, object K) error {
	result := db.Where(object).First(object)
	if result.Error != nil {
		fmt.Printf("Error fetching %v: %v\n", object, result.Error)
	}
	return result.Error
}

// create takes in a pointer to a struct and inserts the data from the struct into the database
func create[K *Word | *WordFrequencyRecord | *Url | *User](db *gorm.DB, object K) error {
	if err := db.Create(object).Error; err != nil {
		log.Printf("Error creating object: %v\n", err)
	}
	return nil
}

// getItemOrCreate takes in a pointer to an object and attempts to fetch the object
// from the database. If it is unsuccessful then it inserts a new object into the database.
func getItemOrCreate[K *Word | *WordFrequencyRecord | *Url | *User](db *gorm.DB, object K) error {
	err := getItem(db, object)
	if err != nil {
		log.Printf("[WARNING] Could not get %v: trying to create\n", object)
		err = create(db, object)
	}
	return err
}

func batchInsertWords(db *gorm.DB, words []*Word, batchSize int) error {
	for i := 0; i < len(words); i += batchSize {
		end := i + batchSize
		if end > len(words) {
			end = len(words)
		}

		batch := words[i:end]
		if err := db.Create(&batch).Error; err != nil {
			return err
		}
	}
	return nil
}

func batchInsertWordFrequencyRecords(db *gorm.DB, wordFrequencyRecords []*WordFrequencyRecord, batchSize int) error {
	if len(wordFrequencyRecords) == 0 {
		return nil
	}

	// Helper function to execute a batch insert
	insertBatch := func(batch []*WordFrequencyRecord) error {
		if err := db.CreateInBatches(batch, batchSize).Error; err != nil {
			return err
		}
		return nil
	}

	// Process the records in batches
	for i := 0; i < len(wordFrequencyRecords); i += batchSize {
		end := i + batchSize
		if end > len(wordFrequencyRecords) {
			end = len(wordFrequencyRecords)
		}

		// Insert the current batch
		if err := insertBatch(wordFrequencyRecords[i:end]); err != nil {
			return err
		}
	}

	return nil
}
