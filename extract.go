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
	cleanedData := strings.TrimFunc(n.Data, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	return n.Type == html.TextNode && len(cleanedData) > 1
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

// invalidNode returns true if the node is a style or script element
func invalidNode(n *html.Node) bool {
	return n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style")
}

// sanitizeHTML takes in an HTML node and removes any script or style content from the tree
func sanitizeHTML(n *html.Node) {
	if n == nil {
		return
	}

	var toRemove []*html.Node
	// Repeat the process for all the sibling nodes
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if invalidNode(c) {
			toRemove = append(toRemove, c)
		} else {
			sanitizeHTML(c)
		}
	}
	for _, node := range toRemove {
		if node.Parent != nil {
			node.Parent.RemoveChild(node)
		}
	}
}

func validTitleNode(n *html.Node, currTitle string) bool {
	return n.Type == html.ElementNode && n.Data == "title" && currTitle == ""
}

func containsOGDescription(n *html.Node) bool {
	for _, attr := range n.Attr {
		if attr.Key == "property" && attr.Val == "og:description" {
			return true
		}
	}
	return false
}
func containsDescription(n *html.Node) bool {
	for _, attr := range n.Attr {
		if attr.Key == "name" && attr.Val == "description" {
			return true
		}
	}
	return false
}

func validMetaDescription(n *html.Node, currDescription string) bool {
	return currDescription == "" && n.Type == html.ElementNode && n.Data == "meta" &&
		(containsOGDescription(n) || containsDescription(n))
}

func getDescription(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "content" {
			return attr.Val
		}
	}
	return ""
}

// extract takes in an array of bytes from an HTML page and returns two slices of type string.
// The first slice returned is the list of words found in the document.
// The second slice returned is the list of hrefs found in the document.
func extract(text []byte) ([]string, []string, string, string) {
	var words []string
	var hrefs []string
	var title string
	var description string

	reader := bytes.NewReader(text)
	doc, err := html.Parse(reader)
	if err != nil {
		log.Fatalf("HTML parse error: %v\n", err)
	}

	var processDocument func(*html.Node)
	processDocument = func(n *html.Node) {
		if validMetaDescription(n, description) { // Pull the description from the meta tags
			description = getDescription(n)
		} else if validTitleNode(n, title) { // Pull the valid title node
			title = n.FirstChild.Data
		} else if validTextNode(n) { // Pull a valid text node
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

	return words, hrefs, title, description
}
