package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Swap swaps the element at idx with the element at idx + 1
func Swap(sli []int, idx int) {
	if idx >= len(sli)-1 {
		return
	}

	a := sli[idx+1]
	sli[idx+1] = sli[idx]
	sli[idx] = a
}

// BubbleSort order a slice using bubble sort algorithm
func BubbleSort(sli []int) {
	for i := 0; i < len(sli); i++ {
		for j := 0; j < len(sli)-i-1; j++ {
			if sli[j] > sli[j+1] {
				Swap(sli, j)
			}
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	mysli := make([]int, 0)

	fmt.Println("Insert a sequence of integers separated by spaces (up to 10)")
	scanner.Scan()
	usrInput := scanner.Text()

	substrings := strings.Split(usrInput, " ")
	nOfInts := len(substrings)
	for i := 0; i < 10 && i < nOfInts; i++ {
		intValue, err := strconv.Atoi(substrings[i])
		if err != nil {
			panic(fmt.Sprint("Error, only accept numbers - ", err.Error()))
		}
		mysli = append(mysli, intValue)
	}
	BubbleSort(mysli)
	fmt.Println(mysli)
}
