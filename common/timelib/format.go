package timelib

import "time"

var FormatEndo = "2006-01-02 15:04:05"

func String(t *time.Time) string {
	if t == nil {
		return ""
	}

	return t.Format(FormatEndo)
}
