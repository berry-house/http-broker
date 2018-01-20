package drivers

import "github.com/berry-house/http-broker/models"

// DatabaseError is a type for driver errors
type DatabaseError string

var _ error = DatabaseError("")

func (e DatabaseError) Error() string { return string(e) }

const (
	// DatabaseInvalidIDError is the default error for an invalid ID
	DatabaseInvalidIDError = DatabaseError("invalid ID")
	// DatabaseUnexpectedError is a general error for unexpected errors
	DatabaseUnexpectedError = DatabaseError("unexpected error")
)

// DatabaseMemory is an in-memory database driver
type DatabaseMemory struct {
	Data map[uint][]*models.TemperatureData
}

var _ Database = (*DatabaseMemory)(nil)

// WriteTemperature writes temperature data into memory
func (d *DatabaseMemory) WriteTemperature(temp *models.TemperatureData) error {
	if d == nil || d.Data == nil || temp == nil {
		return DatabaseUnexpectedError
	}
	list, ok := d.Data[temp.ID]
	if !ok {
		return DatabaseInvalidIDError
	}
	if list == nil {
		return DatabaseUnexpectedError
	}
	d.Data[temp.ID] = append(list, temp)

	return nil
}
