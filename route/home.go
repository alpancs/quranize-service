package route

import (
	"html/template"
	"net/http"
	"os"
)

var data = struct{ Production bool }{os.Getenv("GO_ENV") == "production"}

func Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("view/home.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, data)
}
