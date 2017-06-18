package route

import (
	"html/template"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("view/home.html")
	t.Execute(w, nil)
}
