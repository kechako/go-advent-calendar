package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
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
	now := time.Now()
	var year int

	if r.URL.Path == "/" {
		if now.Month() >= 11 {
			// 11月以降は今年
			year = now.Year()
		} else {
			// 10月以前は去年
			year = now.Year() - 1
		}
	} else {
		// URL から年を取得
		_, err := fmt.Sscanf(r.URL.Path, "/%d", &year)
		if err != nil || year < yearMin || year > yearMax {
			// 見つからない
			w.WriteHeader(404)
			fmt.Fprintln(w, "Not Found")
			return
		}
	}
	err := calendarTmpl.Execute(w, &calendarData{
		Year: year,
	})
	if err != nil {
		w.WriteHeader(500)
	}
}

func main() {
	http.HandleFunc("/", calenerHundler)

	http.ListenAndServe(":8080", nil)
}
