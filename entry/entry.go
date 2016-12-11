package entry

import (
	"bytes"
	"html/template"
	"io"
	"net/http"

	"github.com/kechako/go-advent-calendar/config"
	"github.com/kechako/go-advent-calendar/router"
	"github.com/kechako/go-advent-calendar/store"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
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

func init() {
	router.Router.HandleFunc("/{year:\\d*}/entries/{day:\\d*}", entryHandler).Methods("GET")
}

// エントリーを表示するハンドラー
func entryHandler(w http.ResponseWriter, r *http.Request) {
	// App Engine のコンテキスト取得
	ctx := appengine.NewContext(r)

	conf := config.GetConfig()

	year, err := router.GetYear(r, "year")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	day, err := router.GetDay(r, "day")
	if err != nil {
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
		log.Errorf(ctx, "Template error: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// バッファからレスポンスにコピー
	io.Copy(w, buff)
}