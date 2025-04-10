package controllers

import (
	"banco/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DepositController struct {
	DB *gorm.DB
}

func (d *DepositController) CreateDeposit() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var deposit models.Deposit
	}
}