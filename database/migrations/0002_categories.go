package sql

import migrate "github.com/rubenv/sql-migrate"

func MigrationCategories() *migrate.Migration {
	return &migrate.Migration{
		Id: "0002_categories",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS categories (
				id SERIAL PRIMARY KEY,
				name VARCHAR(100) NOT NULL UNIQUE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				created_by VARCHAR(100) NULL,
				modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				modified_by VARCHAR(100) NULL
			);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS categories CASCADE;`,
		},
	}
}
