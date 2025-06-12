package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID           uuid.UUID `gorm:"column:userid"`
	FirstName        *string   `gorm:"column:firstname"`
	LastName         *string   `gorm:"column:lastname"`
	Email            *string   `gorm:"column:email"`
	PhoneNumber      *string   `gorm:"column:phonenumber"`
	City             *string   `gorm:"column:city"`
	HashedPassword   *string   `gorm:"column:hashedpassword"`
	RegistrationDate time.Time `gorm:"column:registrationdate"`
	AccountStatus    *string   `gorm:"column:accountstatus"`
	UserType         *string   `gorm:"column:usertype"`
	credit           int64     `gorm:"column:credit"`
}

type Ticket struct {
	TicketID          uuid.UUID `gorm:"column:ticketid" json:"ticket_id"`
	VehicleType       string    `gorm:"column:vehicletype" json:"vehicle_type"`
	Origin            string    `gorm:"column:origin" json:"origin"`
	Destination       string    `gorm:"column:destination" json:"destination"`
	DepartureTime     time.Time `gorm:"column:departuretime" json:"departure_time"`
	ArrivalTime       time.Time `gorm:"column:arrivaltime" json:"arrival_time"`
	TicketPrice       float64   `gorm:"column:ticketprice" json:"ticket_price"`
	RemainingCapacity int       `gorm:"column:remainingcapacity" json:"remaining_capacity"`
	CompanyID         uuid.UUID `gorm:"column:companyid" json:"company_id"`
	TravelClass       string    `gorm:"column:travelclass" json:"travel_class"`
}
type TicketDetails struct {
	TicketID          uuid.UUID `json:"ticket_id"`
	VehicleType       string    `json:"vehicle_type"`
	Origin            string    `json:"origin"`
	Destination       string    `json:"destination"`
	DepartureTime     time.Time `json:"departure_time"`
	ArrivalTime       time.Time `json:"arrival_time"`
	TicketPrice       float64   `json:"ticket_price"`
	RemainingCapacity int       `json:"remaining_capacity"`
	TravelClass       string    `json:"travel_class"`

	TrainStarRating  *string `json:"train_star_rating,omitempty"`
	TrainAmenities   *string `json:"train_amenities,omitempty"`
	IsCoupeAvailable *bool   `json:"is_coupe_available,omitempty"`

	AirlineName        *string `json:"airline_name,omitempty"`
	FlightClass        *string `json:"flight_class,omitempty"`
	NumberOfStops      *int    `json:"number_of_stops,omitempty"`
	FlightNumber       *string `json:"flight_number,omitempty"`
	OriginAirport      *string `json:"origin_airport,omitempty"`
	DestinationAirport *string `json:"destination_airport,omitempty"`
	FlightAmenities    *string `json:"flight_amenities,omitempty"`

	BusCompanyName *string `json:"bus_company_name,omitempty"`
	BusType        *string `json:"bus_type,omitempty"`
	SeatsPerRow    *string `json:"seats_per_row,omitempty"`
	BusAmenities   *string `json:"bus_amenities,omitempty"`
}

type TicketList struct {
	Ticket            `json:"ticket"`
	ReservationStatus string    `json:"reservation_status"`
	ReservationTime   time.Time `json:"reservation_time"`
	Number            int       `json:"number"`
}

type Report struct {
	ReportID       uuid.UUID  `json:"report_id"`
	UserID         uuid.UUID  `json:"user_id" binding:"required,uuid"`
	TicketID       *uuid.UUID `json:"ticket_id,omitempty" binding:"omitempty,uuid"`
	PaymentID      *uuid.UUID `json:"payment_id,omitempty" binding:"omitempty,uuid"`
	ReportCategory string     `json:"report_category" binding:"required,oneof=PaymentIssue Delay Cancellation Other"`
	ReportText     string     `json:"report_text,omitempty"`
	ReportStatus   *string    `json:"report_status,omitempty" binding:"omitempty,oneof=Reviewed Pending"`
	ReportTime     time.Time  `json:"report_time,omitempty"`
}

type Penalty struct {
	TicketID          uuid.UUID `json:"ticket_id"`
	UserID            uuid.UUID `json:"user_id"`
	PenaltyAmount     int64     `json:"penalty_amount"`
	PenaltyPercentage float64   `json:"penalty_percentage"`
	RefundAmount      int64     `json:"refund_amount"`
}
