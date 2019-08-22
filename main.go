package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"

	_ "github.com/kechako/go-advent-calendar/calendar"
	_ "github.com/kechako/go-advent-calendar/entry"
)

func main() {
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, runtime.Version())
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
