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

const DEFAULT_TOP_KEYWORD_LIMIT = 100

func UpdateTopKeywords() {
	startTime := time.Now()
	mongodbURL := os.Getenv("MONGODB_HOST")
	session, err := mgo.Dial(mongodbURL)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer session.Close()

	last7Days := bson.M{"timestamp": bson.M{"$gt": time.Now().AddDate(0, 0, -7)}}
	iter := session.DB(os.Getenv("MONGODB_DATABASE")).C("history").Find(last7Days).Iter()
	defer iter.Close()

	TopKeywords = getTopKeywords(iter)
	log.Println("update top keywords elapsed time:", time.Since(startTime))
}

func getTopKeywords(iter *mgo.Iter) []string {
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

	topKeywords := []string{}
	for _, keywordScore := range keywordScores {
		if len(topKeywords) == DEFAULT_TOP_KEYWORD_LIMIT {
			break
		}
		topKeywords = append(topKeywords, keywordScore.Keyword)
	}

	return topKeywords
}
