package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

type History struct {
	Timestamp time.Time
	Keyword   string
}

var historyCollection = newMongodbCollection()

func Log(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	keyword := strings.ToLower(strings.TrimSpace(string(data)))
	if keyword != "" {
		err := historyCollection.Insert(History{time.Now(), keyword})
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			fmt.Println(err.Error())
		}
	}
}

func newMongodbCollection() *mgo.Collection {
	session, err := mgo.Dial(os.Getenv("MONGODB_HOST"))
	if err != nil {
		panic(err)
	}
	return session.DB(os.Getenv("MONGODB_DATABASE")).C("history")
}
