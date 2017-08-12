package job

import (
	"fmt"
	"os"
	"sort"
	"strings"
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

const DEFAULT_TRENDING_KEYWORD_LIMIT = 100

var (
	TrendingKeywords []string

	lastId            string
	last7Days         = bson.M{"timestamp": bson.M{"$gt": time.Now().AddDate(0, 0, -7)}}
	historyCollection *mgo.Collection
)

func Start() {
	session, err := mgo.Dial(os.Getenv("MONGODB_HOST"))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	historyCollection = session.DB(os.Getenv("MONGODB_DATABASE")).C("history")

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
	iter := historyCollection.Find(last7Days).Sort("timestamp").Iter()
	defer iter.Close()
	TrendingKeywords = getTrendingKeywords(iter)
	fmt.Println("trending keywords updated in ", time.Since(startTime))
}

func getTrendingKeywords(iter *mgo.Iter) []string {
	frequency := make(map[string]int)
	var doc Document
	for iter.Next(&doc) {
		frequency[strings.Replace(doc.Keyword, " ", "", -1)]++
		lastId = doc.Id
	}

	keywordScores := []KeywordScore{}
	for keyword, score := range frequency {
		keywordScores = append(keywordScores, KeywordScore{keyword, score})
	}
	sort.Sort(ByTopScore(keywordScores))

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
	last7Days = bson.M{"timestamp": bson.M{"$gt": time.Now().AddDate(0, 0, -7)}}
	historyCollection.Find(last7Days).Sort("-timestamp").Limit(1).One(&lastDoc)
	return lastId != lastDoc.Id
}

type ByTopScore []KeywordScore

func (keywordScores ByTopScore) Len() int {
	return len(keywordScores)
}
func (keywordScores ByTopScore) Less(i, j int) bool {
	if keywordScores[i].Score == keywordScores[j].Score {
		return keywordScores[i].Keyword < keywordScores[j].Keyword
	}
	return keywordScores[i].Score > keywordScores[j].Score
}
func (keywordScores ByTopScore) Swap(i, j int) {
	keywordScores[i], keywordScores[j] = keywordScores[j], keywordScores[i]
}
