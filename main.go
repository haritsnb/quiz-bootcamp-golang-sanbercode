package main

import (
	"app/config"
	"app/database"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env di awal
	_ = godotenv.Load()

	// DB, Migrasi, Seeder
	config.ConnectDatabase()
	database.RunMigrations(config.DB)
	database.RunSeeders(config.DB)

	// Router
	route := gin.Default()

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT")
		if port == "" {
			port = "8080"
		}
	}

	route.Run(":" + port)
}
