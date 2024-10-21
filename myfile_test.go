package main

import (
	"fmt"
	"log"
	"net/http"
)

func init() {
	go func() {
		port := 8080
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			log.Fatalf("Error starting server: %v", err.Error())
		}
	}()

	http.HandleFunc("/documents/top10/", func(w http.ResponseWriter, r *http.Request) {
		corpusHandler(w, r)
	})
}
