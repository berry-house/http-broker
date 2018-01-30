package services_test

import (
	"reflect"
	"testing"

	"github.com/berry-house/http-broker/drivers"

	"github.com/berry-house/http-broker/models"
	"github.com/berry-house/http-broker/services"
)

func TestTemperatureInvalidDataError(t *testing.T) {
	tests := map[string]struct {
		err      services.TemperatureInvalidDataError // error
		expected string                               // expected message
	}{
		"General test": {services.TemperatureInvalidDataError("error message"), "error message"},
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

func TestTemperatureDatabaseDriverError(t *testing.T) {
	tests := map[string]struct {
		err      services.TemperatureDatabaseDriverError // error
		expected string                                  // expected message
	}{
		"General test": {services.TemperatureDatabaseDriverError("error message"), "error message"},
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

func (d *mockDatabaseDriver) WriteTemperature(temp *models.TemperatureData) error {
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

func TestTemperatureWrite(t *testing.T) {
	// Setup
	service := services.TemperatureDatabase{
		Driver: &mockDatabaseDriver{},
	}

	tests := map[string]struct {
		temp     *models.TemperatureData // input
		expected error                   // expected error
	}{
		"Happy path":           {&models.TemperatureData{ID: 1, Timestamp: 1516478286, Temperature: 23.5}, nil},
		"nil data":             {nil, services.TemperatureInvalidDataError("nil data")},
		"Invalid ID":           {&models.TemperatureData{ID: 6, Timestamp: 1516478286, Temperature: 20}, services.TemperatureInvalidID},
		"Temperature too low":  {&models.TemperatureData{ID: 1, Timestamp: 1516478286, Temperature: -50.3}, services.TemperatureInvalidTemperature},
		"Temperature too high": {&models.TemperatureData{ID: 1, Timestamp: 1516478286, Temperature: 56.3}, services.TemperatureInvalidTemperature},
		"Database error":       {&models.TemperatureData{ID: 5, Timestamp: 1516478286, Temperature: 20}, services.TemperatureDatabaseDriverError("mocked error")},
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
