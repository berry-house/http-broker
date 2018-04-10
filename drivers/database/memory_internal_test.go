package database

import (
	"reflect"
	"testing"

	"github.com/berry-house/http_broker/models"
)

func TestNewMemory(t *testing.T) {
	testsSuccessful := map[string]struct {
		data     map[uint][]*models.StatusData // input
		expected Database                      // expected driver
	}{
		"Happy path": {
			data: map[uint][]*models.StatusData{
				1: []*models.StatusData{},
			},
			expected: &Memory{
				data: map[uint][]*models.StatusData{
					1: []*models.StatusData{},
				},
			},
		},
		"nil list": {
			data: map[uint][]*models.StatusData{
				1: nil,
			},
			expected: &Memory{
				data: map[uint][]*models.StatusData{
					1: nil,
				},
			},
		},
	}
	for testName, testCase := range testsSuccessful {
		t.Run(testName, func(t *testing.T) {
			d, err := NewMemory(testCase.data)
			if err != nil {
				t.Errorf("No error expected, got %+v", err)
			}

			if !reflect.DeepEqual(d, testCase.expected) {
				t.Errorf("Expected %+v, got %+v", testCase.expected, d)
			}
		})
	}

	testsFailure := map[string]struct {
		data     map[uint][]*models.StatusData // input
		expected error                         // expected error
	}{
		"nil data": {
			data:     nil,
			expected: DatabaseInvalidDataError("nil data"),
		},
	}
	for testName, testCase := range testsFailure {
		t.Run(testName, func(t *testing.T) {
			_, err := NewMemory(testCase.data)
			if err == nil {
				t.Error("Error expected")
			}

			if !reflect.DeepEqual(err, testCase.expected) {
				t.Errorf("Expected %+v, got %+v", testCase.expected, err)
			}
		})
	}
}
