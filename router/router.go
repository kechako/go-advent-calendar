package router

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kechako/go-advent-calendar/config"
)

var Router = mux.NewRouter()

// Errors
var (
	// URL から年が見つからないエラー。
	ErrYearNotFound = errors.New("Year not found.")
	// URL から日が見つからないエラー。
	ErrDayNotFound = errors.New("Day not found.")
)

func init() {
	http.Handle("/", Router)
}

// リクエストのパスで指定された年を取得します。
// name はパスパターンで指定したパラメーター名です。
func GetYear(r *http.Request, name string) (int, error) {
	conf := config.GetConfig()
	var year int

	vars := mux.Vars(r)
	yearStr := vars[name]
	if yearStr == "" {
		now := time.Now().In(conf.Location)

		if now.Month() >= 11 {
			// 11月以降は今年
			year = now.Year()
		} else {
			// 10月以前は去年
			year = now.Year() - 1
		}
	} else {
		var err error
		// URL から年を取得
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			// 見つからない
			return 0, ErrYearNotFound
		}

		if year < conf.CalendarYearMin || year > conf.CalendarYearMax {
			// 範囲外
			return 0, ErrYearNotFound
		}
	}

	return year, nil
}

// リクエストのパスで指定された日を取得します。
// name はパスパターンで指定したパラメーター名です。
func GetDay(r *http.Request, name string) (int, error) {
	vars := mux.Vars(r)
	dayStr := vars[name]
	if dayStr == "" {
		return 0, ErrDayNotFound
	}

	// URL から日を取得
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		// 見つからない
		return 0, ErrDayNotFound
	}

	if day < 1 || day > 25 {
		// 範囲外
		return 0, ErrDayNotFound
	}

	return day, nil
}
