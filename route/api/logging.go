package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type History struct {
	Timestamp time.Time
	Keyword   string
}

func Log(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	keyword := strings.ToLower(strings.TrimSpace(string(data)))
	if keyword != "" {
		err := HistoryCollection.Insert(History{time.Now(), keyword})
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			fmt.Println(err.Error())
		}
	}
}
