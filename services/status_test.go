package services_test

import (
	"reflect"
	"testing"

	"github.com/berry-house/http_broker/drivers/database"

	"github.com/berry-house/http_broker/models"
	"github.com/berry-house/http_broker/services"
)

func TestStatusInvalidDataError(t *testing.T) {
	tests := map[string]struct {
		err      services.StatusInvalidDataError // error
		expected string                          // expected message
	}{
		"General test": {services.StatusInvalidDataError("error message"), "error message"},
	}
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			errorMsg := testCase.err.Error()
			if errorMsg != testCase.expected {
				t.Errorf("Expected %s, got %s", testCase.expected, errorMsg)
			}
		})
	}
}

func TestStatusDatabaseDriverError(t *testing.T) {
	tests := map[string]struct {
		err      services.StatusDatabaseDriverError // error
		expected string                             // expected message
	}{
		"General test": {services.StatusDatabaseDriverError("error message"), "error message"},
	}
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			errorMsg := testCase.err.Error()
			if errorMsg != testCase.expected {
				t.Errorf("Expected %s, got %s", testCase.expected, errorMsg)
			}
		})
	}
}

type mockDatabaseDriver struct{}

func (d *mockDatabaseDriver) Exists(id uint) (bool, error) {
	if id > 0 && id < 5 {
		return true, nil
	}
	if id == 5 {
		return false, database.DatabaseUnexpectedError("mocked error")
	}

	return false, nil
}

func (d *mockDatabaseDriver) WriteStatus(data *models.StatusData) error {
	if data == nil {
		return database.DatabaseInvalidDataError("nil data")
	}
	// Mocked valid IDs
	if data.ID > 0 && data.ID < 5 {
		return nil
	}
	// Mocked driver error
	if data.ID == 0 || data.ID == 5 {
		return database.DatabaseUnexpectedError("mocked error")
	}

	return database.DatabaseInvalidDataError("invalid id")
}

func TestStatusWrite(t *testing.T) {
	// Setup
	service := services.StatusDatabase{
		Driver: &mockDatabaseDriver{},
	}

	tests := map[string]struct {
		data     *models.StatusData // input
		expected error              // expected error
	}{
		"Happy path":           {&models.StatusData{ID: 1, Timestamp: 1516478286, Temperature: 23}, nil},
		"nil data":             {nil, services.StatusInvalidDataError("nil data")},
		"Invalid ID":           {&models.StatusData{ID: 6, Timestamp: 1516478286}, services.StatusInvalidID},
		"Temperature too low":  {&models.StatusData{ID: 1, Timestamp: 1516478286, Temperature: -50}, services.StatusInvalidData},
		"Temperature too high": {&models.StatusData{ID: 1, Timestamp: 1516478286, Temperature: 56}, services.StatusInvalidData},
		"Light too high":       {&models.StatusData{ID: 1, Timestamp: 1516478286, Light: 153}, services.StatusInvalidData},
		"Humidity too high":    {&models.StatusData{ID: 1, Timestamp: 1516478286, Humidity: 105}, services.StatusInvalidData},
		"Database error":       {&models.StatusData{ID: 5, Timestamp: 1516478286, Temperature: 20}, services.StatusDatabaseDriverError("mocked error")},
	}
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			err := service.Write(testCase.data)
			if !reflect.DeepEqual(err, testCase.expected) {
				t.Errorf("Expected %+v, got %+v", testCase.expected, err)
			}
		})
	}
}
