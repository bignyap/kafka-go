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
	SQLDriver string
	URLConfig DBURLConfig
	// ConnectionURL string
	PoolConfig DBPoolConfig
}

type DBURLConfig struct {
	Username   string
	Password   string
	Database   string
	Properties string
}

type DBPoolConfig struct {
	MaxOpenConnections int
	MaxIdleConnections int
	MaxIdleTime        int
}

func (dbConfig *DBConfig) Connect() (*sql.DB, error) {

	driver := dbConfig.SQLDriver
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.URLConfig.Username,
		dbConfig.URLConfig.Database,
		dbConfig.URLConfig.Database,
	)

	// sql.Open("mysql", "user:password@/dbname")
	db, err := sql.Open(driver, connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(dbConfig.PoolConfig.MaxOpenConnections)
	db.SetMaxIdleConns(dbConfig.PoolConfig.MaxIdleConnections)
	db.SetConnMaxIdleTime(time.Duration(dbConfig.PoolConfig.MaxIdleTime))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil

}

func NewDBConn() (*sql.DB, error) {

	dbconfig := DBConfig{
		SQLDriver: utils.GetEnvString("SQL_DRIVER", "postgres"),
		URLConfig: DBURLConfig{
			Username: utils.GetEnvString("DB_USER", ""),
			Password: utils.GetEnvString("DB_PASSWORD", ""),
			Database: utils.GetEnvString("DB_NAME", ""),
		},
		PoolConfig: DBPoolConfig{
			MaxOpenConnections: utils.GetEnvInt("MAX_OPEN_CONNECTIONS", 10),
			MaxIdleConnections: utils.GetEnvInt("MAX_IDLE_CONNECTIONS", 10),
			MaxIdleTime:        utils.GetEnvInt("MAX_IDLE_TIME", 10),
		},
	}

	return dbconfig.Connect()
}

func WithTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
