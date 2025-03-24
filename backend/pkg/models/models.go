package models

import (
	"time"

	"github.com/google/uuid"
)


type User struct {
	UserID          	uuid.UUID `gorm:"column:userid"`
	FirstName       	*string    `gorm:"column:firstname"`
	LastName        	*string    `gorm:"column:lastname"`
	Email           	*string    `gorm:"column:email"`
	PhoneNumber     	*string    `gorm:"column:phonenumber"`
	City            	*string    `gorm:"column:city"`
	HashedPassword  	*string    `gorm:"column:hashedpassword"`
	RegistrationDate 	time.Time   `gorm:"column:registrationdate"`
	AccountStatus   	*string    `gorm:"column:accountstatus"`
	UserType        	*string    `gorm:"column:usertype"`
}