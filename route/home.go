package route

import (
	"html/template"
	"net/http"
	"os"

	"github.com/pressly/chi"
)

type Data struct {
	Production bool
	Input      string
}

func Home(w http.ResponseWriter, r *http.Request) {
	data := Data{os.Getenv("GO_ENV") == "production", chi.URLParam(r, "input")}
	t, err := template.ParseFiles("view/home.html")
	if err != nil {
		if !data.Production {
			w.Write([]byte(err.Error()))
		}
		panic(err)
	}
	t.Execute(w, data)
}
