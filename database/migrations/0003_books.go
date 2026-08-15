package sql

import migrate "github.com/rubenv/sql-migrate"

func MigrationBooks() *migrate.Migration {
	return &migrate.Migration{
		Id: "0003_books",
		Up: []string{
			`CREATE TABLE IF NOT EXISTS books (
				id SERIAL PRIMARY KEY,
				title VARCHAR(255) NOT NULL,
				description VARCHAR(255) NULL,
				image_url VARCHAR(500) NULL,
				release_year INT NULL,
				price INT NOT NULL DEFAULT 0,
				total_page INT NULL,
				thickness VARCHAR(50) NULL,
				category_id INT NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				created_by VARCHAR(100) NULL,
				modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				modified_by VARCHAR(100) NULL,
				CONSTRAINT fk_books_category FOREIGN KEY (category_id) REFERENCES categories (id) ON UPDATE CASCADE ON DELETE RESTRICT
			);`,
			`CREATE INDEX IF NOT EXISTS idx_books_category_id ON books(category_id);`,
		},
		Down: []string{
			`DROP TABLE IF EXISTS books CASCADE;`,
		},
	}
}
