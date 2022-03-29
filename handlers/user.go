package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"wallet-engine/models"
	"wallet-engine/response"
	"wallet-engine/utilities"
)
type Handler struct {
	WalletService models.WalletService
}

func (handlers *Handler) CreateWallet() gin.HandlerFunc {
	return func(context *gin.Context) {

		var user = &models.User{}
		secretKey, err := utilities.GenerateHashPassword(user.Password)
		if err != nil {
			log.Fatalf("error :%v", err)
		}

		//user.Id = uuid.New().String()
		user.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
		user.SecretKey = string(secretKey)

		if errs := utilities.Decode(context, &user); errs != nil {
			response.JSON(context, http.StatusInternalServerError, nil, errs, "")
			return
		}

		userD, err := handlers.WalletService.CreateWallet(user)

		if err != nil {
			response.JSON(context, http.StatusNotFound, nil, []string{"User email already exists"}, "")
			return
		}
		response.JSON(context, http.StatusCreated, gin.H{"data": userD}, nil, "User created successfully")
		return

	}
}

