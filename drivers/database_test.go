package drivers_test

import (
	"reflect"
	"testing"

	"github.com/berry-house/http-broker/drivers"
	"github.com/berry-house/http-broker/models"
)

func TestDatabaseError(t *testing.T) {
	tests := map[string]struct {
		err      drivers.DatabaseError // error
		expected string                // expected message
	}{
		"Invalid ID error": {drivers.DatabaseInvalidIDError, "invalid ID"},
		"Unexpected error": {drivers.DatabaseUnexpectedError, "unexpected error"},
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

func TestDatabaseMemoryWriteTemperature(t *testing.T) {
	// Setup
	driver := drivers.DatabaseMemory{
		Data: map[uint][]*models.TemperatureData{
			1: []*models.TemperatureData{},
			2: []*models.TemperatureData{},
			3: []*models.TemperatureData{},
			4: []*models.TemperatureData{},
			5: nil,
		},
	}

	tests := map[string]struct {
		temp     *models.TemperatureData // input
		expected error                   // expected error
	}{
		"Happy path":   {&models.TemperatureData{ID: 1, Timestamp: 1516478286, Temperature: 23.5}, nil},
		"nil data":     {nil, drivers.DatabaseUnexpectedError},
		"Invalid ID":   {&models.TemperatureData{ID: 6, Timestamp: 1516478286, Temperature: 20}, drivers.DatabaseInvalidIDError},
		"Driver error": {&models.TemperatureData{ID: 5, Timestamp: 1516478286, Temperature: 20}, drivers.DatabaseUnexpectedError},
	}
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			err := driver.WriteTemperature(testCase.temp)
			if !reflect.DeepEqual(err, testCase.expected) {
				t.Errorf("Expected %+v, got %+v", testCase.expected, err)
			}
		})
	}
}
