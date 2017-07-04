package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alpancs/quranize/route"
	"github.com/go-chi/chi"
)

func init() {
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "7000")
	}
}

func main() {
	log.Println("Quranize is running in port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), setUpRouter())
}

func setUpRouter() http.Handler {
	router := chi.NewRouter()

	FileServer(router, "/file", http.Dir("public"))

	router.Get("/{input}", route.Home)

	router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Use(jsonify)
		apiRouter.Get("/encode/{input}", route.Encode)
		apiRouter.Get("/locate/{input}", route.Locate)
	})

	return router
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func jsonify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}
