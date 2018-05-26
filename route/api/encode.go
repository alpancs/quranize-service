package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	go postToChannel(keyword)
}

func postToChannel(keyword string) {
	reqBody, err := json.Marshal(Response{"@quranize", keyword})
	if err != nil {
		fmt.Println(err)
		return
	}
	url := telegramAPI + "sendMessage"
	res, err := http.Post(url, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		fmt.Println(err)
		return
	}
	resCode := res.StatusCode
	if resCode != 200 {
		resBody, _ := ioutil.ReadAll(res.Body)
		fmt.Printf("URL: %s. request body: %s. response code: %d. response body: %s", url, string(reqBody), resCode, string(resBody))
	}
}
