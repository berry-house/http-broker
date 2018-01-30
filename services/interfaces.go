// Package services holds all services.
// A service interacts with drivers.
// An example of functionality is validating data (in values) and writing it to the database.
// Error returning should be related only to data logical errors.
package services

import "github.com/berry-house/http-broker/models"

// Temperature is an inteface for temperature services
type Temperature interface {
	Write(temp *models.TemperatureData) error
}
