package job

import (
	"log"
	"os"
	"sort"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Document struct {
	Id      string `bson:"_id"`
	Keyword string
}

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

var (
	TrendingKeywords []string

	lastId  string
	history *mgo.Collection
)

func init() {
	session, err := mgo.Dial(os.Getenv("MONGODB_HOST"))
	if err != nil {
		log.Println(err.Error())
		return
	}
	history = session.DB(os.Getenv("MONGODB_DATABASE")).C("history")

	UpdateTrendingKeywords()
	go WatchTrendingKeywords()
}

func WatchTrendingKeywords() {
	for {
		time.Sleep(1 * time.Minute)
		if needToUpdate() {
			UpdateTrendingKeywords()
		}
	}
}

func UpdateTrendingKeywords() {
	startTime := time.Now()
	last7Days := bson.M{"timestamp": bson.M{"$gt": time.Now().AddDate(0, 0, -7)}}
	iter := history.Find(last7Days).Sort("timestamp").Iter()
	defer iter.Close()
	TrendingKeywords = getTrendingKeywords(iter)
	log.Println("update trending keywords elapsed time:", time.Since(startTime))
}

func getTrendingKeywords(iter *mgo.Iter) []string {
	frequency := make(map[string]int)
	var doc Document
	for iter.Next(&doc) {
		frequency[doc.Keyword]++
		lastId = doc.Id
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
	var lastDoc Document
	history.Find(nil).Sort("-timestamp").Limit(1).One(&lastDoc)
	return lastId != lastDoc.Id
}
