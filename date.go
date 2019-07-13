package quickbooks

import (
	"time"
)

const format = "2006-01-02"

// Date represents a Quickbooks date
type Date struct {
	time.Time
}

// UnmarshalJSON removes time from parsed date
func (d *Date) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	// Strip time off
	b = b[:10]

	d.Time, err = time.Parse(format, string(b))

	return err
}

func (d Date) String() string {
	return d.Format(format)
}
