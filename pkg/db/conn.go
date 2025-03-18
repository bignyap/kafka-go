package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func NewDBConn() (*sql.DB, error) {

	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		utils.GetEnvString("DB_USER", ""),
		utils.GetEnvString("DB_PASSWORD", ""),
		utils.GetEnvString("DB_NAME", ""),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
