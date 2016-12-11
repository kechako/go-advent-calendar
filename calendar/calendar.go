package calendar

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/kechako/go-advent-calendar/config"
	"github.com/kechako/go-advent-calendar/router"
	"github.com/kechako/go-advent-calendar/store"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

// テンプレート
var calendarTmpl = template.Must(template.ParseFiles("templates/calendar.tmpl"))

// 日を格納する構造体
type calendarDay struct {
	Day   int
	Entry *store.Entry
	Date  time.Time
}

// 週を格納する構造体
type calendarWeek struct {
	Days []*calendarDay
}

// カレンダーに表示するデータを格納する構造体
type calendarData struct {
	Name  string
	Year  int
	Weeks []*calendarWeek
	Now   time.Time
}

func init() {
	router.Router.HandleFunc("/{year:\\d*}", calendarHandler).Methods("GET")
}

// カレンダーを表示するハンドラー
func calendarHandler(w http.ResponseWriter, r *http.Request) {
	// App Engine のコンテキスト取得
	ctx := appengine.NewContext(r)

	conf := config.GetConfig()

	year, err := router.GetYear(r, "year")
	if err != nil {
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

	// カレンダー用データを作成
	data := &calendarData{
		Name:  conf.CalendarName,
		Year:  year,
		Weeks: makeWeeks(year, entries, conf.Location),
		Now:   time.Now().In(conf.Location),
	}

	// 直接レスポンスに書き込むとエラー時にステータスを変更できないので
	// バッファーに書き込む
	buff := new(bytes.Buffer)
	err = calendarTmpl.Execute(buff, data)
	if err != nil {
		log.Errorf(ctx, "Template error: %s", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// バッファからレスポンスにコピー
	io.Copy(w, buff)
}

// year で指定された年の週データを作成します。
func makeWeeks(year int, entries []*store.Entry, loc *time.Location) []*calendarWeek {
	entMap := makeEntryMap(entries)

	weeks := make([]*calendarWeek, 0, 5)
	t := time.Date(year, 12, 1, 0, 0, 0, 0, loc)
	for wi := 0; ; wi++ {
		week := &calendarWeek{
			Days: make([]*calendarDay, 0, 7),
		}
		weeks = append(weeks, week)
		for di := 0; di < 7; di++ {
			day := &calendarDay{}
			week.Days = append(week.Days, day)

			if t.Weekday() == time.Weekday(di) {
				day.Day = t.Day()
				day.Entry = entMap[day.Day]
				day.Date = t
				t = t.AddDate(0, 0, 1)
				if t.Month() != 12 {
					return weeks
				}
			}
		}
	}
}

// エントリーの日付をキーとするマップを作成します。
func makeEntryMap(entries []*store.Entry) map[int]*store.Entry {
	entMap := make(map[int]*store.Entry)
	for _, e := range entries {
		entMap[e.Day] = e
	}
	return entMap
}
