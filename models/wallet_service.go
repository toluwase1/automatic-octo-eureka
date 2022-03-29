package models

import "fmt"

func (wallet *Wallet) CreditUserWallet(money float64, userID string) {
	fmt.Println("old balance",wallet.Balance)
	wallet.UserID = userID
	wallet.Balance += money
	fmt.Println("new balance",wallet.Balance)
}

func (wallet *Wallet) DebitUserWallet(money float64, userID string) {
	wallet.UserID = userID
	wallet.Balance -= money
}



type WalletService interface {
	CreateWallet(user *User) (* User, error)
	GetUserByEmail(email string) ([]* User, error)
	CheckIfPasswordExists(userReference string) ([]* User, error)
	SaveTransaction(t * Transaction) (* Transaction, error)
	PostToAccount(a * Wallet) (* Wallet, error)
	GetAccountBalance(userID string) (*Wallet, error)
	ChangeUserStatus(isActive bool, userReference string) (interface{}, error)
}


type DefaultWalletService struct {
	repo Repository
}

func NewWalletService(repo Repository) *DefaultWalletService {
	return &DefaultWalletService{
		repo: repo,
	}
}

func (user *DefaultWalletService) CreateWallet(userInfo *User) (*User, error) {
	return user.repo.CreateWallet(userInfo)
}

func (user *DefaultWalletService) GetUserByEmail(email string) ([]*User, error) {
	return user.repo.GetUserByEmail(email)
}

func (user *DefaultWalletService) CheckIfPasswordExists(userReference string) ([]*User, error) {
	return user.repo.CheckIfPasswordExists(userReference)
}

func (user *DefaultWalletService) PostToAccount(a *Wallet) (*Wallet, error) {
	return user.repo.PostToAccount(a)
}

func (user *DefaultWalletService) SaveTransaction(t *Transaction) (*Transaction, error) {
	return user.repo.SaveTransaction(t)
}

func (user *DefaultWalletService) GetAccountBalance(userID string) (*Wallet, error) {
	return user.repo.GetAccountBalance(userID)
}

func (user *DefaultWalletService) ChangeUserStatus(isActive bool, userReference string) (interface{}, error) {
	return user.repo.ChangeUserStatus(isActive, userReference)
}
