package main

import (
	"app/config"
	"app/controllers"
	"app/database"
	"app/middlewares"
	"app/repositories"
	"app/services"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Muat .env di awal
	_ = godotenv.Load()

	// DB, Migrasi, Seeder
	config.ConnectDatabase()
	database.RunMigrations(config.DB)
	database.RunSeeders(config.DB)

	// Repositories
	userRepo := repositories.NewUserRepository(config.DB)
	categoryRepo := repositories.NewCategoryRepository(config.DB)
	bookRepo := repositories.NewBookRepository(config.DB)

	// Services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	bookService := services.NewBookService(bookRepo, categoryRepo)

	// Controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)
	categoryController := controllers.NewCategoryController(categoryService)
	bookController := controllers.NewBookController(bookService)

	// Router
	route := gin.Default()

	// Static route untuk file upload
	route.Static("/storages", "./storages")

	// Routing API
	api := route.Group("/api")
	{

		// Public Auth
		api.POST("/users/login", authController.Login)

		// Protected Routes (JWT Middleware)
		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware())
		{
			// Logout
			protected.POST("/users/logout", authController.Logout)

			// Users
			protected.GET("/users", userController.GetAll)
			protected.POST("/users", userController.Create)
			protected.GET("/users/:id", userController.GetByID)
			protected.PUT("/users/:id", userController.Update)
			protected.DELETE("/users/:id", userController.Delete)
			protected.GET("/categories/:id/books", categoryController.GetBooks)

			// Categories
			protected.GET("/categories", categoryController.GetAll)
			protected.POST("/categories", categoryController.Create)
			protected.GET("/categories/:id", categoryController.GetByID)
			protected.PUT("/categories/:id", categoryController.Update)
			protected.DELETE("/categories/:id", categoryController.Delete)

			// Book
			protected.GET("/books", bookController.GetAll)
			protected.POST("/books", bookController.Create)
			protected.GET("/books/:id", bookController.GetByID)
			protected.PUT("/books/:id", bookController.Update)
			protected.DELETE("/books/:id", bookController.Delete)
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
