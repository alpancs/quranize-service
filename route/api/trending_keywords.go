package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alpancs/quranize-service/job"
)

const DEFAULT_TRENDING_KEYWORDS_LIMIT = 6

func TrendingKeywords(w http.ResponseWriter, r *http.Request) {
	limit := normalizeLimit(r.URL.Query().Get("limit"))
	json.NewEncoder(w).Encode(job.TrendingKeywords[:limit])
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
