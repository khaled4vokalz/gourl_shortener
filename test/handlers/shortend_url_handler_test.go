package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/khaled4vokalz/gourl_shortener/internal/config"
	"github.com/khaled4vokalz/gourl_shortener/internal/db"
	"github.com/khaled4vokalz/gourl_shortener/internal/handlers"
)

// MockStorage for testing
type MockCache struct{}

func (m *MockCache) Set(shortened, originalUrl string, ttl time.Duration) error {
	return nil
}

func (m *MockCache) Get(shortened string) (string, error) {
	return "", nil
}

func (m *MockCache) IsAlive() bool {
	return true
}

func TestShortenUrlHandler(t *testing.T) {
	backup := config.GetConfig
	defer func() { config.GetConfig = backup }()
	config.GetConfig = func() *config.Config {
		return &config.Config{
			Environment: "dev",
		}
	}

	reqBody := map[string]string{
		"url": "https://example.com",
	}

	rr := performPOSTReq(reqBody, t)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 but got %v", status)
	}

	expected := `{"shortened_url":"http://localhost:8080/sRVrAaAE"}`
	if strings.TrimRight(rr.Body.String(), "\n") != expected {
		t.Errorf("Expected body %v but got %v", expected, rr.Body.String())
	}
}

func TestShortenUrlHandlerBadRequest(t *testing.T) {
	backup := config.GetConfig
	defer func() { config.GetConfig = backup }()
	config.GetConfig = func() *config.Config {
		return &config.Config{
			Environment: "dev",
		}
	}

	reqBody := map[string]string{
		"url": "https://",
	}

	rr := performPOSTReq(reqBody, t)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status 400 but got %v", status)
	}
	expected := "Invalid URL"
	if strings.TrimRight(rr.Body.String(), "\n") != expected {
		t.Errorf("Expected body %v but got %v", expected, rr.Body.String())
	}
}

func performPOSTReq(reqBody map[string]string, t *testing.T) *httptest.ResponseRecorder {
	storage, err := db.NewInMemoryDb()

	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(reqBodyJson))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	router := http.NewServeMux()
	router.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		handlers.ShortenUrlHandler(w, r, storage, &MockCache{}, config.ShortenerSettings{Length: 6})
	},
	)

	router.ServeHTTP(rr, req)
	return rr
}
