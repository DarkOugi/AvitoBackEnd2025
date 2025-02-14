package service

import (
	"avito/internal/entity"
	"avito/pkg/auth"
	"avito/pkg/jwt"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	GetUserInfo(ctx context.Context, login string) (*entity.User, bool, error)
	InitUser(ctx context.Context, login, password string) error
	//getMerchInfo(login string) ([]*Merch, error)
	//getUserToUserInfo(sql, login string) ([]*User, error)
	GetInfo(ctx context.Context, login string) (int, []*entity.Merch, []*entity.User, []*entity.User, error)
	BuyItem(ctx context.Context, login, merch string) error
	SendCoin(ctx context.Context, from, to string, cost int) error
}

var (
	ErrBadAuth      = errors.New("bad auth")
	ErrBadPassword  = errors.New("bad password")
	ErrUnCorrectJWT = errors.New("not correct JST")
)

type Service struct {
	rep Repository
}

func NewService(rep Repository) *Service {
	return &Service{
		rep: rep,
	}
}
func (sv *Service) Auth(ctx context.Context, login, password string) (string, error) {
	if ok := auth.CheckLogin(login); !ok {
		return "", fmt.Errorf("%w: not valid login", ErrBadAuth)
	}

	if len(password) < 4 {
		return "", fmt.Errorf("%w: very short password", ErrBadAuth)
	}

	log.Debug().Msgf("%s", password)
	u, ok, err := sv.rep.GetUserInfo(ctx, login)
	if err != nil {
		return "", fmt.Errorf("can't get user info: %w", err)
	}

	if !ok {
		cErr := sv.rep.InitUser(ctx, login, password)
		if cErr != nil {
			return "", fmt.Errorf("can't create user: %w", cErr)
		}

		tokenJWT, errJWT := jwt.GenerateTokenAccess(login)
		if err != nil {
			return "", fmt.Errorf("can't create token: %w", errJWT)
		}
		return tokenJWT, nil
	}

	if password == u.Password {
		tokenJWT, errJWT := jwt.GenerateTokenAccess(login)
		if err != nil {
			return "", fmt.Errorf("can't create token: %w", errJWT)
		}
		return tokenJWT, nil
	}

	return "", fmt.Errorf("uncorrect login/password: %w", ErrBadPassword)
}
func (sv *Service) Info(ctx context.Context, sec string) (int, []*entity.Merch, []*entity.User, []*entity.User, error) {
	// получить токен проверить что он валидный
	// передать токен в функцию
	t, err := jwt.GetInfoFromToken(sec)
	if err != nil {
		return 0, []*entity.Merch{}, []*entity.User{}, []*entity.User{}, fmt.Errorf("uncorrect security: %w", ErrUnCorrectJWT)
	}
	balance, merch, from, to, errSql := sv.rep.GetInfo(ctx, t.User)
	if err != nil {
		return 0, []*entity.Merch{}, []*entity.User{}, []*entity.User{}, fmt.Errorf("can't get correct info: %w", errSql)
	}
	return balance, merch, from, to, nil

}

func (sv *Service) SendCoin(ctx context.Context, sec, toUs string, amount int) error {
	t, err := jwt.GetInfoFromToken(sec)
	if err != nil {
		return ErrUnCorrectJWT
	}

	errTs := sv.rep.SendCoin(ctx, t.User, toUs, amount)

	if errTs != nil {
		return fmt.Errorf("not success transaction: %w", errTs)
	}

	return nil
}

func (sv *Service) BuyItem(ctx context.Context, item string, sec string) error {
	t, err := jwt.GetInfoFromToken(sec)
	if err != nil {
		return ErrUnCorrectJWT
	}

	errBY := sv.rep.BuyItem(ctx, t.User, item)
	if errBY != nil {
		return fmt.Errorf("not success transaction: %w", errBY)
	}
	return nil
}
