package main

import (
	"github.com/emirpasic/gods/sets/hashset"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// serveDocument takes in a file path and a
// *http.ResponseWriter and serves the HTML document to the client
func serveDocument(path string, w *http.ResponseWriter) {
	fileContent, err := openAndReadFile(path)
	if err != nil {
		http.Error(*w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = (*w).Write(fileContent)
	if err != nil {
		log.Fatalf("Could not serve file content %v\n", err)
	}
}

// getParam takes in a *http.Request and a key
// and returns the value associated with that key
func getParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// initSearchHandler creates the http handler and
// handles logic for serving the data on /search requests
func initSearchHandler(index *Index, w *http.ResponseWriter, r *http.Request, stopWords *hashset.Set) {
	// get the searchTerm from the Request and then search the index for the term
	searchTerm := getStemmedWord(getParam(r, "term"), stopWords)
	freq, found := search(index, searchTerm, stopWords)

	// If no search term was found then serve the no_results document
	if !found {
		serveDocument("./static/no_results.html", w)
		return
	}

	// Convert the Frequency found into templateData to be embedded into the html
	templateData := getTemplateData(freq, searchTerm)
	fileContent, _ := openAndReadFile("./static/search.html")
	executeTemplate(*w, string(fileContent), templateData)
}

// executeTemplate creates a new template and then embeds the data within the TemplateData struct into the html
func executeTemplate(w http.ResponseWriter, fileContent string, templateData *TemplateData) {
	tmpl, err := template.New("demo").Parse(fileContent)
	if err != nil {
		log.Fatalf("Erorr parsing file content %v\n", err)
	}
	err = tmpl.Execute(w, templateData)
	if err != nil {
		return
	}
}

// initHandlers creates all the http handlers for the web server
func initHandlers(index *Index, stopWords *hashset.Set) {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/no_results", http.FileServer(http.Dir("static")))
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		initSearchHandler(index, &w, r, stopWords)
	})
	http.HandleFunc("/documents/top11/", func(w http.ResponseWriter, r *http.Request) {
		corpusHandler(&w, r)
	})
}

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
		return
	}
	_, err = (*w).Write(fileContent)
	if err != nil {
		log.Fatalf("Error writing response: %v", err.Error())
	}
}

// startServer starts a server on localhost with the desired port
func startServer(port string) {
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting server:%v\n", err)
	}
}
