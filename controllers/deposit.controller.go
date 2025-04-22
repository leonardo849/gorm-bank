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

var reserve = 0.8 //80% goes to the bank 

func (d *DepositController) CreateDeposit() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var deposit models.Deposit
		var input dto.CreateDepositDTO
		var customer models.Customer
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
		err := d.DB.Transaction(func(tx *gorm.DB) error {
			
			deposit = models.Deposit{
				BankAccountID: customer.BankAccount.ID,
				Amount: input.Amount,
			}
			
			result = tx.Model(&models.BankAccount{}).Where("customer_id = ?", customer.ID).Update("balance", gorm.Expr("balance + ?", input.Amount))
			if result.Error != nil {
				return result.Error
			}
			result = tx.Model(&models.BankAccount{}).Where("id = ?", 1).Update("balance", gorm.Expr("balance + ?", amountToTheBank))
			if result.Error != nil {
				return result.Error
			}
			result = tx.Create(&deposit)
			if result.Error != nil {
				return result.Error
			}
			return nil
		})
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		
		return ctx.Status(200).JSON(fiber.Map{"message": "deposit was created"})
	}
}

