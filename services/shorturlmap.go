package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"sync"
)

type ShortURLMap struct {
	sync.RWMutex
	Urls map[string]string
}

func (s *ShortURLMap) ShortenURL(longURL string) string {
	// Generate a short URL
	b := make([]byte, 5)
	rand.Read(b)
	shortURL := base64.URLEncoding.EncodeToString(b)

	// Store the short URL and its corresponding long URL in the map
	s.Urls[shortURL] = longURL

	// Open the file for appending
	file, err := os.OpenFile("urls.kv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Append the long URL and its corresponding short URL to the file
	_, err = fmt.Fprintf(file, "\"%s\":\"%s\"\n", shortURL, longURL)
	if err != nil {
		log.Fatal(err)
	}

	return shortURL
}

// ResolveURL returns the original long URL given a short URL
func (m *ShortURLMap) ResolveURL(shortURL string) (string, error) {
	m.RLock()
	defer m.RUnlock()

	longURL, ok := m.Urls[shortURL]
	if !ok {
		return "", fmt.Errorf("short URL not found")
	}

	return longURL, nil
}
