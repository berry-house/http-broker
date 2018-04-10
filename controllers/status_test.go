package controllers_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/berry-house/http_broker/controllers"
	"github.com/berry-house/http_broker/models"
	"github.com/berry-house/http_broker/services"
)

// Handler mock
type mockHandlerStatus struct {
	c controllers.Status
}

func (h *mockHandlerStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Routing and methods should be handled by the server, not by the controller
	h.c.Write(w, r)
}

// Service mock
type mockStatusService struct{}

var _ services.Status = (*mockStatusService)(nil)

func (s *mockStatusService) Write(data *models.StatusData) error {
	if data == nil {
		return services.StatusInvalidDataError("nil data")
	}

	// Threshold values
	if data.Temperature < -30 || data.Temperature > 50 ||
		data.Light > 150 ||
		data.Humidity > 100 {
		return services.StatusInvalidData
	}

	// Mocked valid IDs
	if data.ID > 0 && data.ID < 5 {
		return nil
	}
	// Mocked driver error
	if data.ID == 0 || data.ID == 5 {
		return services.StatusDatabaseDriverError("mocked error")
	}

	return services.StatusInvalidID
}

// Utilities
func buildStatusRequest(method, path string, body []byte) *http.Request {
	req, err := http.NewRequest(method, path, bytes.NewReader(body))
	if err != nil {
		panic(err.Error())
	}

	return req
}

func TestWriteStatus(t *testing.T) {
	// Setup
	handler := &mockHandlerStatus{
		c: controllers.Status{
			Service: &mockStatusService{},
		},
	}
	server := httptest.NewServer(handler)
	defer server.Close()

	tests := map[string]struct {
		request            *http.Request // input
		expectedStatus     string        // expected status
		expectedStatusCode int           // expected status code
	}{
		"Happy path": {
			request:            buildStatusRequest("POST", server.URL, []byte(`{"id":1,"timestamp":1516472722,"status":20}`)),
			expectedStatus:     "OK.\n",
			expectedStatusCode: http.StatusOK,
		},
		"No body": {
			request:            buildStatusRequest("POST", server.URL, nil),
			expectedStatus:     "Invalid body.\n",
			expectedStatusCode: http.StatusBadRequest,
		},
		"Invalid ID": {
			request:            buildStatusRequest("POST", server.URL, []byte(`{"id":-1,"timestamp":1516472722}`)),
			expectedStatus:     "Invalid body.\n",
			expectedStatusCode: http.StatusBadRequest,
		},
		"Non-existing ID": {
			request:            buildStatusRequest("POST", server.URL, []byte(`{"id":7,"timestamp":1516472722}`)),
			expectedStatus:     "Invalid ID.\n",
			expectedStatusCode: http.StatusNotFound,
		},
		"Temperature too high": {
			request:            buildStatusRequest("POST", server.URL, []byte(`{"id":1,"timestamp":1516472722,"temperature":163}`)),
			expectedStatus:     "Invalid data.\n",
			expectedStatusCode: http.StatusBadRequest,
		},
		"Temperature too low": {
			request:            buildStatusRequest("POST", server.URL, []byte(`{"id":1,"timestamp":1516472722,"temperature":-80}`)),
			expectedStatus:     "Invalid data.\n",
			expectedStatusCode: http.StatusBadRequest,
		},
		"Humidity too high": {
			request:            buildStatusRequest("POST", server.URL, []byte(`{"id":1,"timestamp":1516472722,"humidity":103}`)),
			expectedStatus:     "Invalid data.\n",
			expectedStatusCode: http.StatusBadRequest,
		},
		"Light too high": {
			request:            buildStatusRequest("POST", server.URL, []byte(`{"id":1,"timestamp":1516472722,"light":165}`)),
			expectedStatus:     "Invalid data.\n",
			expectedStatusCode: http.StatusBadRequest,
		},
		"Database error": {
			request:            buildStatusRequest("POST", server.URL, []byte(`{"id":5,"timestamp":1516472722,"status":20}`)),
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
