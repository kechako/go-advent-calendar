package entry

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/kechako/go-advent-calendar/router"
	"github.com/kechako/go-advent-calendar/store"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	http.Handle("/api/entries/",
		router.IPFilterHandler(
			router.MethodsHandler(map[string]http.Handler{
				"GET":  http.HandlerFunc(entriesAPIHandler),
				"POST": http.HandlerFunc(postEntriesAPIHandler),
			})))
}

// エントリーを JSON で返すAPIハンドラー
func entriesAPIHandler(w http.ResponseWriter, r *http.Request) {
	// App Engine のコンテキスト取得
	ctx := appengine.NewContext(r)

	params, err := router.GetPathParams(r, "api", "entries", ":year")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	year, ok := router.GetYear(params["year"])
	if !ok {
		http.NotFound(w, r)
		return
	}

	s := store.NewStore(ctx)
	entries, err := s.GetEntries(year)
	if err != nil {
		log.Errorf(ctx, "Get entries error : %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 直接レスポンスに書き込むとエラー時にステータスを変更できないので
	// バッファーに書き込む
	buff := new(bytes.Buffer)
	err = json.NewEncoder(buff).Encode(entries)
	if err != nil {
		log.Errorf(ctx, "Template error: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// バッファからレスポンスにコピー
	io.Copy(w, buff)
}

// エントリーをJSON 受け取るAPIハンドラー
func postEntriesAPIHandler(w http.ResponseWriter, r *http.Request) {
	// App Engine のコンテキスト取得
	ctx := appengine.NewContext(r)

	params, err := router.GetPathParams(r, "api", "entries", ":year")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	year, ok := router.GetYear(params["year"])
	if !ok {
		http.NotFound(w, r)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Bad content type", http.StatusBadRequest)
		return
	}

	entries := make([]*store.Entry, 25)
	err = json.NewDecoder(r.Body).Decode(&entries)
	if err != nil {
		log.Errorf(ctx, "JSON parse error : %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	for _, e := range entries {
		if e.Day < 1 || e.Day > 25 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	}

	s := store.NewStore(ctx)
	for _, e := range entries {
		e.Year = year

		err = s.PutEntry(e)
		if err != nil {
			log.Errorf(ctx, "Put entry error : %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
