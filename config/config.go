package config

import (
	"os"
	"strconv"
)

// Config はカレンダーの設定を格納します。
type Config struct {
	CalendarName    string
	CalendarYearMin int
	CalendarYearMax int
}

var config *Config

func init() {
	config = &Config{
		CalendarName:    os.Getenv("CALENDAR_NAME"),
		CalendarYearMin: atoi(os.Getenv("CALENDAR_YEAR_MIN"), 2014),
		CalendarYearMax: atoi(os.Getenv("CALENDAR_YEAR_MAX"), 2099),
	}
}

// カレンダーの設定を取得します。
func GetConfig() *Config {
	return config
}

// 文字列 s を数値に変換します。
// value は変換できない場合の初期値です。
func atoi(s string, value int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return value
	}
	return i
}
