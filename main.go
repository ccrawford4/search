package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	server, exists := os.LookupEnv("SQL_SERVER_NAME")
	if !exists {
		log.Fatalf("Error! Enviornment variable SQL_SERVER_NAME does not exist")
	}
	user, exists := os.LookupEnv("SQL_USER")
	if !exists {
		log.Fatalf("Error! Enviornment variable SQL_USER does not exist")
	}
	password, exists := os.LookupEnv("SQL_USER_PASSWORD")
	if !exists {
		log.Fatalf("Error! Environment variable SQL_USER_PASSWORD does not exist")
	}
	port, exists := os.LookupEnv("SQL_SERVER_PORT")
	if !exists {
		log.Fatalf("Error! Enviornment variable SQL_SERVER_PORT does not exist")
	}
	database, exists := os.LookupEnv("SQL_DATABASE")
	if !exists {
		log.Fatalf("Error! Enviornment variable SQL_DATABASE does not exist")
	}
	testAPI, exists := os.LookupEnv("TEST_API_ENDPOINT")
	if !exists {
		log.Fatalf("Error! Environment variable TEST_API_ENDPOINT does not exist")
	}
	api, exists := os.LookupEnv("API_ENDPOINT")
	if !exists {
		log.Fatalf("Error! Environment variable API_ENDPOINT does not exist")
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		server, user, password, port, database)

	fmt.Printf("connString: %v\n", connString)
	var idx Index
	idx = newDBIndex(connString, true)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{testAPI, api},             // Allow requests from this origin
		AllowMethods:     []string{"GET", "POST", "OPTIONS"}, // Allow these methods
		AllowHeaders:     []string{"Content-Type"},           // Allow these headers
		AllowCredentials: true,                               // Allow cookies or other credentials
	}))

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
	go crawl(&idx, url)
	err = router.Run(":8080")
	if err != nil {
		return
	}
}
