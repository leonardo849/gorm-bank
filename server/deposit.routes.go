package server

import (
	"banco/controllers"
	"banco/database"
	"banco/middlewares"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupDepositRoutes(depositRoute fiber.Router) {
	depositController := controllers.DepositController{DB: database.DB}
	depositRoute.Post("/create", middlewares.CheckJWTMiddleware(), middlewares.CheckIfItIsOwner() , depositController.CreateDeposit())
	log.Println("deposit's routes are working!")
}