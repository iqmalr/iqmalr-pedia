package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SKUCounter struct {
	VendorID uint `gorm:"primaryKey"`
	Counter  uint `gorm:"not null;default:0"`
}

func NextSKUSeq(db *gorm.DB, vendorID uint) (uint, error) {
	var counter SKUCounter

	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			FirstOrCreate(&counter, SKUCounter{VendorID: vendorID}).Error
		if err != nil {
			return err
		}

		counter.Counter++
		return tx.Save(&counter).Error
	})

	return counter.Counter, err
}
