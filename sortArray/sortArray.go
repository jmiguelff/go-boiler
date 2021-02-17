package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func minOf(vars ...int) int {
	min := vars[0]

	for _, i := range vars {
		if min > i {
			min = i
		}
	}
	return min
}

func removeFirstElement(slice []int) []int {
	return slice[1:]
}

func mergeAndSort(a []int, b []int) []int {
	outSli := make([]int, 0)

	for len(a) != 0 && len(b) != 0 {
		min := minOf(a[0], b[0])
		outSli = append(outSli, min)

		if min == a[0] {
			a = removeFirstElement(a)
		} else {
			b = removeFirstElement(b)
		}
	}

	if len(a) != 0 {
		outSli = append(outSli, a...)
	} else {
		outSli = append(outSli, b...)
	}

	return outSli
}

func sortSlice(sli []int, wg *sync.WaitGroup) {
	sort.Ints(sli)
	fmt.Println(sli)
	wg.Done()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	sli := make([]int, 0)

	fmt.Println("Insert a sequence of integers separated by spaces")

	// Rear user input and convert to an int slice
	scanner.Scan()
	usrInput := scanner.Text()
	substrings := strings.Split(usrInput, " ")
	nOfInts := len(substrings)
	if nOfInts < 4 {
		panic(fmt.Sprint("Error, require at least four elements"))
	}

	for i := 0; i < nOfInts; i++ {
		intValue, err := strconv.Atoi(substrings[i])
		if err != nil {
			panic(fmt.Sprint("Error, only accept numbers - ", err.Error()))
		}
		sli = append(sli, intValue)
	}

	// Break slice into 4 slices (last array gets the additional elements)
	nOfElements := nOfInts / 4
	a := sli[0:nOfElements]
	b := sli[nOfElements : nOfElements*2]
	c := sli[nOfElements*2 : nOfElements*3]
	d := sli[nOfElements*3:]

	var wg sync.WaitGroup
	wg.Add(4)
	go sortSlice(a, &wg)
	go sortSlice(b, &wg)
	go sortSlice(c, &wg)
	go sortSlice(d, &wg)
	wg.Wait()

	// Merge all slices
	o := mergeAndSort(a, b)
	o = mergeAndSort(o, c)
	o = mergeAndSort(o, d)

	fmt.Println(o)
}
