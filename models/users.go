package models

type User struct {

	FirstName       string
	LastName        string
	Email           string
	Currency        string
	Password        string
	HashedSecretKey string
	CreatedAt       string
	IsActive        bool
	currencies SupportedCurrencies
}



type SupportedCurrencies struct {
	USDbal float64
	EURbal float64
	NGNbal float64
	GBPbal float64
}