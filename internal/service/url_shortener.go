package service

import "fmt"

func GenerateShortenedURL(original string) string {
	// TODO: starting with just the `length` of the provided url.
	// but this is never a solid solution :D :D
	// plan is to generate some sort of unique key for each URL, some sort of hash
	return fmt.Sprintf("%v", len(original))
}
