package auth

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashPassword(t *testing.T) {
	t.Run("Test eq hash for eq string", func(t *testing.T) {
		assert.Equal(t, HashPassword("hello"), HashPassword("hello"))
	})

	tests := map[string][]string{
		"hello": {
			"HellO",
			"HELlo",
			"hElLO",
		},
	}

	for k, v := range tests {
		one := HashPassword(k)

		for _, el := range v {
			two := HashPassword(el)
			t.Run("Test not eq hash for not eq string", func(t *testing.T) {
				assert.NotEqual(t, one, two)
			})
		}
	}
}

func TestCheckLogin(t *testing.T) {
	logins := []string{
		"dasd.awdq12312@avito.ru",
		"qqssq.bbbqws@avito.ru",
		"zxc.xax9@avito.ru",
	}
	for _, g := range logins {
		t.Run("Test correct login", func(t *testing.T) {
			assert.Truef(t, CheckLogin(g), "CheckLogin eq = %s", g)
		})
	}

	logins = []string{
		"daSd.awdq12312@avito.ru",
		"qqssq.bbb.qws@avito.ru",
		"zxc.xa!x9@avito.ru",
		"zxc.xa!x9@avito.com",
		"zxc.xa!x9@avi.to.ru",
		".@avito.ru",
	}
	for _, g := range logins {
		t.Run("Test not correct login", func(t *testing.T) {

			assert.False(t, CheckLogin(g), "CheckLogin eq = %s", g)

		})
	}
}
