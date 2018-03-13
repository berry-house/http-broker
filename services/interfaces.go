// Package services holds all services.
// A service interacts with drivers.
// An example of functionality is validating data (in values) and writing it to the database.
// Error returning should be related only to data logical errors.
package services

import "github.com/berry-house/http_broker/models"

// Status is an inteface for status services
type Status interface {
	Write(temp *models.StatusData) error
}
