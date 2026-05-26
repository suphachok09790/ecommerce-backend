package main

import (
	"ecommerce/config"
	"ecommerce/internal/handler"
	"ecommerce/internal/model"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	
	
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env 
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	// migrate tables
	config.DB.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.CartItem{},
		&model.Order{},
		&model.OrderItem{},
	)

	// wire layers: repository → service → handler
	userRepo := repository.NewUserRepository(config.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// create fiber app
	app := fiber.New()

	// auth routes — public, no middleware needed
	auth := app.Group("/api/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// start server
	log.Println("Server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

