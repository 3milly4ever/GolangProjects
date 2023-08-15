package types

import "fmt"

type User struct {
	ID         int    `json:"id"`
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
func NewUser(ID int, Name string, Email string, Weight int, Goal string, Regimen string, DateJoined string) *User {
	var user User
	links := CreateUserHypermediaLinks(user.ID)
	user.Links = links

	return &User{
		ID:         ID,
		Name:       Name,
		Email:      user.Email,
		Weight:     user.Weight,
		Goal:       user.Goal,
		Regimen:    user.Regimen,
		DateJoined: user.DateJoined,
		Links:      links,
	}
}

// creates hypermedialinks
func CreateUserHypermediaLinks(id int) UserLinks {
	baseURL := "/users/"
	return UserLinks{
		Self:   fmt.Sprintf("%s%d", baseURL, id),
		Update: fmt.Sprintf("%s%d/update", baseURL, id),
		Delete: fmt.Sprintf("%s%d/delete", baseURL, id),
	}
}
