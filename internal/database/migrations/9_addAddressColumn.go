package migrations

const (
	_9Up = `
	ALTER TABLE rooms
		ADD COLUMN address TEXT DEFAULT '';
	`

	_9Down = `
	ALTER TABLE rooms
		DROP COLUMN address;
	`
)
