package app

import (
	"net/http"

	"github.com/BelyaevEI/shortener/internal/handlers"
)

func Shortener(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		handlers.ReplacePOST(w, r)
	case http.MethodGet:
		handlers.ReplaceGET(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
