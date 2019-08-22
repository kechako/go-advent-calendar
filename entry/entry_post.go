package entry

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/kechako/go-advent-calendar/log"
	"github.com/kechako/go-advent-calendar/store"
	"github.com/kechako/go-advent-calendar/util"
	"go.uber.org/zap"
)

// エントリーを登録するハンドラー
func (h *Handler) entryPostHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	e.Year = year
	e.Day = day
	e.Title = r.FormValue("title")
	e.Url = r.FormValue("url")
	e.Author = r.FormValue("author")
	e.Section = r.FormValue("section")

	err = h.store.PutEntry(ctx, e)
	if err != nil {
		log.Logger.Error("failed to put entry", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%d", year), http.StatusSeeOther)
}
