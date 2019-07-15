package api

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alpancs/quranize-service/db"
)

func Log(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	keyword := strings.ToLower(strings.TrimSpace(string(data)))
	if keyword == "" {
		return
	}

	_, err := db.Exec(
		`INSERT INTO history (keyword, timestamp) VALUES ($1, $2)`,
		keyword,
		time.Now().In(time.UTC),
	)
	if err != nil {
		log.Println(err)
		errorCode := http.StatusInternalServerError
		http.Error(w, http.StatusText(errorCode), errorCode)
	}
}
