package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bignyap/kafka-go/pkg/utils"
	_ "github.com/lib/pq"
)

type DBObject interface {
	connect() (*sql.DB, error)
}

type DBConfig struct {
	ConnectionURL string
	SQLDriver     string
}

// func NewDBConfig() DBConfig {
// 	return DBConfig{
// 		SQLDriver:     "mysql",
// 		ConnectionURL: "user:password@/dbname",
// 	}
// }

func (dbConfig *DBConfig) Connect() (*sql.DB, error) {

	// sql.Open("mysql", "user:password@/dbname")
	return sql.Open(dbConfig.SQLDriver, dbConfig.ConnectionURL)

}

func NewDBConn() *sql.DB {

	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		utils.GetEnvString("DB_USER", ""),
		utils.GetEnvString("DB_PASSWORD", ""),
		utils.GetEnvString("DB_NAME", ""),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening connection to the database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	return db
}
