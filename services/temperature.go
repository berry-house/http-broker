package services

import (
	"github.com/berry-house/http-broker/drivers"
	"github.com/berry-house/http-broker/models"
)

// TemperatureInvalidDataError is an error type for invalid data errors
type TemperatureInvalidDataError string

// TemperatureDatabaseDriverError is an error type for database driver errors
type TemperatureDatabaseDriverError string

func (e TemperatureInvalidDataError) Error() string    { return string(e) }
func (e TemperatureDatabaseDriverError) Error() string { return string(e) }

const (
	// TemperatureInvalidTemperature is the default error for invalid temperatures
	TemperatureInvalidTemperature = TemperatureInvalidDataError("invalid temperature")
	// TemperatureInvalidID is the default error for non-existent IDs
	TemperatureInvalidID = TemperatureInvalidDataError("invalid ID")
)

// TemperatureDatabase is a service for writing temperature data to database
type TemperatureDatabase struct {
	Driver drivers.Database
}

// Write writes temperature data to the database
func (s *TemperatureDatabase) Write(temp *models.TemperatureData) error {
	if temp == nil {
		return TemperatureInvalidDataError("nil data")
	}
	if temp.Temperature < -30 || temp.Temperature > 50 {
		return TemperatureInvalidTemperature
	}

	switch err := s.Driver.WriteTemperature(temp); err.(type) {
	case nil:
		return nil
	case drivers.DatabaseInvalidDataError:
		return TemperatureInvalidID
	default:
		return TemperatureDatabaseDriverError(err.Error())
	}
}
