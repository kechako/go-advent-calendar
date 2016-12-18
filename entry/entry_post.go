package entry

import (
	"fmt"
	"net/http"

	"github.com/kechako/go-advent-calendar/router"
	"github.com/kechako/go-advent-calendar/store"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// エントリーを登録するハンドラー
func entryPostHandler(w http.ResponseWriter, r *http.Request) {
	// App Engine のコンテキスト取得
	ctx := appengine.NewContext(r)

	params, err := router.GetPathParams(r, "entries", ":year", ":day")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	year, ok := router.GetYear(params["year"])
	if !ok {
		http.NotFound(w, r)
		return
	}

	day, ok := router.GetDay(params["day"])
	if !ok {
		http.NotFound(w, r)
		return
	}

	s := store.NewStore(ctx)
	e, err := s.GetEntry(year, day)
	if err != nil {
		log.Errorf(ctx, "Get entry error : %s", err)
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

	err = s.PutEntry(e)
	if err != nil {
		log.Errorf(ctx, "Put entry error : %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%d", year), http.StatusSeeOther)
}
