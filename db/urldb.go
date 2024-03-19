package db

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type UrlDB struct {
}

func (u *UrlDB) GetUrl(key string) (string, error) {
	file, err := os.Open("urls.kv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		if parts[0] == key {
			return parts[1], nil
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return "", fmt.Errorf("short URL not found")
}
