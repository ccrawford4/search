package main

import (
	"fmt"
)

func main() {
	countMap := search("https://cs272-f24.github.io/tests/rnj/sceneI_30.0.html", "Verona")
	fmt.Println("Searching for Verona:")
	for key, value := range countMap {
		fmt.Printf("%q: %d\n", key, value)
	}

	fmt.Println("Searching for Benvolio")
	countMap = search("https://cs272-f24.github.io/tests/rnj/sceneI_30.1.html", "Benvolio")
	for key, value := range countMap {
		fmt.Printf("%q: %d\n", key, value)
	}

	fmt.Println("Searching for Romeo")
	countMap = search("https://cs272-f24.github.io/tests/rnj/", "Romeo")
	for key, value := range countMap {
		fmt.Printf("%q: %d\n", key, value)
	}
}
