package server

import (
	"banco/controllers"
	"banco/database"
	"banco/middlewares"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupWithdrawRoutes(withdrawRoute fiber.Router) {
	withdrawController := controllers.WithdrawController{DB: database.DB}
	withdrawRoute.Post("/create", middlewares.CheckJWTMiddleware(), middlewares.CheckIfItIsOwner(), withdrawController.CreateWithdraw())
	log.Println("withdraw's routes are working!")
}