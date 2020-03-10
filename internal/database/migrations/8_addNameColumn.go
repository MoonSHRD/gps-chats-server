package migrations

const (
	_8Up = `
		ALTER TABLE rooms
			ADD COLUMN name TEXT DEFAULT '';
	`

	_8Down = `
		ALTER TABLE rooms
			DROP COLUMN name;
	`
)
