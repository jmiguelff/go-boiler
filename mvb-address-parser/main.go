package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type MVBLine struct {
	Timestamp uint32
	Address   uint16
	FrameSize uint8
}

func main() {
	// Open the file
	file, err := os.Open("console.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Local variables
	var data []MVBLine
	var lineCount int64

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate over each line
	for scanner.Scan() {
		line := scanner.Text()
		// Process each line here
		lineCount++
		// fmt.Println(line)

		var timestamp int
		var addr int64
		var frameSize int

		// Find the index of "ts="
		index := strings.Index(line, "ts=")
		if index != -1 {
			// Get the substring after "ts="
			substring := line[index+len("ts=") : index+len("ts=")+len("62940941095")]

			// Convert the substring to an integer
			timestamp, err = strconv.Atoi(substring)
			if err != nil {
				log.Fatal(err)
			}

			// Use the integer value here
			//fmt.Println(timestamp)
		}

		// Find the index of "addr="
		index = strings.Index(line, "addr=")
		if index != -1 {
			// Get the substring after "addr="
			substring := line[index+len("addr=") : index+len("addr=")+len("000")]
			// Convert the substring to an integer
			addr, err = strconv.ParseInt(substring, 16, 64)
			if err != nil {
				log.Fatal(err)
			}

			// Use the integer value here
			// fmt.Printf("%03x\n", addr)
		}

		// Find size of data (hacky)
		indexInital := strings.Index(line, "kProcessData")
		indexFinal := strings.Index(line, "Bit,")
		if indexInital != -1 && indexFinal != -1 {
			// Get the substring after "addr="
			substring := line[indexInital+len("KProcessData") : indexFinal]
			// Convert the substring to an integer
			frameSize, err = strconv.Atoi(substring)
			if err != nil {
				log.Fatal(err)
			}

			// Use the integer value here
			//fmt.Println(frameSize)
		}

		// Append the MVBLine struct to the slice
		data = append(data, MVBLine{
			Timestamp: uint32(timestamp),
			Address:   uint16(addr),
			FrameSize: uint8(frameSize),
		})
	}
	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Print the number of lines read
	fmt.Printf("Number of lines read: %d\n", lineCount)

	// Subtract the first timestamp to the last timestamp
	milliseconds := data[len(data)-1].Timestamp - data[0].Timestamp
	delta := time.Duration(milliseconds) * time.Microsecond
	fmt.Printf("Sample duration: %.2fs\n", delta.Seconds())

	// For each address, print the number of times it appears
	fmt.Println("Count the number of times each address appears")
	counts := make(map[uint16]int)
	for _, v := range data {
		counts[v.Address]++
	}
	for k, v := range counts {
		fmt.Printf("%03x: %d\n", k, v)
	}

	fmt.Println("Remove duplicates and print")

	// Copy data slice to other slice
	var dataCopy []MVBLine

	// Copy elements from data to dataCopy without repeating elements with the same address
	for _, v := range data {
		// Check if the address is already in dataCopy
		found := false
		for _, w := range dataCopy {
			if v.Address == w.Address {
				found = true
				break
			}
		}
		if !found {
			dataCopy = append(dataCopy, v)
		}
	}

	// Order dataCopy by address
	sort.Slice(dataCopy, func(i, j int) bool {
		return dataCopy[i].Address < dataCopy[j].Address
	})

	// Print only the address of each element in dataCopy
	for _, v := range dataCopy {
		fmt.Printf("%03x\n", v.Address)
	}
}
