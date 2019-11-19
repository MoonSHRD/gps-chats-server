package migrations

import (
	migrate "github.com/rubenv/sql-migrate"
)

var (
	MigrationsList = &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			&migrate.Migration{
				Id:   "1",
				Up:   []string{_1Up},
				Down: []string{_1Down},
			},
			&migrate.Migration{
				Id:   "2",
				Up:   []string{_2Up},
				Down: []string{_2Down},
			},
		},
	}
)
