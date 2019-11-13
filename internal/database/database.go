package database

import (
	"fmt"

	"github.com/MoonSHRD/sonis/internal/database/migrations"
	"github.com/MoonSHRD/sonis/internal/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // for sqlx
	"github.com/sirupsen/logrus"
)

type Database struct {
	dbConnection *sqlx.DB
}

func New(cfg utils.Config) (*Database, error) {
	var err error
	logger := logrus.New()
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s", cfg.DatabaseUser, cfg.DatabasePassword, cfg.DatabaseName, cfg.DatabaseHost, cfg.DatabasePort, "disable")
	dbConnection, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	db := &Database{
		dbConnection: dbConnection,
	}

	err = migrations.New(db.dbConnection.DB).Migrate()
	if err != nil {
		logger.Errorf("Failed to process the migration. Reason: %s", err.Error())
		return nil, err
	}
	return db, nil
}

func (db *Database) GetDatabaseConnection() *sqlx.DB {
	return db.dbConnection
}

func (db *Database) CloseConnection() {
	db.dbConnection.Close()
	db.dbConnection = nil
}
