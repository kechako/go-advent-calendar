package entry

import (
	"bytes"
	"html/template"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kechako/go-advent-calendar/config"
	"github.com/kechako/go-advent-calendar/log"
	"github.com/kechako/go-advent-calendar/store"
	"github.com/kechako/go-advent-calendar/util"
	"go.uber.org/zap"
)

// テンプレート
var entryTmpl = template.Must(template.ParseFiles("templates/entry.tmpl"))

// エントリーに表示するデータを格納する構造体
type entryData struct {
	Name  string
	Year  int
	Day   int
	Entry *store.Entry
}

type Handler struct {
	store *store.Store
}

func NewHandler(s *store.Store) *Handler {
	return &Handler{
		store: s,
	}
}

func (h *Handler) RegisterHandler(r chi.Router) {
	r.Route("/entries", func(r chi.Router) {
		r.Get("/{year}/{day}", h.entryHandler)
		r.Post("/{year}/{day}", h.entryPostHandler)
	})
	r.Route("/api/entries", func(r chi.Router) {
		r.Get("/{year}", h.entriesAPIHandler)
		r.Post("/{year}", h.postEntriesAPIHandler)
	})
}

// エントリーを表示するハンドラー
func (h *Handler) entryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conf := config.GetConfig()

	year, ok := util.GetYear(chi.URLParam(r, "year"))
	if !ok {
		http.NotFound(w, r)
		return
	}

	day, ok := util.GetDay(chi.URLParam(r, "day"))
	if !ok {
		http.NotFound(w, r)
		return
	}

	e, err := h.store.GetEntry(ctx, year, day)
	if err != nil {
		log.Logger.Error("failed to get entry", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if e == nil {
		// データなし
		e = new(store.Entry)
	}

	// カレンダー用データを作成
	data := &entryData{
		Name:  conf.CalendarName,
		Year:  year,
		Day:   day,
		Entry: e,
	}

	// 直接レスポンスに書き込むとエラー時にステータスを変更できないので
	// バッファーに書き込む
	buff := new(bytes.Buffer)
	err = entryTmpl.Execute(buff, data)
	if err != nil {
		log.Logger.Error("failed to execute template", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// バッファからレスポンスにコピー
	io.Copy(w, buff)
}
