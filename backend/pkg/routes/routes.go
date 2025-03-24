package routes

import (
	"tourism/pkg/controllers"

	"github.com/gorilla/mux"
)

var RegisterRoutes = func(router *mux.Router) {
	router.HandleFunc("/login/", controllers.LoginHandler).Methods("POST", "PUT")
}