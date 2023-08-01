package main

import (
	"fmt"
	"time"
)

// package level variables defined at the top outside all functions
var conferenceName = "Go Conference"

const conferenceTickets = 50

var remainingTickets uint = 50
var bookings = make([]UserData, 0) //a slice/list of strings it's dynamic so size can be 0 and increase by itself

type UserData struct {
	firstName       string
	lastName        string
	email           string
	numberOfTickets uint
	//isOptedInForNewsletter bool
}

//we dont put all the variables here because the general practice is to define variables as local as possible, aka within its own scope where it is used

func main() {

	//%T prints the types of the variables
	//uint can not be negative

	greetUsers()

	for {

		firstName, lastName, email, userTickets := getUserInput()
		isValidName, isValidEmail, isValidTicketNumber := ValidateUserInput(firstName, lastName, email, userTickets, remainingTickets)

		if isValidName && isValidEmail && isValidTicketNumber {
			bookTicket(userTickets, firstName, lastName, email)
			sendTicket(userTickets, firstName, lastName, email)

			firstNames := getFirstNames()
			fmt.Printf("The first names of the bookings are: %v\n", firstNames)

			// fmt.Printf("The whole array: %v\n", bookings)
			// fmt.Printf("The first value type: %v\n", bookings[0])
			// fmt.Printf("Array type: %T\n", bookings)
			// fmt.Printf("Array length: %v\n", len(bookings))

			noTicketsRemaining := remainingTickets == 0

			if noTicketsRemaining {
				//end the program
				fmt.Println("All tickets are sold come back next year, sorry and thank you.")
				break
			}

		} else {
			if !isValidName {
				fmt.Println("Please enter a name that is more than or equal to two characters in length.")
			}
			if !isValidEmail {
				fmt.Println("Please make sure you enter a valid email.")
			}
			if !isValidTicketNumber {
				fmt.Println("Please make sure you're not trying to purchase more tickets than there are available.")
			}
		}

	}
}

//functions

func greetUsers() {
	fmt.Printf("Welcome to %v booking application\n", conferenceName)
	fmt.Printf("We have a total of %v tickets and %v are still available.\n", conferenceTickets, remainingTickets)
	fmt.Println("Get your tickets here to attend.")
}

func getFirstNames() []string {
	firstNames := []string{}
	//to iterate through a slice we need a range expression
	//for arrays and slices, range provides the index and value for each element
	//this is a nested for loop below

	for _, booking := range bookings {
		//strings.Fields() splits the string with white space as a separator
		//var names = strings.Fields(booking) //this names will be an array containing the first name and the last name as separate strings ((to extract from a slice)
		firstNames = append(firstNames, booking.firstName)
	}
	return firstNames
}

func getUserInput() (string, string, string, uint) {

	var firstName string
	var lastName string
	var email string
	var userTickets uint

	//ask user for their name and tickets
	fmt.Print("Enter your first name: ")
	fmt.Scan(&firstName)

	fmt.Print("\nEnter your last name: ")
	fmt.Scan(&lastName)

	fmt.Print("\nEnter your email: ")
	fmt.Scan(&email)

	fmt.Print("\nEnter the amount of tickets to purchase: ")
	fmt.Scan(&userTickets)

	return firstName, lastName, email, userTickets
}

func bookTicket(userTickets uint, firstName string, lastName string, email string) {
	//arrays in go have a fixed size
	remainingTickets = remainingTickets - userTickets

	//create a struct for a user
	var userData = UserData{
		firstName:       firstName,
		lastName:        lastName,
		email:           email,
		numberOfTickets: userTickets,
	}

	bookings = append(bookings, userData)
	fmt.Printf("List of bookings is %v\n", bookings)

	fmt.Printf("Thank you %v %v for booking %v tickets. You will receive a confirmation email at %v\n", firstName, lastName, userTickets, email)
	fmt.Printf("%v tickets remaining for %v\n", remainingTickets, conferenceName)
}

func sendTicket(userTickets uint, firstName string, lastName string, email string) {
	time.Sleep(10 * time.Second)
	var ticket = fmt.Sprintf("%v tickets for %v %v", userTickets, firstName, lastName)
	fmt.Println("--------------")
	fmt.Printf("Sending ticket:\n %v to email address %v", ticket, email)
	fmt.Println("\n--------------")
}
