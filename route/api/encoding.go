package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/alpancs/quranize/service"
	"github.com/go-chi/chi"
)

func Encode(w http.ResponseWriter, r *http.Request) {
	keyword, _ := url.QueryUnescape(chi.URLParam(r, "keyword"))
	json.NewEncoder(w).Encode(service.Encode(keyword))
}
