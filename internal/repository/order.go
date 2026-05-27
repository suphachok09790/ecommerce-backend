package repository

import (
	"ecommerce/internal/model"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) FindByUser(userID uint) ([]model.Order, error) {
	var orders []model.Order
	err := r.db.Preload("OrderItems").
		Where("user_id = ?", userID).
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) CreateWithItems(
	tx *gorm.DB,
	order *model.Order,
	items []model.OrderItem,
) error {
	// insert order row first — we need order.ID for the items
	if err := tx.Create(order).Error; err != nil {
		return err
	}

	// attach order ID to every item
	for i := range items {
		items[i].OrderID = order.ID
	}

	// insert all items in one query
	return tx.Create(&items).Error
}