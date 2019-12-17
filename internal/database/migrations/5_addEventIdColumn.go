package migrations

const (
	_5Up = `
		ALTER TABLE rooms
			ADD COLUMN event_id TEXT DEFAULT '';
	`

	_5Down = `
		ALTER TABLE rooms
			DROP COLUMN event_id;
	`
)
