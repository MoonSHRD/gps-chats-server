package migrations

const (
	_6Up = `
		ALTER TABLE rooms
			ADD COLUMN event_start_date TIMESTAMP WITH TIME ZONE NULL;
	`

	_6Down = `
		ALTER TABLE rooms
			DROP COLUMN event_start_date;
	`
)
