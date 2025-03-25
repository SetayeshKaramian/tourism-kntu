package services

import (
	"errors"
	"log"
	"tourism/pkg/config"
	"tourism/pkg/models"
	"tourism/pkg/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var db *gorm.DB 

func init() {
	config.Connect()
	db = config.GetDB()
}

func GetUserForLogin(phonenumber, password string) (map[string]interface{}, error) {
	var user map[string]interface{}
	result := db.Raw(`SELECT * FROM "User" WHERE phonenumber = ? LIMIT 1`, phonenumber).Scan(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, result.Error
	}

	passwordHash, ok := user["hashedpassword"].(string)
	if !ok {
		return nil, errors.New("password field not found or invalid")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}


func GetUserIDWithPhoneNumber(phoneNumber string) (uuid.UUID, error) {
	var userID string
	result := db.Raw(`SELECT userid FROM "User" WHERE phonenumber = ? LIMIT 1`, phoneNumber).Scan(&userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return uuid.Nil, errors.New("user not found")
		}
		return uuid.Nil, result.Error
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return uuid.Nil, err
	}

	return parsedUserID, nil
}

func CreateNewUser(newUser models.User) (models.User, error){
	query := `INSERT INTO "User" (FirstName, LastName, Email, PhoneNumber, City, HashedPassword, AccountStatus, UserType) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING UserID, FirstName, LastName, Email, PhoneNumber, City, HashedPassword, AccountStatus, UserType, RegistrationDate`

	hasedhPass, err := utils.HashPassword(*newUser.HashedPassword)
	if err != nil {
		return models.User{}, err
	}

	sqlDB, _ := db.DB()
	var createdUser models.User;
	query_error := sqlDB.QueryRow(query,
		newUser.FirstName,
		newUser.LastName,
		newUser.Email,
		newUser.PhoneNumber,
		newUser.City,
		hasedhPass,
		newUser.AccountStatus,
		newUser.UserType,
	).Scan(&createdUser.UserID,
		&createdUser.FirstName,
		&createdUser.LastName,
		&createdUser.Email,
		&createdUser.PhoneNumber,
		&createdUser.City,
		&createdUser.HashedPassword,
		&createdUser.AccountStatus,
		&createdUser.UserType,
		&createdUser.RegistrationDate,)
	if query_error != nil {
		return models.User{}, err
	}

	return createdUser, nil
}

func UpdateUser(updateData models.User, currentUser models.User) error {
	sqlDB, _ := db.DB()
	query := `UPDATE "User" SET `
	argCounter := 1
	userID := currentUser.UserID

	if updateData.FirstName != nil {
		query += `firstname = '` + *updateData.FirstName + `', `
		argCounter++
	}
	if updateData.LastName != nil {
		query += `lastname = '` + *updateData.LastName + `', `
		argCounter++
	}
	if updateData.Email != nil {
		query += `email = '` + *updateData.Email + `', `
		argCounter++
	}
	if updateData.PhoneNumber != nil {
		query += `phonenumber = '` + *updateData.PhoneNumber + `', `
		argCounter++
	}
	if updateData.City != nil {
		query += `city = '` +  *updateData.City + `', `
		argCounter++
	}
	if updateData.AccountStatus != nil {
		query += `accountstatus = '` + *updateData.AccountStatus + `', `
		argCounter++
	}
	if updateData.UserType != nil {
		query += `usertype = '` + *updateData.UserType + `', `
		argCounter++
	}

	if argCounter == 0 {
		return errors.New("no field to update")
	}
	query = query[:len(query)-2] 

	query += ` WHERE userid = '` + userID.String() + `';`

	result, err := sqlDB.Exec(query)
	if err != nil {
		return errors.New("query not executed")
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func GetTicketsWithDetails(origin, destination, vehicleType, departureTime, arrivalTime *string) []*models.Ticket {
	var tickets []*models.Ticket
	query := `SELECT * FROM ticket WHERE `
	
	if origin != nil {
		query += `origin = '` + *origin + `' AND `
	}
	if destination != nil {
		query += `destination = '` + *destination + `' AND `
	}
	if vehicleType != nil {
		query += `vehicletype = '` + *vehicleType + `' AND `
	}
	if departureTime != nil {
		query += `departuretime > '` + *departureTime + `' AND `
	}
	if arrivalTime != nil {
		query += `arrivaltime < '` + *arrivalTime + `' AND `
	}
	query = query[:len(query) - 4]
	db.Raw(query).Scan(&tickets)
	return tickets
}

func GetTicketDetails(ticketid uuid.UUID) *models.TicketDetails {
	var ticket *models.TicketDetails
	query := `
		SELECT 
		t.ticketid AS ticket_id, t.vehicletype AS vehicle_type, t.origin, t.destination, 
		t.departuretime AS departure_time, t.arrivaltime AS arrial_time, t.ticketprice AS ticket_price, 
		t.remainingcapacity AS remaining_capacity, t.travelclass AS travel_class,
		
		tr.starrating AS train_star_rating, tr.amenities AS train_amenities, tr.iscoupeavailable
		AS is_coupe_available,
		
		f.airlinename AS airline_name, f.flightclass AS flight_class, f.numberofstops AS number_of_stops,
		f.flightnumber AS flight_number, 
		f.originairport AS origin_airport, f.destinationairport AS destination_airport, 
		f.amenities AS flight_amenities,
		
		b.buscompanyname AS bus_company_name, b.bustype AS bus_type, b.seatsperrow AS seats_per_row,
		b.amenities AS bus_amenities
		
		FROM ticket t
		LEFT JOIN traindetails tr ON t.ticketid = tr.ticketid AND t.vehicletype = 'Train'
		LEFT JOIN flightdetails f ON t.ticketid = f.ticketid AND t.vehicletype = 'Airplane'
		LEFT JOIN busdetails b ON t.ticketid = b.ticketid AND t.vehicletype = 'Bus'
		WHERE t.ticketid = 
	'` + ticketid.String() + `';`

	log.Println(query)
	db.Raw(query).Scan(&ticket)
	return ticket
}

func GetCities() []string {
	var cities []string
	query := `SELECT DISTINCT destination FROM ticket`
	db.Raw(query).Scan(&cities)
	return cities
}

func GetMyTickets(userid uuid.UUID) []*models.TicketList {
	var tickets []*models.TicketList
	query := `SELECT t.*, r.reservationstatus AS reservation_status, r.reservationtime AS reservation_time
	FROM ticket t JOIN reservation r ON t.ticketid = r.ticketid WHERE r.userid = '` + userid.String() + `';`
	db.Raw(query).Scan(&tickets)
	return tickets
}

func CreateReport(report models.Report) (string, error) {
	query := `
		INSERT INTO reports (userid, ticketid, paymentid, reportcategory, reporttext, reportstatus, reporttime)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
		RETURNING reportid;
	`
	var reportID string
	err := db.Raw(query, report.UserID, report.TicketID, report.PaymentID, report.ReportCategory, 
		report.ReportText, "Pending").Scan(&reportID).Error
	if err != nil {
		return "", errors.New("query not executed")
	}
	return reportID, nil
}

func GetAllReports() []*models.Report {
	var reports []*models.Report
	query := `SELECT reportid AS report_id, userid AS user_id, paymentid AS payment_id,
	ticketid AS ticket_id, reportcategory AS report_category, reportstatus AS report_status,
	reporttext AS resport_text, reporttime AS report_time FROM reports`
	db.Raw(query).Scan(&reports)
	return reports
}