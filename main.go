package main

import (
	"log"
	"net/http"

	"github.com/berry-house/http-broker/controllers"
	"github.com/berry-house/http-broker/drivers"
	"github.com/berry-house/http-broker/models"
	"github.com/berry-house/http-broker/services"
	"github.com/gorilla/mux"
)

func main() {
	// Drivers
	databaseDriver, _ := drivers.NewDatabaseMemory(
		map[uint][]*models.TemperatureData{
			1: []*models.TemperatureData{},
			2: []*models.TemperatureData{},
			3: []*models.TemperatureData{},
			4: []*models.TemperatureData{},
			5: []*models.TemperatureData{},
		},
	)

	// Services
	temperatureService := services.TemperatureDatabase{
		Driver: databaseDriver,
	}

	// Controllers
	temperatureController := controllers.Temperature{
		Service: &temperatureService,
	}

	// Router
	router := mux.NewRouter()
	router.HandleFunc("/broker/temperature", temperatureController.Write).Methods("POST")

	// Server
	server := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8000",
	}
	log.Fatal(server.ListenAndServe())
}
