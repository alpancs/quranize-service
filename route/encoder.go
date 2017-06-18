package route

import (
	"net/http"
	"strings"

	"github.com/alpancs/quranize/service"
	"github.com/pressly/chi"
)

func Encode(w http.ResponseWriter, r *http.Request) {
	text := chi.URLParam(r, "text")
	encodeds := service.Encode(text)
	w.Write([]byte(strings.Join(encodeds, ", ")))
}
