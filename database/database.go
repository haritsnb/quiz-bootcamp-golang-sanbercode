package database

import (
	"database/sql"
	"fmt"
	"log"

	schema "app/database/migrations"
	seeder "app/database/seeders"

	migrate "github.com/rubenv/sql-migrate"
)

func RunMigrations(db *sql.DB) {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			schema.MigrationUsers(),
			schema.MigrationCategories(),
			schema.MigrationBooks(),
		},
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatalf("Gagal menjalankan migrasi: %v", err)
	}
	fmt.Printf("Migrasi berhasil dijalankan: %d migrasi diterapkan.\n", n)
}

func RunSeeders(db *sql.DB) {
	seeder.SeedUser(db)
}
