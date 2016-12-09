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
