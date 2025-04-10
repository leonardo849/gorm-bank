package controllers

import (
	"banco/dto"
	"banco/models"
	"errors"
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type DepositController struct {
	DB *gorm.DB
}

var reserve = 0.8 //80% goes to the bank as well

func (d *DepositController) CreateDeposit() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var deposit models.Deposit
		var input dto.CreateDepositDTO
		var customer models.Customer
		var owner models.Customer
		claims := ctx.Locals("customer").(jwt.MapClaims)
		IDfloat := claims["ID"].(float64)
		if err := ctx.BodyParser(&input); err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validate.Struct(input); err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": "amount needs to be greater than 0.9"})
		}
		amountToTheBank := int(math.Round(reserve * float64(input.Amount)))
		IDuint := uint(IDfloat)
		result := d.DB.Preload("BankAccount").First(&customer, IDuint)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return ctx.Status(404).JSON(fiber.Map{"error": "your user doesn't exist"})
			} else {
				return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
			}
		}
		d.DB.Model(&models.Customer{}).Preload("BankAccount").First(&owner, "name = ?", "OWNER")
		deposit = models.Deposit{
			CustomerID: IDuint,
			Amount: input.Amount,
		}
		result = d.DB.Create(&deposit)
		if result.Error != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}
		result = d.DB.Model(&models.BankAccount{}).Where("customer_id = ?", customer.ID).Update("balance", int(input.Amount) + int(customer.BankAccount.Balance))
		if result.Error != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}
		result = d.DB.Model(&models.BankAccount{}).Where("customer_id = ?", owner.ID).Update("balance", int(owner.BankAccount.Balance) + int(amountToTheBank))
		if result.Error != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}
		return ctx.Status(200).JSON(fiber.Map{"message": "deposit was created"})
	}
}

