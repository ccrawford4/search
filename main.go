package main

import (
	_ "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	_ "net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connString, exists := os.LookupEnv("AZURE_SQL_CONNECTIONSTRING")
	if !exists {
		log.Fatalf("No Connection AZURE_SQL_CONNECTIONSTRING Provided\n")
	}

	var idx Index
	idx = newDBIndex(connString, false)
	router := gin.Default()

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
