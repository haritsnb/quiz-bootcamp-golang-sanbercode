package main

import (
	"app/config"
	"app/controllers"
	"app/database"
	"app/repositories"
	"app/services"
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

	// Repositories
	userRepo := repositories.NewUserRepository(config.DB)
	categoryRepo := repositories.NewCategoryRepository(config.DB)

	// Services
	userService := services.NewUserService(userRepo)
	categoryService := services.NewCategoryService(categoryRepo)

	// Controllers
	userController := controllers.NewUserController(userService)
	categoryController := controllers.NewCategoryController(categoryService)

	// Router
	route := gin.Default()

	// Routing API
	api := route.Group("/api")
	{

		protected := api.Group("/")
		{
			// Users
			protected.GET("/users", userController.GetAll)
			protected.POST("/users", userController.Create)
			protected.GET("/users/:id", userController.GetByID)
			protected.PUT("/users/:id", userController.Update)
			protected.DELETE("/users/:id", userController.Delete)

			// Categories
			protected.GET("/categories", categoryController.GetAll)
			protected.POST("/categories", categoryController.Create)
			protected.GET("/categories/:id", categoryController.GetByID)
			protected.PUT("/categories/:id", categoryController.Update)
			protected.DELETE("/categories/:id", categoryController.Delete)
		}
	}

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
