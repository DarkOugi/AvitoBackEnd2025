package db

import (
	"avito/internal/entity"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
)

// получить хэшированный пароль и баланс пользователя

func (db *PostgresDB) GetUserInfo(ctx context.Context, login string) (*entity.User, bool, error) {
	var u entity.User

	selectUserInfo := "SELECT password,balance FROM Users WHERE login = $1"
	err := db.conn.QueryRow(ctx, selectUserInfo, login).Scan(&u.Password, &u.Cost)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &u, false, nil
		}
		return &u, false, err
	}

	return &u, true, nil
}

// создать пользователя

func (db *PostgresDB) InitUser(ctx context.Context, login, password string) error {
	insertUser := "INSERT INTO Users (login,password,balance) VALUES ($1,$2,1000);"
	_, err := db.conn.Exec(ctx, insertUser, login, password)

	return err
}

// получить название мерча и колличество имеющиеся у пользователя

func (db *PostgresDB) getMerchInfo(ctx context.Context, login string) ([]*entity.Merch, error) {
	infoMerch := "SELECT count(*) AS cnt,merch_id FROM MerchUsers WHERE user_id = $1 GROUP BY merch_id"
	rowsMerch, err := db.conn.Query(ctx, infoMerch, login)

	if err != nil {
		return nil, err
	}
	defer rowsMerch.Close()

	var merch []*entity.Merch
	for rowsMerch.Next() {
		var name string
		var cnt int

		err = rowsMerch.Scan(&cnt, &name)
		if err != nil {
			return merch, err
		}

		m := entity.Merch{
			Cnt:  cnt,
			Name: name,
		}
		merch = append(merch, &m)
	}
	return merch, nil
}

// получить истории трат и зачислений пользователя

func (db *PostgresDB) getUserToUserInfo(ctx context.Context, sql, login string) ([]*entity.User, error) {
	rows, err := db.conn.Query(ctx, sql, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*entity.User
	for rows.Next() {
		var user string
		var cost int
		//var transaction float64
		err = rows.Scan(&user, &cost)
		if err != nil {
			return history, err

		}

		u := entity.User{
			Name: user,
			Cost: cost,
		}
		history = append(history, &u)
	}

	return history, nil
}

// получить полную информацию о пользователе(история покупок + зачисления + списания + текущий баланс)

func (db *PostgresDB) GetInfo(ctx context.Context, login string) (int, []*entity.Merch, []*entity.User, []*entity.User, error) {
	var (
		err error

		balanceInfo *entity.User

		merchInfo []*entity.Merch

		fromUserInfo []*entity.User
		toUserInfo   []*entity.User
	)

	balanceInfo, ok, err := db.GetUserInfo(ctx, login)
	if err != nil && !ok {
		return balanceInfo.Cost, merchInfo, fromUserInfo, toUserInfo, fmt.Errorf("can't get info from table user: %w", err)
	}

	merchInfo, err = db.getMerchInfo(ctx, login)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return balanceInfo.Cost, merchInfo, fromUserInfo, toUserInfo, fmt.Errorf("can't get info from table merchUser: %w", err)
	}

	fromUserInfoQuery := "SELECT to_id,cost FROM UserToUser WHERE from_id = $1"
	fromUserInfo, err = db.getUserToUserInfo(ctx, fromUserInfoQuery, login)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return balanceInfo.Cost, merchInfo, fromUserInfo, toUserInfo, fmt.Errorf("can't get info from table UserToUser: %w", err)
	}

	toUserInfoQuery := "SELECT from_id,cost FROM UserToUser WHERE to_id = $1"
	toUserInfo, err = db.getUserToUserInfo(ctx, toUserInfoQuery, login)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return balanceInfo.Cost, merchInfo, fromUserInfo, toUserInfo, fmt.Errorf("can't get info from table UserToUser: %w", err)
	}

	return balanceInfo.Cost, merchInfo, fromUserInfo, toUserInfo, err
}

// БЛОК ЗАПРОСОВ С ТРАНЗАКЦИЯМИ

func (db *PostgresDB) BuyItem(ctx context.Context, login, merch string) (err error) {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("can't start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	updateBalanceBuyItem := "UPDATE users SET balance = balance - (SELECT cost FROM Merch WHERE name = $1) WHERE login = $2;"
	_, err = tx.Exec(ctx, updateBalanceBuyItem, merch, login)
	if err != nil {
		return fmt.Errorf("can't update table user: %w", err)
	}

	insertHistoryBuyItems := "INSERT INTO merchusers (merch_id,user_id,cost) VALUES  ( $1, $2,(SELECT cost FROM merch WHERE name = $3));"
	_, err = tx.Exec(ctx, insertHistoryBuyItems, merch, login, merch)
	if err != nil {
		return fmt.Errorf("can't update table merchusers: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("can't complete transaction: %w", err)
	}

	return nil
}

func (db *PostgresDB) SendCoin(ctx context.Context, from, to string, cost int) (err error) {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("can't begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	updateFromSendCoin := "UPDATE users SET balance = balance - $1 WHERE login = $2;"
	_, err = tx.Exec(ctx, updateFromSendCoin, cost, from)
	if err != nil {
		return fmt.Errorf("can't update table users(From): %w", err)
	}

	updateToSendCoin := "UPDATE users SET balance = balance + $1 WHERE login = $2;"
	_, err = tx.Exec(ctx, updateToSendCoin, cost, to)
	if err != nil {
		return fmt.Errorf("can't update table users(To): %w", err)
	}

	insertHistorySendCoin := "INSERT INTO UserToUser (from_id,to_id,cost) VALUES  ( $1, $2,$3)"
	_, err = tx.Exec(ctx, insertHistorySendCoin, from, to, cost)
	if err != nil {
		return fmt.Errorf("can't update table UserToUser: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("can't complete transaction: %w", err)
	}
	return nil
}
