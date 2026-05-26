package handler

import (
	"ecommerce/internal/middleware"
	"ecommerce/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var body struct {
		Name string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error":"invalid body"})
	}

	if body.Name == "" || body.Email == "" || body.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "name, email and password are required"})
	}

	if err := h.service.Register(body.Name, body.Email, body.Password); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "email already in use"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "registered successfully"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var body struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid body"})
	}

	user, err := h.service.Login(body.Email, body.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := middleware.GenerateToken(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "could not generate token"})
	}

	return c.JSON(fiber.Map{"token": token})
}