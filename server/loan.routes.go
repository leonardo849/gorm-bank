package server

import (
	"banco/controllers"
	"banco/database"
	"banco/middlewares"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupLoanRoutes(loanRoute fiber.Router) {
	loanController := controllers.LoanController{DB: database.DB}
	loanRoute.Post("/create", middlewares.CheckJWTMiddleware(), middlewares.CheckIfItIsOwner(), loanController.CreateLoan())
	loanRoute.Post("/pay", middlewares.CheckJWTMiddleware(), middlewares.CheckIfItIsOwner(), loanController.PayLoan())
	log.Println("loan's routes are working!")
}