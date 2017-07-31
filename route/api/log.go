package api

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type History struct {
	Timestamp time.Time
	Keyword   string
}

func Log(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword != "" {
		err := HistoryCollection.Insert(History{time.Now(), keyword})
		if err != nil {
			w.WriteHeader(500)
			log.Println(err.Error())
		}
	}
}
