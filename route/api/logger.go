package api

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

type History struct {
	Timestamp time.Time
	Keyword   string
}

func Log(w http.ResponseWriter, r *http.Request) {
	keyword, _ := url.QueryUnescape(chi.URLParam(r, "keyword"))
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return
	}

	err := HistoryCollection.Insert(History{time.Now(), keyword})
	if err != nil {
		w.WriteHeader(500)
		log.Println(err.Error())
	}
}
