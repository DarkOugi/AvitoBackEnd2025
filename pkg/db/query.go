package db

import (
	"context"
	"fmt"
)

// получить хэшированный пароль и баланс пользователя

func (db *PostgresDB) GetUserInfo(login string) (string, int, error) {
	var password string
	var balance int

	selectUserInfo := "SELECT password,balance FROM Users WHERE login = $1"
	err := db.Conn.QueryRow(context.Background(), selectUserInfo, login).Scan(&password, &balance)

	return password, balance, err
}

// создать пользователя

func (db *PostgresDB) InitUser(login, password string) error {
	insertUser := "INSERT INTO Users (login,password,balance) VALUES ($1,$2,1000);"
	_, err := db.Conn.Exec(context.Background(), insertUser, login, password)

	return err
}

// получить название мерча и колличество имеющиеся у пользователя

type Merch struct {
	name string
	cnt  int
}

func (db *PostgresDB) getMerchInfo(login string) ([]*Merch, error) {
	infoMerch := "SELECT count(*) AS cnt,merch_id FROM MerchUsers WHERE user_id = $1 GROUP BY merch_id"
	rowsMerch, err := db.Conn.Query(context.Background(), infoMerch, login)

	if err != nil {
		return nil, err
	}
	defer rowsMerch.Close()

	var merch []*Merch
	for rowsMerch.Next() {
		var name string
		var cnt int
		err := rowsMerch.Scan(&name, &cnt)

		if err == nil {
			m := Merch{
				cnt:  cnt,
				name: name,
			}
			merch = append(merch, &m)
		} else {
			return merch, err
		}
	}
	return merch, nil
}

// получить истории трат и зачислений пользователя

type User struct {
	name string
	cost int
	//transaction float64 // в задании не используется, но в целом позитивное поле
}

func (db *PostgresDB) getUserToUserInfo(sql, login string) ([]*User, error) {
	fmt.Printf("LOGIN %s\n", login)
	rows, err := db.Conn.Query(context.Background(), sql, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*User
	for rows.Next() {
		var user string
		var cost int
		//var transaction float64
		err := rows.Scan(&user, &cost)
		if err == nil {
			u := User{
				name: user,
				cost: cost,
			}
			history = append(history, &u)
		} else {
			return history, err
		}
	}

	return history, nil
}

// получить полную информацию о пользователе(история покупок + зачисления + списания + текущий баланс)

func (db *PostgresDB) GetInfo(login string) (int, []*Merch, []*User, []*User, error) {
	var err error

	var balanceInfo int

	var merchInfo []*Merch

	var fromUserInfo []*User
	var toUserInfo []*User

	_, balanceInfo, err = db.GetUserInfo(login)
	if err != nil {
		return balanceInfo, merchInfo, fromUserInfo, toUserInfo, err
	}

	merchInfo, err = db.getMerchInfo(login)
	if err != nil {
		fmt.Printf(" Merch %s", err)
	}

	fromUserInfoQuery := "SELECT to_id,cost FROM UserToUser WHERE from_id = $1"
	fromUserInfo, err = db.getUserToUserInfo(fromUserInfoQuery, login)
	if err != nil {
		fmt.Printf("From User %s", err)
	}

	toUserInfoQuery := "SELECT from_id,cost FROM UserToUser WHERE to_id = $1"
	toUserInfo, err = db.getUserToUserInfo(toUserInfoQuery, login)
	if err != nil {
		fmt.Printf("To User %s", err)
	}

	return balanceInfo, merchInfo, fromUserInfo, toUserInfo, err
}

// БЛОК ЗАПРОСОВ С ТРАНЗАКЦИЯМИ

func (db *PostgresDB) BuyItem(login, merch string) error {
	tx, err := db.Conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	updateBalanceBuyItem := "UPDATE users SET balance = balance - (SELECT cost FROM merch WHERE name = $1) WHERE login = $2;"
	_, err = tx.Exec(context.Background(), updateBalanceBuyItem, merch, login)
	if err != nil {
		return err
	}

	insertHistoryBuyItems := "INSERT INTO MerchUsers (merch_id,user_id,cost) VALUES  ( $1, $2,cost) SELECT cost FROM merch WHERE name = $1;"
	_, err = tx.Exec(context.Background(), insertHistoryBuyItems, merch, login)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (db *PostgresDB) SendCoin(from, to string, cost int) error {
	tx, err := db.Conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	updateFromSendCoin := "UPDATE users SET balance = balance - $1 WHERE login = $2;"
	_, err = tx.Exec(context.Background(), updateFromSendCoin, cost, from)
	if err != nil {
		return err
	}

	updateToSendCoin := "UPDATE users SET balance = balance + $1 WHERE login = $2;"
	_, err = tx.Exec(context.Background(), updateToSendCoin, cost, to)
	if err != nil {
		return err
	}

	insertHistorySendCoin := "INSERT INTO UserToUser (from_id,to_id,cost) VALUES  ( $1, $2,$3)"
	_, err = tx.Exec(context.Background(), insertHistorySendCoin, from, to, cost)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}
