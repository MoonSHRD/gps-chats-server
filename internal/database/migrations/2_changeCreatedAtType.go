package migrations

const (
	_2Up = `
		ALTER TABLE rooms
			ALTER COLUMN created_at TYPE TIMESTAMP WITH TIME ZONE,
			ALTER COLUMN created_at SET NOT NULL,
			ALTER COLUMN created_at SET DEFAULT now() AT TIME ZONE 'UTC';
	`
	_2Down = `
		ALTER TABLE rooms
			ALTER COLUMN created_at TYPE TIMESTAMP,
			ALTER COLUMN created_at SET NOT NULL,
			ALTER COLUMN created_at SET DEFAULT now();
	`
)
