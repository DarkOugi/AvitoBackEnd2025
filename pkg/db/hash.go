package db

import (
	"crypto/sha512"
	"encoding/hex"
)

// хранить пароль в чистом виде нельзя - так что спрячем )))

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
