package services

import (
	"errors"
	"tourism/pkg/config"

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