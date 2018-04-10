// Package controllers holds all controllers.
// A controller is a request dispatcher, interacting with services.
// An example of functionality is sending data to a database service.
// Error returning is direct to the client.
package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/berry-house/http_broker/models"
	"github.com/berry-house/http_broker/services"
	"github.com/berry-house/http_broker/util"
)

// Status is the controller for status data
type Status struct {
	Service services.Status
}

func (c *Status) Write(w http.ResponseWriter, r *http.Request) {
	// Body extraction
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.LogError(r, err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)

		return
	}
	var temp models.StatusData
	if err = json.Unmarshal(body, &temp); err != nil {
		http.Error(w, "Invalid body.", http.StatusBadRequest)

		return
	}

	// Using service
	switch err = c.Service.Write(&temp); err {
	case nil:
		w.Write([]byte("OK.\n"))
	case services.StatusInvalidID:
		http.Error(w, "Invalid ID.", http.StatusNotFound)
	case services.StatusInvalidData:
		http.Error(w, "Invalid data.", http.StatusBadRequest)
	default:
		util.LogError(r, err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
	}
}
