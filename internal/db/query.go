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

// получить полную информацию о пользователе(история покупок + зачисления + списания + текущий баланс)

func (db *PostgresDB) GetInfo(ctx context.Context, login string) (int, []*entity.Merch, []*entity.User, []*entity.User, error) {
	var (
		balanceInfo int // баланс пользователя

		merchInfo []*entity.Merch // купленный пользователем мерч

		fromUserInfo []*entity.User // отправка денег кому-то со счета пользователя
		toUserInfo   []*entity.User // зачисления на счет пользователя от кого-то
	)
	infoSelect := "SELECT COALESCE(us.balance,0),COALESCE(m.cnt,0),COALESCE(m.merch_id,''),COALESCE(toU.from_id,''),COALESCE(toU.cost,0),COALESCE(fromU.to_id,''),COALESCE(fromU.cost,0) FROM Users us\n"
	infoMerch := "LEFT JOIN (SELECT count(*) AS cnt,merch_id,user_id FROM MerchUsers WHERE user_id = $1 GROUP BY merch_id,user_id) AS m ON m.user_id = us.login\n"
	infoJoinTo := "LEFT JOIN UserToUser AS toU ON toU.to_id = us.login\n"
	infoJoinFrom := "LEFT JOIN UserToUser AS fromU ON fromU.from_id = us.login\n"
	infoWhere := "WHERE us.login = $2;"
	fullInfo := infoSelect + infoMerch + infoJoinFrom + infoJoinTo + infoWhere

	rows, err := db.conn.Query(ctx, fullInfo, login, login)
	if err != nil {
		return 0, nil, nil, nil, err
	}
	defer rows.Close()

	uniqMerchInfo := map[string]int{}
	for rows.Next() {

		var balance int
		var cntMerchInfo int
		var nameMerch string
		var giveName string
		var giveCost int
		var sentName string
		var sentCost int

		errScan := rows.Scan(&balance, &cntMerchInfo, &nameMerch, &giveName, &giveCost, &sentName, &sentCost)
		if errScan != nil {
			return balance, merchInfo, toUserInfo, fromUserInfo, errScan
		}

		balanceInfo = balance

		if nameMerch != "" {
			if _, ok := uniqMerchInfo[nameMerch]; !ok {
				uniqMerchInfo[nameMerch] = 1

				m := entity.Merch{
					Cnt:  cntMerchInfo,
					Name: nameMerch,
				}
				merchInfo = append(merchInfo, &m)
			}
		}

		if giveName != "" {
			fromU := entity.User{
				Name:     giveName,
				Password: "",
				Cost:     giveCost,
			}
			toUserInfo = append(toUserInfo, &fromU)
		}

		if sentName != "" {
			toU := entity.User{
				Name:     sentName,
				Password: "",
				Cost:     sentCost,
			}
			fromUserInfo = append(fromUserInfo, &toU)
		}
	}

	return balanceInfo, merchInfo, fromUserInfo, toUserInfo, nil
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
