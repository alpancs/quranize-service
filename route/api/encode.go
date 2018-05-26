package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/alpancs/quranize/quran"
)

type Response struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

var (
	telegramAPI = fmt.Sprintf("https://api.telegram.org/bot%s/", os.Getenv("QURANIZE_TELEGRAM_TOKEN"))
)

func Encode(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	json.NewEncoder(w).Encode(quran.Encode(keyword))
	go tellOwner(keyword)
}

func tellOwner(keyword string) {
	reqBody, err := json.Marshal(Response{"@alpancs", keyword})
	if err != nil {
		fmt.Println(err)
		return
	}
	url := telegramAPI + "sendMessage"
	_, err = http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		fmt.Println(err)
	}
}
