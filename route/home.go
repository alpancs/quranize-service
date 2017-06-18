package route

import (
	"html/template"
	"net/http"
	"os"
)

var (
	t, _ = template.ParseFiles("view/home.html")
	data = struct{ Production bool }{os.Getenv("GO_ENV") == "production"}
)

func Home(w http.ResponseWriter, r *http.Request) {
	t.Execute(w, data)
}
