package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	index := make(Index)
	urlWordTotals := make(Frequency)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},  // Allow requests from this origin
		AllowMethods:     []string{"GET", "POST", "OPTIONS"}, // Allow these methods
		AllowHeaders:     []string{"Content-Type"},           // Allow these headers
		AllowCredentials: true,                               // Allow cookies or other credentials
	}))

	stopWords := getStopWords()
	router.POST("/search", func(c *gin.Context) {
		type SearchRequestBody struct {
			SearchTerm string
		}

		var searchRequestBody SearchRequestBody
		if err := c.BindJSON(&searchRequestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		// get the searchTerm from the Request and then search the index for the term
		searchTerm := getStemmedWord(searchRequestBody.SearchTerm, stopWords)
		freq, found := search(&index, searchTerm, stopWords)

		// If no search term was found then serve the no_results document
		if !found {
			c.JSON(401, gin.H{})
			return
		}

		// Convert the Frequency found into templateData to be embedded into the html
		templateData := getTemplateData(freq, searchTerm, (float64)(len(index)), &urlWordTotals)
		// fileContent, _ := openAndReadFile("./static/search.html")
		// executeTemplate(*w, string(fileContent), templateData)
		c.IndentedJSON(200, templateData)
	})
	go crawl(&index, &urlWordTotals, parseURL("https://cs272-f24.github.io/top10/"), stopWords)
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
