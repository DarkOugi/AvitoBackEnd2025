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

func setError(ctx *fasthttp.RequestCtx, err error, codeErr int, errMsg string) {
	errStr, errConv := js.ToJSError(errMsg)
	if errConv != nil {
		log.Err(errConv).Msg("Error js.ToJSError")
	}
	log.Info().Err(err).Msgf("Error")
	ctx.SetStatusCode(codeErr)
	ctx.SetBody(errStr)
}

func (sv *Server) Auth(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	data := ctx.Request.Body()
	us, err := js.GetFromJSUser(data)
	if err != nil {
		setError(ctx, err, fasthttp.StatusBadRequest, "Неверный JSON")
		return
	}

	tokenJWT, errAuth := sv.service.Auth(ctx, us.Login, us.Password)
	if errAuth != nil {
		switch {
		case errors.Is(errAuth, service.ErrBadAuth):
			setError(ctx, errAuth, fasthttp.StatusBadRequest, "Невалидный логин/пароль")
			return
		case errors.Is(errAuth, service.ErrBadPassword):
			setError(ctx, errAuth, fasthttp.StatusUnauthorized, "Неверный логин/пароль")
			return
		default:
			setError(ctx, errAuth, fasthttp.StatusInternalServerError, "Ошибка сервера")
			return
		}
	}

	tokenJS, errJS := js.ToJsToken(tokenJWT)
	if errJS != nil {
		setError(ctx, errJS, fasthttp.StatusInternalServerError, "Ошибка сервера")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(tokenJS)
}

func (sv *Server) Info(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	data := ctx.Request.Body()
	sec, err := js.GetFromJSSecurity(data)
	if err != nil {
		setError(ctx, err, fasthttp.StatusBadRequest, "Неверный JSON")
		return
	}

	balance, merch, from, to, errInfo := sv.service.Info(ctx, sec)
	if errInfo != nil {
		switch {
		case errors.Is(errInfo, service.ErrUnCorrectJWT):
			setError(ctx, errInfo, fasthttp.StatusUnauthorized, "Невалидный токен")
			return
		default:
			setError(ctx, errInfo, fasthttp.StatusInternalServerError, "Ошибка сервера")
			return
		}
	}

	jsonData, errJS := js.ToJsInfo(balance, merch, from, to)
	if errJS != nil {
		setError(ctx, errJS, fasthttp.StatusInternalServerError, "Ошибка сервера")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(jsonData)
}

func (sv *Server) SendCoin(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")

	data := ctx.Request.Body()
	utu, err := js.GetFromJsUserToUser(data)
	if err != nil {
		setError(ctx, err, fasthttp.StatusBadRequest, "Неверный JSON")
		return
	}

	errSC := sv.service.SendCoin(ctx, utu.Security, utu.ToUser, utu.Amount)
	if errSC != nil {
		switch {
		case errors.Is(errSC, service.ErrUnCorrectJWT):
			setError(ctx, errSC, fasthttp.StatusUnauthorized, "Невалидный токен")
			return
		case errors.Is(errSC, service.ErrLogic):
			setError(ctx, errSC, fasthttp.StatusBadRequest, "Отправитель не может быть получателем в рамках 1 перевода")
			return
		default:
			setError(ctx, errSC, fasthttp.StatusInternalServerError, "Ошибка сервера")
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
	if !ok {
		setError(ctx, nil, fasthttp.StatusInternalServerError, "Ошибка сервера")
		return
	}

	data := ctx.Request.Body()
	sec, err := js.GetFromJSSecurity(data)
	if err != nil {
		setError(ctx, err, fasthttp.StatusBadRequest, "Неверный JSON")
		return
	}

	errBY := sv.service.BuyItem(ctx, item, sec)
	if errBY != nil {
		switch {
		case errors.Is(errBY, service.ErrUnCorrectJWT):
			setError(ctx, errBY, fasthttp.StatusUnauthorized, "Невалидный токен")
			return
		default:
			setError(ctx, errBY, fasthttp.StatusInternalServerError, "Ошибка сервера")
			return
		}
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody([]byte{})
}
