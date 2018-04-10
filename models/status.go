package models

// StatusData is a model for status information
type StatusData struct {
	ID          uint  `json:"id"`
	Timestamp   int64 `json:"timestamp"`
	Temperature int   `json:"temperature"`
	Humidity    uint  `json:"humidity"`
	Light       uint  `json:"light"`
}
