package driver

import (
	"database/sql"
	"log"
	"time"
)

// DB hold the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const (
	pqMaxOpenDBConn = 10
	pqMaxIdleDBConn = 5
	pqMaxDBLifetime = 5 * time.Minute
)

//ConnectPQSQL will create connection for ElephantSQL
func ConnectPQSQL(dsn string) (*DB, error) {
	d, err := NewDatabase("postgres", dsn)
	if err != nil {
		log.Fatalln("failed to connect ElephantSQL", err)
		return nil, err
	}
	d.SetMaxOpenConns(pqMaxOpenDBConn)
	d.SetMaxIdleConns(pqMaxIdleDBConn)
	d.SetConnMaxIdleTime(pqMaxDBLifetime)

	dbConn.SQL = d
	log.Println("Connected to database...")
	return dbConn, nil
}

//NewDatabase will create the new database connection
func NewDatabase(databaseType, dsn string) (*sql.DB, error) {
	db, err := sql.Open(databaseType, dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
