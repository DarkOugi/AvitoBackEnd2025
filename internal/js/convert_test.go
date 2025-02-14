package js

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFromJSUser(t *testing.T) {
	gotCorrect := []string{
		`{"login" : "denis.zhilin@avito.ru",    "password" : "12345"}`,
	}
	for _, g := range gotCorrect {
		t.Run("Test correct JSON user", func(t *testing.T) {
			var want = &GetUser{
				Login:    "denis.zhilin@avito.ru",
				Password: "12345",
			}

			res, err := GetFromJSUser([]byte(g))

			if err != nil {
				t.Errorf("%s", err.Error())
			}
			assert.Equal(t, want, res, "Parse json to struct eq = %s", g)
		})
	}

	gotNotCorrect := []string{
		`{"login" : "denis.zhilin@avito.ru",    "password" : 1}`,
		`{"login" : ["denis.zhilin@avito.ru"],    "password" : "12345", "notUsed" : "123" }`,
		`{"login" : denis.zhilin@avito.ru,    "password" : "12345", "notUsed"":123}`,
	}

	for _, g := range gotNotCorrect {
		t.Run("Test not correct JSON user", func(t *testing.T) {
			var want = &GetUser{
				Login:    "denis.zhilin@avito.ru",
				Password: "12345",
			}

			res, err := GetFromJSUser([]byte(g))

			if err == nil {
				t.Errorf("%s", err.Error())
			}
			if want == res {
				t.Errorf("Failed Json = %s", g)
			}
			//assert.Equal(t, want, res, "CheckLogin eq = %s", string(g))
		})
	}
}

func TestGetFromJSSecurity(t *testing.T) {
	t.Run("Test correct JSON security", func(t *testing.T) {
		var want = &GetSecurity{
			Security: "123asd",
		}

		got := `{ "security" : "123asd"}`
		res, err := GetFromJSSecurity([]byte(got))

		if err != nil {
			assert.Errorf(t, err, "Error Parse")
		}

		assert.Equal(t, want.Security, res, "From Json Security eq = %s", got)
	})

	//gotNotCorrect := []string{
	//	`{ "security" : "123asd"}`,
	//	`{"security" : "123asd", "notUsed" : "123" }`,
	//	`{"security" : 123}`,
	//}

	//for _, g := range gotNotCorrect {
	//	t.Run("Test not correct JSON security", func(t *testing.T) {
	//		//var want = &GetSecurity{
	//		//	Security: "123asd",
	//		//}
	//
	//		_, err := GetFromJSSecurity([]byte(g))
	//
	//		if err == nil {
	//			t.Errorf("%s", err.Error())
	//		}
	//		//if want == res {
	//		//	t.Errorf("Failed Json = %s", g)
	//		//}
	//		//assert.Equal(t, want, res, "CheckLogin eq = %s", string(g))
	//	})
	//}

}

func TestGetFromJsUserToUser(t *testing.T) {
	t.Run("Test correct JSON userToUser", func(t *testing.T) {
		var want = &GetUserToUser{
			Security: "123asd",
			ToUser:   "den",
			Amount:   100,
		}

		got := `{ "security" : "123asd", "toUser":"den", "amount":100 }`
		res, err := GetFromJsUserToUser([]byte(got))

		assert.Nil(t, err)
		assert.Equal(t, want, res, "From Json UserToUser eq = %s", got)
	})

}

func TestToJSError(t *testing.T) {
	t.Run("Test correct JSON userToUser", func(t *testing.T) {
		got := "123asd"

		_, err := ToJSError(got)

		assert.Nil(t, err)
		//assert.Equal(t, []byte(want), res, "From Json UserToUser eq = %s", got)
	})
}
