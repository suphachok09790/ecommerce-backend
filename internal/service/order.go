package service

import (
	"ecommerce/internal/model"
	"ecommerce/internal/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type OrderService struct {
	db          *gorm.DB
	orderRepo   *repository.OrderRepository
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewOrderService(
	db *gorm.DB,
	orderRepo *repository.OrderRepository,
	cartRepo *repository.CartRepository,
	productRepo *repository.ProductRepository,
) *OrderService {
	return &OrderService{
		db:          db,
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (s *OrderService) CreateOrder(userID uint) (model.Order, error) {
	// load cart BEFORE transaction — read only, no need to lock yet
	cartItems, err := s.cartRepo.FindByUser(userID)
	if err != nil || len(cartItems) == 0 {
		return model.Order{}, errors.New("cart is empty")
	}

	var order model.Order

	// s.db.Transaction opens a transaction
	// return nil  → COMMIT   (all steps saved)
	// return error → ROLLBACK (all steps undone)
	err = s.db.Transaction(func(tx *gorm.DB) error {

		var orderItems []model.OrderItem
		var total float64

		// step 1 — check stock and build order items
		for _, cartItem := range cartItems {
			var product model.Product

			// read product INSIDE transaction for fresh data
			if err := tx.First(&product, cartItem.ProductID).Error; err != nil {
				return fmt.Errorf("product %d not found", cartItem.ProductID)
			}

			// not enough stock → return error → ROLLBACK
			if product.Stock < cartItem.Quantity {
				return fmt.Errorf("insufficient stock for: %s", product.Name)
			}

			total += product.Price * float64(cartItem.Quantity)

			orderItems = append(orderItems, model.OrderItem{
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
				Price:     product.Price, // snapshot price at checkout
			})
		}

		// step 2 — create order + order items
		order = model.Order{
			UserID: userID,
			Total:  total,
			Status: "pending",
		}

		if err := s.orderRepo.CreateWithItems(tx, &order, orderItems); err != nil {
			return err
		}

		// step 3 — deduct stock
		// gorm.Expr does the math IN the database — not in Go
		// this prevents race conditions when two users checkout at the same time
		for _, cartItem := range cartItems {
			result := tx.Model(&model.Product{}).
				Where("id = ?", cartItem.ProductID).
				Update("stock", gorm.Expr("stock - ?", cartItem.Quantity))

			if result.Error != nil {
				return result.Error
			}
		}

		// step 4 — clear the cart
		if err := s.cartRepo.DeleteByUser(tx, userID); err != nil {
			return err
		}

		// return nil = COMMIT
		return nil
	})

	if err != nil {
		return model.Order{}, err
	}
	// reload order from DB with OrderItems populated
	s.db.Preload("OrderItems").First(&order, order.ID)

	return order, nil
}

func (s *OrderService) GetOrders(userID uint) ([]model.Order, error) {
	return s.orderRepo.FindByUser(userID)
}