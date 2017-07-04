package route

import (
	"html/template"
	"net/http"
	"os"

	"github.com/alpancs/quranize/service"
	"github.com/go-chi/chi"
)

type Data struct {
	Production                   bool
	Keyword, Alphabet, QuranText string
}

func Home(w http.ResponseWriter, r *http.Request) {
	go service.UpdateTopKeywords()
	keyword := chi.URLParam(r, "keyword")

	alphabet := "alphabet"
	quranText := "Al-Qur'an"
	encodeds := service.Encode(keyword)
	if len(encodeds) > 0 {
		alphabet = keyword
		quranText = encodeds[0]
	}

	data := Data{os.Getenv("GO_ENV") == "production", keyword, alphabet, quranText}
	t, err := template.ParseFiles("view/home.html")
	if err != nil {
		if !data.Production {
			w.Write([]byte(err.Error()))
		}
		panic(err)
	}
	t.Execute(w, data)
}
