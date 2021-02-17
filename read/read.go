package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"
)

const (
	maxLen = 20
)

// Name struct to hold name
type Name struct {
	fname string
	lname string
}

// SetName Set full name to struct Name object
func (n *Name) SetName(fname string, lname string) {
	var aux []rune

	n.fname = fname
	if utf8.RuneCountInString(n.fname) > maxLen {
		aux = []rune(fname)
		n.fname = string(aux[:maxLen])
	}

	n.lname = lname
	if utf8.RuneCountInString(n.lname) > maxLen {
		aux = []rune(lname)
		n.lname = string(aux[:maxLen])
	}
}

// GetFirstName Get first name from struct Name object
func (n *Name) GetFirstName() string {
	return n.fname
}

// GetLastName Get last name from struct Name object
func (n *Name) GetLastName() string {
	return n.lname
}

// GetFullName Get full name from struct Name object
func (n *Name) GetFullName() string {
	return n.fname + " " + n.lname
}

func main() {
	var nameSlice []Name
	var usrInput string
	var line string

	// Get file name from user
	fmt.Println("Enter file name")
	_, err := fmt.Scan(&usrInput)
	if err != nil {
		panic(fmt.Sprint("Unable to read user input", err.Error()))
	}

	// Open file
	fd, err := os.Open(usrInput)
	if err != nil {
		panic(fmt.Sprint("Unable to open file", err.Error()))
	}
	defer fd.Close()

	// Parse file
	reader := bufio.NewReader(fd)
	for {
		line, err = reader.ReadString('\n')
		// This solves the runtime error if there is a empty line in the middle of the file
		if len(line) > 1 {
			splitRes := strings.Split(line, " ")                // Split line in two strings
			splitRes[1] = strings.TrimSuffix(splitRes[1], "\n") // Remove the line feed character fron second string

			newName := &Name{}
			newName.SetName(splitRes[0], splitRes[1]) // Add both names to Name struct
			nameSlice = append(nameSlice, *newName)   // Append result to slice of Names
		}

		if err != nil {
			break
		}
	}

	// Just to demonstrate going through all elements of the slice (can use direct printLn)
	for i := range nameSlice {
		fmt.Printf("First name: %s, Last name: %s\n", nameSlice[i].GetFirstName(), nameSlice[i].GetLastName())
	}
	// fmt.Println(nameSlice)
}
