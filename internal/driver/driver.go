package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	maxOpenDbConn = 10
	maxIdleDbConn = 5
	maxDbConnLife = 5 * time.Minute
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

// ConnectSQL creates connection pool for Postgres
func ConnectSQL(dsn string) (*DB, error) {
	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDbConnLife)

	dbConn.SQL = d
	err = testDb(d)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// testDb pings the database
func testDb(d *sql.DB) error {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewDatabase creates a new DB connetion for the application
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
