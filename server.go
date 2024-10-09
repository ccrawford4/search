package main

import (
	"log"
	"net/http"
	"strings"
)

// corpusHandler to serve local documents
func corpusHandler(w *http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	if strings.HasSuffix(urlPath, "/") {
		urlPath += "index.html"
	}
	filePath := "." + urlPath
	fileContent, err := openAndReadFile(filePath)
	if err != nil {
		_, err = (*w).Write([]byte("404 No Page found!"))
		if err != nil {
			log.Printf("Could not serve 404 page %v\n", err)
		}
		return
	}
	_, err = (*w).Write(fileContent)
	if err != nil {
		log.Printf("Error writing response: %v", err.Error())
	}
}
