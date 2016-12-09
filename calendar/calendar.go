package calendar

import (
	"bytes"
	"html/template"
	"io"
	"net/http"

	"github.com/kechako/go-advent-calendar/config"
	"github.com/kechako/go-advent-calendar/router"
)

// テンプレート
var calendarTmpl = template.Must(template.ParseFiles("templates/calendar.tmpl"))

// カレンダーに表示するデータを格納する構造体
type calendarData struct {
	Name string
	Year int
}

func init() {
	router.Router.HandleFunc("/{year:\\d*}", calendarHandler).Methods("GET")
}

// カレンダーを表示するハンドラー
func calendarHandler(w http.ResponseWriter, r *http.Request) {
	conf := config.GetConfig()

	year, err := router.GetYear(r, "year")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 直接レスポンスに書き込むとエラー時にステータスを変更できないので
	// バッファーに書き込む
	buff := new(bytes.Buffer)
	err = calendarTmpl.Execute(buff, &calendarData{
		Name: conf.CalendarName,
		Year: year,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// バッファからレスポンスにコピー
	io.Copy(w, buff)
}
