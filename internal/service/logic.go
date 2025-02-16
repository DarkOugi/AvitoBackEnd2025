package service

import (
	"avito/internal/entity"
	"avito/pkg/auth"
	"avito/pkg/jwt"
	"context"
	"errors"
	"fmt"
)

type Repository interface {
	GetUserInfo(ctx context.Context, login string) (*entity.User, bool, error)
	InitUser(ctx context.Context, login, password string) error
	GetInfo(ctx context.Context, login string) (int, []*entity.Merch, []*entity.User, []*entity.User, error)
	BuyItem(ctx context.Context, login, merch string) error
	SendCoin(ctx context.Context, from, to string, cost int) error
}

const minLenPassword = 4

var (
	ErrBadAuth      = errors.New("bad auth")
	ErrBadPassword  = errors.New("bad password")
	ErrUnCorrectJWT = errors.New("not correct JWT")
	ErrLogic        = errors.New("logic error")
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
	// логика авторизации
	// 1 проверить логин и пароль на валидность
	// Проверяем существует ли такой пользователь
	// Если существует то
	// а) получаем его хэшированный пароль и сравниваем с переданным пользователейм в форме
	// (нужно его перед этим захэшировать)
	// если пароли не равны - ошибка, иначе вернем jwt токен
	// б) создаем нового пользователя и возвращаем jwt токен

	if ok := auth.CheckLogin(login); !ok {
		return "", fmt.Errorf("%w: not valid login", ErrBadAuth)
	}

	if len(password) < minLenPassword {
		return "", fmt.Errorf("%w: very short password", ErrBadAuth)
	}
	password = auth.HashPassword(password)

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
		if errJWT != nil {
			return "", fmt.Errorf("can't create token: %w", errJWT)
		}
		return tokenJWT, nil
	}

	if password == u.Password {
		tokenJWT, errJWT := jwt.GenerateTokenAccess(login)
		if errJWT != nil {
			return "", fmt.Errorf("can't create token: %w", errJWT)
		}
		return tokenJWT, nil
	}

	return "", fmt.Errorf("uncorrect login/password: %w %s", ErrBadPassword, password)
}

// Info
//
// balance : баланс пользователя
//
// merch   : купленный пользователем мерч
//
// from    : переводы пользователя другим людям
//
// to      : переводы на счет пользователя
func (sv *Service) Info(ctx context.Context, sec string) (int, []*entity.Merch, []*entity.User, []*entity.User, error) {
	t, err := jwt.GetInfoFromToken(sec)
	if err != nil {
		return 0, nil, nil, nil, fmt.Errorf("uncorrect security: %w", ErrUnCorrectJWT)
	}
	balance, merch, from, to, errSQL := sv.rep.GetInfo(ctx, t.User)
	if errSQL != nil {
		return 0, nil, nil, nil, fmt.Errorf("can't get correct info: %w", errSQL)
	}
	return balance, merch, from, to, nil
}

func (sv *Service) SendCoin(ctx context.Context, sec, toUs string, amount int) error {
	t, err := jwt.GetInfoFromToken(sec)
	if err != nil {
		return ErrUnCorrectJWT
	}

	if t.User == toUs {
		return ErrLogic
	}

	errTS := sv.rep.SendCoin(ctx, t.User, toUs, amount)

	if errTS != nil {
		return fmt.Errorf("don't send coin: %w", errTS)
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
		return fmt.Errorf("don't buy item: %w", errBY)
	}
	return nil
}
