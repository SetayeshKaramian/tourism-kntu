package controllers

import (
	"encoding/json"
	"fmt"

	"net/http"
	"tourism/pkg/models"
	"tourism/pkg/services"
	"tourism/pkg/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	// "github.com/google/uuid"
	// "github.com/gorilla/mux"
)


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		type LoginRequest struct {
			PhoneNumber string `json:"phonenumber"`
			Password string `json:"password"`
		}
		var loginReq LoginRequest
			err := json.NewDecoder(r.Body).Decode(&loginReq)
		if err != nil {
			http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
			return
		}
		user , err := services.GetUserForLogin(loginReq.PhoneNumber, loginReq.Password)
		if err != nil {
			fmt.Print(user)
			http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusForbidden)
			return
		}
		
		otp := utils.GenerateOTP(loginReq.PhoneNumber, utils.Purposes["login"])
	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"otp": "%s", "phonenumber": "%s"}`, otp, loginReq.PhoneNumber)))

	case "PUT":
		var otpReq struct {
			PhoneNumber string `json:"phonenumber"`
			OTP 		string `json:"otp"`
		}
		err := json.NewDecoder(r.Body).Decode(&otpReq)
		if err != nil {
			http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
			return 
		}
		if !utils.ValidateOTP(otpReq.PhoneNumber, utils.Purposes["login"], otpReq.OTP) {
			http.Error(w, `{"error": "Invalid OTP"}`, http.StatusUnauthorized)
			return
		}
		userID, err := services.GetUserIDWithPhoneNumber(otpReq.PhoneNumber)
		if err != nil {
			http.Error(w, `{"error": "Invalid phone number"}`, http.StatusForbidden)
			return
		}
		token, err := utils.GenerateJWT(userID)
		if err != nil {
			http.Error(w, "token could not be generated", http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))
	}
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUserInfo models.User
	err := json.NewDecoder(r.Body).Decode(&newUserInfo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := services.CreateNewUser(newUserInfo)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}
	token, err := utils.GenerateJWT(user.UserID)
	if err != nil {
		http.Error(w, "token could not be generated", http.StatusForbidden)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"token": "%s"}`, token)))
}

func UpdateUser(w http.ResponseWriter, r *http.Request) { 
	vars := mux.Vars(r)
	userID, exists := vars["id"]  
	if !exists {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	userid, _ := uuid.Parse(userID)
	user, _ := r.Context().Value(utils.UserContextKey).(*models.User)

	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	if user.UserID != userid {
		http.Error(w, "you dont have permission", http.StatusBadRequest)
		return
	}
	
	var updateData models.User
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	err := services.UpdateUser(updateData, *user)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "User updated successfully"}`))
}