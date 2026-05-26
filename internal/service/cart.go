package service

import (
	"ecommerce/internal/model"
	"ecommerce/internal/repository"
	"errors"

	"gorm.io/gorm"
)

type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewCartService(cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *CartService {
	return &CartService{cartRepo: cartRepo, productRepo: productRepo}
}

func (s *CartService) GetCart(userID uint) ([]model.CartItem, error) {
	return s.cartRepo.FindByUser(userID)
}

func (s *CartService) AddItem(userID, productID uint, quantity int) error {
	if quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}

	_, err := s.productRepo.FindByID(productID)
	if err != nil {
		return errors.New("product not found")
	}

	existing, err := s.cartRepo.FindItem(userID, productID)

	if err == nil {
		existing.Quantity += quantity
		return s.cartRepo.Update(&existing)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		item := model.CartItem{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
		}
		return s.cartRepo.Create(&item)
	}

	return err
}

func (s *CartService) UpdateQuantity(userID, productID uint, quantity int) error {
	item, err := s.cartRepo.FindItem(userID, productID)
	if err != nil {
		return errors.New("item not in cart")
	}

	if quantity <= 0 {
		return s.cartRepo.Delete(userID, productID)
	}

	item.Quantity = quantity
	return s.cartRepo.Update(&item)
}

func (s *CartService) RemoveItem(userID, productID uint) error {
	_, err := s.cartRepo.FindItem(userID, productID)
	if err != nil {
		return errors.New("item not in cart")
	}
	return s.cartRepo.Delete(userID, productID)
}