package route

import (
	"encoding/json"
	"net/http"

	"github.com/alpancs/quranize/service"
	"github.com/pressly/chi"
)

func Encode(w http.ResponseWriter, r *http.Request) {
	text := chi.URLParam(r, "text")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(service.Encode(text))
}
