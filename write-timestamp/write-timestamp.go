package main

import (
	"bufio"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Create file
	f, err := os.Create("timestamp.log")
	check(err)

	// After being used auto close
	defer f.Close()

	for {
		// Get time from system
		t := time.Now()

		// Write to file
		w := bufio.NewWriter(f)
		_, err := w.WriteString(t.String() + "\n")
		check(err)

		w.Flush()

		// Wait 1 second
		time.Sleep(10 * time.Second)
	}
}
