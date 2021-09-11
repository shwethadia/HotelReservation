package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

//DB holds the connection pool
type DB struct {

	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenConn = 10

const maxIdleDbConn = 5

const maxDbLifeTime = 5 * time.Minute

//ConnectSQL creates database pool for postgress
func ConnectSQL(dsn string) (*DB, error) {

	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	d.SetMaxOpenConns(maxOpenConn)

	d.SetConnMaxIdleTime(maxIdleDbConn)
	
	d.SetConnMaxLifetime(maxDbLifeTime)

	dbConn.SQL = d

	err = testDB(d)
	if err != nil {

		return nil, err
	}

	return dbConn, nil

}

//Tries to ping the database
func testDB(d *sql.DB) error {

	err := d.Ping()
	if err != nil {
		return err
	}

	return nil
}

//NewDatabase created a new database for the application
func NewDatabase(dsn string) (*sql.DB, error) {

	db, err := sql.Open("pgx", dsn)

	if err != nil {

		return nil, err
	}

	if err = db.Ping(); err != nil {

		return nil, err
	}

	return db, nil

}
