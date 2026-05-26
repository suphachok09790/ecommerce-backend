package handler

import (
	"ecommerce/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	service *service.CartService
}

func NewCartHandler(s *service.CartService) *CartHandler {
	return &CartHandler{service: s}
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	items, err := h.service.GetCart(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

func (h *CartHandler) AddItem(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var body struct {
		ProductID uint `json:"product_id"`
		Quantity  int `json:"quantity"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	if err := h.service.AddItem(userID, body.ProductID, body.Quantity); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "item added to cart"})
}

func (h *CartHandler) UpdateQuantity(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	productID, err := strconv.ParseUint(c.Params("product_id"), 10 , 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid product_id"})
	}

	var body struct {
		Quantity int `json:"quantity"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}
	
	if err := h.service.UpdateQuantity(userID, uint(productID), body.Quantity); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "cart updated"})
}

func (h *CartHandler) RemoveItem(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	productID, err := strconv.ParseUint(c.Params("product_id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid product_id"})
	}

	if err := h.service.RemoveItem(userID, uint(productID)); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "item removed"})
}