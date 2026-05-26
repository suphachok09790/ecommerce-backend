package repository

import (
	"ecommerce/internal/model"
	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) FindByUser(userID uint) ([]model.CartItem, error) {
	var items []model.CartItem
	err := r.db.Preload("Product").Where("user_id = ?", userID).Find(&items).Error
	return items, err
}

func (r *CartRepository) FindItem(userID, productID uint) (model.CartItem, error) {
	var item model.CartItem
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&item).Error
	return item, err
}

func (r *CartRepository) Create(item *model.CartItem) error {
	return r.db.Create(item).Error
}

func (r *CartRepository) Update(item *model.CartItem) error {
	return r.db.Save(item).Error
}

func (r *CartRepository) Delete(userID, productID uint) error {
	return r.db.
	Where("user_id =? AND product_id = ?", userID, productID).
	Delete(&model.CartItem{}).Error
}

func (r *CartRepository) DeleteByUser(tx *gorm.DB, userID uint) error {
	return tx.
	Where("user_id = ?", userID).
	Delete(&model.CartItem{}).Error
}