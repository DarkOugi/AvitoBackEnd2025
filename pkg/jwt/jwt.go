package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

var jwtSecretKey = []byte("testJWT")

func GenerateTokenAccess(user string) (string, error) {
	payload := jwt.MapClaims{
		"user": user,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwtSecretKey)

	return t, err
}

type UserClaims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

func GetInfoFromToken(token string) (*UserClaims, error) {

	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if claims, ok := t.Claims.(*UserClaims); ok && t.Valid {
		return claims, nil
	} else {
		return claims, err
	}
}
