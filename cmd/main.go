package main

import (
	"log"

	app "github.com/peterkuchinov/The-link-shortener-on-Golang/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
