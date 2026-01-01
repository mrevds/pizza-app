package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CvvHash(cvv string) string {
	hashedCvv, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Sprintf("failed to hash cvv: %v", err)
	}
	return string(hashedCvv)
}
