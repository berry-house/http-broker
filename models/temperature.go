package models

// TemperatureData is a model for temperature information
type TemperatureData struct {
	ID          uint    `json:"id"`
	Timestamp   int64   `json:"timestamp"`
	Temperature float64 `json:"temperature"`
}
