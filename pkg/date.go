package quickbooks

import (
	"time"
)

const secondFormat = "2006-01-02"
const format = "2006-01-02T15:04:05-07:00"

// Date represents a Quickbooks date
type Date struct {
	time.Time
}

// UnmarshalJSON removes time from parsed date
func (d *Date) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	d.Time, err = time.Parse(format, string(b))
	if err != nil {
		d.Time, err = time.Parse(secondFormat, string(b))
	}

	return err
}

func (d Date) String() string {
	return d.Format(format)
}
