package middlewares

import (
	"banco/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CheckIfItIsOwner() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims := ctx.Locals("customer").(jwt.MapClaims)
		role := claims["role"].(float64)
		if int(role) == utils.OWNER {
			return ctx.Status(401).JSON(fiber.Map{"error": "the owner account can't do that"})
		}
		return ctx.Next()
	}
}