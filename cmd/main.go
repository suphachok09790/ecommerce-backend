package main

import (
	"ecommerce/config"
	"ecommerce/internal/model"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	config.DB.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.CartItem{},
		&model.Order{},
		&model.OrderItem{},
	)

	log.Println("Tables migrated successfully")

	log.Println("Server is ready")
}

