package middlewares

import (
	"banco/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thoas/go-funk"
)

func CheckRole() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims := ctx.Locals("customer").(jwt.MapClaims)
		role := claims["role"].(float64)
		exists := funk.Contains([]float64{utils.OWNER, utils.MANAGER}, role)
		if exists {
			return ctx.Next()
		}
		return ctx.Status(401).JSON(fiber.Map{"message": "you are a customer"})
	}
}