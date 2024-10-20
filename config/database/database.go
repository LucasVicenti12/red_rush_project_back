package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type DatabaseConfig struct {
	username string
	password string
	host     string
	port     string
	dbname   string
}

func ConnectDatabase() (*sql.DB, error) {
	config := DatabaseConfig{
		username: os.Getenv("DB_USERNAME"),
		password: os.Getenv("DB_PASSWORD"),
		host:     os.Getenv("DB_HOST"),
		port:     os.Getenv("DB_PORT"),
		dbname:   os.Getenv("DB_NAME"),
	}

	conn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		config.username,
		config.password,
		config.host,
		config.port,
		config.dbname,
	)

	db, err := sql.Open("mysql", conn)

	if err != nil {
		return nil, err
	}

	return db, nil
}
