package main

import (
	initdb "avito/initDB"
	"avito/pkg/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/valyala/fasthttp"
)

type JsUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var conn *pgx.Conn

// для инита базы данных
//func main() {
//	db, err := pg.Connect("localhost", "5432", "avito", "0000", "avitodb")
//	if err != nil {
//		fmt.Printf("Connect to db %s\n", err.Error())
//	}
//	err = initDB.CreateTables(db)
//	if err != nil {
//		fmt.Printf("Error create table  %s\n", err.Error())
//	}
//	err = initDB.CreateIndex(db)
//	if err != nil {
//		fmt.Printf("Error create index  %s\n", err.Error())
//	}
//	err = initDB.InsertValue(db)
//	if err != nil {
//		fmt.Printf("Error insert to table  %s\n", err.Error())
//	}
//
//}

//	func main() {
//		// create, insert := initdb.Init_mech()
//		// db, err := pg.Connect("localhost", "5432", "avito", "0000", "avitodb")
//		// if err != nil {
//		// 	fmt.Printf("Connect to db %s\n", err.Error())
//		// } else {
//		// 	_, err = db.Exec(create)
//		// 	if err != nil {
//		// 		fmt.Printf("Create table err: %s\n", err.Error())
//		// 	} else {
//		// 		fmt.Printf("Create table success\n")
//		// 	}
//		// 	in, err := db.Exec(insert)
//		// 	if err != nil {
//		// 		fmt.Printf("Create rows err: %s\n", err.Error())
//		// 		panic(err)
//		// 	} else {
//		// 		fmt.Printf("Insert table %d rows success\n", in)
//		// 	}
//		// 	defer db.Close()
//		// }
//		// q := token.MetaJWT("Denis.Zhilin")
//		// fmt.Println(q)
//		fmt.Println(pg.HashPassword("asdadqwq212e1d2wd"))
//
// }
func handler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/":
		ctx.SetContentType("text/plain; charset=utf-8")
		ctx.SetBody([]byte("Welcome to the fasthttp server"))
	case "/api/info":
		return
	case "/api/sendCoin":
		return
	case "/api/buy/{item}":
		return
	case "/api/auth":
		body := ctx.Request.Body()
		var user JsUser
		// как будто нужно ограничивать длинну
		if len(body) > 512 {
			fmt.Println("Слишком большой запрос")
		}

		err := json.Unmarshal(body, &user)
		if err != nil {
			fmt.Println(err.Error())
		}

		if user.Login != "" && user.Password != "" {
			pass := db.HashPassword(user.Password)
			fmt.Printf("HASH PASSWORD %s\n", pass)

			dbPass, _, errC := db.GetUserInfo(conn, user.Login)
			fmt.Printf("new Pass %s", errC)
			if dbPass != "" && dbPass == pass {
				fmt.Println("YES YES YES")
			} else if dbPass == "" { // сюда более умную проверку
				err := db.InitUser(conn, user.Login, pass)
				if err == nil {
					fmt.Println("User Create")
				}
			}
		}

	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

func main() {
	var err error
	if conn == nil {

		conn, err = db.Connect("localhost", "5432", "avito", "0000", "avitodb")
		if err != nil {
			panic(err)
		}
	}
	defer conn.Close(context.Background())
	err = initdb.CreateTables(conn)
	if err != nil {
		fmt.Printf("Error create table  %s\n", err.Error())
	}
	if err := fasthttp.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
