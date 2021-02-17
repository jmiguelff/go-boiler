package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// GenDisplayFn generates a function that calculates the displacement at a give time
func GenDisplayFn(accel, vInit, sInit float64) func(float64) float64 {
	fn := func(t float64) float64 {
		return (0.5*accel*math.Pow(t, 2) + vInit*t + sInit)
	}
	return fn
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	sli := make([]float64, 0)

	fmt.Println("Insert values for acceleration, initial velocity and initial displacement separated by white spaces")
	scanner.Scan()
	usrInput := scanner.Text()

	substrings := strings.Split(usrInput, " ")
	if len(substrings) != 3 {
		panic(fmt.Sprint("Need three numbers!"))
	}

	for i := 0; i < len(substrings); i++ {
		n, err := strconv.ParseFloat(substrings[i], 64)
		if err != nil {
			panic(fmt.Sprint("Error, only accept numbers - ", err.Error()))
		}
		sli = append(sli, n)
	}

	fmt.Println("Enter a value for time to compute the displacement")
	var t float64
	_, err := fmt.Scan(&t)
	if err != nil {
		panic(fmt.Sprint("Unable to read user input", err.Error()))
	}

	fn := GenDisplayFn(sli[0], sli[1], sli[2])
	fmt.Println(fn(t))
}
