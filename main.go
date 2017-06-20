package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alpancs/quranize/route"
	"github.com/alpancs/quranize/service"
	"github.com/pressly/chi"
)

func init() {
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "7000")
	}
}

func dummy() {
	runeSet := make(map[rune]bool)
	for _, sura := range service.Quran.Suras {
		for _, aya := range sura.Ayas {
			for _, r := range []rune(aya.Text) {
				runeSet[r] = true
			}
		}
	}
	fmt.Println(len(runeSet))
	for r := range runeSet {
		fmt.Printf("%c %d\n", r, r)
	}
}

func main() {
	dummy()
	return
	log.Println("Quranize is running in port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), setUpRouter())
}

func setUpRouter() http.Handler {
	router := chi.NewRouter()

	router.Get("/", route.Home)

	router.Route("/api", func(apiRouter chi.Router) {
		apiRouter.Use(jsonify)
		apiRouter.Get("/encode/:input", route.Encode)
		apiRouter.Get("/locate/:input", route.Locate)
	})

	router.FileServer("/", http.Dir("public"))

	return router
}

func jsonify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}
