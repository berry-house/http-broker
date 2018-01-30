package drivers_test

import (
	"reflect"
	"testing"

	"github.com/berry-house/http-broker/drivers"
	"github.com/berry-house/http-broker/models"
)

func TestDatabaseInvalidDataErrorError(t *testing.T) {
	tests := map[string]struct {
		err      drivers.DatabaseInvalidDataError // error
		expected string                           // expected message
	}{
		"General test": {drivers.DatabaseInvalidDataError("some message"), "some message"},
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

func TestDatabaseUnexpectedError(t *testing.T) {
	tests := map[string]struct {
		err      drivers.DatabaseUnexpectedError // error
		expected string                          // expected message
	}{
		"General test": {drivers.DatabaseUnexpectedError("some message"), "some message"},
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

func TestDatabaseMemoryExists(t *testing.T) {
	// Setup
	driver, _ := drivers.NewDatabaseMemory(
		map[uint][]*models.TemperatureData{
			1: []*models.TemperatureData{},
			2: []*models.TemperatureData{},
			3: nil,
		},
	)

	testsSuccess := map[string]struct {
		id       uint // ID
		expected bool // expected result
	}{
		"Happy path":      {1, true},
		"Non-existent ID": {4, false},
	}
	for testName, testCase := range testsSuccess {
		t.Run(testName, func(t *testing.T) {
			exists, err := driver.Exists(testCase.id)
			if err != nil {
				t.Errorf("No error expected, gt %+v", err)
			}
			if exists != testCase.expected {
				t.Errorf("Expected %t, got %t", testCase.expected, exists)
			}
		})
	}
}

func TestDatabaseMemoryWriteTemperature(t *testing.T) {
	// Setup
	driver, _ := drivers.NewDatabaseMemory(
		map[uint][]*models.TemperatureData{
			1: []*models.TemperatureData{},
			2: []*models.TemperatureData{},
			3: []*models.TemperatureData{},
			4: []*models.TemperatureData{},
			5: nil,
		},
	)

	tests := map[string]struct {
		temp     *models.TemperatureData // input
		expected error                   // expected error
	}{
		"Happy path":   {&models.TemperatureData{ID: 1, Timestamp: 1516478286, Temperature: 23.5}, nil},
		"nil data":     {nil, drivers.DatabaseInvalidDataError("nil data")},
		"Invalid ID":   {&models.TemperatureData{ID: 6, Timestamp: 1516478286, Temperature: 20}, drivers.DatabaseInvalidDataError("invalid ID")},
		"Driver error": {&models.TemperatureData{ID: 5, Timestamp: 1516478286, Temperature: 20}, drivers.DatabaseUnexpectedError("nil list")},
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
