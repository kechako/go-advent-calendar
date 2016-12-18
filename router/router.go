package router

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kechako/go-advent-calendar/config"
)

// URL からパラメーターが見つからないエラー
var ErrNotFound = errors.New("parameters not found.")

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

// HTTP メソッドをチェックするハンドラー。
func MethodsHandler(handlers map[string]http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handler, ok := handlers[r.Method]; ok {
			handler.ServeHTTP(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
}

func GetPathParams(r *http.Request, params ...string) (map[string]string, error) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	segs := strings.Split(path, "/")
	if path != "" {
		segs = segs[1:]
	}
	if len(segs) != len(params) {
		return nil, ErrNotFound
	}

	paramMap := make(map[string]string)

	for i, p := range params {
		seg := segs[i]
		if strings.HasPrefix(p, ":") {
			paramMap[p[1:]] = seg
		} else {
			if p != seg {
				return nil, ErrNotFound
			}
		}
	}

	return paramMap, nil
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
