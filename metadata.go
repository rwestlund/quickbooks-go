package quickbooks

import "time"

// MetaData is a timestamp of genesis and last change of a Quickbooks object
type MetaData struct {
	CreateTime      time.Time
	LastUpdatedTime time.Time
}
