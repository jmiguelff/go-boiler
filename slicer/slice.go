package main

import (
	"fmt"
	"sort"
	"strconv"
)

func main() {
	// I assumed by length the problem meant capacity...
	mysli := make([]int, 0, 3)
	var usrInput string
	for usrInput != "X" {
		fmt.Println("Write an integer:")
		_, err := fmt.Scan(&usrInput)
		if err != nil {
			panic(fmt.Sprint("Unable to read user input", err.Error()))
		}

		intValue, err := strconv.Atoi(usrInput)
		if err != nil {
			continue
		}

		mysli = append(mysli, intValue)
		sort.Ints(mysli)
		fmt.Println(mysli)
	}
}
