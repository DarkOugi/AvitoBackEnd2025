package auth

import (
	"crypto/sha512"
	"encoding/hex"
	"regexp"
)

func CheckLogin(login string) bool {
	standart := `^[a-z]+\.[a-z]+\d*@avito\.ru$`

	match, _ := regexp.MatchString(standart, login)

	return match
}

func HashPassword(password string) string {
	sha512 := sha512.New()

	passwordBytes := []byte(password)
	salt := []byte("VetyStrongSalt")
	passwordBytes = append(passwordBytes, salt...)

	sha512.Write(passwordBytes)
	hashedPasswordBytes := sha512.Sum(nil)

	hashedPasswordHex := hex.EncodeToString(hashedPasswordBytes)
	return hashedPasswordHex
}
