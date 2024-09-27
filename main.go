package main

import "github.com/gin-gonic/gin"

func main() {
	index := make(Index)
	router := gin.Default()
	stopWords := getStopWords()
	router.GET("/search", func(c *gin.Context) {
		term, found := c.Params.Get("term")
		if !found {
			c.JSON(400, gin.H{})
		}
		// get the searchTerm from the Request and then search the index for the term
		searchTerm := getStemmedWord(term, stopWords)
		freq, found := search(&index, searchTerm, stopWords)

		// If no search term was found then serve the no_results document
		if !found {
			c.JSON(401, gin.H{})
			return
		}

		// Convert the Frequency found into templateData to be embedded into the html
		templateData := getTemplateData(freq, searchTerm)
		// fileContent, _ := openAndReadFile("./static/search.html")
		// executeTemplate(*w, string(fileContent), templateData)
		c.IndentedJSON(200, templateData)
	})
	go router.Run(":8080")
	crawl(&index, parseURL("https://cs272-f24.github.io/top10/"), stopWords)
}
