package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/berry-house/http_broker/controllers"
	"github.com/berry-house/http_broker/drivers"
	"github.com/berry-house/http_broker/models"
	"github.com/berry-house/http_broker/services"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
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

	// Logging
	loggerJSON, err := ioutil.ReadFile("conf/logger.json")
	if err != nil {
		panic(err)
	}

	var loggerConfig zap.Config
	err = json.Unmarshal(loggerJSON, &loggerConfig)
	if err != nil {
		panic(err)
	}

	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	// Context
	ctx := context.Background()
	ctx = context.WithValue(ctx, "logger", logger)

	// Router
	router := mux.NewRouter()
	router.HandleFunc("/broker/temperature", temperatureController.Write).Methods("POST")
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rctx := r.WithContext(ctx)
			next.ServeHTTP(w, rctx)
		})
	})

	// Server
	server := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8000",
	}

	log.Fatal(server.ListenAndServe())
}
