package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

type Migration struct {
	dbConn *sql.DB
}

func New(dbConn *sql.DB) *Migration {
	return &Migration{
		dbConn: dbConn,
	}
}

func (m *Migration) Migrate() error {
	goose.SetDialect("postgres")
	goose.AddNamedMigration("1_initial.go", InitialUp, nil)

	if m.dbConn != nil {
		err := goose.Up(m.dbConn, ".")
		if err != nil {
			return err
		}
	}

	return nil
}
