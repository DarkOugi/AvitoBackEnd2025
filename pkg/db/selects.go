package db

import (
	"context"
	"github.com/jackc/pgx/v5"
)

// получить хэшированный пароль и баланс пользователя
var selectUserInfo = `
	SELECT password,balance FROM Users
	WHERE login = $1
`

func GetUserInfo(db *pgx.Conn, login string) (string, int, error) {
	var password string
	var balance int
	err := db.QueryRow(context.Background(), selectUserInfo, login).Scan(&password, &balance)

	return password, balance, err
}

// осздать пользователя

func InitUser(db *pgx.Conn, login, password string) error {
	insertUser := "INSERT INTO Users (login,password,balance) VALUES ($1,$2,1000);"
	_, err := db.Exec(context.Background(), insertUser, login, password)

	return err
}

// получить название мерча его прайс и колличество
var infoMerch = `
	SELECT count(*) AS cnt,merch_id FROM MerchUsers
	WHERE user_id = $1
	GROUP BY merch_id
`

type Merch struct {
	name string
	cnt  int
}

func getMerchInfo(db *pgx.Conn, login string) []*Merch {
	rowsMerch, err := db.Query(context.Background(), infoMerch, login)
	if err != nil {
		return nil
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
			//merchCostCount := map[string]int{
			//	"cnt":  cnt,
			//	"cost": cost,
			//}
			//merchInfo[name] = merchCostCount
			merch = append(merch, &m)
		}
	}
	return merch
}

// получить истории трат и зачислений пользователя
var fromUserInfo = `
	SELECT to_id,cost,transactionTime FROM UserToUser
	WHERE from_id = $1
`

var toUserInfo = `
	SELECT to_id,cost,transactionTime FROM UserToUser
	WHERE to_id = $1
`

type User struct {
	name string
	cost int
	//transaction float64 // тут скорее всего будет тип время
}

func getUserToUserInfo(db *pgx.Conn, sql, login string) []*User {
	rows, err := db.Query(context.Background(), sql, login)
	if err != nil {
		return nil
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
		}
	}
	return history
}

// получить полную информацию о пользователе(история покупок + зачисления + списания + текущий баланс)
func GetInfo(db *pgx.Conn, login string) ([]*Merch, int, []*User, []*User) {
	merchInfo := getMerchInfo(db, login)        //right
	_, balanceInfo, _ := GetUserInfo(db, login) //right
	fromUserInfo := getUserToUserInfo(db, fromUserInfo, login)
	toUserInfo := getUserToUserInfo(db, toUserInfo, login)
	return merchInfo, balanceInfo, fromUserInfo, toUserInfo
}

// БЛОК ЗАПРОСОВ С ТРАНЗАКЦИЯМИ

// var buyItemUpdate = `
//
//		BEGIN;
//		UPDATE users SET balance = balance - (
//			SELECT cost FROM merch
//			WHERE name = '%s'
//		)
//	   WHERE login = '%s';
//		INSERT INTO MerchUsers (merch_id,user_id,cost)
//		VALUES  ( '%s',
//	             '%s',
//				  cost
//				)
//		SELECT cost FROM merch
//	   WHERE name = '%s'
//		;
//		COMMIT;
//
// `
var updateBalaceBuyItme = `
	UPDATE users SET balance = balance - (
			SELECT cost FROM merch
			WHERE name = $1
	)
	WHERE login = $2;
`
var insertHistroyBuyItems = `
	INSERT INTO MerchUsers (merch_id,user_id,cost) 
	VALUES  ( $1, 
              $2,
			  cost
			)
	SELECT cost FROM merch
    WHERE name = $1
	;
`

func BuyItem(db *pgx.Conn, login, merch string) error {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), updateBalaceBuyItme, merch, login)
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), insertHistroyBuyItems, merch, login)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

// var sendCoin = `
//
//		BEGIN;
//		UPDATE users SET balance = balance - %s
//	   WHERE login = '%s';
//		UPDATE users SET balance = balance + %s
//	   WHERE login = '%s';
//		INSERT INTO UserToUser (from_id,to_id,cost)
//		VALUES  ( '%s', '%s',%s)
//		COMMIT;
//
// `
var updateFromSendCoin = `
	UPDATE users SET balance = balance - $1
    WHERE login = $2;
`
var updateToSendCoin = `
	UPDATE users SET balance = balance + $1
    WHERE login = $2;
`
var insertHistorySendCoin = `
	INSERT INTO UserToUser (from_id,to_id,cost)
	VALUES  ( $1, $2,$3)
`

func SendCoin(db *pgx.Conn, from, to string, cost int) error {
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	_, err = tx.Exec(context.Background(), updateFromSendCoin, cost, from)
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), updateToSendCoin, cost, to)
	if err != nil {
		return err
	}
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
