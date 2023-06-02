package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", replaceURLPOST)    //Обработчик для POST запроса
	http.HandleFunc("/{id}", replaceURLGET) //Обработчик для GET зарпоса

	err := http.ListenAndServe(`:8080`, nil)
	if err != nil {
		panic(err)
	}
}

func replaceURLPOST(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := make([]byte, r.ContentLength)
	_, err := r.Body.Read(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURL := "http://localhost:8080/EwHXdJfB"

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s\n", shortURL)
}

func replaceURLGET(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	originalURL := "https://practicum.yandex.ru/"

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
