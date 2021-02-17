package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Animal struct type to hold animal properties
type Animal struct {
	food       string
	locomotion string
	sound      string
}

// Eat prints what an animal eats
func (a *Animal) Eat() {
	fmt.Println(a.food)
}

// Move prints how an animal moves
func (a *Animal) Move() {
	fmt.Println(a.locomotion)
}

// Speak prints the sound an animal makes
func (a *Animal) Speak() {
	fmt.Println(a.sound)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Data base of animals
	animalsDb := make(map[string]*Animal)
	animalsDb["cow"] = &Animal{food: "grass", locomotion: "walk", sound: "moo"}
	animalsDb["bird"] = &Animal{food: "worms", locomotion: "fly", sound: "peep"}
	animalsDb["snake"] = &Animal{food: "mice", locomotion: "slither", sound: "hsss"}

	for {
		fmt.Print(">")
		scanner.Scan()
		usrInput := scanner.Text()
		substrings := strings.Split(usrInput, " ")
		if len(substrings) != 2 {
			fmt.Println("Invalid request, try again (need two strings")
			continue
		}

		// Check if animal exists in the db
		if _, ok := animalsDb[substrings[0]]; !ok {
			fmt.Println("That animal does not exist in the database select between [cow][bird][snake")
			continue
		}

		switch substrings[1] {
		case "eat":
			animalsDb[substrings[0]].Eat()
		case "move":
			animalsDb[substrings[0]].Move()
		case "speak":
			animalsDb[substrings[0]].Speak()
		default:
			fmt.Println("Field not available please select between [eat][move][speak]")
		}
	}
}
