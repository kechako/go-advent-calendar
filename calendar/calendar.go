package calendar

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/kechako/go-advent-calendar/config"
	"github.com/kechako/go-advent-calendar/log"
	"github.com/kechako/go-advent-calendar/store"
	"github.com/kechako/go-advent-calendar/util"
	"go.uber.org/zap"
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

type Handler struct {
	store *store.Store
}

func NewHandler(s *store.Store) *Handler {
	return &Handler{
		store: s,
	}
}

func (h *Handler) RegisterHandler(r chi.Router) {
	r.Get("/", h.calendarHandler)
	r.Get("/{year}", h.calendarHandler)
}

// カレンダーを表示するハンドラー
func (h *Handler) calendarHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	conf := config.GetConfig()

	y := chi.URLParam(r, "year")

	var year int
	if y == "" {
		year = util.CurrentYear()
	} else {
		var ok bool
		year, ok = util.GetYear(y)
		if !ok {
			http.NotFound(w, r)
			return
		}
	}

	entries, err := h.store.GetEntries(ctx, year)
	if err != nil {
		log.Logger.Error("failed to get entry", zap.Error(err))
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
		log.Logger.Error("failed to execute template", zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// バッファからレスポンスにコピー
	io.Copy(w, buff)
}

// year で指定された年の週データを作成します。
func makeWeeks(year int, entries []*store.Entry, loc *time.Location) []*calendarWeek {
	entMap := makeEntryMap(entries)

	weeks := make([]*calendarWeek, 0, 6)
	t := time.Date(year, 12, 1, 0, 0, 0, 0, loc)
	for wi := 0; ; wi++ {
		week := &calendarWeek{
			Days: make([]*calendarDay, 0, 7),
		}
		weeks = append(weeks, week)
		for di := 0; di < 7; di++ {
			day := &calendarDay{}
			week.Days = append(week.Days, day)

			if t.Month() == 12 && t.Weekday() == time.Weekday(di) {
				day.Day = t.Day()
				day.Entry = entMap[day.Day]
				day.Date = t
				t = t.AddDate(0, 0, 1)
			}
		}
		if t.Month() != 12 {
			break
		}
	}
	return weeks
}

// エントリーの日付をキーとするマップを作成します。
func makeEntryMap(entries []*store.Entry) map[int]*store.Entry {
	entMap := make(map[int]*store.Entry)
	for _, e := range entries {
		entMap[e.Day] = e
	}
	return entMap
}
