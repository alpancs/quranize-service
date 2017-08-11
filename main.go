package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alpancs/quranize/job"
	"github.com/alpancs/quranize/route"
	"github.com/alpancs/quranize/route/api"
	"github.com/alpancs/quranize/route/webhook"
	"github.com/go-chi/chi"
)

func main() {
	job.RunInBackground()
	log.Println("Quranize is running in port " + getPort())
	http.ListenAndServe(":"+getPort(), newRouter())
}

func getPort() string {
	if os.Getenv("PORT") == "" {
		return "7000"
	}
	return os.Getenv("PORT")
}

func newRouter() http.Handler {
	router := chi.NewRouter()

	router.Route("/", func(homeRouter chi.Router) {
		homeRouter.Use(vary)
		homeRouter.Get("/", route.Home)
		homeRouter.Get("/{keyword:^([A-Za-z' ]|%20)+$}", route.Home)
		fileServer(homeRouter, "/", http.Dir("public"))
	})

	router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Use(jsonify)
		apiRouter.Get("/encode", api.Encode)
		apiRouter.Get("/locate", api.Locate)
		apiRouter.Get("/translate/{sura}/{aya}", api.Translate)
		apiRouter.Get("/tafsir/{sura}/{aya}", api.Tafsir)
		apiRouter.Get("/trending_keywords", api.TrendingKeywords)
		apiRouter.Post("/log", api.Log)
	})

	router.Route("/webhook", func(webhookRouter chi.Router) {
		webhookRouter.Post("/"+os.Getenv("QURANIZE_TELEGRAM_TOKEN"), webhook.Telegram)
	})

	return router
}

func fileServer(router chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		router.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	router.With(oneYearCache).Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func jsonify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func vary(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		next.ServeHTTP(w, r)
	})
}

func oneYearCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=31536000")
		next.ServeHTTP(w, r)
	})
}
