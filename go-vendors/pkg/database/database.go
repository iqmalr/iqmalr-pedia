package database

import (
	"log"
	_ "os"

	"github.com/iqmalr-pedia/go-vendors/internal/config"
	"github.com/iqmalr-pedia/go-vendors/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	DB *gorm.DB
}

var Database DBInstance

func ConnectDB() {
	db, err := gorm.Open(postgres.Open(config.AppConfig.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	err = db.AutoMigrate(
		&models.Vendor{},
		&models.VendorUser{},
		&models.VendorSetting{},
		&models.VendorBankAccount{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	log.Println("Vendor database connected successfully")
	Database = DBInstance{DB: db}
}

func GetDB() *gorm.DB {
	return Database.DB
}
