package api

import (
	"encoding/json"
	"net/http"

	"github.com/alpancs/quranize/quran"
)

func Encode(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	json.NewEncoder(w).Encode(quran.Encode(keyword))
}
