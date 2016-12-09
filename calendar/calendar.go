package calendar

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kechako/go-advent-calendar/router"
)

// 許容する年の範囲
const (
	yearMin = 2014
	yearMax = 2099
)

// テンプレート
var calendarTmpl = template.Must(template.ParseFiles("templates/calendar.tmpl"))

// カレンダーに表示するデータを格納する構造体
type calendarData struct {
	Year int
}

func init() {
	router.Router.HandleFunc("/{year:\\d*}", calendarHandler).Methods("GET")
}

// カレンダーを表示するハンドラー
func calendarHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	now := time.Now()
	var year int

	vars := mux.Vars(r)
	yearStr := vars["year"]
	if yearStr == "" {
		if now.Month() >= 11 {
			// 11月以降は今年
			year = now.Year()
		} else {
			// 10月以前は去年
			year = now.Year() - 1
		}
	} else {
		// URL から年を取得
		year, err = strconv.Atoi(yearStr)
		if err != nil || year < yearMin || year > yearMax {
			// 見つからない
			w.WriteHeader(404)
			fmt.Fprintln(w, "Not Found")
			return
		}
	}
	err = calendarTmpl.Execute(w, &calendarData{
		Year: year,
	})
	if err != nil {
		w.WriteHeader(500)
	}
}
