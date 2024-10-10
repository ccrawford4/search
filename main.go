package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func init() {
	// Only load .env file in development environment
	if os.Getenv("ENV") != "production" {
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
	err = createSearchResult(rsClient, "romeo", &SearchResult{
		Frequency{
			"https://example1.com": 32,
			"https://example2.com": 13,
			"https://example3.com": 4,
		},
		10,
		false,
	})
	if err != nil {
		log.Fatalf("Error creating search result: %v\n", err)
	}
	log.Printf("Successfully created romeo object.\n")
	result, err := getSearchResult(rsClient, "romeo")
	if err != nil {
		log.Fatalf("Error getting search result for word romeo %v\n", err)
	}
	fmt.Printf("Result: %v\n", result)
	//var idx Index
	//idx = newDBIndex(connString, false)
	//router := gin.Default()
	//
	//router.POST("/search", func(c *gin.Context) {
	//	type SearchRequestBody struct {
	//		SearchTerm string
	//	}
	//
	//	var searchRequestBody SearchRequestBody
	//	if err := c.BindJSON(&searchRequestBody); err != nil {
	//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	}
	//
	//	// get the searchTerm from the Request and then search the index for the term
	//	result := getTemplateData(&idx, searchRequestBody.SearchTerm)
	//	c.IndentedJSON(200, result)
	//})
	//
	//url, err := parseURL("https://cs272-f24.github.io/top10/")
	//if err != nil {
	//	log.Fatalf("Could not parse seed url: %v", err)
	//}
	//go crawl(&idx, url)
	//err = router.Run(":8080")
	//if err != nil {
	//	return
	//}
}
