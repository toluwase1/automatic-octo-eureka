package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"wallet-engine/models"
	"wallet-engine/response"
	"wallet-engine/utilities"
)

func (handlers *Handler) GetTransactionHistoryByUserId() gin.HandlerFunc {
	return func(context *gin.Context) {

		userID := context.Param("id")
		transaction := &models.Transaction{}

		transaction.UserID = userID

		transaction.Id = uuid.New().String()


		if err := utilities.Decode(context, &transaction); err != nil {
			response.JSON(context, http.StatusNotFound, nil, []string{"cannot decode transaction"}, "")
			return
		}


		if transaction.Amount < 1000 {
			response.JSON(context, http.StatusNotFound, nil, []string{"sorry you can't deposit less than N1000.00"}, "")
			return
		}

		userDB, err := handlers.WalletService.CheckIfPasswordExists(userID)
		if err != nil {
			response.JSON(context, http.StatusNotFound, nil, []string{"could not fetch user reference "}, "")
			return
		}

		var hashedPassword string
		var activationStatus bool
		for _, user := range userDB {
			hashedPassword = user.SecretKey
			activationStatus = user.IsActive
		}

		if correct := utilities.CheckPasswordHash(transaction.Password, []byte(hashedPassword)); correct {
			response.JSON(context, http.StatusNotFound, nil, []string{"Invalid password"}, "")
			return
		}

		account := &models.Wallet{}
		currentUser := &models.User{}
		currentUser.IsActive = activationStatus


		if currentUser.IsActive == false {
			response.JSON(context, http.StatusNotFound, nil, []string{"sorry, please activate your account"}, "")
			return
		}

		t, err := handlers.WalletService.GetAccountBalance(userID)
		if err != nil {
			log.Fatalln(err)
		}

		account.Balance = t.Balance

		account.CreditUserWallet(transaction.Amount, transaction.UserID)

		userTransaction, err := handlers.WalletService.SaveTransaction(transaction)
		if err != nil {
			response.JSON(context, http.StatusNotFound, nil, []string{"could not fetch userid "}, "")
			return
		}

		currentAccount, err := handlers.WalletService.PostToAccount(account)
		if err != nil {
			response.JSON(context, http.StatusNotFound, nil, []string{"could not post to account "}, "")
			return
		}

		response.JSON(context, http.StatusCreated, gin.H{
			"transaction id": userTransaction.Id,
			"amount credited to user account ": userTransaction.Amount,
			"New account balance": currentAccount.Balance,
		},
			nil,
			"account has been successfully credited with " +fmt.Sprintf("%f", userTransaction.Amount)+"your new account balance is "+ fmt.Sprintf("%f", currentAccount.Balance))
	}
}
