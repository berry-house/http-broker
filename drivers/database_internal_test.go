package drivers

import (
	"reflect"
	"testing"

	"github.com/berry-house/http-broker/models"
)

func TestNewDatabaseMemory(t *testing.T) {
	testsSuccessful := map[string]struct {
		data     map[uint][]*models.TemperatureData // input
		expected Database                           // expected driver
	}{
		"Happy path": {
			data: map[uint][]*models.TemperatureData{
				1: []*models.TemperatureData{},
			},
			expected: &DatabaseMemory{
				data: map[uint][]*models.TemperatureData{
					1: []*models.TemperatureData{},
				},
			},
		},
		"nil list": {
			data: map[uint][]*models.TemperatureData{
				1: nil,
			},
			expected: &DatabaseMemory{
				data: map[uint][]*models.TemperatureData{
					1: nil,
				},
			},
		},
	}
	for testName, testCase := range testsSuccessful {
		t.Run(testName, func(t *testing.T) {
			d, err := NewDatabaseMemory(testCase.data)
			if err != nil {
				t.Errorf("No error expected, got %+v", err)
			}

			if !reflect.DeepEqual(d, testCase.expected) {
				t.Errorf("Expected %+v, got %+v", testCase.expected, d)
			}
		})
	}

	testsFailure := map[string]struct {
		data     map[uint][]*models.TemperatureData // input
		expected error                              // expected error
	}{
		"nil data": {
			data:     nil,
			expected: DatabaseInvalidDataError("nil data"),
		},
	}
	for testName, testCase := range testsFailure {
		t.Run(testName, func(t *testing.T) {
			_, err := NewDatabaseMemory(testCase.data)
			if err == nil {
				t.Error("Error expected")
			}

			if !reflect.DeepEqual(err, testCase.expected) {
				t.Errorf("Expected %+v, got %+v", testCase.expected, err)
			}
		})
	}
}
