package database

import (
	"database/sql"

	"github.com/MoonSHRD/sonis/internal/database/migrations"
	migrate "github.com/rubenv/sql-migrate"
)

type Migration struct {
	dbConn *sql.DB
}

func NewMigrations(dbConn *sql.DB) *Migration {
	return &Migration{
		dbConn: dbConn,
	}
}

func (m *Migration) Migrate() error {
	_, err := migrate.Exec(m.dbConn, "postgres", migrations.MigrationsList, migrate.Up)

	return err
}
