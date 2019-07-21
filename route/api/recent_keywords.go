package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/alpancs/quranize-service/db"
)

const DEFAULT_RECENT_KEYWORDS_LIMIT = 6

func RecentKeywords(w http.ResponseWriter, r *http.Request) {
	limit := normalizeLimit(r.URL.Query().Get("limit"), DEFAULT_RECENT_KEYWORDS_LIMIT)
	json.NewEncoder(w).Encode(recentKeywords(limit))
}

func recentKeywords(limit int) []string {
	keywords := []string{}
	rows, err := db.Query(`
		SELECT regexp_replace(keyword, '[^a-z'']', '', 'g') FROM history
		WHERE timestamp >= $1
		GROUP BY 1
		ORDER BY max(timestamp) DESC
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
