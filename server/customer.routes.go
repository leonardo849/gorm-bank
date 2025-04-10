package server

import (
	"banco/controllers"
	"banco/database"
	"banco/middlewares"
	"log"

	"github.com/gofiber/fiber/v2"
)

func SetupCustomerRoutes(customerRoute fiber.Router) {
	customerController := controllers.CustomerController{DB: database.DB}
	customerRoute.Get("/all", middlewares.CheckJWTMiddleware(), middlewares.CheckRole(), customerController.FindAllCustomers())
	customerRoute.Post("/login", customerController.LoginCustomer())
	customerRoute.Post("/create", customerController.CreateCustomer())
	customerRoute.Get("/one/:id", middlewares.CheckJWTMiddleware(), customerController.FindOneCustomer())
	customerRoute.Patch("/update/:role/:id", middlewares.CheckJWTMiddleware(), middlewares.CheckRole(), customerController.ChangeCustomerRole())
	log.Println("Customer's routes are working!")
}