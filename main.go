package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
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

var (
	port             int
	runningMode      string
	loggerConfigFile string
	databaseAddress  string
	databaseName     string
	databaseUsername string
	databasePassword string
)

func init() {
	flag.IntVar(&port, "port", 8000, "Port in which the service listens")
	flag.StringVar(&runningMode, "runningMode", "", "Running mode of the server (either \"prod\" or \"test\")")
	flag.StringVar(&loggerConfigFile, "loggerConfigFile", "", "Path of JSON file for logging configuration.")
	flag.StringVar(&databaseAddress, "databaseAddress", "", "Address for the database")
	flag.StringVar(&databaseName, "databaseName", "", "Name of the database")
	flag.StringVar(&databaseUsername, "databaseUsername", "", "Username for the database")
	flag.StringVar(&databasePassword, "databasePassword", "", "Password for the database")
}

func main() {
	flag.Parse()

	var temperatureController controllers.Temperature

	switch runningMode {
	case "prod":
		// Drivers
		databaseDriver, err := drivers.NewDatabaseMySQL(
			fmt.Sprintf("%s:%s@tcp(%s)/%s", databaseUsername, databasePassword, databaseAddress, databaseName),
		)
		if err != nil {
			panic(err.Error())
		}
		// Services
		temperatureService := services.TemperatureDatabase{
			Driver: databaseDriver,
		}

		// Controllers
		temperatureController = controllers.Temperature{
			Service: &temperatureService,
		}
	case "test":
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
		temperatureController = controllers.Temperature{
			Service: &temperatureService,
		}
	default:
		panic("Invalid running mode. Use http_broker -h.")
	}

	// Logging
	loggerJSON, err := ioutil.ReadFile(loggerConfigFile)
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
