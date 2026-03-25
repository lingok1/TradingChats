package utils

import (
	"time"
)

var BeijingLocation *time.Location

func init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.FixedZone("CST", 8*60*60)
	}
	BeijingLocation = loc
}

// NowString 获取当前北京时间字符串 (YYYY-MM-DD HH:mm:ss)
func NowString() string {
	return time.Now().In(BeijingLocation).Format("2006-01-02 15:04:05")
}

// FormatTime 将 time.Time 格式化为北京时间字符串
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.In(BeijingLocation).Format("2006-01-02 15:04:05")
}
