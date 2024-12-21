package tests

import (
	"testing"

	"github.com/khaled4vokalz/gourl_shortener/internal/utils"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"http://example.com", true},
		{"https://example.com", true},
		{"ftp://example.com", false},
		{"example.com", false},
		{"http://", false},
		{"https://www.google.com", true},
		{"invalid-url", false},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			result := utils.IsValidURL(test.url)
			if result != test.expected {
				t.Errorf("For URL %s: expected %v, got %v", test.url, test.expected, result)
			}
		})
	}
}
