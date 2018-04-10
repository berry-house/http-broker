package database_test

import (
	"reflect"
	"testing"

	"github.com/berry-house/http_broker/drivers/database"
	"github.com/berry-house/http_broker/models"
)

func TestDatabaseInvalidDataErrorError(t *testing.T) {
	tests := map[string]struct {
		err      database.DatabaseInvalidDataError // error
		expected string                            // expected message
	}{
		"General test": {database.DatabaseInvalidDataError("some message"), "some message"},
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

func TestUnexpectedError(t *testing.T) {
	tests := map[string]struct {
		err      database.DatabaseUnexpectedError // error
		expected string                           // expected message
	}{
		"General test": {database.DatabaseUnexpectedError("some message"), "some message"},
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

func TestMemoryExists(t *testing.T) {
	// Setup
	driver, _ := database.NewMemory(
		map[uint][]*models.StatusData{
			1: []*models.StatusData{},
			2: []*models.StatusData{},
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

func TestMemoryWriteStatus(t *testing.T) {
	// Setup
	driver, _ := database.NewMemory(
		map[uint][]*models.StatusData{
			1: []*models.StatusData{},
			2: []*models.StatusData{},
			3: []*models.StatusData{},
			4: []*models.StatusData{},
			5: nil,
		},
	)

	tests := map[string]struct {
		data     *models.StatusData // input
		expected error              // expected error
	}{
		"Happy path":   {&models.StatusData{ID: 1, Timestamp: 1516478286}, nil},
		"nil data":     {nil, database.DatabaseInvalidDataError("nil data")},
		"Invalid ID":   {&models.StatusData{ID: 6, Timestamp: 1516478286}, database.DatabaseInvalidDataError("invalid ID")},
		"Driver error": {&models.StatusData{ID: 5, Timestamp: 1516478286}, database.DatabaseUnexpectedError("nil list")},
	}
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			err := driver.WriteStatus(testCase.data)
			if !reflect.DeepEqual(err, testCase.expected) {
				t.Errorf("Expected %+v, got %+v", testCase.expected, err)
			}
		})
	}
}
