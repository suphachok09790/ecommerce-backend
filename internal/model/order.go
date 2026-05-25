package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID uint `json:"user_id"`
	Status string `json:"status" gorm:"default:pending"`
	Total float64 `json:"total"`
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID uint `json:"id" gorm:"primaryKey"`
	OrderID uint `json:"order_id"`
	ProductID uint `json:"product_id"`
	Quantity int `json:"quantity"`
	Price float64 `json:"price"`
}