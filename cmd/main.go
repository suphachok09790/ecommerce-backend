package main

import (
	"ecommerce/config"
	"ecommerce/internal/handler"
	"ecommerce/internal/model"
	"ecommerce/internal/repository"
	"ecommerce/internal/service"
	"ecommerce/internal/middleware"
	
	
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

	// wire auth layers: repository → service → handler
	userRepo := repository.NewUserRepository(config.DB)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	// wire products
	productRepo := repository.NewProductRepository(config.DB)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	// cart
	cartRepo := repository.NewCartRepository(config.DB)
	cartService := service.NewCartService(cartRepo, productRepo)
	cartHandler := handler.NewCartHandler(cartService)

	// orders
	orderRepo := repository.NewOrderRepository(config.DB)
	orderService := service.NewOrderService(config.DB, orderRepo, cartRepo, productRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	// create fiber app
	app := fiber.New()

	// auth routes — public, no middleware needed
	auth := app.Group("/api/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// product routes — public read
	public := app.Group("/api")
	public.Get("/products", productHandler.GetAll)
	public.Get("/products/:id", productHandler.GetByID)

	// protected routes — login required
	api := app.Group("/api", middleware.RequireAuth)
	api.Get("/cart", cartHandler.GetCart)
	api.Post("/cart", cartHandler.AddItem)
	api.Put("/cart/:product_id", cartHandler.UpdateQuantity)
	api.Delete("/cart/:product_id", cartHandler.RemoveItem)
	api.Post("/orders", orderHandler.CreateOrder)
	api.Get("/orders", orderHandler.GetOrders)

	// admin routes — RequireAuth + RequireAdmin both must pass
	admin := app.Group("/api/admin", middleware.RequireAuth, middleware.RequireAdmin)
	admin.Post("/products",       productHandler.Create)
	admin.Put("/products/:id",    productHandler.Update)
	admin.Delete("/products/:id", productHandler.Delete)

	// start server
	log.Println("Server running on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

