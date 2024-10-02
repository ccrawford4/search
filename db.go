package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type Word struct {
	gorm.Model
	Name string `gorm:"unique"`
}
type Url struct {
	gorm.Model
	Name  string `gorm:"unique"`
	Count uint
}
type WordFrequencyRecord struct {
	gorm.Model
	Count uint
	Word  Word
	Url   Url
}

func migrateTables(db *gorm.DB) {
	err := db.AutoMigrate(&Word{}, &Url{}, &WordFrequencyRecord{})
	if err != nil {
		log.Fatalf("Error creating tables: %v\n", err)
	}
}

func connectToDB(dbName string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	migrateTables(db)
	return db
}

func insertRow[K Word | Url | WordFrequencyRecord](db *gorm.DB, object K) {
	result := db.Create(&object)
	if result.Error != nil {
		log.Fatalf("Error inserting row: %v\n", result.Error)
	}
	fmt.Printf("Inserted new row: %v\n", object)
}

func fillDB(db *gorm.DB, index *Index) {
	// word -> [url] -> count

}
