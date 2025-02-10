package db

var selectUser = `
	SELECT password FROM users
	WHERE login = '%s'
`

// нужно возвращать еще и баланс из checkUser чтоб потом проверять могу ли я купить вещь
var buyItem = `
	BEGIN;
	UPDATE users SET balance = balance - (
		SELECT cost FROM merch
		WHERE name = '%s'
	)
    WHERE login = '%s';
	INSERT INTO MerchUsers (merch_id,user_id) 
	VALUES  ( '%s', '%s');
	COMMIT;
`
var sendCoin = `
	BEGIN;
	UPDATE users SET balance = balance - %s
    WHERE login = '%s';
	UPDATE users SET balance = balance + %s
    WHERE login = '%s';
	INSERT INTO UserToUser (from_id,to_id) 
	VALUES  ( '%s', '%s')
	COMMIT;
`
var infoBalance = `
	SELECT balance FROM users
	WHERE login = '%s'
`
var infoMerch = `
	WITH mu AS (
	SELECT count(*) AS cnt,merch_id FROM MerchUsers
	WHERE user_id = '%s'
	GROUP BY merch_id
	)
	SELECT name,cost,cnt FROM merch
	JOIN mu on mu.merch_id = name 
`

//func SelectUser(login string) {
//	fmt.Sprintf(select_user, login)
//}
