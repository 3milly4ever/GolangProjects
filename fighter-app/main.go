package main

import (
	"fmt"
)

//var orgName = "UFC"

//const maxNumOfFighters = 40

//var curNumOfFighters = 0

//var userOption ints

var bantamWeight = make([]FighterInfo, 0)
var featherWeight = make([]FighterInfo, 0)
var lightWeight = make([]FighterInfo, 0)
var welterWeight = make([]FighterInfo, 0)

type FighterInfo struct {
	firstName string
	lastName  string
	weight    int
	southpaw  string
	reach     uint
}

func main() {

	greetUsers()

	//firstName, lastName, weight, southpaw, reach := getUserInput()

	var firstName string
	var lastName string
	var weight int
	var southpaw string
	var reach uint
	var option int = 0

	for option != 7 {

		fmt.Println("Press 1 to enter fighter information ")
		fmt.Println("Press 2 to print all Bantamweight fighters")
		fmt.Println("Press 3 to print all Featherweight fighters")
		fmt.Println("Press 4 to print all Lightweight fighters")
		fmt.Println("Press 5 to print all Welterweight fighters")
		fmt.Println("Press 6 to print all fighter info")
		fmt.Println("Press 7 to end the program")
		fmt.Scanf("%d", &option)

		switch option {

		case 1:
			fmt.Print("Please enter fighter first name: ")
			fmt.Scan(&firstName)
			fmt.Print("\nPlease enter fighter last name: ")
			fmt.Scan(&lastName)
			fmt.Print("\nPlease enter fighter weight: ")
			fmt.Scan(&weight)
			fmt.Print("\nPlease enter please enter fighter reach: ")
			fmt.Scanln(&reach)
			fmt.Print("\nPlease enter if fighter is southpaw or orthodox: ")
			fmt.Scan(&southpaw)
			var discardNewline string
			fmt.Scan(&discardNewline)
			//getUserInput(firstName, lastName, weight, reach, southpaw)
			isValidName, isValidWeight, isValidReach, isValidSouthPawInput := validateInput(firstName, lastName, weight, reach, southpaw)

			if !isValidName {
				fmt.Println("Please enter a name that is more than or equal to two characters in length")
				continue
			}
			if !isValidWeight {
				fmt.Println("Please enter a weight that is 135, 145, 155 or 170")
				continue
			}
			if !isValidReach {
				fmt.Println("Please enter a reach that's between or is 65 and 85")
				continue
			}
			if !isValidSouthPawInput {
				fmt.Println("Error, please enter southpaw or orthodox")
				continue
			}
			if isValidName && isValidWeight && isValidReach && isValidSouthPawInput {
				populateFighters(firstName, lastName, weight, southpaw, reach)
			}
		case 2:
			printAllBantamweights()
		case 3:
			printAllFeatherweights()
		case 4:
			printAllLightweights()
		case 5:
			printAllWelterweights()
		case 6:
			printAllFighterInfo()
		case 7:
			fmt.Println("Thanks for using the program, have a good one")
		default:
			fmt.Println("Error, please enter 1, 2, 3, 4, 5, 6, or 7")
		}
	}
	//return firstName, lastName, weight, southpaw, reach
}

// func getUserInput(firstName string, lastName string, weight uint, reach uint, southpaw string) {

// }

func validateInput(firstName string, lastName string, weight int, reach uint, southpaw string) (bool, bool, bool, bool) {

	isValidName := len(firstName) >= 2 && len(lastName) >= 2
	isValidWeight := weight >= 135 && weight <= 170
	isValidReach := reach >= 65 && reach <= 85
	isValidSouthPawInput := southpaw == "southpaw" || southpaw == "orthodox"

	return isValidName, isValidWeight, isValidReach, isValidSouthPawInput

}

// func deleteFighters() {

// }

func greetUsers() {
	fmt.Println("Welcome to the fighter information storing application.")
}

func populateFighters(firstName string, lastName string, weight int, southpaw string, reach uint) {

	var fighter = FighterInfo{
		firstName: firstName,
		lastName:  lastName,
		weight:    weight,
		southpaw:  southpaw,
		reach:     reach,
	}

	switch fighter.weight {
	case 135:
		bantamWeight = append(bantamWeight, fighter)
	case 145:
		featherWeight = append(featherWeight, fighter)
	case 155:
		lightWeight = append(lightWeight, fighter)
	case 170:
		welterWeight = append(welterWeight, fighter)
	default:
		fmt.Println("Invalid weight.")
	}
}

func printAllFighterInfo() {

	for _, fighter := range bantamWeight {
		fmt.Println("First Name:", fighter.firstName)
		fmt.Println("Last Name:", fighter.lastName)
		fmt.Println("Weight:", fighter.weight)
		fmt.Println("Southpaw:", fighter.southpaw)
		fmt.Println("Reach:", fighter.reach)
		fmt.Println("--------------------")
	}

	for _, fighter := range featherWeight {
		fmt.Println("First Name:", fighter.firstName)
		fmt.Println("Last Name:", fighter.lastName)
		fmt.Println("Weight:", fighter.weight)
		fmt.Println("Southpaw:", fighter.southpaw)
		fmt.Println("Reach:", fighter.reach)
		fmt.Println("--------------------")
	}

	for _, fighter := range lightWeight {
		fmt.Println("First Name:", fighter.firstName)
		fmt.Println("Last Name:", fighter.lastName)
		fmt.Println("Weight:", fighter.weight)
		fmt.Println("Southpaw:", fighter.southpaw)
		fmt.Println("Reach:", fighter.reach)
		fmt.Println("--------------------")
	}

	for _, fighter := range welterWeight {
		fmt.Println("First Name:", fighter.firstName)
		fmt.Println("Last Name:", fighter.lastName)
		fmt.Println("Weight:", fighter.weight)
		fmt.Println("Southpaw:", fighter.southpaw)
		fmt.Println("Reach:", fighter.reach)
		fmt.Println("--------------------")
	}
}

func printAllBantamweights() {
	for _, fighter := range bantamWeight {
		fmt.Println("First Name:", fighter.firstName)
		fmt.Println("Last Name:", fighter.lastName)
		fmt.Println("Weight:", fighter.weight)
		fmt.Println("Southpaw:", fighter.southpaw)
		fmt.Println("Reach:", fighter.reach)
		fmt.Println("--------------------")
	}
}

func printAllFeatherweights() {
	for _, fighter := range featherWeight {
		fmt.Println("First Name:", fighter.firstName)
		fmt.Println("Last Name:", fighter.lastName)
		fmt.Println("Weight:", fighter.weight)
		fmt.Println("Southpaw:", fighter.southpaw)
		fmt.Println("Reach:", fighter.reach)
		fmt.Println("--------------------")
	}
}

func printAllLightweights() {
	for _, fighter := range lightWeight {
		fmt.Println("First Name:", fighter.firstName)
		fmt.Println("Last Name:", fighter.lastName)
		fmt.Println("Weight:", fighter.weight)
		fmt.Println("Southpaw:", fighter.southpaw)
		fmt.Println("Reach:", fighter.reach)
		fmt.Println("--------------------")
	}
}

func printAllWelterweights() {
	for _, fighter := range welterWeight {
		fmt.Println("First Name:", fighter.firstName)
		fmt.Println("Last Name:", fighter.lastName)
		fmt.Println("Weight:", fighter.weight)
		fmt.Println("Southpaw:", fighter.southpaw)
		fmt.Println("Reach:", fighter.reach)
		fmt.Println("--------------------")
	}
}
