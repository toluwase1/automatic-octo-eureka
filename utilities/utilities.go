package utilities

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
	"wallet-engine/models"
)

func GenerateHashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CheckPasswordHash(password, hashedPassword string) error {
	log.Println("hashedPassword: ", hashedPassword)
	log.Println("password: ", password)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	log.Println(err)
	log.Println(err == nil)
	return err // nil means it is a match
}

//func iniciarSesion(w http.ResponseWriter, r *http.Request) {
//
//	w.Header().Add("Content-Type", "application/json")
//	w.Header().Add("Access-Control-Allow-Origin", "*")
//	w.Header().Add("Access-Control-Allow-Methods:", "POST")
//	var usuarios []
//
//	Usuarioerr := r.ParseForm()
//
//	if err != nil {
//		return
//	}
//
//	body, err := ioutil.ReadAll(r.Body)
//
//	if err != nil {
//		panic(err.Error())
//	}
//	datosPost := make(map[string]string)
//	json.Unmarshal(body, &datosPost)
//	email    := datosPost["email"]
//	password := datosPost["password"]
//	result, err := db.Query("SELECT email, password FROM usuario WHERE email = ?", &email)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	for result.Next() {
//
//		var usuario Usuario
//		err := bcrypt.CompareHashAndPassword([]byte(*&usuario.Password), []byte(password))
//		if err != nil {
//			panic(err.Error())
//		}
//
//		defer result.Close()
//		errr := result.Scan(&usuario.Email, &usuario.Nickname, &usuario.Password, &usuario.Roscos, &usuario.Roscosperfectos, &usuario.Rol, &usuario.Foto)
//
//		if errr != nil {
//			panic(err.Error())
//		}
//
//		usuarios = append(usuarios, usuario)
//	}
//	json.NewEncoder(w).Encode(usuarios)
//}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

}

func Decode(c *gin.Context, v interface{}) []string {
	err := c.ShouldBindJSON(v)
	log.Println(err)
	if err != nil {
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
		fmt.Println("here.. ", DBURL)
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

func AutoMigrate(Database *gorm.DB) {
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

// four sum
