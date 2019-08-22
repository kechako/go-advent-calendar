package entry

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kechako/go-advent-calendar/store"
	"github.com/kechako/go-advent-calendar/util"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// エントリーを JSON で返すAPIハンドラー
func (h *Handler) entriesAPIHandler(w http.ResponseWriter, r *http.Request) {
	// App Engine のコンテキスト取得
	ctx := appengine.NewContext(r)

	year, ok := util.GetYear(chi.URLParam(r, "year"))
	if !ok {
		http.NotFound(w, r)
		return
	}

	entries, err := h.store.GetEntries(ctx, year)
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
func (h *Handler) postEntriesAPIHandler(w http.ResponseWriter, r *http.Request) {
	// App Engine のコンテキスト取得
	ctx := appengine.NewContext(r)

	year, ok := util.GetYear(chi.URLParam(r, "year"))
	if !ok {
		http.NotFound(w, r)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Bad content type", http.StatusBadRequest)
		return
	}

	entries := make([]*store.Entry, 25)
	err := json.NewDecoder(r.Body).Decode(&entries)
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

	for _, e := range entries {
		e.Year = year

		err = h.store.PutEntry(ctx, e)
		if err != nil {
			log.Errorf(ctx, "Put entry error : %s", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
