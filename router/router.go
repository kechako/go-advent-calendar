package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

var Router = mux.NewRouter()

func init() {
	http.Handle("/", Router)
}
