package database

import (
	"banco/models"
	"banco/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		return nil, fmt.Errorf("there isn't dsn")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	DB = db
	err = migrateModels(db)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = createOwner(db)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return db, nil
}

func migrateModels(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Customer{}, &models.BankAccount{}, &models.Deposit{}, &models.Loan{}, &models.BankTransfer{})
	if err != nil {
		return err
	}
	log.Println("Tables are ok")
	return nil
}

func createOwner(db *gorm.DB) error {
	var owner models.Customer
	ownerPassword := os.Getenv("OWNER_PASSWORD")
	if ownerPassword == "" {
		return fmt.Errorf("there isn't owner's password")
	}
	result := db.Where("name = ?", "OWNER").First(&owner)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			err := db.Transaction(func(tx *gorm.DB) error {
				var bankAccount models.BankAccount
				var reserve uint = 100000
				owner = models.Customer{
					Name:          "OWNER",
					Password:      ownerPassword,
					Role:          utils.OWNER,
					RoleUpdatedAt: time.Now(),
				}
				result = db.Create(&owner)
				if result.Error != nil {
					return result.Error
				}
				bankAccount = models.BankAccount{
					CustomerID: owner.ID,
					Balance:    reserve,
				}
				result = db.Create(&bankAccount)
				if result.Error != nil {
					return result.Error
				}
				log.Println("owner was created")
				return nil
			})
			return err
			// var bankAccount models.BankAccount
			// var reserve uint = 100000
			// owner = models.Customer{
			// 	Name: "OWNER",
			// 	Password: ownerPassword,
			// 	Role: utils.OWNER,
			// 	RoleUpdatedAt: time.Now(),
			// }
			// result = db.Create(&owner)
			// if result.Error != nil {
			// 	return result.Error
			// }
			// bankAccount = models.BankAccount{
			// 	CustomerID: owner.ID,
			// 	Balance: reserve,
			// }
			// result = db.Create(&bankAccount)
			// if result.Error != nil {
			// 	return result.Error
			// }
			// log.Println("owner was created")
			// return nil
		} else {
			return result.Error
		}

	}
	log.Println("owner already exists")
	return nil
}
