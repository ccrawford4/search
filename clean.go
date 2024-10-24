package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

// ensureLeadingSlash ensures the href starts with a '/'.
func ensureLeadingSlash(href string) string {
	if !strings.HasPrefix(href, "/") {
		return "/" + href
	}
	return href
}

// clean takes a host URL and a href, and returns the fully formatted URL.
func clean(host, href string) (string, error) {
	relativeURL, err := parseURL(href)
	if err != nil {
		log.Printf("Could not parseHREF: %q\n", href)
		return href, err
	}

	// Ignore fragment urls
	if relativeURL.Fragment != "" || relativeURL.RawFragment != "" {
		return "", fmt.Errorf("ERROR! Url cannot be a fragment\n")
	}

	if relativeURL.RawQuery != "" {
		return "", fmt.Errorf("ERROR! URL With Query Params Cannot Be Cleaned %q\n", href)
	}

	// Return the href if it is already a full URL.
	if relativeURL.Scheme != "" {
		return href, nil
	}

	hostUrl, err := url.Parse(host)
	if err != nil {
		log.Printf("Could not parse host URL: %q\n", host)
		return href, err
	}

	// Ensure the href starts with a '/' if it's a relative path.
	href = ensureLeadingSlash(href)

	// Construct the full URL using the host's scheme and host.
	var source = "%s://%s%s"
	return fmt.Sprintf(source, hostUrl.Scheme, hostUrl.Host, href), nil
}
