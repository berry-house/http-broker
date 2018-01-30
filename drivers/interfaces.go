// Package drivers holds all drivers.
// A driver is the lowest functionality layer, interacting with resource sources.
// An example of functionality is accessing a database.
// Error returning should be related only to the sources.
package drivers

import "github.com/berry-house/http-broker/models"

// Database is an interface for database drivers
type Database interface {
	Exists(id uint) (bool, error)
	WriteTemperature(t *models.TemperatureData) error
}
