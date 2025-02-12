package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

//type Repository interface {
//	GetUserInfo(login string) (string, int, error)
//	InitUser(login, password string) error
//	getMerchInfo(login string) ([]*Merch, error)
//	getUserToUserInfo(sql, login string) ([]*User, error)
//	GetInfo(login string) (int, []*Merch, []*User, []*User, error)
//	BuyItem(login, merch string) error
//	SendCoin(from, to string, cost int) error
//}

type PostgresDB struct {
	Conn *pgx.Conn
}

func GetConnect(host, port, username, password, dbname string) (*PostgresDB, error) {
	strСonn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbname)
	conn, err := pgx.Connect(context.Background(), strСonn)
	if err == nil {
		return &PostgresDB{Conn: conn}, nil
	}

	return nil, err
}
