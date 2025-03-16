package db

import "database/sql"

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
