package models

import (
	"fighter-management-app/pkg/config"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB

type Organization struct {
	gorm.Model
	Name     string     `gorm:""json:"name"`
	Networth int        `json:"networth"`
	Fighters []*Fighter `json:"fighters,omitempty" gorm:"foreignKey:OrganizationRefer"`
}

type Fighter struct {
	gorm.Model
	Name              string        `gorm:""json:"name"`
	Reach             int           `json:"reach"`
	Age               int           `json:"age"`
	Weight            int           `json:"weight"`
	OrganizationRefer *int          `json:"organizationrefer"`
	Organization      *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationRefer;references:ID;belongsTo"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&Organization{}, &Fighter{})
}

// func (f Fighter) CreateFighter(organizationID int) *Fighter {
// 	//assign the OrganizationRefer to the provided OrganizationFighter
// 	f.OrganizationRefer = &organizationID
// 	// Create the new fighter
// 	db.Create(f)
// 	fmt.Printf("Saved Fighter: %+v\n", f)
// 	return &f

// }

func (o *Organization) CreateOrganization() *Organization {
	var existingOrg Organization
	res := db.Where("name = ?", o.Name).First(&existingOrg)
	//if the res ^ is found(if it already exists) the code below finds it and returns it and a new one will not be created
	if res.Error == nil {
		return &existingOrg
	}

	db.Create(o)
	return o
}

func DeleteFighter(ID int64) Fighter {
	var fighter Fighter
	db.Where("ID=?", ID).Delete(&fighter)
	return fighter
}

func DeleteOrganization(ID int64) Organization {
	var organization Organization
	db.Where("ID=?", ID).Delete(&organization)
	return organization
}

func GetAllFighters() []Fighter {
	var Fighters []Fighter //declaring a slice of Fighter structs
	db.Find(&Fighters)
	return Fighters
}

func GetAllOrganizations() []Organization {
	var Organizations []Organization
	db.Find(&Organizations)
	return Organizations
}

func GetFighterById(Id int64) (*Fighter, *gorm.DB) {
	var getFighter Fighter
	db := db.Where("ID=?", Id).Find(&getFighter)
	return &getFighter, db
}

func GetOrganizationById(Id int64) (*Organization, *gorm.DB) {
	var getOrganization Organization
	db := db.Where("ID=?", Id).Find(&getOrganization)
	return &getOrganization, db
}

// func (f *Fighter) GetOrganization() error {
// 	err := db.Model(f).Related(&f.Organization, "OrganizationID").Error
// 	return err
// }

// func AssociateFighterWithOrg() ([]Organization, error) {

// 	organizationIDS := []uint{}
// 	if err := db.Model(&Fighter{}).Pluck("organization_id", &organizationIDS).Error; err != nil {
// 		return nil, err
// 	}

// 	organizations := []Organization{}
// 	if err := db.Where(organizationIDS).Find(&organizations).Error; err != nil {
// 		return nil, err
// 	}

// 	// Preload the associated Fighters for each Organization
// 	for i := range organizations {
// 		if err := db.Model(&organizations[i]).Preload("Fighters").Error; err != nil {
// 			return nil, err
// 		}
// 	}

// 	// fighters := []Fighter{}
// 	// if err := db.Model(&Fighter{}).Find(&fighters).Error; err != nil {
// 	// 	return nil, err
// 	// }

// 	// for i := range fighters {
// 	// 	if err := fighters[i].GetOrganization(); err != nil {
// 	// 		return nil, err
// 	// 	}
// 	// }

// 	return organizations, nil
// }

// func AssociateOrgWithFighter() ([]Fighter, error) {
// 	fighters := []Fighter{}
// 	if err := db.Model(&Fighter{}).Find(&fighters).Error; err != nil {
// 		return nil, err
// 	}

// 	for i := range fighters {
// 		if err := fighters[i].GetOrganization(); err != nil {
// 			return nil, err
// 		}
// 	}

// 	return fighters, nil

// }
