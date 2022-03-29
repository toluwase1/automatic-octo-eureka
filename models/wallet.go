package models

type Transaction struct {
	UserID               string
	Id 					 string
	Amount               float64
	TransactionType      string
	Password             string
}

type TransactionHistory struct {
	userId 			uint
	AmountCredited float64
	AmountDebited float64
}

type Wallet struct {
	UserID  string
	Balance float64
}


type SupportedCurrencies struct {
	UserID  string
	USDBal float64
	EURBal float64
	NGNBal float64
	GBPBal float64
}
