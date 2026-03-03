package database

import (
	"fmt"
	"log"
	"os"

	"github.com/iqmalr-pedia/go-transactions/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBInstance struct {
	DB *gorm.DB
}

var Database DBInstance

func ConnectDB() {
	var dsn string

	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "iqmalr")
	dbName := getEnv("DB_NAME", "transaction_db")
	dbPort := getEnv("DB_PORT", "5432")
	dbSSLMode := getEnv("DB_SSL_MODE", "disable")
	dbTimeZone := getEnv("DB_TIME_ZONE", "Asia/Jakarta")

	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode, dbTimeZone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	err = db.AutoMigrate(
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.OrderPayment{},
		&models.OrderStatusHistory{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	log.Println("Database connected successfully")
	Database = DBInstance{DB: db}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetDB() *gorm.DB {
	return Database.DB
}

func WithTransaction(db *gorm.DB, fn func(*gorm.DB) error) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Transaction rolled back due to panic: %v", r)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback().Error; rbErr != nil {
			log.Printf("Failed to rollback transaction: %v", rbErr)
		}
		return err
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}

	return nil
}
