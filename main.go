package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alpancs/quranize/route"
	"github.com/pressly/chi"
)

func main() {
	if os.Getenv("PORT") == "" {
		os.Setenv("PORT", "8000")
	}
	log.Println("Linguist is running in port " + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), setUpRouter())
}

func setUpRouter() http.Handler {
	router := chi.NewRouter()

	router.Get("/api/encode/:text", route.Encode)

	return router
}
