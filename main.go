package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
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
	idx = newDBIndex(connString, false, rsClient)
	router := gin.Default()
	if os.Getenv("ENV") == "development" {
		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"http://127.0.0.1:3000", "http://localhost:3000"},
			AllowMethods:     []string{"POST", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Add any other required headers here
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))
	}

	router.POST("/search", func(c *gin.Context) {
		type SearchRequestBody struct {
			SearchTerm string
		}

		var searchRequestBody SearchRequestBody
		if err := c.BindJSON(&searchRequestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		// get the searchTerm from the Request and then search the index for the term
		result := getTemplateData(&idx, searchRequestBody.SearchTerm)
		c.IndentedJSON(200, result)
	})

	url, err := parseURL("https://cs272-f24.github.io/top10/")
	if err != nil {
		log.Fatalf("Could not parse seed url: %v", err)
	}
	// Use the test crawl flag to avoid parsing robots.txt and delaying overall crawl time
	go crawl(&idx, url, true)
	err = router.Run(":8080")
	if err != nil {
		return
	}
}
