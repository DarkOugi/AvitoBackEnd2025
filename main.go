package main

import (
	"avito/pkg/db"
	"avito/pkg/jwt"
	"context"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
)

type JsUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type JsUserToUser struct {
	JWT    string `json:"security"`
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
type OnlyToken struct {
	JWT string `json:"security"`
}

//type JsInfo struct {
//	balance     int
//	merch       []db.Merch
//	coinHistory map[string]*JsUserInfo
//}
//
//type JsUserInfo struct {
//}

var pSQL *db.PostgresDB

func handler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/api/info":
		body := ctx.Request.Body()
		var to OnlyToken
		// как будто нужно ограничивать длинну
		if len(body) > 512 {
			fmt.Println("Слишком большой запрос")
		}

		err := json.Unmarshal(body, &to)
		if err != nil {
			fmt.Println(err.Error())
		}
		t, _ := jwt.GetInfoFromToken(to.JWT)
		balance, merch, from, to1, err1 := pSQL.GetInfo(t.User)
		if err1 == nil {
			fmt.Println(balance)
			fmt.Println(merch)
			fmt.Println(from)
			fmt.Println(to1)
		}

	case "/api/sendCoin":
		body := ctx.Request.Body()
		var to JsUserToUser
		// как будто нужно ограничивать длинну
		if len(body) > 512 {
			fmt.Println("Слишком большой запрос")
		}

		err := json.Unmarshal(body, &to)
		if err != nil {
			fmt.Println(err.Error())
		}
		t, _ := jwt.GetInfoFromToken(to.JWT)
		err = pSQL.SendCoin(t.User, to.ToUser, to.Amount)
		if err == nil {
			fmt.Println("TRANSACTION SUCCESS")
		} else {
			fmt.Printf("%s\n", err)
		}

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

			dbPass, _, errC := pSQL.GetUserInfo(user.Login)
			fmt.Printf("new Pass %s", errC)
			if dbPass != "" && dbPass == pass {
				token, err := jwt.GenerateTokenAccess(user.Login)
				if err == nil {
					fmt.Println(token)
				}
			} else if dbPass == "" { // сюда более умную проверку
				err := pSQL.InitUser(user.Login, pass)
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
	if pSQL == nil {
		pSQL, err = db.GetConnect("localhost", "5432", "avito", "0000", "avitodb")
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
	defer pSQL.Conn.Close(context.Background())
	if err != nil {
		fmt.Printf("Error create table  %s\n", err.Error())
	}
	if err := fasthttp.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
	//decode := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiYmFkZGF5In0.wr_9TpqhbB9ASJ01UbaXirypTHkEuZHta-15YoaR2Xg"
	//t, _ := jwt.GetInfoFromToken(decode)
	//fmt.Printf(t.User)

}
