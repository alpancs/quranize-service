package route

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/alpancs/quranize/job"
	"github.com/alpancs/quranize/service"
	"github.com/go-chi/chi"
)

type Location struct {
	Sura, Aya, Begin, End int
	SuraName, AyaText     string
}

type History struct {
	Timestamp time.Time
	Keyword   string
}

const DEFAULT_TRENDING_KEYWORDS_LIMIT = 6

var history *mgo.Collection

func init() {
	session, err := mgo.Dial(os.Getenv("MONGODB_HOST"))
	if err == nil {
		history = session.DB(os.Getenv("MONGODB_DATABASE")).C("history")
	} else {
		log.Println(err.Error())
	}
}

func Encode(w http.ResponseWriter, r *http.Request) {
	keyword, _ := url.QueryUnescape(chi.URLParam(r, "keyword"))
	json.NewEncoder(w).Encode(service.Encode(keyword))
}

func Locate(w http.ResponseWriter, r *http.Request) {
	keyword := chi.URLParam(r, "keyword")
	locations := []Location{}
	for _, loc := range service.Locate(keyword) {
		suraName := service.QuranEnhanced.Suras[loc.Sura].Name
		ayaText := service.QuranEnhanced.Suras[loc.Sura].Ayas[loc.Aya].Text
		ayaTextRune := []rune(ayaText)
		begin := indexAfterSpaces(ayaTextRune, loc.SliceIndex)
		end := begin + indexAfterSpaces(ayaTextRune[begin:], strings.Count(keyword, " ")+1) - 1
		locations = append(locations, Location{loc.Sura, loc.Aya, begin, end, suraName, ayaText})
	}
	json.NewEncoder(w).Encode(locations)
}

func TrendingKeywords(w http.ResponseWriter, r *http.Request) {
	limit := normalizeLimit(r.URL.Query().Get("limit"))
	json.NewEncoder(w).Encode(job.TrendingKeywords[:limit])
}

func Log(w http.ResponseWriter, r *http.Request) {
	keyword, _ := url.QueryUnescape(chi.URLParam(r, "keyword"))
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	if keyword == "" {
		return
	}

	err := history.Insert(History{bson.Now(), keyword})
	if err != nil {
		w.WriteHeader(500)
		log.Println(err.Error())
	}
}

func Translate(w http.ResponseWriter, r *http.Request) {
	sura, errSura := strconv.Atoi(chi.URLParam(r, "sura"))
	aya, errAya := strconv.Atoi(chi.URLParam(r, "aya"))
	if errSura == nil && errAya == nil && validIndex(sura, aya) {
		json.NewEncoder(w).Encode(service.QuranTranslationID.Suras[sura-1].Ayas[aya-1].Text)
	} else {
		w.WriteHeader(400)
	}
}

func indexAfterSpaces(text []rune, remain int) int {
	for i, r := range text {
		if remain == 0 {
			return i
		}
		if r == ' ' {
			remain--
		}
	}
	return len(text) + 1
}

func normalizeLimit(queryLimit string) int {
	limit, err := strconv.Atoi(queryLimit)
	if err != nil {
		limit = DEFAULT_TRENDING_KEYWORDS_LIMIT
	}
	if limit < 0 {
		limit = 0
	}
	if limit > len(job.TrendingKeywords) {
		limit = len(job.TrendingKeywords)
	}
	return limit
}

func validIndex(sura, aya int) bool {
	if sura < 1 || sura > len(service.QuranTranslationID.Suras) {
		return false
	}
	if aya < 1 || aya > len(service.QuranTranslationID.Suras[sura-1].Ayas) {
		return false
	}
	return true
}
