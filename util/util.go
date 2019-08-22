package util

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/kechako/go-advent-calendar/config"
)

// IP フィルタリングを行うハンドラーを返します。
func IPFilterHandler(next http.Handler) http.Handler {
	ipList := config.GetConfig().IPWhiteList

	if len(ipList) == 0 {
		// 設定なし
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler := next.ServeHTTP

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}

		found := false
		for _, ip := range ipList {
			if host == ip {
				found = true
				break
			}
		}
		if !found {
			handler = http.NotFound
		}

		handler(w, r)
	})
}

// 年の文字列から年の値を取得します。
func GetYear(s string) (int, bool) {
	conf := config.GetConfig()

	var year int
	if s == "" {
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
		year, err = strconv.Atoi(s)
		if err != nil {
			// 見つからない
			return 0, false
		}

		if year < conf.CalendarYearMin || year > conf.CalendarYearMax {
			// 範囲外
			return 0, false
		}
	}

	return year, true
}

func GetDay(s string) (int, bool) {
	if s == "" {
		return 0, false
	}

	// URL から日を取得
	day, err := strconv.Atoi(s)
	if err != nil {
		// 見つからない
		return 0, false
	}

	if day < 1 || day > 25 {
		// 範囲外
		return 0, false
	}

	return day, true
}
