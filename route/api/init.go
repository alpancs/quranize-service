package api

import (
	"log"
	"os"

	"gopkg.in/mgo.v2"
)

var HistoryCollection *mgo.Collection

func init() {
	session, err := mgo.Dial(os.Getenv("MONGODB_HOST"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	HistoryCollection = session.DB(os.Getenv("MONGODB_DATABASE")).C("history")
}
