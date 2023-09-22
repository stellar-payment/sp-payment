package timeutil

import (
	"time"
)

func ParseLocaltime(t string) (local time.Time) {
	loc, _ := time.LoadLocation("Asia/Makassar")
	local, _ = time.ParseInLocation("2006-01-02 15:04:05", t, loc)
	return
}

func ParseLocaltimeOnly(t string) (local time.Time) {
	loc, _ := time.LoadLocation("Asia/Makassar")
	local, _ = time.ParseInLocation("15:04:05", t, loc)
	return
}

func ParseDate(t string) (local time.Time) {
	local, _ = time.Parse("2006-01-02", t)
	return
}

func ConvertLocalTime(t time.Time) (local time.Time) {
	loc, _ := time.LoadLocation("Asia/Makassar")
	return t.In(loc)
}

func FormatLocaltime(t time.Time) (str string) {
	loc, _ := time.LoadLocation("Asia/Makassar")
	return t.In(loc).Format("2006-01-02 15:04:05")
}

func FormatLocaltimeOnly(t time.Time) (str string) {
	loc, _ := time.LoadLocation("Asia/Makassar")
	return t.In(loc).Format("15:04:05")
}

func FormatDate(t time.Time) (str string) {
	if t.IsZero() {
		return ""
	}

	return t.Format("2006-01-02")
}

func FormatVerboseTime(t time.Time) (str string) {
	loc, _ := time.LoadLocation("Asia/Makassar")
	return t.In(loc).Format("2006-01-02T15:04:05-07:00")
}

func FormatVerboseLogTime(t time.Time) (str string) {
	loc, _ := time.LoadLocation("Asia/Makassar")
	return t.In(loc).Format("2006-01-02T150405")
}

func FormatFlattenTime(t time.Time) (str string) {
	loc, _ := time.LoadLocation("Asia/Makassar")
	return t.In(loc).Format("2006010215040507")
}

func GetStartEndMonth(t time.Time) (first time.Time, last time.Time) {
	first = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	last = first.AddDate(0, 1, -1)
	return
}

func Truncate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
