package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alpancs/quranize-service/quran"
	"github.com/mssola/user_agent"
)

type Response struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

var (
	telegramAPI      = fmt.Sprintf("https://api.telegram.org/bot%s/", os.Getenv("QURANIZE_TELEGRAM_TOKEN"))
	hasTelegramToken = os.Getenv("QURANIZE_TELEGRAM_TOKEN") != ""
)

func Encode(w http.ResponseWriter, r *http.Request) {
	if pythonRequest(r) {
		json.NewEncoder(w).Encode([]string{"who is python?", "please contact the developer via Telegram @alpancs or email alpancs@gmail.com"})
		return
	}

	startTime := time.Now()
	encodeds := quran.Encode(r.URL.Query().Get("keyword"))
	if hasTelegramToken {
		go postToChannel(r, encodeds, time.Since(startTime))
	}

	json.NewEncoder(w).Encode(encodeds)
}

func pythonRequest(r *http.Request) bool {
	browserName, _ := user_agent.New(r.UserAgent()).Browser()
	return strings.Contains(strings.ToLower(browserName), "python")
}

func postToChannel(req *http.Request, resp []string, duration time.Duration) {
	keyword := req.URL.Query().Get("keyword")
	ua := user_agent.New(req.UserAgent())
	browserName, browserVersion := ua.Browser()
	browser := browserName + " " + browserVersion
	if ua.Mobile() {
		browser += " (mobile)"
	}
	msg := strings.Join([]string{
		fmt.Sprintf("keyword: %s, response: %s, duration: %d ms.", keyword, strings.Join(resp, ", "), duration.Milliseconds()),
		fmt.Sprintf("browser: %s, OS: %s.", browser, ua.OS()),
	}, "\n")

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
