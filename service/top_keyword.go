package service

import (
	"log"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func WatchTopKeywords() {
	for {
		TopKeywords = getTopKeywords()
		time.Sleep(1 * time.Hour)
	}
}

func getTopKeywords() []string {
	mongodbURL := os.Getenv("MONGODB_HOST")
	session, err := mgo.Dial(mongodbURL)
	if err != nil {
		log.Println(err.Error())
		return []string{}
	}
	defer session.Close()

	lastDay := bson.M{"timestamp": bson.M{"$gte": time.Now().AddDate(0, 0, -1)}}
	iter := session.DB(os.Getenv("MONGODB_DATABASE")).C("history").Find(lastDay).Iter()
	defer iter.Close()

	return rankKeywords(iter)
}

func rankKeywords(iter *mgo.Iter) []string {
	frequency := make(map[string]int)
	var doc struct{ Keyword string }
	for iter.Next(&doc) {
		frequency[doc.Keyword]++
	}

	keywords := []string{}
	for key := range frequency {
		keywords = append(keywords, key)
	}

	return keywords
}
