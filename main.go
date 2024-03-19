package main

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "log"
    "net/http"
    "sync"

    "github.com/gofiber/fiber/v2"
)

// ShortURLMap holds mappings between short and long URLs
type ShortURLMap struct {
    sync.RWMutex
    urls map[string]string
}

// ShortenURL generates a short URL and maps it to the original URL
func (m *ShortURLMap) ShortenURL(longURL string) string {
    m.Lock()
    defer m.Unlock()

    shortURL := generateShortURL()
    m.urls[shortURL] = longURL

    return shortURL
}

// ResolveURL returns the original long URL given a short URL
func (m *ShortURLMap) ResolveURL(shortURL string) (string, error) {
    m.RLock()
    defer m.RUnlock()

    longURL, ok := m.urls[shortURL]
    if !ok {
        return "", fmt.Errorf("short URL not found")
    }

    return longURL, nil
}

func generateShortURL() string {
    b := make([]byte, 6)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}

func main() {
    app := fiber.New()

    // Initialize ShortURLMap
    urlMap := &ShortURLMap{
        urls: make(map[string]string),
    }

    // POST /shorten
    app.Post("/shorten", func(c *fiber.Ctx) error {
        type request struct {
            LongURL string `json:"long_url"`
        }

        var req request
        if err := c.BodyParser(&req); err != nil {
            return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
        }

        shortURL := urlMap.ShortenURL(req.LongURL)
        return c.JSON(fiber.Map{"short_url": shortURL})
    })

    // GET /:shortURL
    app.Get("/:shortURL", func(c *fiber.Ctx) error {
        shortURL := c.Params("shortURL")
        longURL, err := urlMap.ResolveURL(shortURL)
        if err != nil {
            return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Short URL not found"})
        }

        return c.Redirect(longURL, http.StatusMovedPermanently)
    })

    log.Fatal(app.Listen(":3000"))
}
