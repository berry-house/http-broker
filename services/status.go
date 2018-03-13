package services

import (
	"github.com/berry-house/http_broker/drivers"
	"github.com/berry-house/http_broker/models"
)

// StatusInvalidDataError is an error type for invalid data errors
type StatusInvalidDataError string

// StatusDatabaseDriverError is an error type for database driver errors
type StatusDatabaseDriverError string

func (e StatusInvalidDataError) Error() string    { return string(e) }
func (e StatusDatabaseDriverError) Error() string { return string(e) }

const (
	// StatusInvalidStatus is the default error for invalid statuss
	StatusInvalidStatus = StatusInvalidDataError("invalid status")
	// StatusInvalidID is the default error for non-existent IDs
	StatusInvalidID = StatusInvalidDataError("invalid ID")
)

// StatusDatabase is a service for writing status data to database
type StatusDatabase struct {
	Driver drivers.Database
}

// Write writes status data to the database
func (s *StatusDatabase) Write(temp *models.StatusData) error {
	if temp == nil {
		return StatusInvalidDataError("nil data")
	}
	if temp.Temperature < -30 || temp.Temperature > 50 {
		return StatusInvalidStatus
	}

	switch err := s.Driver.WriteStatus(temp); err.(type) {
	case nil:
		return nil
	case drivers.DatabaseInvalidDataError:
		return StatusInvalidID
	default:
		return StatusDatabaseDriverError(err.Error())
	}
}
