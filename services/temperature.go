package services

import (
	"github.com/berry-house/http-broker/drivers"
	"github.com/berry-house/http-broker/models"
)

// TemperatureError is a type for service errors
type TemperatureError string

func (e TemperatureError) Error() string { return string(e) }

var _ error = TemperatureError("")

const (
	// TemperatureNilDataError is the default error for nil temperature data
	TemperatureNilDataError = TemperatureError("nil temperature data")
	// TemperatureInvalidIDError is the default error for invalid ID
	TemperatureInvalidIDError = TemperatureError("invalid ID")
	// TemperatureInvalidTemperatureError is the default error for too high or too low temperatures
	TemperatureInvalidTemperatureError = TemperatureError("invalid temperature")
	// TemperatureDatabaseDriverError is the default error for unidentified errors drivers
	TemperatureDatabaseDriverError = TemperatureError("database driver error")
)

var _ Temperature = (*TemperatureDatabase)(nil)

// TemperatureDatabase is a service for writing temperature data to database
type TemperatureDatabase struct {
	Driver drivers.Database
}

// Write writes temperature data to the database
func (s *TemperatureDatabase) Write(temp *models.TemperatureData) error {
	// Validating data
	if temp == nil {
		return TemperatureNilDataError
	}
	if temp.Temperature < -30 || temp.Temperature > 50 {
		return TemperatureInvalidTemperatureError
	}

	switch err := s.Driver.WriteTemperature(temp); err {
	case nil:
		return nil
	case drivers.DatabaseInvalidIDError:
		return TemperatureInvalidIDError
	default:
		return TemperatureDatabaseDriverError
	}
}
