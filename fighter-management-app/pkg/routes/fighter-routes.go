package routes

import (
	"fighter-management-app/pkg/controllers"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func createFighterHandlerWithDB(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateFighterHandler(w, r, db)
	}
}

func createOrganizationHandlerWithDB(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateOrganizationHandler(w, r, db)
	}
}

func getAllFightersHandlerWithDB(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllFightersHandler(w, r, db)
	}
}

func getAllOrganizationsHandlerWithDB(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controllers.GetAllOrganizationsHandler(w, r, db)
	}
}

var RegisterFighterRoutes = func(router *mux.Router, db *gorm.DB) {
	//creates a new fighter and a new org and posts the info
	router.HandleFunc("/fighters/", createFighterHandlerWithDB(db)).Methods("POST")
	router.HandleFunc("/organization/", createOrganizationHandlerWithDB(db)).Methods("POST")
	//gets all
	router.HandleFunc("/fighters/", getAllFightersHandlerWithDB(db)).Methods("GET")

	router.HandleFunc("/organization/", getAllOrganizationsHandlerWithDB(db)).Methods("GET")
	//gets specific one by id
	router.HandleFunc("/fighters/{fighterId}", controllers.GetFighterByIdHandler).Methods("GET")
	router.HandleFunc("/organization/{organizationId}", controllers.GetOrganizationByIdHandler).Methods("GET")
	//update specific fighter/orgs fields
	router.HandleFunc("/fighters/{fighterId}", controllers.UpdateFighter).Methods("PUT")
	router.HandleFunc("/organization/{organizationId}", controllers.UpdateOrganization).Methods("PUT")
	//delete specific fighter/org
	router.HandleFunc("/fighters/{fighterId}", controllers.DeleteFighterHandler).Methods("DELETE")
	router.HandleFunc("/organization/{organizationId}", controllers.DeleteOrganizationHandler).Methods("DELETE")

}
