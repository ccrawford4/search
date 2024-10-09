package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
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
	server, exists := os.LookupEnv("server")
	if !exists {
		log.Fatalf("Error! Enviornment variable 'server' does not exist")
	}
	user, exists := os.LookupEnv("user")
	if !exists {
		log.Fatalf("Error! Enviornment variable 'user' does not exist")
	}
	password, exists := os.LookupEnv("password")
	if !exists {
		log.Fatalf("Error! Environment variable 'password' does not exist")
	}
	port, exists := os.LookupEnv("port")
	if !exists {
		log.Fatalf("Error! Enviornment variable 'port' does not exist")
	}
	database, exists := os.LookupEnv("database")
	if !exists {
		log.Fatalf("Error! Enviornment variable 'database' does not exist")
	}
	api, exists := os.LookupEnv("api")
	if !exists {
		log.Fatalf("Error! Environment variable 'api' does not exist")
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		server, user, password, port, database)

	// playground(connString)

	fmt.Printf("connString: %v\n", connString)
	var idx Index
	idx = newDBIndex(connString, false)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{api},                      // Allow requests from this origin
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
