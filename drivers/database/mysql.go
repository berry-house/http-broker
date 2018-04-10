package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/berry-house/http_broker/models"
	_ "github.com/go-sql-driver/mysql" // MySQL
)

const (
	plantQuery   = `SELECT COUNT(*) FROM plant WHERE id = ?;`
	statusInsert = `REPLACE INTO
						conditions(plantID, time, lightIntensity, soilHumidity, airTemperature)
					VALUES(?, ?, ?, ?, ?);`
)

// MySQL is a MySQL database driver
type MySQL struct {
	database *sql.DB
}

// NewMySQL creates a new MySQL driver
func NewMySQL(conn string) (*MySQL, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println(err)
		return nil, DatabaseUnexpectedError(err.Error())
	}
	err = db.Ping()
	if err != nil {
		return nil, DatabaseUnexpectedError(err.Error())
	}

	return &MySQL{database: db}, nil
}

// Exists checks if current id exists
func (d *MySQL) Exists(id uint) (bool, error) {
	var rowsNumber int
	err := d.database.QueryRow(plantQuery, id).Scan(&rowsNumber)
	if err != nil {
		return false, DatabaseUnexpectedError(err.Error())
	}

	return rowsNumber != 0, nil
}

// WriteStatus writes status data into memory
func (d *MySQL) WriteStatus(temp *models.StatusData) error {
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
	insert, err := d.database.Prepare(statusInsert)
	if err != nil {
		return DatabaseUnexpectedError(err.Error())
	}
	_, err = insert.Exec(temp.ID, timestampString, temp.Light, temp.Humidity, temp.Temperature)
	if err != nil {
		return DatabaseUnexpectedError(err.Error())
	}

	return nil
}
