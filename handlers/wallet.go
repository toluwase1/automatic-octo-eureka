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




func (handlers *Handler) ActivateWallet() gin.HandlerFunc {
	return func(context *gin.Context) {
		userReference := context.Param("id")

		activate := context.Query("activate")

		user := &models.User{}

		// setting message to be sent as response
		var message string
		var status bool
		if activate == "true" {
			message = "account activation is successful"
			status = true
		} else {
			message = "account deactivation is successfu1"
			status = false

		}

		// Handles activation and deactivation of the wallet
		user.ActivateWallet(status)
		_, err := handlers.WalletService.ChangeUserStatus(user.IsActive, userReference)
		if err != nil {
			response.JSON(context, http.StatusNotFound, nil, []string{"could not change user status "}, "")
			return
		}

		response.JSON(context, http.StatusCreated, gin.H{
			"message": message,
		}, nil, "user activation status changed")
	}
}


func (handlers *Handler) CreditWallet() gin.HandlerFunc {
	return func(context *gin.Context) {

		fmt.Println("context ",context)

		userID := context.Param("id")
		transaction := &models.Transaction{}

		transaction.UserID = userID
		transaction.TransactionType = "credit"
		//transaction.Amount= context

		transaction.Id = uuid.New().String()

		//transaction.Type = "Credit"
		log.Println("transaction id" ,transaction.Id)
		log.Println("transaction user" ,transaction.UserID)
		log.Println("transaction amount" ,transaction.Amount)
		log.Println("transaction pw" ,transaction.Password)
		//log.Println("transaction type" ,transaction.Type)
		// Binding the json
		log.Println("transaction details" ,transaction)
		fmt.Println("context ",context)

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
		log.Println("getting id",t)
		if err != nil {
			log.Fatalln(err)
		}

		account.Balance = t.Balance

		account.CreditUserWallet(transaction.Amount, transaction.UserID)

		log.Println("transaction here: ",transaction)

		userTransaction, err := handlers.WalletService.SaveTransaction(transaction)
		log.Println("userTransaction: ",userTransaction)
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



func (handlers *Handler) DebitWallet() gin.HandlerFunc {
	return func(context *gin.Context) {

		userID := context.Param("id")
		transaction := &models.Transaction{}
		transaction.Id = uuid.New().String()
		transaction.UserID = userID
		transaction.TransactionType = "debit"

		if err := utilities.Decode(context, &transaction); err != nil {
			response.JSON(context, http.StatusNotFound, nil, []string{"transaction could not be decoded"}, "")
			return
		}

		if transaction.Amount < 1000 {
			response.JSON(context, http.StatusNotFound, nil, []string{"sorry you can't debit less than N1000.00"}, "")
			return
		}

		userDB, err := handlers.WalletService.CheckIfPasswordExists(userID)
		if err != nil {
			response.JSON(context, http.StatusNotFound, nil, []string{"could not fetch user reference "}, "")
			return
		}

		var hashedPassword string
		var checkIfUserActive bool
		for _, user := range userDB {
			hashedPassword = user.SecretKey
			checkIfUserActive = user.IsActive
		}

		// Confirming the password provided by the user
		if correct := utilities.CheckPasswordHash(transaction.Password, []byte(hashedPassword)); correct {
			response.JSON(context, http.StatusNotFound, nil, []string{"Invalid password"}, "")
			return
		}

		account := &models.Wallet{}
		wUser := &models.User{}
		wUser.IsActive = checkIfUserActive

		// Checking if the user is active
		if wUser.IsActive == false {
			response.JSON(context, http.StatusNotFound, nil, []string{"Sorry, your account is not active"}, "")
			return
		}

		t, err := handlers.WalletService.GetAccountBalance(userID)
		if err != nil {
			log.Print(err)
		}


		account.Balance = t.Balance


		if account.Balance <= 0 {
			response.JSON(context, http.StatusNotFound, nil, []string{"Your balance is insufficient to perform the specified operation"}, "")
			return
		}
		if account.Balance < transaction.Amount {
			response.JSON(context, http.StatusNotFound, nil, []string{"Sorry, your account is insufficient for this transaction"}, "")
			return
		}
		account.DebitUserWallet(transaction.Amount, transaction.UserID)

		userTransaction, err := handlers.WalletService.SaveTransaction(transaction)
		if err != nil {
			response.JSON(context, http.StatusNotFound, nil, []string{"could not fetch user reference "}, "")
			return
		}
		newBal, err := handlers.WalletService.PostToAccount(account)
		if err != nil {
			response.JSON(context, http.StatusNotFound, nil, []string{"could not debit to user account "}, "")
			return
		}
		response.JSON(context, http.StatusCreated, gin.H{
			"transaction id": userTransaction.Id,
			"amount debited from specified user account ":       userTransaction.Amount,
			"account balance":       newBal.Balance,
		},
			nil,
			"account has been successfuly debited with " +fmt.Sprintf("%f", userTransaction.Amount)+" your new account balance is "+ fmt.Sprintf("%f", newBal.Balance))
	}
}

