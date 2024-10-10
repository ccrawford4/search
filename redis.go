package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

const keyPrefix = "term:"

func getRSClient(redisHost, redisPassword string) (*redis.Client, error) {
	op := &redis.Options{
		Addr:     redisHost,
		Password: redisPassword,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12},
		WriteTimeout: 5 * time.Second}
	client := redis.NewClient(op)
	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("failed to connect with redis instance at %s - %v", redisHost, err)
		return client, err
	}
	log.Printf("successfully connected with redis instance at %s", redisHost)
	return client, nil
}

func createSearchResult(rsClient *redis.Client, searchTerm string, searchResult *SearchResult) error {
	// Marshal the searchResult to JSON
	data, err := json.Marshal(searchResult)
	if err != nil {
		return err
	}

	// Set the JSON data as a value for the search term key
	err = rsClient.Set(context.Background(), keyPrefix+searchTerm, data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func getSearchResult(rsClient *redis.Client, searchTerm string) (*SearchResult, error) {
	result, err := rsClient.Get(context.Background(), keyPrefix+searchTerm).Result()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		log.Printf("Empty response from cache for key %q: ", searchTerm)
		return nil, nil
	}
	fmt.Printf("Result: %v\n", result)
	var searchResult SearchResult
	err = json.Unmarshal([]byte(result), &searchResult)
	if err != nil {
		return nil, err
	}
	return &searchResult, nil
}
