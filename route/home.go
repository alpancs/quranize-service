package route

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/alpancs/quranize/core"
	"github.com/go-chi/chi"
)

type Data struct {
	IsProduction    bool
	Keyword         string
	Transliteration string
	QuranText       string
	CssVersion      string
	JsVersion       string
}

var (
	isProduction = os.Getenv("ENV") == "production"
	cssVersion   = getVersion("/home.css")
	jsVersion    = getVersion("/home.js")
)

func Home(w http.ResponseWriter, r *http.Request) {
	keyword, _ := url.QueryUnescape(chi.URLParam(r, "keyword"))

	transliteration := "transliteration"
	quranText := "Alquran"
	encodeds := core.Encode(keyword)
	if len(encodeds) > 0 {
		transliteration = keyword
		quranText = encodeds[0]
	}

	data := Data{isProduction, keyword, transliteration, quranText, cssVersion, jsVersion}
	t, err := template.ParseFiles("view/home.html")
	if err != nil {
		if !data.IsProduction {
			w.Write([]byte(err.Error()))
		}
		panic(err)
	}
	t.Execute(w, data)
}

func getVersion(filePath string) string {
	raw, err := ioutil.ReadFile("public" + filePath)
	if err != nil {
		panic(err)
	}
	hash := md5.New()
	hash.Write(raw)
	return hex.EncodeToString(hash.Sum(nil))
}
