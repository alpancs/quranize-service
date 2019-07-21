package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alpancs/quranize-service/db"
)

const DEFAULT_TRENDING_KEYWORDS_LIMIT = 6

func TrendingKeywords(w http.ResponseWriter, r *http.Request) {
	limit := normalizeLimit(r.URL.Query().Get("limit"), DEFAULT_TRENDING_KEYWORDS_LIMIT)
	json.NewEncoder(w).Encode(trendingKeywords(limit))
}

func normalizeLimit(queryLimit string, defaultLimit int) int {
	limit, err := strconv.Atoi(queryLimit)
	if err != nil {
		limit = defaultLimit
	}
	if limit < 0 {
		limit = 0
	}
	return limit
}

func trendingKeywords(limit int) []string {
	keywords := []string{}
	rows, err := db.Query(`
		SELECT regexp_replace(keyword, '[^a-z'']', '', 'g') FROM history
		WHERE timestamp >= $1
		GROUP BY 1
		ORDER BY count(1) DESC
		LIMIT $2`,
		time.Now().Add(-30*24*time.Hour).In(time.UTC),
		limit,
	)
	if err != nil {
		log.Println(err)
		return keywords
	}
	defer rows.Close()
	for rows.Next() {
		var keyword string
		if err := rows.Scan(&keyword); err != nil {
			log.Println(err)
			return keywords
		}
		keywords = append(keywords, keyword)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
	}
	return keywords
}
