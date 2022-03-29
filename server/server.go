package server

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"wallet-engine/handlers"
	"wallet-engine/models"
	"wallet-engine/utilities"
)

func Start() {
	router := initializeRouter()
	database := utilities.Initialize()
	handler := handlers.Handler{
		WalletService: models.NewWalletService(models.NewWalletRepositoryDB(database)),
	}
	DefineRouter(router, &handler)
	PORT := fmt.Sprintf(":%s", os.Getenv("API_PORT"))
	if PORT == ":" {
		PORT += "8080"
	}
	s := &http.Server{
		Handler: router,
		Addr:    PORT,
	}
	wait := make(chan os.Signal)

	log.Printf("Server Started at Port%s", PORT)

	go startServer(s)

	signal.Notify(wait, os.Interrupt)

	<-wait
	log.Printf("Shutting down the server...")

	time.Sleep(time.Second * 2) // sleep for 1 second

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("An error occurred: %s", err)
	}
	log.Printf("Server shut down")
}

func startServer(s *http.Server)  {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("An error occurred with the server: %s", err)
			return
		}
}
func initializeRouter() *gin.Engine {
	router := gin.Default()
	if os.Getenv("GIN_MODE") == "testing" {
		return router
	}

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,

		MaxAge: 12 * time.Hour,
	}))
	return router
}

func DefineRouter(r *gin.Engine, handler *handlers.Handler) {
	router := r.Group("/api/v1")
	router.POST("/create-user-wallet", handler.CreateWallet())
	router.PUT("/activate-deactivate-wallet/:id", handler.ActivateWallet())
	router.POST("/credit-user-wallet/:id", handler.CreditWallet())
	router.POST("/debit-user-wallet/:id", handler.DebitWallet())
	router.GET("/transaction-history/:id", handler.GetTransactionHistoryByUserid)
	router.GET("/credit-history/:id", handler.GetCreditTransactionHistoryByUserid)
	router.GET("/debit-history/:id", handler.GetDebitTransactionHistoryByUserid)

}