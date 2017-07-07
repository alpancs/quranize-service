package service

import (
	"log"
	"os"
	"sort"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type KeywordScore struct {
	Keyword string
	Score   int
}
type KeywordScores []KeywordScore

func (keywordScores KeywordScores) Len() int {
	return len(keywordScores)
}
func (keywordScores KeywordScores) Less(i, j int) bool {
	if keywordScores[i].Score == keywordScores[j].Score {
		return keywordScores[i].Keyword < keywordScores[j].Keyword
	}
	return keywordScores[i].Score > keywordScores[j].Score
}
func (keywordScores KeywordScores) Swap(i, j int) {
	keywordScores[i], keywordScores[j] = keywordScores[j], keywordScores[i]
}

const DEFAULT_TRENDING_KEYWORD_LIMIT = 100

var lastId string

func WatchTrendingKeywords() {
	for {
		UpdateTrendingKeywords()
		time.Sleep(5 * time.Minute)
	}
}

func UpdateTrendingKeywords() {
	startTime := time.Now()
	if !needToUpdate() {
		return
	}

	session, err := mgo.Dial(os.Getenv("MONGODB_HOST"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer session.Close()

	last7Days := bson.M{"timestamp": bson.M{"$gt": time.Now().AddDate(0, 0, -7)}}
	iter := session.DB(os.Getenv("MONGODB_DATABASE")).C("history").Find(last7Days).Iter()
	defer iter.Close()

	TrendingKeywords = getTrendingKeywords(iter)
	log.Println("update trending keywords elapsed time:", time.Since(startTime))
}

func getTrendingKeywords(iter *mgo.Iter) []string {
	frequency := make(map[string]int)
	var doc struct{ Keyword string }
	for iter.Next(&doc) {
		frequency[doc.Keyword]++
	}

	keywordScores := KeywordScores{}
	for keyword, score := range frequency {
		keywordScores = append(keywordScores, KeywordScore{keyword, score})
	}
	sort.Sort(keywordScores)

	trendingKeywords := []string{}
	for _, keywordScore := range keywordScores {
		if len(trendingKeywords) == DEFAULT_TRENDING_KEYWORD_LIMIT {
			break
		}
		trendingKeywords = append(trendingKeywords, keywordScore.Keyword)
	}

	return trendingKeywords
}

func needToUpdate() bool {
	session, err := mgo.Dial(os.Getenv("MONGODB_HOST"))
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer session.Close()

	var doc struct {
		Id string `bson:"_id"`
	}
	session.DB(os.Getenv("MONGODB_DATABASE")).C("history").Find(nil).Sort("-timestamp").Limit(1).One(&doc)

	defer func() { lastId = doc.Id }()
	return lastId != doc.Id
}
