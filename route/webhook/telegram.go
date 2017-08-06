package webhook

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Telegram(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
}
