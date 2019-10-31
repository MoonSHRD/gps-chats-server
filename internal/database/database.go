package database

import (
	"fmt"

	"github.com/MoonSHRD/sonis/internal/database/migrations"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // for sqlx
	"github.com/sirupsen/logrus"
)

type Database struct {
	dbConnection *sqlx.DB
}

func New() (*Database, error) {
	var err error
	logger := logrus.New()
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s", "postgres", "postgres", "sonis", "localhost", 15432, "disable") // FIXME make options dynamic
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
