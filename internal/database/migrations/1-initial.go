package migrations

import "database/sql"

const (
	createRoomsTable = `
		CREATE TABLE rooms
		(
			id SERIAL
				CONSTRAINT rooms_pk
					PRIMARY KEY,
			latitude REAL NOT NULL,
			longitude REAL NOT NULL,
			ttl INT DEFAULT 1 NOT NULL,
			room_id TEXT NOT NULL,
			category TEXT,
			created_at TIMESTAMP DEFAULT now() NOT NULL
		);
	`
)

func InitialUp(tx *sql.Tx) error {
	_, err := tx.Exec(createRoomsTable)
	return err
}
