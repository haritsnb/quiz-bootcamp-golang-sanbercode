package sql

import migrate "github.com/rubenv/sql-migrate"

func MigrationUsers() *migrate.Migration {
	return &migrate.Migration{
		Id: "0001_users",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS users (
				id SERIAL PRIMARY KEY,
				username VARCHAR(100) NOT NULL UNIQUE,
				password VARCHAR(255) NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				created_by VARCHAR(100) NULL,
				modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				modified_by VARCHAR(100) NULL
			);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS users CASCADE;`,
		},
	}
}
