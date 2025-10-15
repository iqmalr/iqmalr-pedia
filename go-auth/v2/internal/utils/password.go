package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidatePasswordStrength(password string) bool {
	if len(password) < 8 {
		return false
	}

	//var (
	//	hasUpper   bool
	//	hasLower   bool
	//	hasNumber  bool
	//	hasSpecial bool
	//)
	//
	//for _, char := range password {
	//	switch {
	//	case char >= 'A' && char <= 'Z':
	//		hasUpper = true
	//	case char >= 'a' && char <= 'z':
	//		hasLower = true
	//	case char >= '0' && char <= '9':
	//		hasNumber = true
	//	case char == '!' || char == '@' || char == '#' || char == '$' || char == '%' || char == '^' || char == '&' || char == '*':
	//		hasSpecial = true
	//	}
	//}
	//
	//return hasUpper && hasLower && hasNumber && hasSpecial

	return true
}
