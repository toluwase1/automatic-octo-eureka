package models

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"strconv"
)

type Repository interface {
	CreateWallet(user *User) (*User, error)
	GetUserByEmail(email string) ([]*User, error)
	CheckIfPasswordExists(userReference string) ([]*User, error)
	PostToAccount(user *Wallet) (*Wallet, error)
	SaveTransaction(transaction *Transaction) (*Transaction, error)
	GetAccountBalance(userID string) (*Wallet, error)
	ChangeUserStatus(isActive bool, userReference string) (interface{}, error)
}

// Mysql struct
type Mysql struct {
	DB *gorm.DB
}

func NewWalletRepositoryDB(db *gorm.DB) *Mysql {
	return &Mysql{
		DB: db,
	}
}


func (mysql *Mysql) GetUserByEmail(email string) ([]*User, error) {
	var user []*User
	userEmail := mysql.DB.Where("email = ?", email).First(&user)
	return user, userEmail.Error
}


func (mysql *Mysql) CreateWallet(user *User) (*User, error) {
	var wallet  Wallet

	err := mysql.DB.Create(user).Error
	wallet.UserID= strconv.Itoa(int(user.Id))
	err2 := mysql.DB.Create(wallet).Error
	log.Println(err2)
	return user, err
}

func (mysql *Mysql) CheckIfPasswordExists(userReference string) ([]*User, error) {
	var user []*User
	userFound := mysql.DB.Where("id = ?", userReference).First(&user)
	fmt.Println("user record does not exist: ",userFound.Error)
	return user, userFound.Error
}


func (mysql *Mysql) SaveTransaction(transaction *Transaction) (*Transaction, error) {
	err := mysql.DB.Create(transaction).Error
	log.Println("error in saving transaction: ",err)
	return transaction, err
}


func (mysql *Mysql) PostToAccount(account *Wallet) (*Wallet, error) {
	var newAccount Wallet
	mysql.DB.First(&newAccount)
	newAccount.UserID=account.UserID
	newAccount.Balance = account.Balance

	fmt.Println("new account details",newAccount)
	err := mysql.DB.Model(&newAccount).Where("user_id = ?", account.UserID).Update("balance", account.Balance).Error
	log.Print(err)
	return account, err
}


func (mysql *Mysql) ChangeUserStatus(isActive bool, id string) (interface{}, error) {

	result :=
		mysql.DB.Model(User{}).
			Where("id = ?", id).
			Updates(
				User{
					IsActive: isActive,
				},
			)
	return result, result.Error
}

//func (mysql *Mysql) GetAllTransactionsById(id string) ([]*Wallet, error) {
//	var history []*Wallet
//	historyFound := mysql.DB.Where("user_id = ?", id).First(&history)
//	fmt.Println(historyFound.Error)
//	return historyFound, historyFound.Error
//	//mysql.DB.First(&history, "id = ?", "string_primary_key")
//}

func (mysql *Mysql) GetAccountBalance(userID string) (*Wallet, error) {
	var user *Wallet
	userEmail := mysql.DB.Where("user_id = ?", userID).Take(&user)
	return user, userEmail.Error
}
