package jwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateTokenAccess(t *testing.T) {
	users := []string{
		"Denis",
		"denis.zhili",
		"1231@ASDASaczxc",
	}
	for _, u := range users {
		t.Run("Test correct create JWT", func(t *testing.T) {
			token, err := GenerateTokenAccess(u)
			t.Logf("%s\n", token)
			assert.Nil(t, err)
		})
	}
}

func TestGetInfoFromToken(t *testing.T) {
	users := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiRGVuaXMifQ.fXfpc4ZIVQM7f_BF-E0zIDNPHPXqDIK3H87VX4wLhu8",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiZGVuaXMuemhpbGkifQ.PUH6RxRKQ4rFoQvQqHumjFzrQTiES_AGZEww5camtHE",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiMTIzMUBBU0RBU2FjenhjIn0.Bln1anwnkruAlmvCEqEzzLb3vlJ4pjjch57XZ0o-pLs",
	}
	for _, u := range users {
		t.Run("Test correct create JWT", func(t *testing.T) {
			user, err := GetInfoFromToken(u)
			t.Logf("%s\n", user.User)
			assert.Nil(t, err)
		})
	}
}
