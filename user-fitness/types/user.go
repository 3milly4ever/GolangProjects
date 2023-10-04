package types

import (
	"fmt"

	"github.com/google/uuid"
)

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

type UserCredentials struct {
	ID       uuid.UUID
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLinks struct {
	Self   string `json:"self"`
	Update string `json:"update"`
	Delete string `json:"delete"`
}

type PaginatedUserResponse struct {
	Users       []User          `json:"users"`
	TotalUsers  int             `json:"totalUsers"`
	TotalPages  int             `json:"TotalPages"`
	CurrentPage int             `json:"CurrentPage"`
	PageSize    int             `json:"pageSize"`
	Links       PaginationLinks `json:"links"`
}

type PaginationLinks struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
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

// creates pagination links
func CreatePaginationLinks(baseURL string, currentPage, pageSize, totalUsers int) PaginationLinks {
	var paginationLinks PaginationLinks

	//Calculate total number of page
	totalPages := (totalUsers + pageSize - 1) / pageSize

	//Generate links for the first page, last page, and next page
	paginationLinks.First = fmt.Sprintf("%s?page=1&pageSize=%d", baseURL, pageSize)
	paginationLinks.Last = fmt.Sprintf("%s?page=%d&pageSize=%d", baseURL, totalPages, pageSize)
	//if it is not page 0 then it has a previous page
	if currentPage > 1 {
		paginationLinks.Prev = fmt.Sprintf("%s?page=%d&pageSize=%d", baseURL, currentPage-1, pageSize)
	}
	//if it is not equal to total pages, it is not the last page so it has a next. otherwise the next field would be omitted
	if currentPage < totalPages {
		paginationLinks.Next = fmt.Sprintf("%s?page=%d&pageSize=%d", baseURL, currentPage+1, pageSize)
	}

	return paginationLinks
}
