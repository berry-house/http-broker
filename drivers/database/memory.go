package database

import "github.com/berry-house/http_broker/models"

// Memory is an in-memory database driver
type Memory struct {
	data map[uint][]*models.StatusData
}

// NewMemory creates a new DatabaseMemory driver
func NewMemory(data map[uint][]*models.StatusData) (*Memory, error) {
	if data == nil {
		return nil, DatabaseInvalidDataError("nil data")
	}

	return &Memory{data: data}, nil
}

// Exists checks if current ID exists
func (d *Memory) Exists(id uint) (bool, error) {
	_, ok := d.data[id]
	if !ok {
		return false, nil
	}

	return true, nil
}

// WriteStatus writes status data into memory
func (d *Memory) WriteStatus(temp *models.StatusData) error {
	if d == nil {
		return DatabaseUnexpectedError("nil driver")
	}
	if temp == nil {
		return DatabaseInvalidDataError("nil data")
	}

	list, ok := d.data[temp.ID]
	if !ok {
		return DatabaseInvalidDataError("invalid ID")
	}
	if list == nil {
		return DatabaseUnexpectedError("nil list")
	}
	d.data[temp.ID] = append(list, temp)

	return nil
}
