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
	router.Handle("/get/my/tickets/", utils.JWTMiddleware(http.HandlerFunc(controllers.GetMyTickets))).Methods("GET")
	router.Handle("/create/report/", utils.JWTMiddleware(http.HandlerFunc(controllers.MakeReport))).Methods("POST")
	router.Handle("/reports/", utils.JWTMiddleware(http.HandlerFunc(controllers.GetAllReports))).Methods("GET")
	router.Handle("/cancelled/tickets/", utils.JWTMiddleware(http.HandlerFunc(controllers.GetAllCancelledTickets))).Methods("GET")
}