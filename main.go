package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alpancs/quranize-service/route"
	"github.com/alpancs/quranize-service/route/api"
	"github.com/alpancs/quranize-service/route/webhook"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
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

	router.Route("/", func(compressedRouter chi.Router) {
		compressedRouter.Use(middleware.DefaultCompress)
		compressedRouter.Use(header("Vary", "Accept-Encoding"))

		homeRouter := compressedRouter.With(header("Content-Type", "text/html; charset=utf-8"))
		homeRouter.Get("/", route.Home)
		homeRouter.Get("/{keyword:^([A-Za-z' ]|%20)+$}", route.Home)

		cacheControl := "no-store"
		if isProduction {
			cacheControl = "public, max-age=31536000"
		}
		cachedRouter := compressedRouter.With(header("Cache-Control", cacheControl))
		fileServer(cachedRouter, "/", http.Dir("public"))
	})

	router.Route("/api", func(apiRouter chi.Router) {
		if !isProduction {
			apiRouter.Use(delay(1200 * time.Millisecond))
		}
		apiRouter.Use(header("Content-Type", "application/json; charset=utf-8"))
		apiRouter.Use(header("Cache-Control", "no-store"))
		cachedRouter := apiRouter
		if isProduction {
			cachedRouter = apiRouter.With(header("Cache-Control", "public, max-age=43200"))
		}
		cachedRouter.Get("/encode", api.Encode)
		cachedRouter.Get("/locate", api.Locate)
		cachedRouter.Get("/aya/{sura}/{aya}", api.Aya)
		cachedRouter.Get("/translation/{sura}/{aya}", api.Translation)
		cachedRouter.Get("/tafsir/{sura}/{aya}", api.Tafsir)
		apiRouter.Get("/trending_keywords", api.TrendingKeywords)
		apiRouter.Get("/recent_keywords", api.RecentKeywords)
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

func delay(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(duration)
			next.ServeHTTP(w, r)
		})
	}
}
