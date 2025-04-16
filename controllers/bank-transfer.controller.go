package controllers

import (
	"banco/dto"
	"banco/models"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thoas/go-funk"
	"gorm.io/gorm"
)

type BankTransferController struct {
	DB *gorm.DB
}

func (b *BankTransferController) CreateTransfer() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var input dto.CreateBankTransferDTO
		var receiver models.Customer
		var sender models.Customer
		claims := ctx.Locals("customer").(jwt.MapClaims)
		IdTokenFloat := claims["ID"].(float64)
		IdTokenUint := uint(IdTokenFloat)
		if err := ctx.BodyParser(&input); err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if err := validate.Struct(input); err != nil {
			return ctx.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		result := b.DB.Preload("BankAccount").First(&sender, IdTokenUint)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return ctx.Status(404).JSON(fiber.Map{"error": "sender wasn't found"})
			} else {
				return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
			}
		}

		result = b.DB.Preload("BankAccount").First(&receiver, input.ReceiverID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return ctx.Status(404).JSON(fiber.Map{"error": "receiver wasn't found"})
			} else {
				return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
			}
		}

		if sender.BankAccount.Balance < input.Amount {
			return ctx.Status(401).JSON(fiber.Map{"error": "you don't have enough money"})
		}
		

		err := b.DB.Transaction(func(tx *gorm.DB) error {
			var sender models.Customer
			var receiver models.Customer

			if err := tx.Preload("BankAccount").First(&sender, IdTokenUint).Error; err != nil {
				return err
			}
			if err := tx.Preload("BankAccount").First(&receiver, input.ReceiverID).Error; err != nil {
				return err
			}

			if sender.BankAccount.Balance < input.Amount {
				return fiber.NewError(fiber.StatusUnauthorized, "you don't have enough money")
			}

			if err := tx.Model(&models.BankAccount{}).
				Where("customer_id = ?", sender.ID).
				Update("balance", gorm.Expr("balance - ?", input.Amount)).Error; err != nil {
				return err
			}

			if err := tx.Model(&models.BankAccount{}).
				Where("customer_id = ?", receiver.ID).
				Update("balance", gorm.Expr("balance + ?", input.Amount)).Error; err != nil {
				return err
			}

			transfer := models.BankTransfer{
				Amount:     input.Amount,
				ReceiverID: receiver.ID,
				SenderID:   sender.ID,
			}
			if err := tx.Create(&transfer).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		
		return ctx.Status(200).JSON(fiber.Map{"message": "transfer was created!"})
	}
}

func (b *BankTransferController) FindAllTransfers() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var transfers []models.BankTransfer
		result := b.DB.Find(&transfers)
		if result.Error != nil {
			return ctx.Status(500).JSON(fiber.Map{"error": result.Error.Error()})
		}
		mapped := funk.Map(transfers, func(transfer models.BankTransfer) dto.FindBankTransferDTO {
			return dto.FindBankTransferDTO{
				ID: transfer.ID,
				Amount: transfer.Amount,
				ReceiverID: transfer.ReceiverID,
				SenderID: transfer.SenderID,
				
			}
		}).([]dto.FindBankTransferDTO)
		return ctx.Status(200).JSON(mapped)
	}
}