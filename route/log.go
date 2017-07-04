package route

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/go-chi/chi"
)

type History struct {
	Timestamp time.Time
	Input     string
}

func Log(w http.ResponseWriter, r *http.Request) {
	mongodbURL := os.Getenv("MONGODB_HOST")
	session, err := mgo.Dial(mongodbURL)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err.Error())
		return
	}

	defer session.Close()
	input, _ := url.QueryUnescape(chi.URLParam(r, "input"))
	err = session.DB(os.Getenv("MONGODB_DATABASE")).C("history").Insert(History{time.Now(), input})
	if err != nil {
		w.WriteHeader(500)
		log.Println(err.Error())
	}
}
