package routes

import (
	"net/http"
	"tourism/pkg/controllers"
	"tourism/pkg/utils"

	"github.com/gorilla/mux"
)

var RegisterRoutes = func(router *mux.Router) {
	router.HandleFunc("/login/", controllers.LoginHandler).Methods("POST", "PUT")
	router.Handle("/new/user/", http.HandlerFunc(controllers.RegisterUser)).Methods("POST")
	router.Handle("/update/user/{id}/", utils.JWTMiddleware(http.HandlerFunc(controllers.UpdateUser))).Methods("PUT")
	router.HandleFunc("/get/tickets/", controllers.GetTickets).Methods("GET")
	router.HandleFunc("/get/ticket/{id}/", controllers.GetTicket).Methods("GET")
	router.HandleFunc("/get/cities/", controllers.GetCities).Methods("GET")
}