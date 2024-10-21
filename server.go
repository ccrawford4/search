package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
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
