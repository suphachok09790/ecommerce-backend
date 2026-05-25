package middleware

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Claims is the data we embed inside the token
type Claims struct {
	UserID uint `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken — called after successful login
func GenerateToken(userID uint, role string) (string, error) {
	claims := Claims{
		UserID:  userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// ParseToken — called inside middleware to verify the token
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token)(interface{}, error){
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// RequireAuth — runs before any protected handler
func RequireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "missing token"})
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{"error": "invalid token format"})
	}

	claims, err := ParseToken(parts[1])
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "invalid or expired token"})
	}

	// store in context so handlers can read it
	c.Locals("user_id", claims.UserID)
	c.Locals("role", claims.Role)

	return c.Next()
}

// RequireAdmin — always stacked after RequireAuth
func RequireAdmin(c *fiber.Ctx) error {
	role, ok := c.Locals("role").(string)
	if !ok || role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "admin access required"})
	}
	return c.Next()
}


