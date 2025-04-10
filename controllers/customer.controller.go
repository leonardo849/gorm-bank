package controllers

import (
	"banco/dto"
	"banco/functionscrypto"
	"banco/models"
	"banco/utils"
	"errors"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

var validate = validator.New()

type CustomerController struct {
	DB            *gorm.DB
}

func (c *CustomerController) FindAllCustomers() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var customers []models.Customer
		var array []dto.FindCustomerDTO
		result := c.DB.Preload("BankAccount").Find(&customers)
		if result.Error != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}
		array = funk.Map(customers, func(customer models.Customer) dto.FindCustomerDTO {
			return dto.FindCustomerDTO{ID: customer.ID, Name: customer.Name, CreatedAt: customer.CreatedAt, UpdatedAt: customer.UpdatedAt, RoleUpdatedAt: customer.RoleUpdatedAt,
				BankAccount: dto.FindBankAccountDTO{
					ID:         customer.BankAccount.ID,
					CustomerID: customer.BankAccount.CustomerID,
					Balance:    customer.BankAccount.Balance,
				}}
		}).([]dto.FindCustomerDTO)

		return ctx.Status(200).JSON(array)
	}
}

func (c *CustomerController) CreateCustomer() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var input dto.CreateCustomerDTO
		if err := ctx.BodyParser(&input); err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validate.Struct(input); err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		customer := models.Customer{
			Name:          input.Name,
			Password:      input.Password,
			Role:          utils.CUSTOMER,
			RoleUpdatedAt: time.Now(),
		}

		result := c.DB.Create(&customer)
		if result.Error != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}

		bankAccount := models.BankAccount{
			CustomerID: customer.ID,
			Balance:    0,
		}
		result = c.DB.Create(&bankAccount)
		if result.Error != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}

		return ctx.Status(200).JSON(fiber.Map{"message": "customer was created"})
	}
}

func (c *CustomerController) LoginCustomer() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var input dto.LoginCustomerDTO
		var customer models.Customer
		if err := ctx.BodyParser(&input); err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validate.Struct(input); err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		result := c.DB.First(&customer, "name = ?", input.Name)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return ctx.Status(404).JSON(fiber.Map{"error": "that customer doesn't exist"})
			} else {
				return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
			}
		}
		comparePassword := functionscrypto.CompareHash(input.Password, customer.Password)
		if !comparePassword {
			return ctx.Status(401).JSON(fiber.Map{"error": "password is wrong"})
		}
		jwtToken, err := functionscrypto.GenerateJWT(customer.ID, utils.Role(customer.Role), customer.RoleUpdatedAt)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(200).JSON(fiber.Map{"token": jwtToken})
	}
}

func (c *CustomerController) FindOneCustomer() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims := ctx.Locals("customer").(jwt.MapClaims)
		role := claims["role"].(float64)
		var searchedID int
		var err error
		var customer models.Customer
		var customerDTO dto.FindCustomerDTO
		if int(role) == utils.CUSTOMER {
			idFloat := claims["ID"].(float64)
			searchedID = int(idFloat)
		} else {
			IdParam := ctx.Params("id")
			if IdParam == "" {
				return ctx.Status(400).JSON(fiber.Map{"error": "there isn't ID"})
			}
			searchedID, err = strconv.Atoi(IdParam)
			if err != nil {
				return ctx.Status(400).JSON(fiber.Map{"error": "the ID isn't a number"})
			}
		}
		result := c.DB.Preload("BankAccount").First(&customer, searchedID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return ctx.Status(404).JSON(fiber.Map{"error": "the customer doesn't exist"})
			} else {
				return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
			}
		}
		customerDTO = dto.FindCustomerDTO{
			ID: customer.ID,
			Name: customer.Name,
			CreatedAt: customer.CreatedAt,
			UpdatedAt: customer.UpdatedAt,
			RoleUpdatedAt: customer.RoleUpdatedAt,
			BankAccount: dto.FindBankAccountDTO{
				ID: customer.BankAccount.ID,
				Balance: customer.BankAccount.Balance,
				CustomerID: customer.BankAccount.CustomerID,
			},
		}
		return ctx.Status(200).JSON(customerDTO)
	}
}

