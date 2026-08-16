package main

import (
	"app/config"
	"app/controllers"
	"app/database"
	"app/docs"
	"app/middlewares"
	"app/repositories"
	"app/services"
	"os"

	_ "app/docs"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Book & Category Management REST API
// @version         1.0
// @description     Dokumentasi interaktif RESTful API Pengelolaan User, Kategori, dan Buku beserta Upload Cover Gambar.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Harits Nala Barrun
// @contact.email   developer.haritsnb@gmail.com

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath        /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer <TOKEN_JWT_ANDA>
func main() {
	/**
	 * Setup
	 */
	// Kosongkan aga Swagger otomatis mengikuti port yang sedang dibuka pada browser (8081 / 8080 / RAILWAY)
	docs.SwaggerInfo.Host = ""

	// Muat .env di awal
	_ = godotenv.Load()

	// DB, Migrasi, Seeder
	config.ConnectDatabase()
	database.RunMigrations(config.DB)
	database.RunSeeders(config.DB)

	/**
	 * Layer
	 */
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

	/**
	 * Router
	 */
	route := gin.Default()

	// Endpoint Swagger UI Docs
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Static route untuk file upload
	route.Static("/storages", "./storages")

	// Routing API
	api := route.Group("/api")
	{
		// Public Routes
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

	/**
	 * Running Project
	 */
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT")
		if port == "" {
			port = "8080"
		}
	}

	route.Run(":" + port)
}
