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

	router.Get("/{keyword}", route.Home)

	router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Use(jsonify)
		apiRouter.Get("/encode/{keyword}", route.Encode)
		apiRouter.Get("/locate/{keyword}", route.Locate)
		apiRouter.Get("/top-keywords", route.TopKeywords)
	})

	router.Get("/log/{keyword}", route.Log)

	fileServer(router, "/file", http.Dir("public"))

	return router
}

func fileServer(router chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		router.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	router.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func jsonify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}
