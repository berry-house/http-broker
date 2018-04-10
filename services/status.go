package services

import (
	"github.com/berry-house/http_broker/drivers/database"
	"github.com/berry-house/http_broker/models"
)

// StatusInvalidDataError is an error type for invalid data errors
type StatusInvalidDataError string

// StatusDatabaseDriverError is an error type for database driver errors
type StatusDatabaseDriverError string

func (e StatusInvalidDataError) Error() string    { return string(e) }
func (e StatusDatabaseDriverError) Error() string { return string(e) }

const (
	// StatusInvalidData is the default error for invalid data
	StatusInvalidData = StatusInvalidDataError("invalid data")
	// StatusInvalidID is the default error for non-existent IDs
	StatusInvalidID = StatusInvalidDataError("invalid ID")
)

// StatusDatabase is a service for writing status data to database
type StatusDatabase struct {
	Driver database.Database
}

// Write writes status data to the database
func (s *StatusDatabase) Write(data *models.StatusData) error {
	if data == nil {
		return StatusInvalidDataError("nil data")
	}

	// Threshold values
	if data.Temperature < -30 || data.Temperature > 50 ||
		data.Light > 150 ||
		data.Humidity > 100 {
		return StatusInvalidData
	}

	switch err := s.Driver.WriteStatus(data); err.(type) {
	case nil:
		return nil
	case database.DatabaseInvalidDataError:
		return StatusInvalidID
	default:
		return StatusDatabaseDriverError(err.Error())
	}
}
