package controllers

import (
	"banco/dto"
	"banco/models"
	"banco/utils"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type LoanController struct {
	DB *gorm.DB
}

var percentageOfMoneyAvailableForLoan float64 = 0.02

func (l *LoanController) CreateLoan() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var input dto.CreateLoanDTO
		var loan models.Loan
		var owner models.Customer
		// var customer models.Customer
		var bankAccount models.BankAccount
		claims := ctx.Locals("customer").(jwt.MapClaims)
		IDfloat := claims["ID"].(float64)
		if err := ctx.BodyParser(&input); err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		if err := validate.Struct(input); err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		if input.MaturityDate.Before(time.Now()) {
			return ctx.Status(401).JSON(fiber.Map{"error": "the maturity date is older than the current date"})
		}

		if input.MaturityDate.Sub(time.Now()) < time.Hour * 24 * 30 {
			return ctx.Status(400).JSON(fiber.Map{"error": "the maturity date is so much close"})
		}

		interestRate := utils.SetInterestRate(input.MaturityDate.Sub(time.Now()))

		result := l.DB.Preload("BankAccount").First(&owner, "name = ?", "OWNER")

		if result.Error != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}


		if float64(input.Amount) > float64(owner.BankAccount.Balance)*percentageOfMoneyAvailableForLoan {
			return ctx.Status(500).JSON(fiber.Map{"error": "we can't lend that value"})
		}

		result = l.DB.First(&bankAccount, "customer_id = ?", uint(IDfloat))
	
		if result.Error != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}

		if bankAccount.Loan.ID != 0 {
			return ctx.Status(401).JSON(fiber.Map{"error": "you already have a loan"})
		}
		
		totalAmount := uint(float64(input.Amount) * (1 + (interestRate / 100)))

		if totalAmount > bankAccount.Balance {
			return ctx.Status(500).JSON(fiber.Map{"error": "sorry. You can't take a loan with an amount greater than your balance"})
		}

		err := l.DB.Transaction(func(tx *gorm.DB) error {
			loan = models.Loan{
				BaseAmount: input.Amount,
				InterestRate: interestRate,
				TotalAmount: totalAmount,
				MaturityDate: input.MaturityDate,
				BankAccountID: bankAccount.ID,
			}
	
			result = tx.Create(&loan)
			if result.Error != nil {
				return result.Error
			}

			result = tx.Model(&models.BankAccount{}).Where("customer_id = ?", bankAccount.CustomerID).Update("balance", gorm.Expr("balance + ?", loan.BaseAmount))
			if result.Error != nil {
				return result.Error
			}

			result = tx.Model(&models.BankAccount{}).Where("customer_id = ?", owner.ID).Update("balance", gorm.Expr("balance - ?", loan.BaseAmount))


			return nil
		})

		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		
	
		message := "your loan was created. You'll have to pay:" + " " + strconv.Itoa(int(totalAmount)) + " " + "interest rate:" + " "+ strconv.FormatFloat(interestRate, 'f', 2, 64)



		return ctx.Status(200).JSON(fiber.Map{"message": message})
	}
}

func (l *LoanController) PayLoan() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var customer models.Customer
		var owner models.Customer
		claims := ctx.Locals("customer").(jwt.MapClaims)
		IDfloat := claims["ID"].(float64)

		result := l.DB.Preload("BankAccount.Loan").First(&customer, uint(IDfloat))
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return ctx.Status(404).JSON(fiber.Map{"error": "your customer doesn't exist"})
			} else {
				return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
			}
		}

		if customer.BankAccount.Loan.ID == 0 {
			return ctx.Status(404).JSON(fiber.Map{"error": "you don't have any loan"})
		}

		l.DB.Preload("BankAccount").First(&owner, "name = ?", "OWNER")

		err := l.DB.Transaction(func(tx *gorm.DB) error {
			result = tx.Model(&models.BankAccount{}).Where("customer_id = ?", customer.ID).Update("balance", gorm.Expr("balance - ?", customer.BankAccount.Loan.TotalAmount))
			if result.Error != nil {
				return result.Error
			}
			result = tx.Model(&models.BankAccount{}).Where("customer_id", owner.ID).Update("balance", gorm.Expr("balance + ?", customer.BankAccount.Loan.TotalAmount))
			if result.Error != nil {
				return result.Error
			}

			result = tx.Model(&models.Loan{}).Delete(&customer.BankAccount.Loan)
			if result.Error != nil {
				return result.Error
			}

			return nil
		})

		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return ctx.Status(200).JSON(fiber.Map{"message": "your loan was paid!"})
	}
}
