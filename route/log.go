package route

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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

	mongodbURL := os.Getenv("MONGODB_HOST")
	session, err := mgo.Dial(mongodbURL)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err.Error())
		return
	}

	defer session.Close()
	err = session.DB(os.Getenv("MONGODB_DATABASE")).C("history").Insert(History{bson.Now(), keyword})
	if err != nil {
		w.WriteHeader(500)
		log.Println(err.Error())
	}
}
