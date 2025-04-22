package controllers

import (
	"banco/dto"
	"banco/functionscrypto"
	"banco/models"
	"banco/utils"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

var validate = validator.New()

type CustomerController struct {
	DB *gorm.DB
}

func (c *CustomerController) FindAllCustomers() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		page := ctx.Params("page")
		pageSize := ctx.Params("page_size")
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		pageSizeInt, err := strconv.Atoi(pageSize)
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		offset := (pageInt - 1) * pageSizeInt
		var customers []models.Customer
		var array []dto.FindCustomerDTO
		result := c.DB.Preload("BankAccount").Limit(pageSizeInt).Offset(offset).Find(&customers)
		if result.Error != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}
		array = funk.Map(customers, func(customer models.Customer) dto.FindCustomerDTO {
			return dto.FindCustomerDTO{
				ID: customer.ID, Name: customer.Name, CreatedAt: customer.CreatedAt, UpdatedAt: customer.UpdatedAt, RoleUpdatedAt: customer.RoleUpdatedAt,
				Role: customer.Role,
				BankAccount: dto.FindBankAccountDTO{
					ID:         customer.BankAccount.ID,
					CustomerID: customer.BankAccount.CustomerID,
					Balance:    customer.BankAccount.Balance,
					SentTransfers: nil,
					ReceivedTransfers: nil,
					Loan: nil,
					Deposits: nil,
					Withdrawals: nil,
				},
			}
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
		if strings.ToUpper(input.Name) == "OWNER" {
			return ctx.Status(401).JSON(fiber.Map{"error": "you can't create a user with that name"})
		}
		err := c.DB.Transaction(func(tx *gorm.DB) error {
			customer := models.Customer{
				Name:          input.Name,
				Password:      input.Password,
				Role:          utils.CUSTOMER,
				RoleUpdatedAt: time.Now(),
			}

			result := tx.Create(&customer)
			if result.Error != nil {
				return result.Error
			}

			bankAccount := models.BankAccount{
				CustomerID: customer.ID,
				Balance:    0,
			}
			result = tx.Create(&bankAccount)
			if result.Error != nil {
				return result.Error
			}
			return nil
		})
		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
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

		result := c.DB.Preload("BankAccount").First(&customer, "name = ?", input.Name)
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
		jwtToken, err := functionscrypto.GenerateJWT(customer.ID, customer.Role, customer.RoleUpdatedAt, customer.BankAccount.ID)
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
		result := c.DB.Preload("BankAccount.SentTransfers").Preload("BankAccount.ReceivedTransfers").Preload("BankAccount.Deposits").Preload("BankAccount.Loan").Preload("BankAccount.Withdrawals").First(&customer, searchedID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return ctx.Status(404).JSON(fiber.Map{"error": "the customer doesn't exist"})
			} else {
				return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
			}
		}

		depositsDTO := funk.Map(customer.BankAccount.Deposits, func(deposit models.Deposit) dto.FindDepositDTO {
			return dto.FindDepositDTO{
				ID:         deposit.ID,
				Amount:     deposit.Amount,
				BankAccountID: deposit.BankAccountID,
				CreatedAt:  deposit.CreatedAt,
			}
		}).([]dto.FindDepositDTO)

		sentTransfersDTO := funk.Map(customer.BankAccount.SentTransfers, func(sentTransfer models.BankTransfer) dto.FindBankTransferDTO {
			return dto.FindBankTransferDTO{
				ID: sentTransfer.ID,
				Amount: sentTransfer.Amount,
				ReceiverBankAccountID: sentTransfer.ReceiverBankAccountID,
				SenderBankAccountID: sentTransfer.SenderBankAccountID,
				CreatedAt: sentTransfer.CreatedAt,
			}
		}).([]dto.FindBankTransferDTO)

		receivedTransfersDTO := funk.Map(customer.BankAccount.ReceivedTransfers, func(receivedTransfer models.BankTransfer) dto.FindBankTransferDTO {
			return dto.FindBankTransferDTO{
				ID: receivedTransfer.ID,
				Amount: receivedTransfer.Amount,
				ReceiverBankAccountID: receivedTransfer.ReceiverBankAccountID,
				SenderBankAccountID: receivedTransfer.SenderBankAccountID,
				CreatedAt: receivedTransfer.CreatedAt,
			}
		}).([]dto.FindBankTransferDTO)

		withdrawalsDTO := funk.Map(customer.BankAccount.Withdrawals, func(withdraw models.Withdraw) dto.FindWithdrawDTO {
			return dto.FindWithdrawDTO{
				ID: withdraw.ID,
				Amount: withdraw.Amount,
				BankAccountID: withdraw.BankAccountID,
				CreatedAt: withdraw.CreatedAt,
			}
		}).([]dto.FindWithdrawDTO)

		

		findLoanDTO := dto.FindLoanDTO{
			ID: customer.BankAccount.Loan.ID,
			Amount: customer.BankAccount.Loan.BaseAmount,
			TotalAmount: customer.BankAccount.Loan.TotalAmount,
			InterestRate: customer.BankAccount.Loan.InterestRate,
			CreatedAt: customer.BankAccount.Loan.CreatedAt,
			MaturityDate: customer.BankAccount.Loan.MaturityDate,
		}


		

		customerDTO = dto.FindCustomerDTO{
			ID:            customer.ID,
			Name:          customer.Name,
			CreatedAt:     customer.CreatedAt,
			UpdatedAt:     customer.UpdatedAt,
			RoleUpdatedAt: customer.RoleUpdatedAt,
			Role:          customer.Role,
			BankAccount: dto.FindBankAccountDTO{
				ID:         customer.BankAccount.ID,
				Balance:    customer.BankAccount.Balance,
				CustomerID: customer.BankAccount.CustomerID,
				SentTransfers: sentTransfersDTO,
				ReceivedTransfers: receivedTransfersDTO,
				Loan: &findLoanDTO,
				Deposits:      depositsDTO,
				Withdrawals: withdrawalsDTO,
			},
		}
		return ctx.Status(200).JSON(customerDTO)
	}
}

func (c *CustomerController) ChangeCustomerRole() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		role := ctx.Params("role")
		customerID := ctx.Params("id")
		customerIdInt, err := strconv.Atoi(customerID)
		if err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": "id isn't valid"})
		}
		roleNumber, err := strconv.Atoi(role)
		if err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": "param isn't a number"})
		}
		if roleNumber == utils.OWNER || roleNumber <= 0 && roleNumber >= 3 {
			return ctx.Status(400).JSON(fiber.Map{"error": "you can't give that role to a customer or that role isn't valid"})
		}

		result := c.DB.Model(&models.Customer{}).Where("id = ?", customerIdInt).Updates(map[string]interface{}{
			"role": roleNumber,
			"role_updated_at": time.Now(),
		})
		if result.Error != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}

		return ctx.Status(200).JSON(fiber.Map{"message": "customer was updated!"})
	}
}

func (c *CustomerController) DeleteCustomer() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		claims := ctx.Locals("customer").(jwt.MapClaims)
		role := claims["role"].(float64)
		var searchedID int
		var searchedCustomer models.Customer
		var err error
		if funk.Contains([2]int{utils.MANAGER, utils.OWNER}, int(role)) {
			idParam := ctx.Params("id")
			searchedID, err = strconv.Atoi(idParam)
			if err != nil {
				return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
		} else {
			idReq := claims["ID"].(float64)
			searchedID = int(idReq)
		}


		result := c.DB.Preload("BankAccount").First(&searchedCustomer, searchedID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return ctx.Status(500).JSON(fiber.Map{"error": "that customer doesn't exist"})
			} else {
				return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
			}
		}

		if searchedCustomer.BankAccount.Balance > 0 {
			return ctx.Status(401).JSON(fiber.Map{"error": "that bank account has balance"})
		}


		result = c.DB.Delete(&searchedCustomer)
		if result.Error != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}
		return ctx.Status(200).JSON(fiber.Map{"message": "customer was deleted!"})
	}
}