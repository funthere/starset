package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(p string) string {
	salt := 8
	hash, _ := bcrypt.GenerateFromPassword([]byte(p), salt)

	return string(hash)
}

func ComparePassword(h, p []byte) bool {
	return bcrypt.CompareHashAndPassword(h, p) == nil
}
