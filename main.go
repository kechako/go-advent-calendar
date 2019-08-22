package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/go-chi/chi"
	"github.com/kechako/go-advent-calendar/calendar"
	"github.com/kechako/go-advent-calendar/entry"
	"github.com/kechako/go-advent-calendar/store"
	"github.com/kechako/go-advent-calendar/util"
)

func main() {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Fatal("GOOGLE_CLOUD_PROJECT environment variable not set")
	}

	ctx := context.Background()

	r := chi.NewRouter()
	r.Use(util.IPFilterHandler)
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, runtime.Version())
	})

	s, err := store.NewStore(ctx, projectID)
	if err != nil {
		log.Fatal("Failed to initialize Cloud Datastore client: ", err)
	}
	defer s.Close()

	c := calendar.NewHandler(s)
	c.RegisterHandler(r)

	e := entry.NewHandler(s)
	e.RegisterHandler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
