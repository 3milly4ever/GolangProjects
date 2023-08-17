package types

import "fmt"

type User struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Weight     int    `json:"weight"`
	Goal       string `json:"goal"`
	Regimen    string `json:"regimen"`
	DateJoined string `json:"date_joined"`

	Links UserLinks `json:"links"`
}

type UserLinks struct {
	Self   string `json:"self"`
	Update string `json:"update"`
	Delete string `json:"delete"`
}

// creates a new user
func NewUser(ID int64, Name string, Email string, Weight int, Goal string, Regimen string, DateJoined string, links UserLinks) *User {
	//	var user User
	//links := CreateUserHypermediaLinks(ID)
	//	user.Links = links

	return &User{
		ID:         ID,
		Name:       Name,
		Email:      Email,
		Weight:     Weight,
		Goal:       Goal,
		Regimen:    Regimen,
		DateJoined: DateJoined,
		Links:      links,
	}
}

// creates hypermedialinks
func CreateUserHypermediaLinks(id int64) UserLinks {
	baseURL := "/users/"
	return UserLinks{
		Self:   fmt.Sprintf("%s%d", baseURL, id),
		Update: fmt.Sprintf("%s%d", baseURL, id),
		Delete: fmt.Sprintf("%s%d", baseURL, id),
	}
}
