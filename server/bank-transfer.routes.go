package server

import (
	"banco/controllers"
	"banco/database"
	"banco/middlewares"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupBankTransferRoutes(BankTransferRoute fiber.Router) {
	bankController := controllers.BankTransferController{DB: database.DB}
	BankTransferRoute.Post("/create", middlewares.CheckJWTMiddleware(),middlewares.CheckIfItIsOwner(), bankController.CreateTransfer())
	BankTransferRoute.Get("/all", middlewares.CheckJWTMiddleware(), middlewares.CheckRole(), bankController.FindAllTransfers())
	log.Println("bank transfer's routes are working")
}