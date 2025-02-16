package js

import (
	"avito/internal/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFromJSUser(t *testing.T) {
	jsCorrect := []string{
		`{"login" : "denis.zhilin@avito.ru",    "password" : "12345"}`,
	}
	for _, js := range jsCorrect {
		t.Run("Test correct JSON user", func(t *testing.T) {
			var want = &GetUser{
				Login:    "denis.zhilin@avito.ru",
				Password: "12345",
			}

			res, err := GetFromJSUser([]byte(js))

			assert.Nil(t, err)
			assert.Equal(t, want, res, "Parse json to struct eq = %s", js)
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

		got, err := GetFromJSSecurity([]byte(`{ "security" : "123asd"}`))

		assert.Nil(t, err, "Error Parse")
		assert.Equal(t, want.Security, got, "From Json Security eq = %s", got)
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

		js := `{ "security" : "123asd", "toUser":"den", "amount":100 }`
		got, err := GetFromJsUserToUser([]byte(js))

		assert.Nil(t, err)
		assert.Equal(t, want, got, "From Json UserToUser eq = %s", got)
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

func TestToJsToken(t *testing.T) {
	t.Run("Test correct JSON userToUser", func(t *testing.T) {
		got := "123asd"

		_, err := ToJsToken(got)

		assert.Nil(t, err)
		//assert.Equal(t, []byte(want), res, "From Json UserToUser eq = %s", got)
	})
}

func TestToJsMerch(t *testing.T) {
	t.Run("Test correct JSON userToUser", func(t *testing.T) {
		got := []*entity.Merch{
			{
				Name: "t-shirt",
				Cnt:  80,
			},
		}

		merch := ToJsMerch(got)

		assert.NotNil(t, merch)
		//assert.Equal(t, []byte(want), res, "From Json UserToUser eq = %s", got)
	})
}

func TestToJsFromUser(t *testing.T) {
	t.Run("Test correct JSON userToUser", func(t *testing.T) {
		got := []*entity.User{
			{
				Name:     "t-shirt",
				Password: "123",
				Cost:     80,
			},
		}

		user := ToJsFromUser(got)

		assert.NotNil(t, user)
		//assert.Equal(t, []byte(want), res, "From Json UserToUser eq = %s", got)
	})
}

func TestToJsToUser(t *testing.T) {
	t.Run("Test correct JSON userToUser", func(t *testing.T) {
		got := []*entity.User{
			{
				Name:     "t-shirt",
				Password: "123",
				Cost:     80,
			},
		}

		user := ToJsToUser(got)

		assert.NotNil(t, user)
		//assert.Equal(t, []byte(want), res, "From Json UserToUser eq = %s", got)
	})
}

func TestToJsInfo(t *testing.T) {
	t.Run("Test correct JSON userToUser", func(t *testing.T) {
		gotU := []*entity.User{
			{
				Name:     "t-shirt",
				Password: "123",
				Cost:     80,
			},
		}

		gotM := []*entity.Merch{
			{
				Name: "t-shirt",
				Cnt:  80,
			},
		}

		_, err := ToJsInfo(100, gotM, gotU, gotU)

		assert.Nil(t, err)
		//assert.Equal(t, []byte(want), res, "From Json UserToUser eq = %s", got)
	})
}
