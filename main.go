package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/henimbola/shlk/db"
	"github.com/henimbola/shlk/services"
)

func main() {
	_, err := os.Stat("urls.kv")
	if err == nil {
		fmt.Println("File exists")
	} else if os.IsNotExist(err) {
		file, err := os.Create("urls.kv")

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()
	} else {
		fmt.Println("Error:", err)
	}

	db := &db.UrlDB{}
	app := fiber.New()

	// Initialize ShortURLMap
	urlMap := &services.ShortURLMap{
		Urls: make(map[string]string),
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

		longURL, err := db.GetUrl(shortURL)
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Short URL not found"})
		}

		return c.Redirect(longURL, http.StatusMovedPermanently)
	})

	app.Static("/", "./public")

	log.Fatal(app.Listen(":3000"))
}
