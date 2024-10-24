package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// executeTemplate creates a new template and then embeds the data within the TemplateData struct into the html
func executeTemplate(w http.ResponseWriter, fileContent string, templateData *TemplateData) {
	tmpl, err := template.New("demo").Parse(fileContent)
	if err != nil {
		log.Printf("Error parsing file content %v\n", err)
		return
	}
	if err = tmpl.Execute(w, templateData); err != nil {
		log.Printf("Error executing template %v\n", err)
	}
}

func searchHandler(idx *Index, c *gin.Context) {
	type SearchRequestBody struct {
		SearchTerm string
		Email      string
	}

	var searchRequestBody SearchRequestBody
	if err := c.BindJSON(&searchRequestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	// get the searchTerm from the Request and then search the index for the term
	result := getTemplateData(idx, searchRequestBody.SearchTerm)
	if result == nil {
		c.IndentedJSON(404, gin.H{"error": "No results found"})
	} else {
		c.IndentedJSON(200, result)
	}
}

// initHandlers initializes the handlers for the server
func initHandlers(idx *Index, router *gin.Engine) {
	router.GET("/documents/top10/*any", func(c *gin.Context) {
		corpusHandler(c.Writer, c.Request)
	})

	router.POST("/search", func(c *gin.Context) {
		searchHandler(idx, c)
	})
}

// startServer starts a gin server
func startServer(idx *Index) {
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
	initHandlers(idx, router)
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("[ERROR] Starting Server.\n")
	}
}

// corpusHandler to serve local documents
func corpusHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	if strings.HasSuffix(urlPath, "/") {
		urlPath += "index.html"
	}
	filePath := "." + urlPath
	fileContent, err := openAndReadFile(filePath)
	if err != nil {
		_, err = w.Write([]byte("404 No Page found!"))
		if err != nil {
			log.Printf("Could not serve 404 page %v\n", err)
		}
		return
	}
	_, err = w.Write(fileContent)
	if err != nil {
		log.Printf("Error writing response: %v", err.Error())
	}
}
