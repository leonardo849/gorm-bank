package controllers

import (
	"banco/dto"
	"banco/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type WithdrawController struct {
	DB *gorm.DB
}

func (w *WithdrawController) CreateWithdraw() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var input dto.CreateWithdrawDTO
		
		var bankAccount models.BankAccount
		claims := ctx.Locals("customer").(jwt.MapClaims)
		bankAccountID := claims["bank-account_id"].(float64)
		if err := ctx.BodyParser(&input); err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validate.Struct(input); err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		

		result := w.DB.Model(&models.BankAccount{}).Preload("Loan").First(&bankAccount, "id = ?", uint(bankAccountID))
		if result.Error != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}
		if bankAccount.Loan.ID != 0 && bankAccount.Balance - input.Amount < bankAccount.Loan.TotalAmount ||
		bankAccount.Balance < input.Amount {
			return ctx.Status(401).JSON(fiber.Map{"error": "you can't make a withdraw with that value"})
		}

		err := w.DB.Transaction(func(tx *gorm.DB) error {
			// var customer models.Customer
			// var owner models.Customer
			

			// tx.Model(&models.Customer{}).Preload("BankAccount").First(&owner, "name = ?", "OWNER")

			result = tx.Model(&models.BankAccount{}).Where("id = ?", 1).Update("balance", gorm.Expr("balance - ?", input.Amount))
			if result.Error != nil {
				return result.Error
			}

			result = tx.Model(&models.BankAccount{}).Where("id = ?", uint(bankAccountID)).Update("balance", gorm.Expr("balance - ?", input.Amount))
			if result.Error != nil {
				return result.Error
			}

			withdraw := models.Withdraw{
				BankAccountID: uint(bankAccountID),
				Amount: input.Amount,
			}

			result = tx.Create(&withdraw)
			if result.Error != nil {
				return result.Error
			}

			return nil
		})

		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
		}


		

		return ctx.Status(200).JSON(fiber.Map{"message": "withdraw was created!"})

	}
}