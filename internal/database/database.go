package database

import (
	"github.com/MoonSHRD/sonis/internal/models"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

type Database struct {
	dbConnection *pg.DB
}

func New() (*Database, error) {
	var err error
	db := &Database{
		dbConnection: pg.Connect(&pg.Options{
			User:     "postgres",
			Password: "postgres",
			Database: "sonis",
			Addr:     "localhost:15432",
		}), // FIXME make options dynamic
	}
	err = db.initializeDbSchema()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *Database) GetDatabaseConnection() *pg.DB {
	return db.dbConnection
}

func (db *Database) CloseConnection() error {
	err := db.dbConnection.Close()
	if err != nil {
		return err
	}
	db.dbConnection = nil

	return nil
}

func (db *Database) initializeDbSchema() error {
	for _, model := range []interface{}{(*models.Room)(nil)} {
		err := db.dbConnection.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
