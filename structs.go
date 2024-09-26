package main

// Frequency can be used to map URLs to TF-IDF scores or
// to map words to their frequency counts
type Frequency[V float64 | int] map[string]V

// Index map[word][URL]TFIDF score
type Index map[string]Frequency[float64]

// Hit struct is used to format the template response
// URL = the URL that a word is pulled from
// TFIDF = the TF-IDF score (relevancy) of the document for the search
type Hit struct {
	URL   string
	TFIDF float64
}

// Hits is an array of Hit objects
type Hits []Hit

// Len computes the length of the Hits array
func (hits Hits) Len() int {
	return len(hits)
}

// Less compares two Hit objects within
// the Hits array
func (hits Hits) Less(i, j int) bool {
	hitA, hitB := hits[i], hits[j]
	if hitA.TFIDF == hitB.TFIDF {
		return hitB.URL > hitA.URL
	}
	return hitA.TFIDF > hitB.TFIDF
}

// Swap swaps two Hit items in a Hits object
func (hits Hits) Swap(i, j int) {
	hits[i], hits[j] = hits[j], hits[i]
}

// TemplateData is the data structure executed on the
// template for the web server. HITS is a Hits object
// containing the relevant hits for a search. TERM is the search term.
type TemplateData struct {
	HITS Hits
	TERM string
}
