package main

import (
	"io"
	"net/http"
	"runtime"

	_ "github.com/kechako/go-advent-calendar/calendar"
	_ "github.com/kechako/go-advent-calendar/entry"
)

func init() {
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, runtime.Version())
	})
}
