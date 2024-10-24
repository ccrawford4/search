package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	// Only load .env file in development environment
	if os.Getenv("ENV") == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v\n", err)
		}
	}
}

func main() {
	connString, exists := os.LookupEnv("AZURE_SQL_CONNECTIONSTRING")
	if !exists {
		log.Fatalf("No Connection AZURE_SQL_CONNECTIONSTRING Provided\n")
	}
	log.Printf("Connecting string loaded successfully: %v\n", connString)
	redisHost, exists := os.LookupEnv("REDIS_HOST")
	if !exists {
		log.Fatalf("No Redis Host Provided\n")
	}
	redisPassword, exists := os.LookupEnv("REDIS_PASSWORD")
	if !exists {
		log.Fatalf("No Redis Password Provided\n")
	}
	rsClient, err := getRSClient(redisHost, redisPassword)
	if err != nil {
		log.Fatalf("Error getting Redis Client: %v\n", err)
	}

	var idx Index
	// For Production
	idx = newDBIndex(connString, false, rsClient)

	// for testing
	// idx = newDBIndex("dev.db", true, nil)
	go startServer(&idx)

	// For production
	go crawl(&idx, "https://usfca.edu/", false)
}
