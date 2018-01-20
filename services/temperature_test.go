package services_test

import (
	"reflect"
	"testing"

	"github.com/berry-house/http-broker/drivers"

	"github.com/berry-house/http-broker/models"
	"github.com/berry-house/http-broker/services"
)

func TestDatabaseError(t *testing.T) {
	tests := map[string]struct {
		err      services.TemperatureError // error
		expected string                    // expected message
	}{
		"nil temperature data":      {services.TemperatureNilDataError, "nil temperature data"},
		"Invalid ID error":          {services.TemperatureInvalidIDError, "invalid ID"},
		"Invalid temperature error": {services.TemperatureInvalidTemperatureError, "invalid temperature"},
		"Database driver error":     {services.TemperatureDatabaseDriverError, "database driver error"},
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

var _ drivers.Database = (*mockDatabaseDriver)(nil)

func (d *mockDatabaseDriver) WriteTemperature(temp *models.TemperatureData) error {
	if temp == nil {
		return services.TemperatureNilDataError
	}
	// Mocked valid IDs
	if temp.ID < 1 || temp.ID > 5 {
		return drivers.DatabaseInvalidIDError
	}
	// Mocked driver error
	if temp.ID == 5 {
		return drivers.DatabaseUnexpectedError
	}

	return nil
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
		"nil data":             {nil, services.TemperatureNilDataError},
		"Invalid ID":           {&models.TemperatureData{ID: 6, Timestamp: 1516478286, Temperature: 20}, services.TemperatureInvalidIDError},
		"Temperature too low":  {&models.TemperatureData{ID: 1, Timestamp: 1516478286, Temperature: -50.3}, services.TemperatureInvalidTemperatureError},
		"Temperature too high": {&models.TemperatureData{ID: 1, Timestamp: 1516478286, Temperature: 56.3}, services.TemperatureInvalidTemperatureError},
		"Database error":       {&models.TemperatureData{ID: 5, Timestamp: 1516478286, Temperature: 20}, services.TemperatureDatabaseDriverError},
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
