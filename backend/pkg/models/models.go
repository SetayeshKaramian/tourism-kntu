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

type Ticket struct {
	TicketID         uuid.UUID  `gorm:"column:ticketid"`
	VehicleType      string  `gorm:"column:vehicletype"`
	Origin           string  `gorm:"column:origin"`
	Destination      string  `gorm:"column:destination"`
	DepartureTime    time.Time  `gorm:"column:departuretime"`
	ArrivalTime      time.Time  `gorm:"column:arrivaltime"`
	TicketPrice      float64 `gorm:"column:ticketprice"`
	RemainingCapacity int    `gorm:"column:remainingcapacity"`
	CompanyID        uuid.UUID `gorm:"column:companyid"`
	TravelClass      string  `gorm:"column:travelclass"`
}