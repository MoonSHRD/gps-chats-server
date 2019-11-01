package migrations

import (
	"database/sql"
)

const (
	sqlQuery = `
		ALTER TABLE rooms
			ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE,
			ALTER COLUMN created_at SET NOT NULL,
			ALTER COLUMN created_at SET DEFAULT now() AT TIME ZONE 'UTC';
	`
)

func Up002(tx *sql.Tx) error {
	_, err := tx.Exec(sqlQuery)
	return err
}
