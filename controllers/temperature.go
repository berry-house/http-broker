// Package controllers holds all controllers.
// A controller is a request dispatcher, interacting with services.
// An example of functionality is sending data to a database service.
// Error returning is direct to the client.
package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/berry-house/http-broker/models"
	"github.com/berry-house/http-broker/services"
)

// Temperature is the controller for temperature data
type Temperature struct {
	Service services.Temperature
}

func (c *Temperature) Write(w http.ResponseWriter, r *http.Request) {
	// Body extraction
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal server error.", http.StatusInternalServerError)

		return
	}
	var temp models.TemperatureData
	if err = json.Unmarshal(body, &temp); err != nil {
		http.Error(w, "Invalid body.", http.StatusBadRequest)

		return
	}

	// Using service
	switch err = c.Service.Write(&temp); err {
	case nil:
		w.Write([]byte("OK.\n"))
	case services.TemperatureInvalidID:
		http.Error(w, "Invalid ID.", http.StatusNotFound)
	case services.TemperatureInvalidTemperature:
		http.Error(w, "Invalid temperature.", http.StatusBadRequest)
	default:
		// TODO: log error
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
	}
}
