package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Fighter struct {
	FirstName string
	LastName  string
	Weight    int
	Southpaw  string
	Reach     int
}

var fighters []Fighter

func main() {
	fmt.Println("Welcome to the fighter information storing application.")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nPress 1 to enter fighter information")
		fmt.Println("Press 2 to print all Bantamweight fighters")
		fmt.Println("Press 3 to print all Featherweight fighters")
		fmt.Println("Press 4 to print all Lightweight fighters")
		fmt.Println("Press 5 to print all Welterweight fighters")
		fmt.Println("Press 6 to print all fighter info")
		fmt.Println("Press 7 to end the program")

		choice, _ := reader.ReadString('\n')

		switch choice {
		case "1\n":
			fighter := Fighter{}

			fmt.Print("Please enter fighter first name: ")
			fighter.FirstName, _ = reader.ReadString('\n')
			//error check
			fighter.FirstName = strings.TrimSuffix(fighter.FirstName, "\n") // Trim the newline character to get the right length
			if len(fighter.FirstName) < 2 {
				fmt.Print("Please enter a name longer than two characters")
				continue
			}

			fmt.Print("\nPlease enter fighter last name: ")
			fighter.LastName, _ = reader.ReadString('\n')
			//error check
			fighter.LastName = strings.TrimSuffix(fighter.LastName, "\n")
			if len(fighter.LastName) < 2 {
				fmt.Print("Error, please enter a last name longer than two characters")
				continue
			}

			fmt.Print("\nPlease enter fighter weight: ")
			fmt.Scanf("%d", &fighter.Weight)
			//error check
			if fighter.Weight < 135 || fighter.Weight > 170 {
				fmt.Print("Error, please enter a weight either between or exactly 135 and 170")
				continue
			}

			fmt.Print("\nPlease enter fighter reach: ")
			fmt.Scanf("%d", &fighter.Reach)
			//error check
			if fighter.Reach < 65 || fighter.Reach > 85 {
				fmt.Print("Error, please enter a reach either between or exactly 65 and 85")
				continue
			}

			fmt.Print("\nPlease enter if fighter is southpaw or orthodox: ")
			fighter.Southpaw, _ = reader.ReadString('\n')
			fighter.Southpaw = strings.TrimSuffix(fighter.Southpaw, "\n")
			//error check
			if fighter.Southpaw != "orthodox" && fighter.Southpaw != "southpaw" {
				fmt.Print("Error, please enter either southpaw or orthodox")
				continue
			}

			fighters = append(fighters, fighter)

		case "2\n":
			printByWeight("135")
		case "3\n":
			printByWeight("145")
		case "4\n":
			printByWeight("155")
		case "5\n":
			printByWeight("170")
		case "6\n":
			printAllFighters()
		case "7\n":
			fmt.Println("Exiting the program.")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func printByWeight(weightClass string) {
	weightClass = strings.TrimSpace(weightClass)
	fmt.Printf("\nFighters in the %s division:\n", weightClass)
	for _, fighter := range fighters {
		weightStr := strconv.Itoa(fighter.Weight)
		weightStr = strings.TrimSpace(weightStr)
		if weightStr == weightClass {
			fmt.Printf("First Name: %s\n", fighter.FirstName)
			fmt.Printf("Last Name: %s\n", fighter.LastName)
			fmt.Printf("Weight: %d\n", fighter.Weight)
			fmt.Printf("Southpaw: %s\n", fighter.Southpaw)
			fmt.Printf("Reach: %d\n", fighter.Reach)
			fmt.Println("--------------------")
		}
	}
}

func printAllFighters() {
	fmt.Println("\nAll Fighters:")
	for _, fighter := range fighters {
		fmt.Printf("First Name: %s\n", fighter.FirstName)
		fmt.Printf("Last Name: %s\n", fighter.LastName)
		fmt.Printf("Weight: %d\n", fighter.Weight)
		fmt.Printf("Southpaw: %s\n", fighter.Southpaw)
		fmt.Printf("Reach: %d\n", fighter.Reach)
		fmt.Println("--------------------")
	}
}
