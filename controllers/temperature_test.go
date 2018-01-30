package controllers_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/berry-house/http-broker/controllers"
	"github.com/berry-house/http-broker/models"
	"github.com/berry-house/http-broker/services"
)

// Handler mock
type mockHandlerTemperature struct {
	c controllers.Temperature
}

func (h *mockHandlerTemperature) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Routing and methods should be handled by the server, not by the controller
	h.c.Write(w, r)
}

// Service mock
type mockTemperatureService struct{}

var _ services.Temperature = (*mockTemperatureService)(nil)

func (s *mockTemperatureService) Write(temp *models.TemperatureData) error {
	if temp == nil {
		return services.TemperatureInvalidDataError("nil data")
	}
	if temp.Temperature < -30 || temp.Temperature > 50 {
		return services.TemperatureInvalidTemperature
	}
	// Mocked valid IDs
	if temp.ID > 0 && temp.ID < 5 {
		return nil
	}
	// Mocked driver error
	if temp.ID == 0 || temp.ID == 5 {
		return services.TemperatureDatabaseDriverError("mocked error")
	}

	return services.TemperatureInvalidID
}

// Utilities
func buildTemperatureRequest(method, path string, body []byte) *http.Request {
	req, err := http.NewRequest(method, path, bytes.NewReader(body))
	if err != nil {
		panic(err.Error())
	}

	return req
}

func TestWriteTemperature(t *testing.T) {
	// Setup
	handler := &mockHandlerTemperature{
		c: controllers.Temperature{
			Service: &mockTemperatureService{},
		},
	}
	server := httptest.NewServer(handler)
	tests := map[string]struct {
		request            *http.Request // input
		expectedStatus     string        // expected status
		expectedStatusCode int           // expected status code
	}{
		"Happy path": {
			request:            buildTemperatureRequest("POST", server.URL, []byte(`{"id":1,"timestamp":1516472722,"temperature":20}`)),
			expectedStatus:     "OK.\n",
			expectedStatusCode: http.StatusOK,
		},
		"No body": {
			request:            buildTemperatureRequest("POST", server.URL, nil),
			expectedStatus:     "Invalid body.\n",
			expectedStatusCode: http.StatusBadRequest,
		},
		"Invalid ID": {
			request:            buildTemperatureRequest("POST", server.URL, []byte(`{"id":-1,"timestamp":1516472722,"temperature":20}`)),
			expectedStatus:     "Invalid body.\n",
			expectedStatusCode: http.StatusBadRequest,
		},
		"Non-existing ID": {
			request:            buildTemperatureRequest("POST", server.URL, []byte(`{"id":7,"timestamp":1516472722,"temperature":20}`)),
			expectedStatus:     "Invalid ID.\n",
			expectedStatusCode: http.StatusNotFound,
		},
		"Temperature too high": {
			request:            buildTemperatureRequest("POST", server.URL, []byte(`{"id":1,"temperature":163.4}`)),
			expectedStatus:     "Invalid temperature.\n",
			expectedStatusCode: http.StatusBadRequest,
		},
		"Temperature too low": {
			request:            buildTemperatureRequest("POST", server.URL, []byte(`{"id":1,"temperature":-80.4}`)),
			expectedStatus:     "Invalid temperature.\n",
			expectedStatusCode: http.StatusBadRequest,
		},
		"Database error": {
			request:            buildTemperatureRequest("POST", server.URL, []byte(`{"id":5,"timestamp":1516472722,"temperature":20}`)),
			expectedStatus:     "Internal server error.\n",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for testName, testCase := range tests {
		t.Run(testName, func(t *testing.T) {
			client := http.DefaultClient
			response, err := client.Do(testCase.request)
			if err != nil {
				t.Errorf("No error expected, got %+v", err)
			}
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Errorf("No error expected, got %+v", err)
			}
			status := string(body)
			if status != testCase.expectedStatus ||
				response.StatusCode != testCase.expectedStatusCode {
				t.Errorf("Expected %d: %q, got %d: %q", testCase.expectedStatusCode, testCase.expectedStatus, response.StatusCode, status)
			}
		})
	}
}
