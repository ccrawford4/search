package main

import (
	"bytes"
	"golang.org/x/net/html"
	"log"
	"strings"
	"unicode"
)

// validAnchorElement returns true if the node is a valid anchor element
func validAnchorElement(n *html.Node) bool {
	return n.Data == "a"
}

// validTextNode returns true if the node is a valid text node
func validTextNode(n *html.Node) bool {
	return n.Type == html.TextNode
}

// getWords takes in an HTMl text node and returns a slice of strings containing each word in its data.
func getWords(n *html.Node) []string {
	var words []string
	text := n.Data

	// Don't accept characters that are punctuation or spaces
	f := func(c rune) bool {
		return unicode.IsPunct(c) || unicode.IsSpace(c)
	}

	cleanedWords := strings.FieldsFunc(text, f)
	for _, word := range cleanedWords {
		words = append(words, word)
	}
	return words
}

// getHref takes in an anchor HTML node and returns the url from its href if it exists.
// If a href was not found in the anchor tag then the function returns false.
func getHref(n *html.Node) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			return attr.Val, true
		}
	}
	return "", false
}

// sanitizeHTML takes in an HTML node and removes any script or style content from the tree
func sanitizeHTML(n *html.Node) {
	if n == nil {
		return
	}

	// If the element is a script or style node then remove it
	if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
		if n.Parent != nil {
			n.Parent.RemoveChild(n)
		}
		return
	}

	// Repeat the process for all the sibling nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sanitizeHTML(c)
	}
}

// extract takes in an array of bytes from an HTML page and returns two slices of type string.
// The first slice returned is the list of words found in the document.
// The second slice returned is the list of hrefs found in the document.
func extract(text []byte) ([]string, []string) {
	var words []string
	var hrefs []string

	reader := bytes.NewReader(text)
	doc, err := html.Parse(reader)
	if err != nil {
		log.Fatalf("HTML parse error: %v\n", err)
	}

	var processDocument func(*html.Node)
	processDocument = func(n *html.Node) {
		// For text nodes, extract the words from the data
		if validTextNode(n) {
			words = append(words, getWords(n)...)
		} else if validAnchorElement(n) {
			// For anchor elements, try and get the href from the attributes
			href, foundHref := getHref(n)
			if foundHref {
				hrefs = append(hrefs, href)
			}
		}

		// Repeat the process for the sibling nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processDocument(c)
		}
	}
	// Sanitize the HTML to get rid of style and script content, and then perform the processing
	sanitizeHTML(doc)
	processDocument(doc)

	return words, hrefs
}
