package db

import (
	"context"
	"github.com/jackc/pgx/v5"
)

// получить хэшированный пароль и баланс пользователя

func GetUserInfo(db *pgx.Conn, login string) (string, int, error) {
	var password string
	var balance int

	selectUserInfo := "SELECT password,balance FROM Users WHERE login = $1"
	err := db.QueryRow(context.Background(), selectUserInfo, login).Scan(&password, &balance)

	return password, balance, err
}

// создать пользователя

func InitUser(db *pgx.Conn, login, password string) error {
	insertUser := "INSERT INTO Users (login,password,balance) VALUES ($1,$2,1000);"
	_, err := db.Exec(context.Background(), insertUser, login, password)

	return err
}

// получить название мерча и колличество имеющиеся у пользователя

type Merch struct {
	name string
	cnt  int
}

func getMerchInfo(db *pgx.Conn, login string) ([]*Merch, error) {
	infoMerch := "SELECT count(*) AS cnt,merch_id FROM MerchUsers WHERE user_id = $1 GROUP BY merch_id"
	rowsMerch, err := db.Query(context.Background(), infoMerch, login)

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
	//transaction float64 // тут скорее всего будет тип время
}

func getUserToUserInfo(db *pgx.Conn, sql, login string) ([]*User, error) {
	rows, err := db.Query(context.Background(), sql, login)
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

func GetInfo(db *pgx.Conn, login string) ([]*Merch, int, []*User, []*User) {
	merchInfo, _ := getMerchInfo(db, login) //right

	_, balanceInfo, _ := GetUserInfo(db, login) //right

	fromUserInfoQuery := "SELECT to_id,cost,transactionTime FROM UserToUser WHERE from_id = $1"
	fromUserInfo, _ := getUserToUserInfo(db, fromUserInfoQuery, login)

	toUserInfoQuery := "SELECT to_id,cost,transactionTime FROM UserToUser WHERE to_id = $1"
	toUserInfo, _ := getUserToUserInfo(db, toUserInfoQuery, login)

	return merchInfo, balanceInfo, fromUserInfo, toUserInfo
}

// БЛОК ЗАПРОСОВ С ТРАНЗАКЦИЯМИ

func BuyItem(db *pgx.Conn, login, merch string) error {
	tx, err := db.Begin(context.Background())
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

func SendCoin(db *pgx.Conn, from, to string, cost int) error {
	tx, err := db.Begin(context.Background())
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
