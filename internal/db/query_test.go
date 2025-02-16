package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPostgresDB(t *testing.T) {
	// важно запускать последовательно  !!!

	t.Run("create test user1", func(t *testing.T) {
		err := testDB.InitUser(context.Background(), "test1", "0000")

		require.NoError(t, err)
	})

	t.Run("create test user2", func(t *testing.T) {
		err := testDB.InitUser(context.Background(), "test2", "0000")

		require.NoError(t, err)
	})

	t.Run("get user", func(t *testing.T) {
		user, ok, err := testDB.GetUserInfo(context.Background(), "test1")

		require.True(t, ok)
		require.NoError(t, err)
		require.Equal(t, user.Password, "0000")
	})

	t.Run("get un init user", func(t *testing.T) {
		_, ok, err := testDB.GetUserInfo(context.Background(), "test")

		require.False(t, ok)
		require.Nil(t, err)
		require.False(t, ok)
	})

	t.Run("buy Pen", func(t *testing.T) {
		err := testDB.BuyItem(context.Background(), "test1", "pen")
		require.Nil(t, err)
	})

	t.Run("send to test1 from test2", func(t *testing.T) {
		err := testDB.SendCoin(context.Background(), "test1", "test2", 1)
		require.Nil(t, err)
	})

	t.Run("send to test2 from test1", func(t *testing.T) {
		err := testDB.SendCoin(context.Background(), "test2", "test1", 5)
		require.Nil(t, err)
	})

	t.Run("info", func(t *testing.T) {
		_, _, _, _, err := testDB.GetInfo(context.Background(), "test1")
		require.Nil(t, err)
	})
}
