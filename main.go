package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
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

// カレンダーを表示するハンドラー
func calenerHundler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{year:\\d*}", calenerHundler).Methods("GET")

	http.ListenAndServe(":8080", r)
}
