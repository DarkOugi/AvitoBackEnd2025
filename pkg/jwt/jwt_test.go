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
			_, err := GenerateTokenAccess(u)
			//t.Logf("%s\n", token)
			assert.Nil(t, err)
		})
	}
}

func TestGetInfoFromToken(t *testing.T) {
	users := []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiRGVuaXMifQ.Xd2wHTWC2a4F0EsdwTLOZInC5fk9BHbi2UbR4XgyjlY",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiZGVuaXMuemhpbGkifQ.tIV14f8XbXvg8eMdoLpx94P5fVFqvjZUyqByl3WmDT0",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiMTIzMUBBU0RBU2FjenhjIn0.Yr07okh7J7rjBEOKr0Msuz1KobW5CSR1d7NmIpBDeuc",
	}
	for _, u := range users {
		t.Run("Test get info from JWT", func(t *testing.T) {
			_, err := GetInfoFromToken(u)

			assert.Nil(t, err)
		})
	}
}
