package models

// StatusData is a model for temperature information
type StatusData struct {
	ID          uint    `json:"id"`
	Timestamp   int64   `json:"timestamp"`
	Temperature float64 `json:"temperature"`
}
