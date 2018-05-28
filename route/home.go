package route

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/alpancs/quranize-service/quran"
	"github.com/go-chi/chi"
)

type HomeData struct {
	IsProduction    bool
	Keyword         string
	Transliteration string
	QuranText       string
	CssVersion      string
	JsVersion       string
}

var (
	isProduction = os.Getenv("ENV") == "production"
	homeTemplate = parseTemplate("view/home.html")
	cssVersion   = getVersion("/home.css")
	jsVersion    = getVersion("/home.js")
)

func Home(w http.ResponseWriter, r *http.Request) {
	keyword, _ := url.QueryUnescape(chi.URLParam(r, "keyword"))

	transliteration := "alquran"
	quranText := "القرآن"
	encodeds := quran.Encode(keyword)
	if len(encodeds) > 0 {
		transliteration = keyword
		quranText = encodeds[0]
	}

	homeData := HomeData{isProduction, keyword, transliteration, quranText, cssVersion, jsVersion}
	homeTemplate.Execute(w, homeData)
}

func parseTemplate(filePath string) *template.Template {
	parsedTemplate, err := template.ParseFiles(filePath)
	if err != nil {
		panic(err)
	}
	return parsedTemplate
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
