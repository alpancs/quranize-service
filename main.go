package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/alpancs/quranize/job"
	"github.com/alpancs/quranize/route"
	"github.com/alpancs/quranize/route/api"
	"github.com/alpancs/quranize/route/webhook"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	job.Start()
	port := getPort()
	fmt.Println("Quranize is listening port", port)
	http.ListenAndServe(":"+port, newRouter())
}

func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "7000"
}

func newRouter() http.Handler {
	router := chi.NewRouter()

	isProduction := os.Getenv("ENV") == "production"
	if !isProduction {
		router.Use(middleware.Logger)
	}

	router.Route("/", func(compressedRoute chi.Router) {
		compressedRoute.Use(middleware.DefaultCompress)
		compressedRoute.Use(header("Vary", "Accept-Encoding"))

		homeRouter := compressedRoute.With(header("Content-Type", "text/html; charset=utf-8"))
		homeRouter.Get("/", route.Home)
		homeRouter.Get("/{keyword:^([A-Za-z' ]|%20)+$}", route.Home)

		cachedRouter := compressedRoute.With(header("Cache-Control", "public, max-age=31536000"))
		fileServer(cachedRouter, "/", http.Dir("public"))
	})

	router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Use(header("Content-Type", "application/json; charset=utf-8"))
		cachedRouter := apiRouter.With(header("Cache-Control", "public, max-age=3600"))
		if !isProduction {
			cachedRouter = apiRouter
		}
		cachedRouter.Get("/encode", api.Encode)
		cachedRouter.Get("/locate", api.Locate)
		cachedRouter.Get("/aya/{sura}/{aya}", api.Aya)
		cachedRouter.Get("/translation/{sura}/{aya}", api.Translation)
		cachedRouter.Get("/tafsir/{sura}/{aya}", api.Tafsir)
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
	router.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func header(key, value string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(key, value)
			next.ServeHTTP(w, r)
		})
	}
}
