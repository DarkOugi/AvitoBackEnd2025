package server

import (
	"avito/internal/entity"
	"avito/internal/service"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestRep struct {
	correct bool
	newUser bool
}

func (r *TestRep) GetUserInfo(ctx context.Context, login string) (*entity.User, bool, error) {
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
		return nil, false, service.ErrBadAuth
	}
}
func (r *TestRep) InitUser(ctx context.Context, login, password string) error {
	if r.correct {
		return nil
	}
	return service.ErrBadAuth
}

func (r *TestRep) GetInfo(ctx context.Context, login string) (int, []*entity.Merch, []*entity.User, []*entity.User, error) {
	if r.correct {
		return 0, nil, nil, nil, nil
	}
	return 0, nil, nil, nil, service.ErrBadAuth
}
func (r *TestRep) BuyItem(ctx context.Context, login, merch string) error {
	if r.correct {
		return nil
	}
	return service.ErrBadAuth
}
func (r *TestRep) SendCoin(ctx context.Context, from, to string, cost int) error {
	if r.correct {
		return nil
	}
	return service.ErrBadAuth
}
func TestNewServer(t *testing.T) {
	t.Run("Test create mock server", func(t *testing.T) {
		sv := service.NewService(&TestRep{
			correct: true,
		})
		ser := NewServer(sv)
		assert.NotNil(t, ser)
	})
}
