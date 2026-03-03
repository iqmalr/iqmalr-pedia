package main

import (
	"fmt"
	"log"

	"github.com/iqmalr-pedia/go-auth/v2/internal/models"
	"github.com/iqmalr-pedia/go-auth/v2/pkg/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func main() {
	database.ConnectDB()
	db := database.GetDB()

	err := db.Transaction(func(tx *gorm.DB) error {
		fmt.Println("🌱 Memulai proses seeding massal untuk go-auth...")
		if err := seedUsers(tx); err != nil {
			return err // Jika gagal, kembalikan error untuk trigger rollback
		}
		if err := seedUserAddresses(tx); err != nil {
			return err // Jika gagal, kembalikan error untuk trigger rollback
		}
		fmt.Println("✅ Semua proses seeding berhasil.")
		return nil
	})

	if err != nil {
		log.Fatalf("❌ Seeding gagal! Transaksi di-rollback. Error: %v", err)
	} else {
		fmt.Println("✅ Seeding go-auth selesai dan berhasil di-commit ke database!")
	}
}

// seedUsers akan membuat 1000 customer dan 500 vendor
func seedUsers(tx *gorm.DB) error {
	fmt.Println("- Menyiapkan data Users (1000 Customer, 500 Vendor)...")

	// Hash password sekali di luar loop untuk efisiensi
	// Password untuk semua user adalah "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal melakukan hash password: %w", err)
	}

	const customerCount = 1000
	const vendorCount = 500
	var users []models.User

	// Generate 1000 Customer
	for i := 1; i <= customerCount; i++ {
		users = append(users, models.User{
			Name:      fmt.Sprintf("Customer %d", i),
			Email:     fmt.Sprintf("customer%d@example.com", i),
			Password:  string(hashedPassword),
			Role:      "customer",
			Phone:     "+628133333333", // Nomor HP sama untuk semua
			AvatarUrl: "",              // URL dikosongkan
			IsActive:  true,
		})
	}

	// Generate 500 Vendor
	for i := 1; i <= vendorCount; i++ {
		users = append(users, models.User{
			Name:      fmt.Sprintf("Vendor %d", i),
			Email:     fmt.Sprintf("vendor%d@example.com", i),
			Password:  string(hashedPassword),
			Role:      "vendor",
			Phone:     "+628133333333", // Nomor HP sama untuk semua
			AvatarUrl: "",              // URL dikosongkan
			IsActive:  true,
		})
	}

	// Lakukan bulk insert dengan OnConflict untuk menjadikannya idempoten
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&users).Error; err != nil {
		return fmt.Errorf("gagal melakukan bulk insert users: %w", err)
	}

	fmt.Printf("  - Berhasil memproses %d data user.\n", customerCount+vendorCount)
	return nil
}

// seedUserAddresses akan membuat alamat untuk sebagian user yang sudah dibuat
func seedUserAddresses(tx *gorm.DB) error {
	fmt.Println("- Menyiapkan data User Addresses (untuk 100 Customer dan 50 Vendor pertama)...")

	const customerAddressCount = 100
	const vendorAddressCount = 50
	var addresses []models.UserAddress

	// Buat alamat untuk 100 customer pertama (ID 1 - 100)
	for i := 1; i <= customerAddressCount; i++ {
		addresses = append(addresses, models.UserAddress{
			UserID:        uint(i), // Asumsi ID customer dimulai dari 1
			Label:         "Rumah",
			RecipientName: fmt.Sprintf("Customer %d", i),
			Phone:         "+628133333333",
			AddressLine1:  fmt.Sprintf("Jl. Contoh Alamat No. %d", i),
			City:          "Jakarta",
			State:         "DKI Jakarta",
			PostalCode:    fmt.Sprintf("%d", 10000+i),
			Country:       "Indonesia",
			IsDefault:     true,
		})
	}

	// Buat alamat untuk 50 vendor pertama (ID 1001 - 1050)
	// Asumsi ID vendor dimulai setelah customer terakhir (1000)
	for i := 1; i <= vendorAddressCount; i++ {
		vendorID := 1000 + i
		addresses = append(addresses, models.UserAddress{
			UserID:        uint(vendorID),
			Label:         "Kantor",
			RecipientName: fmt.Sprintf("Vendor %d", i),
			Phone:         "+628133333333",
			AddressLine1:  fmt.Sprintf("Jl. Pusat Bisnis No. %d", i),
			City:          "Surabaya",
			State:         "Jawa Timur",
			PostalCode:    fmt.Sprintf("%d", 60000+i),
			Country:       "Indonesia",
			IsDefault:     true,
		})
	}

	// Lakukan bulk insert untuk alamat
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&addresses).Error; err != nil {
		return fmt.Errorf("gagal melakukan bulk insert addresses: %w", err)
	}

	fmt.Printf("  - Berhasil memproses %d data alamat user.\n", customerAddressCount+vendorAddressCount)
	return nil
}
