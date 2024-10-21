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

	var source = "%s://%s%s"
	if strings.Contains(source, "://www") {
		source = "%s://www.%s%s"
	}
	// Construct the full URL using the host's scheme and host.
	return fmt.Sprintf(source, hostUrl.Scheme, hostUrl.Host, href), nil
}
