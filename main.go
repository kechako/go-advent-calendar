package main

import (
	"net/http"

	_ "github.com/kechako/go-advent-calendar/calendar"
	"github.com/kechako/go-advent-calendar/router"
)

func main() {
	http.ListenAndServe(":8080", router.Router)
}
