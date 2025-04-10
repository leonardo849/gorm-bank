package middlewares

import (
	"banco/database"
	"banco/models"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)


func CheckJWTMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token is missing",
			})
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token format",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signature")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid claims",
			})
		}
		customerIdFloat, okID := claims["ID"].(float64)
		roleUpdatedAt, okRoleUpdatedAt := claims["role_updated_at"].(float64)
		if !okID || !okRoleUpdatedAt {
			return ctx.Status(401).JSON(fiber.Map{"error": "invalid token"})
		} 
		customerID := uint(customerIdFloat)
		tokenUpdatedAt := int64(roleUpdatedAt)
		
		var customer models.Customer
		if err := database.DB.First(&customer, customerID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ctx.Status(404).JSON(fiber.Map{"error": "your user doesn't exist"})
			} else {
				return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		}

		dbUpdatedAt := customer.RoleUpdatedAt.Unix()

		if tokenUpdatedAt < dbUpdatedAt {
			return ctx.Status(401).JSON(fiber.Map{"error": "please login again"})
		}

		ctx.Locals("customer", claims)
		return ctx.Next()
	}
}