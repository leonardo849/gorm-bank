package server

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func RunServer() error {
	app := fiber.New()
	port := os.Getenv("PORT")
	if port == "" {
		return fmt.Errorf("there isn't port")
	}

	app.Get("/", func (c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"message": "hi!"})
	})
	customerGroup := app.Group("/customer")
	depositGroup := app.Group("/deposit")
	SetupCustomerRoutes(customerGroup)
	SetupDepositRoutes(depositGroup)
	return app.Listen(port)
}