package controllers

import (
	"encoding/json"
	"fighter-management-app/pkg/models"
	"fighter-management-app/pkg/utils"
	f "fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var NewFighter models.Fighter
var NewOrganization models.Organization

func getAssociatedFighters(db *gorm.DB, organizationID int) []*models.Fighter {
	var fighters []*models.Fighter
	db.Find(&fighters, "organization_refer = ?", organizationID)
	return fighters
}

func CreateFighterHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Parse the request body
	newFighter := &models.Fighter{}
	utils.ParseBody(r, newFighter)
	organizationIDStr := r.FormValue("organizationID")
	organizationID, err := strconv.ParseInt(organizationIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Organization ID", http.StatusBadRequest)
		return
	}

	tx := db.Begin()

	//Fetch the organization with the above id
	var organization models.Organization
	if tx.First(&organization, organizationID).RecordNotFound() {
		tx.Rollback()
		http.Error(w, "Organization with the provided ID does not exist", http.StatusNotFound)
		return
	}

	// Create the new fighter
	orgID := int(organizationID)
	newFighter.OrganizationRefer = &orgID

	//Save the fighter within the transaction
	if err := tx.Create(newFighter).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create fighter", http.StatusInternalServerError)
		return
	}

	fighters := getAssociatedFighters(db, orgID)

	//append the fighter to the Organization fighter's slice
	organization.Fighters = fighters

	//Save organization within the transaction
	if err := tx.Save(&organization).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update organization", http.StatusInternalServerError)
		return
	}

	tx.Commit()
	//db.Save(&organization)

	//db.Preload("Fighters").First(newFighter, newFighter.ID) incorrect
	db.Preload("Organization").First(newFighter, newFighter.ID)

	// Fetch the fighter again to include the associated organization
	res, _ := json.Marshal(newFighter)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateOrganizationHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	CreateOrganization := &models.Organization{}
	utils.ParseBody(r, CreateOrganization)
	o := CreateOrganization.CreateOrganization()
	res, _ := json.Marshal(o)
	w.WriteHeader(http.StatusOK)
	w.Write(res)
	f.Printf("Parsed JSON: %+v\n", CreateOrganization)
}

func GetAllFightersHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	//fighters := models.GetAllFighters()
	var fighters []*models.Fighter
	db.Preload("Organization").Find(&fighters)

	res, err := json.Marshal(fighters)
	if err != nil {
		http.Error(w, "Failed to marshal fighter to JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetAllOrganizationsHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	//	organizations := models.GetAllOrganizations()
	var organizations []*models.Organization
	db.Preload("Fighters").Find(&organizations)

	res, err := json.Marshal(organizations)
	if err != nil {
		http.Error(w, "Failed to marshal organization to JSON", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetFighterByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fighterId := vars["fighterId"]
	ID, err := strconv.ParseInt(fighterId, 0, 0)
	if err != nil {
		f.Println("Error while parsing")
	}
	fighterDetails, _ := models.GetFighterById(ID) // this ID is now the converted integer from strconv.ParseInt, we're getting the other fields from this
	res, _ := json.Marshal(fighterDetails)         //now we are converting the fighter details to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetOrganizationByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	organizationId := vars["organizationId"]
	ID, err := strconv.ParseInt(organizationId, 0, 0)
	if err != nil {
		f.Println("Error while parsing")
	}
	organizationDetails, _ := models.GetOrganizationById(ID)
	res, _ := json.Marshal(organizationDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateFighter(w http.ResponseWriter, r *http.Request) {
	updateFighter := &models.Fighter{}           //create an empty instance of a fighter struct
	utils.ParseBody(r, updateFighter)            //we populate the (empty) updateFighter object with data received from the request body
	vars := mux.Vars(r)                          //extract the url via the route pattern we have defined/fighter
	fighterId := vars["fighterId"]               //we specify which value parameter we want fighter/fighterId
	ID, err := strconv.ParseInt(fighterId, 0, 0) //we parse the ID from a string to an int
	if err != nil {
		f.Println("Error while parsing")
	}
	fighterDetails, db := models.GetFighterById(ID) //we get the fighterDetails that the ID has specified
	if updateFighter.Name != "" {
		fighterDetails.Name = updateFighter.Name
	}
	if updateFighter.Age != 0 {
		fighterDetails.Age = updateFighter.Age
	}
	if updateFighter.Organization != nil {
		fighterDetails.Organization = updateFighter.Organization
	}
	if updateFighter.Reach != 0 {
		fighterDetails.Reach = updateFighter.Reach
	}
	if updateFighter.Weight != 0 {
		fighterDetails.Weight = updateFighter.Weight
	}

	db.Save(&fighterDetails)               //we save the fighterDetails to the database
	res, _ := json.Marshal(fighterDetails) //for the client we return the data we just saved in the database
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	updateOrganization := &models.Organization{}
	utils.ParseBody(r, updateOrganization)

	vars := mux.Vars(r)
	organizationId := vars["organizationId"]
	ID, err := strconv.ParseInt(organizationId, 0, 0)
	if err != nil {
		http.Error(w, "Error while parsing ID request", http.StatusInternalServerError)
		return
	}

	organizationDetails, _ := models.GetOrganizationById(ID)
	if organizationDetails.Name != "" {
		organizationDetails.Name = updateOrganization.Name
	}
	if organizationDetails.Networth != 0 {
		organizationDetails.Networth = updateOrganization.Networth
	}

	res, err := json.Marshal(organizationDetails)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteFighterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fighterId := vars["fighterId"]
	ID, err := strconv.ParseInt(fighterId, 0, 0)
	if err != nil {
		f.Println("Error parsing ID to delete fighter")
	}
	fighter := models.DeleteFighter(ID)
	res, _ := json.Marshal(fighter)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func DeleteOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	organizationId := vars["organizationId"]
	ID, err := strconv.ParseInt(organizationId, 0, 0)
	if err != nil {
		f.Println("Error parsing ID to delete organization")
	}
	organization := models.DeleteOrganization(ID)
	res, _ := json.Marshal(organization)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

//
// func UpdateOrganization(w http.ResponseWriter, r *http.Request) {
// 	updateOrganization := &models.Organization{}
// 	utils.ParseBody(r, updateOrganization)
// 	vars := mux.Vars(r)
// 	organizationId := vars["organizationId"]
// 	ID, err := strconv.ParseInt(organizationId, 0, 0)
// 	if err != nil {
// 		f.Println("Error while parsing ID request")
// 	}
// 	organizationDetails, db := models.GetOrganizationById(ID)
// 	if organizationDetails.Name != "" {
// 		organizationDetails.Name = updateOrganization.Name
// 	}
// 	if organizationDetails.Networth != nil {
// 		organizationDetails.Networth = updateOrganization.Networth
// 	}

// 	db.Save(&organizationDetails)
// 	res, _ := json.Marshal(organizationDetails)
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(res)
// }
