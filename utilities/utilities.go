package utilities

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"wallet-engine/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
	"os"
	"strings"
)


func GenerateHashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CheckPasswordHash(password string, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	return err == nil
}

func Decode(c *gin.Context, v interface{}) []string {
	err := c.ShouldBindJSON(v)
	log.Println(err)
	if  err != nil {
		log.Println("1")
		var errs []string
		validation, ok := err.(validator.ValidationErrors)
		log.Println("2")
		if ok {
			log.Println("3")
			for _, fieldErr := range validation {
				errs = append(errs, NewFieldError(fieldErr).String())
				log.Println(errs)
			}
		} else {
			log.Println("errorr")
			errs = append(errs, "internal server error: "+err.Error())
		}
		return errs
	}
	return nil
}

//func Init() *gorm.DB {
//
//	mongoURL := fmt.Sprintf("%s://%s:%s", os.Getenv("db_type"), os.Getenv("mongo_db_host"), os.Getenv("mongo_db_port"))
//
//	timeout := time.Minute * 15
//
//	// using go mongo-driver  to connect to mongoDB
//	ctx, cancel := context.WithTimeout(context.Background(), timeout)
//	defer cancel()
//	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
//	if err != nil {
//		log.Fatalf("error %v", err)
//	}
//
//	log.Println("Database Connected Successfully...")
//	return client
//}

func Initialize() *gorm.DB {
	var Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string
	Dbdriver = os.Getenv("DB_DRIVER")
	DbUser = os.Getenv("DB_USER")
	DbPassword = os.Getenv("DB_PASSWORD")
	DbHost = os.Getenv("DB_HOST")
	DbName = os.Getenv("DB_NAME")
	DbPort = os.Getenv("DB_PORT")

	// If you want to use postgres, i added support for you in the else block, (dont forgot to edit the .env file)
	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		fmt.Println("here.. ",DBURL)
		dsn := "root:toluwase@tcp(127.0.0.1:3306)/wallet?charset=utf8mb4&parseTime=True&loc=Local"

		Database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
		AutoMigrate(Database)
		return Database
	} else if Dbdriver == "postgres" {
		//DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		//Database, err :=  gorm.Open(postgres.Open(DBURL), &gorm.Config{})
		//AutoMigrate(Database)
		//if err != nil {
		//	fmt.Printf("Cannot connect to %s database", Dbdriver)
		//	log.Fatal("This is the error connecting to postgres:", err)
		//} else {
		//	fmt.Printf("We are connected to the %s database", Dbdriver)
		//}
	} else {
		fmt.Println("Unknown Driver")
	}
	return nil
	//database migration

	//
	//server.Router = gin.Default()
	//server.Router.Use(middlewares.CORSMiddleware())
	//
	//server.initializeRoutes()

}

func AutoMigrate(Database *gorm.DB)  {
	Database.Debug().AutoMigrate(
		&models.User{},
		&models.Transaction{},
		&models.Wallet{},
	)
}




type FieldError struct {
	err validator.FieldError
}

func (f FieldError) String() string {
	var sb strings.Builder

	sb.WriteString("validation failed on field '" + f.err.Field() + "'")
	sb.WriteString(", condition: " + f.err.ActualTag())

	// Print condition parameters, e.g. one_of=red blue -> { red blue }
	if f.err.Param() != "" {
		sb.WriteString(" { " + f.err.Param() + " }")
	}

	if f.err.Value() != nil && f.err.Value() != "" {
		sb.WriteString(fmt.Sprintf(", actual: %v", f.err.Value()))
	}

	return sb.String()
}

func NewFieldError(err validator.FieldError) FieldError {
	return FieldError{err: err}
}