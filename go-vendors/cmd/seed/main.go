package main

import (
	"fmt"
	"log"
	"os"

	"github.com/iqmalr-pedia/go-vendors/internal/models"
	"github.com/iqmalr-pedia/go-vendors/pkg/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AuthUser adalah struct sementara untuk mewakili user dari database go-auth
type AuthUser struct {
	ID    uint
	Email string
	Role  string
}

// TableName memberitahu GORM untuk menggunakan tabel 'users'
func (AuthUser) TableName() string {
	return "users"
}

// Salin fungsi getEnv dari go-auth untuk memastikan cara pembacaan env sama
func getAuthEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// 1. Inisialisasi koneksi ke database vendor
	database.ConnectDB()
	dbVendor := database.GetDB()

	// 2. Bangun koneksi ke database auth dengan cara yang SAMA PERSIS seperti go-auth
	log.Println("Membangun koneksi ke database AUTH...")
	authDbHost := getAuthEnv("DB_HOST", "localhost")
	authDbUser := getAuthEnv("DB_USER", "postgres")
	authDbPassword := getAuthEnv("DB_PASSWORD", "iqmalr")
	authDbName := getAuthEnv("DB_NAME", "auth_db_v2")
	authDbPort := getAuthEnv("DB_PORT", "5432")
	authDbSSLMode := getAuthEnv("DB_SSL_MODE", "disable")
	authDbTimeZone := getAuthEnv("DB_TIME_ZONE", "Asia/Jakarta")

	authDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		authDbHost, authDbUser, authDbPassword, authDbName, authDbPort, authDbSSLMode, authDbTimeZone)

	// CETAK DSN UNTUK DEBUGGING. Ini sangat penting!
	log.Printf("Mencoba terhubung ke database AUTH dengan DSN: %s", authDSN)

	dbAuth, err := gorm.Open(postgres.Open(authDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal koneksi ke database AUTH dengan DSN di atas. Pastikan database dan tabel sudah ada. Error: %v", err)
	}
	log.Println("✅ Koneksi ke database AUTH berhasil (read-only).")

	// 3. Jalankan seluruh proses seeding di dalam satu transaksi di DB vendor
	err = dbVendor.Transaction(func(tx *gorm.DB) error {
		fmt.Println("🌱 Memulai proses seeding massal untuk go-vendors...")

		// Cari semua user dengan role 'vendor' dari database auth
		var vendorUsers []AuthUser
		if err := dbAuth.Where("role = ?", "vendor").Find(&vendorUsers).Error; err != nil {
			return fmt.Errorf("gagal mengambil data vendor dari database auth. Pastikan tabel 'users' ada dan ada data dengan role 'vendor'. Error: %w", err)
		}
		if len(vendorUsers) == 0 {
			return fmt.Errorf("tidak ada user dengan role 'vendor' ditemukan di database auth. Jalankan seeder go-auth terlebih dahulu.")
		}
		fmt.Printf("- Ditemukan %d user vendor di database auth.\n", len(vendorUsers))

		// Jalankan fungsi seeder dengan data yang sudah ditemukan
		if err := seedVendors(tx, vendorUsers); err != nil {
			return err
		}
		if err := seedVendorUsers(tx, vendorUsers); err != nil {
			return err
		}
		if err := seedVendorBankAccounts(tx); err != nil {
			return err
		}
		fmt.Println("✅ Semua proses seeding berhasil.")
		return nil
	})

	// 4. Cek hasil transaksi
	if err != nil {
		log.Fatalf("❌ Seeding gagal! Transaksi di-rollback. Error: %v", err)
	} else {
		fmt.Println("✅ Seeding go-vendors selesai dan berhasil di-commit ke database!")
	}
}

// ... (fungsi seedVendors, seedVendorUsers, seedVendorBankAccounts tidak berubah) ...
func seedVendors(tx *gorm.DB, vendorUsers []AuthUser) error {
	fmt.Println("- Menyiapkan data Vendors...")
	var vendors []models.Vendor
	for _, user := range vendorUsers {
		vendors = append(vendors, models.Vendor{
			OwnerID:      user.ID,
			Name:         fmt.Sprintf("PT. Vendor %s", user.Email),
			Slug:         fmt.Sprintf("pt-vendor-%d", user.ID),
			Description:  fmt.Sprintf("Vendor yang dimiliki oleh %s", user.Email),
			LogoUrl:      "",
			BannerUrl:    "",
			ContactEmail: fmt.Sprintf("contact@%s", user.Email),
			ContactPhone: "+628133333333",
			Status:       "approved",
			BusinessType: "Retail",
			AddressLine1: fmt.Sprintf("Jl. Pergudangan No. %d", user.ID),
			City:         "Jakarta",
			State:        "DKI Jakarta",
			PostalCode:   fmt.Sprintf("%d", 10000+user.ID),
			Country:      "Indonesia",
		})
	}
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&vendors).Error; err != nil {
		return fmt.Errorf("gagal melakukan bulk insert vendors: %w", err)
	}
	fmt.Printf("  - Berhasil memproses %d data vendor.\n", len(vendors))
	return nil
}

func seedVendorUsers(tx *gorm.DB, vendorUsers []AuthUser) error {
	fmt.Println("- Menyiapkan data Vendor Users...")
	var createdVendors []models.Vendor
	ownerIDs := make([]uint, len(vendorUsers))
	for i, user := range vendorUsers {
		ownerIDs[i] = user.ID
	}
	if err := tx.Where("owner_id IN ?", ownerIDs).Find(&createdVendors).Error; err != nil {
		return fmt.Errorf("gagal mengambil data vendor yang baru dibuat: %w", err)
	}
	var vendorUsersToInsert []models.VendorUser
	for _, vendor := range createdVendors {
		vendorUsersToInsert = append(vendorUsersToInsert, models.VendorUser{
			VendorID: vendor.ID,
			UserID:   vendor.OwnerID,
			Role:     "owner",
			IsActive: true,
		})
	}
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&vendorUsersToInsert).Error; err != nil {
		return fmt.Errorf("gagal melakukan bulk insert vendor users: %w", err)
	}
	fmt.Printf("  - Berhasil memproses %d data hubungan vendor-user.\n", len(vendorUsersToInsert))
	return nil
}

func seedVendorBankAccounts(tx *gorm.DB) error {
	fmt.Println("- Menyiapkan data Vendor Bank Accounts (untuk 100 Vendor pertama)...")
	var vendors []models.Vendor
	if err := tx.Limit(100).Find(&vendors).Error; err != nil {
		return fmt.Errorf("gagal mengambil 100 vendor pertama: %w", err)
	}
	bankNames := []string{"Bank Central Asia", "Bank Mandiri", "Bank Negara Indonesia"}
	var bankAccounts []models.VendorBankAccount
	for i, vendor := range vendors {
		bankAccounts = append(bankAccounts, models.VendorBankAccount{
			VendorID:      vendor.ID,
			BankName:      bankNames[i%len(bankNames)],
			AccountNumber: fmt.Sprintf("%010d", vendor.ID),
			AccountHolder: fmt.Sprintf("Pemilik %s", vendor.Name),
			IsPrimary:     i == 0,
		})
	}
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&bankAccounts).Error; err != nil {
		return fmt.Errorf("gagal melakukan bulk insert bank accounts: %w", err)
	}
	fmt.Printf("  - Berhasil memproses %d data akun bank vendor.\n", len(bankAccounts))
	return nil
}
