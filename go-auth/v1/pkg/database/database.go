package database

import (
	"log"

	"github.com/iqmalr-pedia/go-auth/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBInstance struct {
	DB *gorm.DB
}

var Database DBInstance

func ConnectDB() {
	dsn := "host=localhost user=auth_user password=auth_pass dbname=auth_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.VerificationRequest{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	Database = DBInstance{DB: db}
}
