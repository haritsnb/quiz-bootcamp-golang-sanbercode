package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	DB *sql.DB
)

func ConnectDatabase() {
	// Membaca file .env jika ada
	_ = godotenv.Load()

	var connStr string

	// 1. Prioritas 1: Jika ada DATABASE_URL (Railway)
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		connStr = dbURL
	} else {
		// 2. Prioritas 2 & 3: Variabel terpisah
		host := getEnv("PGHOST", getEnv("DB_HOST", "localhost"))
		port := getEnv("PGPORT", getEnv("DB_PORT", "5432"))
		user := getEnv("PGUSER", getEnv("DB_USER", "postgres"))
		password := getEnv("PGPASSWORD", getEnv("DB_PASSWORD", "root"))
		dbname := getEnv("PGDATABASE", getEnv("DB_NAME", "local_golangquiz"))
		sslmode := getEnv("PGSSLMODE", getEnv("DB_SSLMODE", "disable"))

		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	}

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Gagal membuat koneksi database: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Gagal terhubung ke database: ", err)
	}

	fmt.Println("Berhasil terhubung ke database PostgreSQL!")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
