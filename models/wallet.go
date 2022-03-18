package models

type Transaction struct {
	UserID               string
	TransactionReference string
	Amount               float64
	PhoneNumber          string
	Password             string
}

// Wallet struct
type Wallet struct {
	UserID  string
	Balance float64
}

// CreditUserWallet Credits the user data
func (w *Wallet) CreditUserWallet(money float64, userID string) {
	w.UserID = userID
	w.Balance += money
}

// DebitUserWallet Debits the user data
func (w *Wallet) DebitUserWallet(money float64, userID string) {
	w.UserID = userID
	w.Balance -= money
}
