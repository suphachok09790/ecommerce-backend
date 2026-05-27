package handler

import (
	"ecommerce/internal/service"
	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	order, err := h.service.CreateOrder(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(order)
}

func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	orders, err := h.service.GetOrders(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(orders)
}