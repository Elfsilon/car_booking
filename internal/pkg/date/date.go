package date

import (
	"strings"
	"time"
)

type Date struct {
	time.Time
}

func (d *Date) MarshalJSON() ([]byte, error) {
	str := "\"" + d.Time.UTC().Format(time.DateOnly) + "\""
	return []byte(str), nil
}

func (d *Date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return
	}
	d.Time, err = time.Parse(time.DateOnly, s)
	return
}
