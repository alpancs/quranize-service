package route

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"io/ioutil"
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
	JsVersion       string
}

var jsVersion string

func init() {
	jsVersion = getVersion("public/home.js")
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

func getVersion(filePath string) string {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	hash := md5.New()
	hash.Write(raw)
	return hex.EncodeToString(hash.Sum(nil))
}
