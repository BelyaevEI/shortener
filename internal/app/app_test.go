package app_test

import (
	"log"

	"github.com/BelyaevEI/shortener/internal/app"
)

func Example() {

	// This function init and run shortner service
	if err := app.RunServer(); err != nil {
		log.Fatal(err)
	}
}
