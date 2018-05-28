package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/alpancs/quranize-service/quran"
)

type Update struct {
	InlineQuery `json:"inline_query"`
}
type InlineQuery struct{ Id, Query string }

type InlineQueryAnswer struct {
	InlineQueryId string              `json:"inline_query_id"`
	Results       []InlineQueryResult `json:"results"`
}
type InlineQueryResult struct {
	Type                string `json:"type"`
	Id                  string `json:"id"`
	Title               string `json:"title"`
	InputMessageContent `json:"input_message_content"`
}
type InputMessageContent struct {
	MessageText string `json:"message_text"`
}

func Telegram(w http.ResponseWriter, r *http.Request) {
	var update Update
	json.NewDecoder(r.Body).Decode(&update)
	go reply(update.InlineQuery)
}

func reply(inlineQuery InlineQuery) {
	telegramURL := "https://api.telegram.org/bot" + os.Getenv("QURANIZE_TELEGRAM_TOKEN") + "/answerInlineQuery"
	results := []InlineQueryResult{}
	for i, encoded := range quran.Encode(inlineQuery.Query) {
		id := strconv.Itoa(i)
		results = append(results, InlineQueryResult{"article", id, encoded, InputMessageContent{encoded}})
	}
	answerBuffer := bytes.NewBuffer(nil)
	json.NewEncoder(answerBuffer).Encode(InlineQueryAnswer{inlineQuery.Id, results})
	http.Post(telegramURL, "application/json", answerBuffer)
}
