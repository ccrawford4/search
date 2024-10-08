package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

type Word struct {
	gorm.Model
	Name string `gorm:"unique"`
}
type Url struct {
	gorm.Model
	Name  string `gorm:"unique"`
	Count int
}
type WordFrequencyRecord struct {
	gorm.Model
	Count      int
	WordID     uint
	Word       Word
	UrlID      uint
	Url        Url
	IdxWordUrl string `gorm:"uniqueIndex:idx_word_url"`
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
func connectToDB(dbName string, resetDB bool) *gorm.DB {
	if resetDB {
		dropDatabase(dbName)
	}
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	migrateTables(db)
	return db
}

// getItem takes in a pointer to a struct, and fills the
// struct with data from the first entry of the respective table that matches the filter
func getItem[K *Word | *WordFrequencyRecord | *Url](db *gorm.DB, object K) error {
	result := db.Where(object).First(object)
	if result.Error != nil {
		fmt.Printf("Error fetching %v: %v\n", object, result.Error)
	}
	return result.Error
}

// create takes in a pointer to a struct and inserts the data from the struct into the database
func create[K *Word | *WordFrequencyRecord | *Url](db *gorm.DB, object K) error {
	if err := db.Create(object).Error; err != nil {
		log.Printf("Error creating object: %v\n", err)
	}
	return nil
}

// getItemOrCreate takes in a pointer to an object and attempts to fetch the object
// from the database. If it is unsuccessful then it inserts a new object into the database.
func getItemOrCreate[K *Word | *WordFrequencyRecord | *Url](db *gorm.DB, object K) error {
	err := getItem(db, object)
	if err != nil {
		err = create(db, object)
	}
	return err
}
