package migrations

const (
	_7Up = `
		ALTER TABLE rooms
			RENAME COLUMN event_id TO parent_group_id;
	`

	_7Down = `
		ALTER TABLE rooms
			RENAME COLUMN parent_group_id TO event_id;
	`
)
