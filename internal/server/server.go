package server

import (
	"avito/internal/js"
	"avito/internal/service"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

type Server struct {
	service *service.Service
}

func NewServer(sv *service.Service) *Server {
	return &Server{service: sv}
}

func setError(ctx *fasthttp.RequestCtx, codeErr int, errMsg string) {
	errStr, errConv := js.ToJSError(errMsg)
	if errConv != nil {
		log.Err(errConv).Msg("Error js.ToJSError")
	}
	ctx.SetStatusCode(codeErr)
	ctx.SetBody(errStr)
}

func checkSecurity(ctx *fasthttp.RequestCtx) {

}
func (sv *Server) Auth(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	data := ctx.Request.Body()
	us, err := js.GetFromJSUser(data)

	if err != nil {
		setError(ctx, fasthttp.StatusBadRequest, "Неверный JSON")
	}

	tokenJWT, errAuth := sv.service.Auth(ctx, us.Login, us.Password)
	if errAuth != nil {
		switch {
		case errors.Is(errAuth, service.ErrBadAuth):
			setError(ctx, fasthttp.StatusBadRequest, "Не валидный логин/пароль")
			return
		case errors.Is(errAuth, service.ErrBadPassword):
			setError(ctx, fasthttp.StatusUnauthorized, "Не верный логин/пароль")
			return
		default:
			setError(ctx, fasthttp.StatusInternalServerError, "Ошибка сервера")
			return
		}
	}

	tokenJS, errJS := js.ToJsToken(tokenJWT)
	if errJS != nil {
		setError(ctx, fasthttp.StatusInternalServerError, "Ошибка сервера")
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(tokenJS)

}

func (sv *Server) Info(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	data := ctx.Request.Body()
	sec, err := js.GetFromJSSecurity(data)
	if err != nil {
		setError(ctx, fasthttp.StatusBadRequest, "Неверный JSON")
	}

	balance, merch, from, to, errInfo := sv.service.Info(ctx, sec)
	if errInfo != nil {
		switch {
		case errors.Is(errInfo, service.ErrUnCorrectJWT):
			setError(ctx, fasthttp.StatusUnauthorized, "Не валидный токен")
			return
		default:
			setError(ctx, fasthttp.StatusInternalServerError, "Ошибка сервера")
			return
		}
	}

	jsonData, errJS := js.ToJsInfo(balance, merch, from, to)
	if errJS != nil {
		setError(ctx, fasthttp.StatusInternalServerError, "Ошибка сервера")
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(jsonData)
}

func (sv *Server) SendCoin(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	data := ctx.Request.Body()
	utu, err := js.GetFromJsUserToUser(data)
	if err != nil {
		setError(ctx, fasthttp.StatusBadRequest, "Неверный JSON")
	}

	errSC := sv.service.SendCoin(ctx, utu.Security, utu.ToUser, utu.Amount)
	if errSC != nil {
		switch {
		case errors.Is(errSC, service.ErrUnCorrectJWT):
			setError(ctx, fasthttp.StatusUnauthorized, "Не валидный токен")
			return
		default:
			setError(ctx, fasthttp.StatusInternalServerError, "Ошибка сервера")
			return
		}
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte{})
}

func (sv *Server) BuyItem(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	itemAny := ctx.UserValue("item")
	item, ok := itemAny.(string)
	if ok == false {
		setError(ctx, fasthttp.StatusInternalServerError, "Ошибка сервера")
	}

	data := ctx.Request.Body()
	sec, err := js.GetFromJSSecurity(data)
	if err != nil {
		setError(ctx, fasthttp.StatusBadRequest, "Неверный JSON")
	}

	errBY := sv.service.BuyItem(ctx, item, sec)
	if errBY != nil {
		switch {
		case errors.Is(errBY, service.ErrUnCorrectJWT):
			setError(ctx, fasthttp.StatusUnauthorized, "Не валидный токен")
			return
		default:
			setError(ctx, fasthttp.StatusInternalServerError, "Ошибка сервера")
			return
		}
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte{})

}
