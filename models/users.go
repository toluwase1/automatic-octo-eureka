package models

type User struct {
	Id         		uint 	`gorm:"primary_key;auto_increment" json:"id"`
	FirstName       string  `gorm:"size:255;not null;" json:"first_name"`
	LastName        string	`gorm:"size:255;not null;" json:"last_name"`
	Email           string	`gorm:"size:255;not null;unique" json:"email"`
	Password        string  `gorm:"size:100;not null;" json:"password"`
	SecretKey       string  `gorm:"size:100;not null;" json:"secret_key"`
	CreatedAt       string  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	IsActive        bool     `gorm:"size:100;not null;" json:"is_active"`
	Balance			float64
	//Currencies      SupportedCurrencies `gorm:"embedded" json:"currencies"`
}



