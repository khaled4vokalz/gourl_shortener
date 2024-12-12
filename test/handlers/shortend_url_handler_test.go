package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/khaled4vokalz/gourl_shortener/internal/db"
	"github.com/khaled4vokalz/gourl_shortener/internal/handlers"
)

func TestShortenUrlHandler(t *testing.T) {
	reqBody := map[string]string{
		"url": "https://example.com",
	}

	rr := performPOSTReq(reqBody, t)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200 but got %v", status)
	}

	expected := `{"shortened_url":"http://localhost:8080/19"}`
	if strings.TrimRight(rr.Body.String(), "\n") != expected {
		t.Errorf("Expected body %v but got %v", expected, rr.Body.String())
	}
}

func TestShortenUrlHandlerBadRequest(t *testing.T) {
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
	storage := db.NewInMemoryDb()

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
		handlers.ShortenUrlHandler(w, r, storage)
	})

	router.ServeHTTP(rr, req)
	return rr
}
