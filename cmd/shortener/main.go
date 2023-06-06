package main

import (
	"net/http"

	"github.com/BelyaevEI/shortener/internal/app"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, app.Shortener)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
