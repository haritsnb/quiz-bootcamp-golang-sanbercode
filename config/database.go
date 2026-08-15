package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	DB  *sql.DB
	err error
)

func ConnectDatabase() {
	host := getEnv("PGHOST", "localhost")
	port := getEnv("PGPORT", "5432")
	user := getEnv("PGUSER", "postgres")
	password := getEnv("PGPASSWORD", "root")
	dbname := getEnv("PGNAME", "local_golangquiz")
	sslmode := getEnv("PGSSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal membuat koneksi database:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	fmt.Println("Berhasil terhubung ke database PostgreSQL!")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
