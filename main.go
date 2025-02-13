package main

import (
	auth2 "avito/pkg/auth"
	"avito/pkg/db"
	"avito/pkg/js"
	"avito/pkg/jwt"
	"context"
	"fmt"
	"github.com/fasthttp/router"

	//routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var pSQL *db.PostgresDB

//func handler(ctx *fasthttp.RequestCtx) {
//	switch string(ctx.Path()) {
//	case "/api/info":
//		body := ctx.Request.Body()
//		var to OnlyToken
//		// как будто нужно ограничивать длинну
//
//		err := json.Unmarshal(body, &to)
//		if err != nil {
//			fmt.Println(err.Error())
//		}
//		t, _ := jwt.GetInfoFromToken(to.JWT)
//		balance, merch, from, to1, err1 := pSQL.GetInfo(t.User)
//		if err1 == nil {
//			jsonData, err := js.ToJsInfo(balance, merch, from, to1)
//			//jsonData, err := json.Marshal(data)
//			if err != nil {
//				fmt.Println("Ошибка маршаллинга JSON данных:", err)
//
//			} else {
//				ctx.SetContentType("application/json")
//				ctx.SetStatusCode(fasthttp.StatusOK)
//				ctx.SetBody(jsonData)
//			}
//		}
//
//	case "/api/sendCoin":
//		body := ctx.Request.Body()
//		var to JsUserToUser
//		// как будто нужно ограничивать длинну
//
//		err := json.Unmarshal(body, &to)
//		if err != nil {
//			fmt.Println(err.Error())
//		}
//		t, _ := jwt.GetInfoFromToken(to.JWT)
//		err = pSQL.SendCoin(t.User, to.ToUser, to.Amount)
//		if err == nil {
//			fmt.Println("TRANSACTION SUCCESS")
//		} else {
//			fmt.Printf("%s\n", err)
//		}
//
//	case "/api/buy/{item}":
//		return
//	case "/api/auth":
//		body := ctx.Request.Body()
//		var user JsUser
//		// как будто нужно ограничивать длинну
//
//		err := json.Unmarshal(body, &user)
//		if err != nil {
//			fmt.Println(err.Error())
//		}
//
//		if user.Login != "" && user.Password != "" {
//			pass := db.HashPassword(user.Password)
//			fmt.Printf("HASH PASSWORD %s\n", pass)
//
//			dbPass, _, errC := pSQL.GetUserInfo(user.Login)
//			fmt.Printf("new Pass %s", errC)
//			if dbPass != "" && dbPass == pass {
//				token, err := jwt.GenerateTokenAccess(user.Login)
//				if err == nil {
//					fmt.Println(token)
//				}
//				data := map[string]string{
//					"token": token,
//				}
//				jsonData, err := json.Marshal(data)
//				if err != nil {
//					fmt.Println("Ошибка маршаллинга JSON данных:", err)
//
//				} else {
//					ctx.SetContentType("application/json")
//					ctx.SetStatusCode(fasthttp.StatusOK)
//					ctx.SetBody(jsonData)
//				}
//			} else if dbPass == "" { // сюда более умную проверку
//				err := pSQL.InitUser(user.Login, pass)
//				if err == nil {
//					fmt.Println("User Create")
//				}
//			}
//		}
//
//	default:
//		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
//	}
//}

func info(ctx *fasthttp.RequestCtx) {
	body := ctx.Request.Body()
	token, err := js.GetFromJSSecurity(body)
	if err != nil {
		fmt.Println(err.Error())
	}
	t, _ := jwt.GetInfoFromToken(token)
	balance, merch, from, to1, _ := pSQL.GetInfo(t.User)
	jsonData, err := js.ToJsInfo(balance, merch, from, to1)
	//jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Ошибка маршаллинга JSON данных:", err)

	} else {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBody(jsonData)
	}
}
func buyItem(ctx *fasthttp.RequestCtx) {
	itemAny := ctx.UserValue("item")
	item, ok := itemAny.(string)
	if ok == false {
		fmt.Println("ERR")
	}
	body := ctx.Request.Body()
	token, err := js.GetFromJSSecurity(body)
	if err != nil {
		fmt.Println(err.Error())
	}
	t, _ := jwt.GetInfoFromToken(token)
	err = pSQL.BuyItem(t.User, item)
	if err == nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)
	}

}
func auth(ctx *fasthttp.RequestCtx) {
	body := ctx.Request.Body()
	us, err := js.GetFromJSUser(body)
	if err != nil {
		fmt.Println(err.Error())
	}
	if flag := auth2.CheckLogin(us.Login); flag {
		if us.Password != "" {
			pass := auth2.HashPassword(us.Password)
			passSql, _, _ := pSQL.GetUserInfo(us.Login)

			if pass == passSql {
				jwtlog, _ := jwt.GenerateTokenAccess(us.Login)
				jsToken, _ := js.ToJsToken(jwtlog)
				ctx.SetContentType("application/json")
				ctx.SetStatusCode(fasthttp.StatusOK)
				ctx.SetBody(jsToken)

			} else {
				cerr := pSQL.InitUser(us.Login, pass)
				if cerr == nil {
					jwtlog, _ := jwt.GenerateTokenAccess(us.Login)
					jsToken, _ := js.ToJsToken(jwtlog)
					ctx.SetContentType("application/json")
					ctx.SetStatusCode(fasthttp.StatusOK)
					ctx.SetBody(jsToken)
				}
			}
		}
	}
}
func sendCoin(ctx *fasthttp.RequestCtx) {
	body := ctx.Request.Body()
	utu, err := js.GetFromJsUserToUser(body)
	if err != nil {
		fmt.Println(err.Error())
	}
	t, _ := jwt.GetInfoFromToken(utu.Security)
	err = pSQL.SendCoin(t.User, utu.ToUser, utu.Amount)
	if err == nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)
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

	r := router.New()
	r.POST("/api/auth", auth)
	r.GET("/api/buy/{item}", buyItem)
	r.POST("/api/sendCoin", sendCoin)
	r.GET("/api/info", info)

	if err := fasthttp.ListenAndServe(":8080", r.Handler); err != nil {
		panic(err)
	}
	//router := routing.New()
	//decode := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiYmFkZGF5In0.wr_9TpqhbB9ASJ01UbaXirypTHkEuZHta-15YoaR2Xg"
	//t, _ := jwt.GetInfoFromToken(decode)
	//fmt.Printf(t.User)

}
