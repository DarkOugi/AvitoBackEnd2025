package db

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPostgresDB_GetUserInfo(t *testing.T) {
	t.Run("get user", func(t *testing.T) {
		_, _, err := testDB.GetUserInfo(context.Background(), "denis.zhilin@avito,ru")
		assert.NoError(t, err)
	})
}

func TestPostgresDB_GetInfo(t *testing.T) {
	t.Run("get user", func(t *testing.T) {
		_, _, _, _, err := testDB.GetInfo(context.Background(), "denis.zhilin@avito.ru")
		assert.NoError(t, err)
	})
}

//func TestPostgresDB_InitUser(t *testing.T) {
//
//	t.Run("create new user", func(t *testing.T) {
//		err := testDB.InitUser(context.Background(), "asd234565", "qqqq123")
//		assert.NoError(t, err)
//
//	})
//
//
//}
