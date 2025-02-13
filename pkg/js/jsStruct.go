package js

// get JSON

type GetUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetUserToUser struct {
	Security string `json:"security"`
	ToUser   string `json:"toUser"`
	Amount   int    `json:"amount"`
}

type GetSecurity struct {
	Security string `json:"security"`
}

// parse to JSON

type ToToken struct {
	Token string `json:"token"`
}

type ToMerch struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type ToFromUser struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type ToToUser struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type ToCoinHistory struct {
	Received []*ToToUser   `json:"received"`
	Sent     []*ToFromUser `json:"sent"`
}

type ToInfo struct {
	Coins       int           `json:"coins"`
	Inventory   []*ToMerch    `json:"inventory"`
	CoinHistory ToCoinHistory `json:"coinHistory"`
}
