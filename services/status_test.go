package services_test

import (
	"reflect"
	"testing"

	"github.com/berry-house/http_broker/drivers"

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
		return false, drivers.DatabaseUnexpectedError("mocked error")
	}

	return false, nil
}

func (d *mockDatabaseDriver) WriteStatus(temp *models.StatusData) error {
	if temp == nil {
		return drivers.DatabaseInvalidDataError("nil data")
	}
	// Mocked valid IDs
	if temp.ID > 0 && temp.ID < 5 {
		return nil
	}
	// Mocked driver error
	if temp.ID == 0 || temp.ID == 5 {
		return drivers.DatabaseUnexpectedError("mocked error")
	}

	return drivers.DatabaseInvalidDataError("invalid id")
}

func TestStatusWrite(t *testing.T) {
	// Setup
	service := services.StatusDatabase{
		Driver: &mockDatabaseDriver{},
	}

	tests := map[string]struct {
		temp     *models.StatusData // input
		expected error              // expected error
	}{
		"Happy path":      {&models.StatusData{ID: 1, Timestamp: 1516478286, Temperature: 23.5}, nil},
		"nil data":        {nil, services.StatusInvalidDataError("nil data")},
		"Invalid ID":      {&models.StatusData{ID: 6, Timestamp: 1516478286, Temperature: 20}, services.StatusInvalidID},
		"Status too low":  {&models.StatusData{ID: 1, Timestamp: 1516478286, Temperature: -50.3}, services.StatusInvalidStatus},
		"Status too high": {&models.StatusData{ID: 1, Timestamp: 1516478286, Temperature: 56.3}, services.StatusInvalidStatus},
		"Database error":  {&models.StatusData{ID: 5, Timestamp: 1516478286, Temperature: 20}, services.StatusDatabaseDriverError("mocked error")},
	}
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			err := service.Write(testCase.temp)
			if !reflect.DeepEqual(err, testCase.expected) {
				t.Errorf("Expected %+v, got %+v", testCase.expected, err)
			}
		})
	}
}
