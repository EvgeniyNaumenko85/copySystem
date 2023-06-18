package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var database *sql.DB

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// Initiation of connection with DB
func initDB(cfg Config) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	//if err != nil {
	//	panic(err)
	//}
	if err != nil {
		log.Fatal("Couldn't connect to database", err.Error())
	}

	return db
}

// StartDbConnection Creates connection to database
func StartDbConnection(cfg Config) {
	database = initDB(cfg)
}

// GetDBConn func for getting db conn globally
func GetDBConn() *sql.DB {
	return database
}

func CloseDbConnection() error {
	if err := GetDBConn().Close(); err != nil {
		fmt.Errorf("error occurred on database connection closing: %s", err.Error())
	}
	return nil
}
