package main

import (
	"fighter-management-app/pkg/config"
	"fighter-management-app/pkg/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	db := config.GetDB()
	routes.RegisterFighterRoutes(r, db)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:9090", r))

}
