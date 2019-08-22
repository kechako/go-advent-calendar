package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/go-chi/chi"
	"github.com/kechako/go-advent-calendar/calendar"
	"github.com/kechako/go-advent-calendar/entry"
	"github.com/kechako/go-advent-calendar/util"
)

func main() {

	r := chi.NewRouter()
	r.Use(util.IPFilterHandler)
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, runtime.Version())
	})

	c := calendar.NewHandler()
	c.RegisterHandler(r)
	e := entry.NewHandler()
	e.RegisterHandler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
