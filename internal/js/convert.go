//nolint:wrapcheck // тише тише тише
package js

import (
	"avito/internal/entity"
	"encoding/json"
)

// get from JSON
func GetFromJSUser(user []byte) (*GetUser, error) {
	ujs := GetUser{}

	err := json.Unmarshal(user, &ujs)

	return &ujs, err
}

func GetFromJSSecurity(security []byte) (string, error) {
	sjs := GetSecurity{}

	err := json.Unmarshal(security, &sjs)

	return sjs.Security, err
}

func GetFromJsUserToUser(userToUser []byte) (*GetUserToUser, error) {
	utujs := GetUserToUser{}

	err := json.Unmarshal(userToUser, &utujs)

	return &utujs, err
}

// parse to JSON
func ToJSError(err string) ([]byte, error) {
	e := ToError{Errors: err}

	be, errM := json.Marshal(e)

	return be, errM
}
func ToJsToken(token string) ([]byte, error) {
	t := ToToken{Token: token}

	bt, err := json.Marshal(t)

	return bt, err
}

func ToJsMerch(merch []*entity.Merch) []*ToMerch {
	merchJs := []*ToMerch{}
	for _, el := range merch {
		merchJs = append(merchJs, &ToMerch{
			Type:     el.Name,
			Quantity: el.Cnt,
		})
	}

	// bm, err := json.Marshal(merchJs)

	return merchJs
}

func ToJsFromUser(from []*entity.User) []*ToFromUser {
	fromJs := []*ToFromUser{}
	for _, el := range from {
		fromJs = append(fromJs, &ToFromUser{
			ToUser: el.Name,
			Amount: el.Cost,
		})
	}

	// bf, err := json.Marshal(fromJs)

	return fromJs
}

func ToJsToUser(from []*entity.User) []*ToToUser {
	fromJs := []*ToToUser{}
	for _, el := range from {
		fromJs = append(fromJs, &ToToUser{
			FromUser: el.Name,
			Amount:   el.Cost,
		})
	}

	// bf, err := json.Marshal(fromJs)

	return fromJs
}

func ToJsInfo(balance int, merch []*entity.Merch, from, to []*entity.User) ([]byte, error) {
	merchJs := ToJsMerch(merch)

	fromJs := ToJsFromUser(from)
	toJs := ToJsToUser(to)
	coinHistory := ToCoinHistory{
		toJs,
		fromJs,
	}

	infoJs := ToInfo{
		balance,
		merchJs,
		coinHistory,
	}
	bi, err := json.Marshal(infoJs)
	return bi, err
}
