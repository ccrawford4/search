package main

import (
	"context"
	"database/sql"
	_ "errors"
	"fmt"
	_ "github.com/microsoft/go-mssqldb"
	_ "gorm.io/driver/sqlite"
	_ "gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	var db *sql.DB
	var server = "devapisqlserver.database.windows.net"
	var port = 1433
	var user = "sqladmin"
	var password = "u?g'U$nj5^rKxeGJjJ9l"
	var database = "devdb1"

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	var err error
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Printf("Connected!")

	var idx Index
	idx = newDBIndex(connString, db)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://calm-field-07a2a211e.5.azurestaticapps.net"}, // Allow requests from this origin
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},                                                      // Allow these methods
		AllowHeaders:     []string{"Content-Type"},                                                                // Allow these headers
		AllowCredentials: true,                                                                                    // Allow cookies or other credentials
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
