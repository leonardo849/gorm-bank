package server

import (
	"banco/database"
	"banco/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RunServer() error {
	app := fiber.New()
	port := os.Getenv("PORT")
	if port == "" {
		return fmt.Errorf("there isn't port")
	}
	startLoanChecker(database.DB)
	app.Get("/", func (c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"message": "hi!"})
	})
	customerGroup := app.Group("/customer")
	depositGroup := app.Group("/deposit")
	bankTransferGroup := app.Group("/bank-transfer")
	loanGroup := app.Group("/loan")
	withdrawGroup := app.Group("/withdraw")
	SetupCustomerRoutes(customerGroup)
	SetupDepositRoutes(depositGroup)
	SetupBankTransferRoutes(bankTransferGroup)
	SetupLoanRoutes(loanGroup)
	SetupWithdrawRoutes(withdrawGroup)
	
	return app.Listen(port)
}




func startLoanChecker(db *gorm.DB) {
	ticker := time.NewTicker(24 * time.Hour)

	go func() {
		for {
			<-ticker.C
			checkOverdueLoans(db)
		}
	}()
}

func checkOverdueLoans(db *gorm.DB) {
	var overdueLoans []models.Loan
	result := db.Where("maturity_date < ?", time.Now()).Find(&overdueLoans)
	if result.Error != nil {
		log.Println(result.Error.Error())
		return
	}

	if len(overdueLoans) > 0 {
		var customer models.Customer
		var bankAccount models.BankAccount
		
		log.Println("there are overdue loans")
		for _, loan := range overdueLoans {
			db.Model(&models.BankAccount{}).Preload("Loan").First(&bankAccount, loan.BankAccountID)
			db.Model(&models.Customer{}).First(&customer, bankAccount.CustomerID)
			err := db.Transaction(func(tx *gorm.DB) error {
				result := tx.Model(&models.BankAccount{}).Where("id = ?", bankAccount.ID).Update("balance", gorm.Expr("balance - ?", loan.TotalAmount))
				if result.Error != nil {
					return result.Error
				}
				result = tx.Model(&models.BankAccount{}).Where("customer_id = ?", 1).Update("balance", gorm.Expr("balance + ?", loan.TotalAmount))
				if result.Error != nil {
					return result.Error
				}

				result = tx.Model(&models.Loan{}).Delete(&loan)
				if result.Error != nil {
					return result.Error
				}

				return nil
			})

			if err != nil {
				fmt.Println("error:", err.Error())
			}
		}
	}
}