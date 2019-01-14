package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/alpancs/quranize-service/quran"
	"github.com/mssola/user_agent"
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
	go postToChannel(r, keyword)

	if pythonRequest(r) {
		json.NewEncoder(w).Encode([]string{"who is python-requests?", "please contact the developer via Telegram @alpancs or email alpancs@gmail.com"})
		return
	}

	json.NewEncoder(w).Encode(quran.Encode(keyword))
}

func pythonRequest(r *http.Request) bool {
	browserName, _ := user_agent.New(r.UserAgent()).Browser()
	return browserName == "python-requests"
}

func postToChannel(r *http.Request, keyword string) {
	if os.Getenv("ENV") != "production" {
		return
	}

	ua := user_agent.New(r.UserAgent())
	browserName, browserVersion := ua.Browser()
	browser := browserName + " " + browserVersion
	if ua.Mobile() {
		browser += " (mobile)"
	}
	msg := fmt.Sprintf("keyword: %s\nbrowser: %s, OS: %s", keyword, browser, ua.OS())

	reqBody, err := json.Marshal(Response{"@quranize", msg})
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
		fmt.Printf("URL: %s. request body: %s. response code: %d. response body: %s\n", url, string(reqBody), resCode, string(resBody))
	}
}
