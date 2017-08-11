package route

import (
	"html/template"
	"net/http"
	"net/url"
	"os"

	"github.com/alpancs/quranize/service"
	"github.com/go-chi/chi"
)

type Data struct {
	Production      bool
	Keyword         string
	Transliteration string
	QuranText       string
	JsVersion       int64
}

var jsVersion int64

func init() {
	fileInfo, _ := os.Stat("public/home.js")
	jsVersion = fileInfo.ModTime().Unix()
}

func Home(w http.ResponseWriter, r *http.Request) {
	keyword, _ := url.QueryUnescape(chi.URLParam(r, "keyword"))

	transliteration := "transliteration"
	quranText := "Alquran"
	encodeds := service.Encode(keyword)
	if len(encodeds) > 0 {
		transliteration = keyword
		quranText = encodeds[0]
	}

	data := Data{os.Getenv("GO_ENV") == "production", keyword, transliteration, quranText, jsVersion}
	t, err := template.ParseFiles("view/home.html")
	if err != nil {
		if !data.Production {
			w.Write([]byte(err.Error()))
		}
		panic(err)
	}
	t.Execute(w, data)
}
