package db

import (
	"copySys/pkg/logger"
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

func initDB(cfg Config) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Error.Println(err)
		log.Fatal("Couldn't connect to database", err.Error())
	}

	return db
}

func StartDbConnection(cfg Config) {
	database = initDB(cfg)
}

func GetDBConn() *sql.DB {
	return database
}

func CloseDbConnection() error {
	if err := GetDBConn().Close(); err != nil {
		fmt.Errorf("error occurred on database connection closing: %s", err.Error())
		logger.Error.Println(err)
	}
	return nil
}
