package jwt

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecretKey = []byte("testJWT")

func MetaJWT(user string) string {
	payload := jwt.MapClaims{
		"user": user,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		fmt.Printf("Error jwt token created: %s", err.Error())
		return ""
	} else {
		return t
	}
}
