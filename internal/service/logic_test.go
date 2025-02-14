package service

import (
	"avito/internal/entity"
	"avito/pkg/jwt"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var email = "denis.zhilin@avito.ru"

type TestRep struct {
	correct bool
	newUser bool
}

func (r *TestRep) GetUserInfo(ctx context.Context, login string) (*entity.User, bool, error) {
	//fmt.Println(correctUser.Password)
	if r.correct {
		if r.newUser {
			return &entity.User{
				login,
				0,
				"qweqweq",
			}, false, nil
		}

		return &entity.User{
			login,
			0,
			"qweqweq",
		}, true, nil
	} else {
		return nil, false, ErrBadAuth
	}
}
func (r *TestRep) InitUser(ctx context.Context, login, password string) error {
	if r.correct {
		return nil
	}
	return ErrBadAuth
}

func (r *TestRep) GetInfo(ctx context.Context, login string) (int, []*entity.Merch, []*entity.User, []*entity.User, error) {
	if r.correct {
		return 0, nil, nil, nil, nil
	}
	return 0, nil, nil, nil, ErrUnCorrectJWT
}
func (r *TestRep) BuyItem(ctx context.Context, login, merch string) error {
	if r.correct {
		return nil
	}
	return ErrUnCorrectJWT
}
func (r *TestRep) SendCoin(ctx context.Context, from, to string, cost int) error {
	if r.correct {
		return nil
	}
	return ErrUnCorrectJWT
}

func TestNewService(t *testing.T) {
	t.Run("Create NewService", func(t *testing.T) {
		sv := NewService(&TestRep{
			correct: true,
		})
		assert.IsType(t, sv, &Service{})
	})
}
func TestService_Auth(t *testing.T) {
	t.Run("not correct email", func(t *testing.T) {
		sv := NewService(&TestRep{
			correct: true,
			newUser: false,
		})
		_, errAuth := sv.Auth(context.Background(), "denis.2.zhilin@avito.ru", "qweqweq")
		//passJWT, _ := jwt.GenerateTokenAccess("denis.zhilin@avito.ru")

		//assert.Equal(t, passJWT, hashPasswordTest)
		assert.ErrorIs(t, errAuth, ErrBadAuth)
	})
	t.Run("not correct pass", func(t *testing.T) {
		sv := NewService(&TestRep{
			correct: true,
			newUser: false,
		})
		_, errAuth := sv.Auth(context.Background(), "denis.zhilin@avito.ru", "qwe")
		//passJWT, _ := jwt.GenerateTokenAccess("denis.zhilin@avito.ru")

		//assert.Equal(t, passJWT, hashPasswordTest)
		assert.ErrorIs(t, errAuth, ErrBadAuth)
	})
	t.Run("correct Auth User", func(t *testing.T) {
		sv := NewService(&TestRep{
			correct: true,
			newUser: false,
		})
		hashPasswordTest, errAuth := sv.Auth(context.Background(), "denis.zhilin@avito.ru", "qweqweq")
		passJWT, _ := jwt.GenerateTokenAccess("denis.zhilin@avito.ru")

		assert.Equal(t, passJWT, hashPasswordTest)
		assert.Nil(t, errAuth)
	})
	t.Run("correct Auth new User", func(t *testing.T) {
		sv := NewService(&TestRep{
			correct: true,
			newUser: true,
		})

		hashPasswordTest, errAuth := sv.Auth(context.Background(), "denis.zhilin@avito.ru", "unHashPassword")
		passJWT, _ := jwt.GenerateTokenAccess("denis.zhilin@avito.ru")

		assert.Equal(t, passJWT, hashPasswordTest)
		assert.Nil(t, errAuth)
		//assert.IsType(t, sv, &Service{})
	})
}

func TestService_SendCoin(t *testing.T) {
	t.Run("correct Send", func(t *testing.T) {
		sv := NewService(&TestRep{
			correct: true,
			newUser: false,
		})
		passJWT, _ := jwt.GenerateTokenAccess("denis.zhilin@avito.ru")
		errSC := sv.SendCoin(context.Background(), passJWT, "asd", 2)
		assert.Nil(t, errSC)
	})

	t.Run("not correct Send", func(t *testing.T) {
		sv := NewService(&TestRep{
			correct: false,
			newUser: false,
		})
		passJWT, _ := jwt.GenerateTokenAccess("denis.zhilin@avito.ru")
		errSC := sv.SendCoin(context.Background(), passJWT, "asd", 2)
		assert.ErrorIs(t, errSC, ErrUnCorrectJWT)
	})

}

func TestService_BuyItem(t *testing.T) {
	t.Run("correct Buy", func(t *testing.T) {
		sv := NewService(&TestRep{
			correct: true,
			newUser: false,
		})
		passJWT, _ := jwt.GenerateTokenAccess("denis.zhilin@avito.ru")
		errSC := sv.BuyItem(context.Background(), "t-shirt", passJWT)
		assert.Nil(t, errSC)
	})

	t.Run("not correct Buy", func(t *testing.T) {
		sv := NewService(&TestRep{
			correct: false,
			newUser: false,
		})
		passJWT, _ := jwt.GenerateTokenAccess("denis.zhilin@avito.ru")
		errSC := sv.BuyItem(context.Background(), "t-shirt", passJWT)
		assert.ErrorIs(t, errSC, ErrUnCorrectJWT)
	})

}

func TestService_Info(t *testing.T) {
	t.Run("correct Info", func(t *testing.T) {
		sv := NewService(&TestRep{
			correct: true,
			newUser: false,
		})
		passJWT, _ := jwt.GenerateTokenAccess("denis.zhilin@avito.ru")
		_, merch, from, to, errSC := sv.Info(context.Background(), passJWT)
		assert.IsType(t, merch, []*entity.Merch{})
		assert.IsType(t, from, []*entity.User{})
		assert.IsType(t, to, []*entity.User{})

		assert.Nil(t, errSC)
	})

	t.Run("not correct Info", func(t *testing.T) {
		sv := NewService(&TestRep{
			correct: false,
			newUser: false,
		})
		passJWT, _ := jwt.GenerateTokenAccess("denis.zhilin@avito.ru")
		_, _, _, _, errSC := sv.Info(context.Background(), passJWT)
		t.Logf("\nERROR %s \n", errSC)
		assert.ErrorIs(t, errSC, ErrUnCorrectJWT)
	})
}
