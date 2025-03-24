package main

import (
	"log"
	"net/http"
	"tourism/pkg/routes"

	"github.com/gorilla/mux"
)



func main() {
	r := mux.NewRouter()
	routes.RegisterRoutes(r)
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err == nil {
			log.Println("Registered route:", path)
		} else {
			log.Println("Error registering route:", err)
		}
		return nil
	})
	log.Fatal(http.ListenAndServe(":8000", r))
}