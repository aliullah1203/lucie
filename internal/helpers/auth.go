package helpers

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a plain password
func HashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hash)
}

// VerifyPassword compares hashed and plain passwords
//
//	func VerifyPassword(hashedPassword, password string) bool {
//		err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
//		return err == nil
//	}
func VerifyPassword(stored, provided string) bool {
	return stored == provided
}
