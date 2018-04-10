// Package database holds database drivers.
// A driver is the lowest functionality layer, interacting with resource sources.
// An example of functionality is accessing a database.
// Error returning should be related only to the sources.
package database

import "github.com/berry-house/http_broker/models"

// Database is an interface for database drivers
type Database interface {
	Exists(id uint) (bool, error)
	WriteStatus(data *models.StatusData) error
}

// DatabaseInvalidDataError is an error type for invalid data errors
type DatabaseInvalidDataError string

// DatabaseUnexpectedError is an error type for unhandled errors
type DatabaseUnexpectedError string

func (e DatabaseInvalidDataError) Error() string { return string(e) }
func (e DatabaseUnexpectedError) Error() string  { return string(e) }
