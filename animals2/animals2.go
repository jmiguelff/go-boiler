package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Animal interface has 3 methods
type Animal interface {
	Eat()
	Move()
	Speak()
}

// Cow type declaration
type Cow struct {
	name string
}

// Bird type declaration
type Bird struct {
	name string
}

// Snake type declaration
type Snake struct {
	name string
}

// Eat prints what a cow eats
func (c Cow) Eat() {
	fmt.Print("grass")
}

// Move prints how a cow moves
func (c Cow) Move() {
	fmt.Print("walk")
}

// Speak prints the sound a cow makes
func (c Cow) Speak() {
	fmt.Print("moo")
}

// Eat prints what a bird eats
func (b Bird) Eat() {
	fmt.Print("worms")
}

// Move prints how a bird moves
func (b Bird) Move() {
	fmt.Print("fly")
}

// Speak prints the sound a bird makes
func (b Bird) Speak() {
	fmt.Print("peep")
}

// Eat prints what a snake eats
func (s Snake) Eat() {
	fmt.Print("mice")
}

// Move prints how a snake moves
func (s Snake) Move() {
	fmt.Print("slither")
}

// Speak prints the sound a snake makes
func (s Snake) Speak() {
	fmt.Print("hsss")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	animals := make(map[string]Animal)

	for {
		fmt.Print(">")
		scanner.Scan()
		usrInput := scanner.Text()
		substrings := strings.Split(usrInput, " ")
		if len(substrings) != 3 {
			fmt.Println("Invalid request, try again (need three strings")
			continue
		}

		if substrings[0] == "newanimal" {
			// Check if name is taken
			if _, ok := animals[substrings[1]]; ok {
				fmt.Println("Name is taken choose other name")
				continue
			}
			switch strings.ToLower(substrings[2]) {
			case "cow":
				animals[substrings[1]] = Cow{substrings[1]}
				fmt.Println("Created it!")
			case "bird":
				animals[substrings[1]] = Bird{substrings[1]}
				fmt.Println("Created it!")
			case "snake":
				animals[substrings[1]] = Snake{substrings[1]}
				fmt.Println("Created it!")
			default:
				fmt.Println("Invalid type of animal, please choose between cow bird or snake")
			}
		} else if substrings[0] == "query" {
			// Check if name exists
			if _, ok := animals[substrings[1]]; !ok {
				fmt.Println("There is no animal with that name on the database")
				continue
			}
			switch strings.ToLower(substrings[2]) {
			case "eat":
				animals[substrings[1]].Eat()
				fmt.Println()
			case "move":
				animals[substrings[1]].Move()
				fmt.Println()
			case "speak":
				animals[substrings[1]].Speak()
				fmt.Println()
			default:
				fmt.Println("Action not available, please choose between eat, move or speak")
			}
		} else {
			fmt.Println("Invalid command, use newanimal or query commands")
		}
	}
}
