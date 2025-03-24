package services

import (
	"errors"
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

func GetTicketsWithOriginAndDestination(origin, destination string) []*models.Ticket {
	var tickets []*models.Ticket
	query := `SELECT * FROM ticket WHERE origin = ? AND destination = ?`
	db.Raw(query, origin, destination).Scan(&tickets)
	return tickets
}