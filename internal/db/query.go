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
		return &u, false, fmt.Errorf("select user error: %w", err)
	}

	return &u, true, nil
}

// создать пользователя

func (db *PostgresDB) InitUser(ctx context.Context, login, password string) error {
	insertUser := "INSERT INTO Users (login,password,balance) VALUES ($1,$2,1000);"
	_, err := db.conn.Exec(ctx, insertUser, login, password)

	if err != nil {
		return fmt.Errorf("error create user: %w", err)
	}
	return nil
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
func (db *PostgresDB) GetInfo(ctx context.Context, login string) (int,
	[]*entity.Merch, []*entity.User, []*entity.User, error) {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return 0, nil, nil, nil, fmt.Errorf("can't start transaction: %w", err)
	}
	//nolint:errcheck // тише тише тише
	defer tx.Rollback(ctx)

	balanceInfo := 0
	balanceSQL := "SELECT balance FROM Users WHERE login = $1;"
	errBalanceSQL := tx.QueryRow(ctx, balanceSQL, login).Scan(&balanceInfo)
	if errBalanceSQL != nil {
		return 0, nil, nil, nil, fmt.Errorf("can't get info balance: %w", errBalanceSQL)
	}

	merchInfo := []*entity.Merch{}
	merchSQL := "SELECT count(*) AS cnt,merch_id FROM MerchUsers WHERE user_id = $1 GROUP BY merch_id;"
	merchRows, errMerchSQL := tx.Query(ctx, merchSQL, login)
	if errMerchSQL != nil {
		return 0, nil, nil, nil, fmt.Errorf("can't select from merch: %w", errMerchSQL)
	}

	for merchRows.Next() {
		var marchName string
		var merchCnt int

		errScan := merchRows.Scan(&merchCnt, &marchName)
		if errScan != nil {
			return 0, nil, nil, nil, fmt.Errorf("error scan merch info: %w", errScan)
		}

		merchInfo = append(merchInfo, &entity.Merch{
			Name: marchName,
			Cnt:  merchCnt,
		})
	}

	toInfo := []*entity.User{}
	toSQL := "SELECT from_id,cost FROM usertouser WHERE to_id = $1;"
	toRows, errToSQL := tx.Query(ctx, toSQL, login)
	if errToSQL != nil {
		return 0, nil, nil, nil, fmt.Errorf("can't select to_id: %w", errToSQL)
	}

	for toRows.Next() {
		var from string
		var cost int

		errScan := toRows.Scan(&from, &cost)
		if errScan != nil {
			return 0, nil, nil, nil, fmt.Errorf("error scan to info: %w", errScan)
		}

		toInfo = append(toInfo, &entity.User{
			Name:     from,
			Cost:     cost,
			Password: "",
		})
	}

	fromInfo := []*entity.User{}
	fromSQL := "SELECT to_id,cost FROM usertouser WHERE from_id = $1;"
	fromRows, errFromSQL := tx.Query(ctx, fromSQL, login)
	if errFromSQL != nil {
		return balanceInfo, nil, nil, nil, fmt.Errorf("can't select from_id: %w", errFromSQL)
	}

	for fromRows.Next() {
		var to string
		var cost int

		errScan := fromRows.Scan(&to, &cost)
		if errScan != nil {
			return 0, nil, nil, nil, fmt.Errorf("error scan from info: %w", errScan)
		}

		fromInfo = append(fromInfo, &entity.User{
			Name:     to,
			Cost:     cost,
			Password: "",
		})
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, nil, nil, nil, fmt.Errorf("can't complete transaction: %w", err)
	}

	return balanceInfo, merchInfo, fromInfo, toInfo, nil
}

// БЛОК ЗАПРОСОВ С ТРАНЗАКЦИЯМИ

func (db *PostgresDB) BuyItem(ctx context.Context, login, merch string) (err error) {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("can't start transaction: %w", err)
	}
	//nolint:errcheck // тише тише тише
	defer tx.Rollback(ctx)

	updateBalanceBuyItem := `
		UPDATE users 
		SET balance = balance - (SELECT cost FROM Merch WHERE name = $1) 
		WHERE login = $2;
	`
	_, err = tx.Exec(ctx, updateBalanceBuyItem, merch, login)
	if err != nil {
		return fmt.Errorf("can't update table user: %w", err)
	}

	insertHistoryBuyItems := `
		INSERT INTO merchusers (merch_id,user_id,cost) 
		VALUES  ( $1, $2,(SELECT cost FROM merch WHERE name = $3));
	`
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
	//nolint:errcheck // тише тише тише
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

	insertHistorySendCoin := "INSERT INTO UserToUser (from_id,to_id,cost) VALUES  ( $1, $2,$3);"
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
