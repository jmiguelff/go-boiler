package main

import (
	"fmt"
)

func main() {
	var usrInput float64

	fmt.Println("Insert float number")
	_, err := fmt.Scan(&usrInput)
	if err != nil {
		panic(fmt.Sprint("unable to read user input", err.Error()))
	}
	fmt.Println(int64(usrInput))
}
