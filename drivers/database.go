package drivers

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/berry-house/http_broker/models"
	_ "github.com/go-sql-driver/mysql" // MySQL
)

// DatabaseInvalidDataError is an error type for invalid data errors
type DatabaseInvalidDataError string

// DatabaseUnexpectedError is an error type for unhandled errors
type DatabaseUnexpectedError string

func (e DatabaseInvalidDataError) Error() string { return string(e) }
func (e DatabaseUnexpectedError) Error() string  { return string(e) }

// DatabaseMemory is an in-memory database driver
type DatabaseMemory struct {
	data map[uint][]*models.TemperatureData
}

// NewDatabaseMemory creates a new DatabaseMemory driver
func NewDatabaseMemory(data map[uint][]*models.TemperatureData) (Database, error) {
	if data == nil {
		return nil, DatabaseInvalidDataError("nil data")
	}

	return &DatabaseMemory{data: data}, nil
}

// Exists checks if current ID exists
func (d *DatabaseMemory) Exists(id uint) (bool, error) {
	_, ok := d.data[id]
	if !ok {
		return false, nil
	}

	return true, nil
}

// WriteTemperature writes temperature data into memory
func (d *DatabaseMemory) WriteTemperature(temp *models.TemperatureData) error {
	if d == nil {
		return DatabaseUnexpectedError("nil driver")
	}
	if temp == nil {
		return DatabaseInvalidDataError("nil data")
	}

	list, ok := d.data[temp.ID]
	if !ok {
		return DatabaseInvalidDataError("invalid ID")
	}
	if list == nil {
		return DatabaseUnexpectedError("nil list")
	}
	d.data[temp.ID] = append(list, temp)

	return nil
}

const (
	plantQuery        = `SELECT COUNT(*) FROM plant WHERE id = ?;`
	temperatureInsert = `REPLACE INTO conditions(plantID, time, airTemperature) VALUES(?, ?, ?);`
)

// DatabaseMySQL is a MySQL database driver
type DatabaseMySQL struct {
	database *sql.DB
}

// NewDatabaseMySQL creates a new DatabaseMySQL driver
func NewDatabaseMySQL(conn string) (*DatabaseMySQL, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println(err)
		return nil, DatabaseUnexpectedError(err.Error())
	}
	err = db.Ping()
	if err != nil {
		return nil, DatabaseUnexpectedError(err.Error())
	}

	return &DatabaseMySQL{database: db}, nil
}

// Exists checks if current id exists
func (d *DatabaseMySQL) Exists(id uint) (bool, error) {
	var rowsNumber int
	err := d.database.QueryRow(plantQuery, id).Scan(&rowsNumber)
	if err != nil {
		return false, DatabaseUnexpectedError(err.Error())
	}

	return rowsNumber != 0, nil
}

// WriteTemperature writes temperature data into memory
func (d *DatabaseMySQL) WriteTemperature(temp *models.TemperatureData) error {
	if d == nil {
		return DatabaseUnexpectedError("nil driver")
	}
	if temp == nil {
		return DatabaseInvalidDataError("nil data")
	}

	// Check if ID is valid
	exists, err := d.Exists(temp.ID)
	if err != nil {
		return DatabaseUnexpectedError(err.Error())
	}
	if !exists {
		return DatabaseInvalidDataError("non-existent ID")
	}

	// Insert
	timestamp := time.Unix(temp.Timestamp, 0)
	timestampString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		timestamp.Year(), int(timestamp.Month()), timestamp.Day(), timestamp.Hour(), timestamp.Minute(), timestamp.Second())
	insert, err := d.database.Prepare(temperatureInsert)
	if err != nil {
		return DatabaseUnexpectedError(err.Error())
	}
	_, err = insert.Exec(temp.ID, timestampString, temp.Temperature)
	if err != nil {
		return DatabaseUnexpectedError(err.Error())
	}

	return nil
}
