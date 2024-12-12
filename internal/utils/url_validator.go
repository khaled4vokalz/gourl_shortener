package utils

import (
	"net/url"
	"regexp"
)

func IsValidURL(urlStr string) bool {
	parsedURL, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return false
	}

	validScheme := regexp.MustCompile("^(http|https)$")
	if !validScheme.MatchString(parsedURL.Scheme) {
		return false
	}

	return parsedURL.Host != ""
}
