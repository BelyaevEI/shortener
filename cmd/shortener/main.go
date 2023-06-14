package main

import (
	"log"

	"github.com/BelyaevEI/shortener/internal/app"
)

func main() {

	if err := app.RunServer(); err != nil {
		log.Fatal(err)
	}
}
