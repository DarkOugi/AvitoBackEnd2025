package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestBuyItem(t *testing.T) {
	t.Run("create new user", func(t *testing.T) {
		type User struct {
			login    string
			password string
		}
		user := User{
			login:    "denis.zhilin@avito.ru",
			password: "00000",
		}
		userB, err := json.Marshal(user)
		assert.Nil(t, err, "TO JSON")
		//time.Sleep(30 * time.Second)
		fmt.Println(1)
		resp, err := http.Post("http://localhost:8080/api/auth", "application/json", bytes.NewReader(userB))
		//err := testDB.InitUser(context.Background(), "test1", "0000")
		assert.Nil(t, err, "POST")
		fmt.Println(resp)
		//require.NoError(t, err, "POST")

	})
}
