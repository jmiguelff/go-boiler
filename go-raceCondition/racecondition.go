package main

import (
	"fmt"
	"time"
)

// x global value subject to race condition
var x int

// incrementGlobal add value to global variable x
func incrementGlobal(incValue int) {
	for i := 0; i < incValue; i++ {
		x = x + 1
	}
}

// divideGlobal divide the global variable by a value and assign it to the global
func divideGlobal(div int) {
	x = x / div
}

func main() {
	// All the following routines will manipulate the same global variable
	// because they are being run without any synchronization mechanism
	// the final value of the global variable (x) is impossible to determine.
	// If you run the code using "go run -race racecondition.go" you will see
	// that golang is able to detect the race conditions.
	go incrementGlobal(10)
	go divideGlobal(10)
	go incrementGlobal(2)

	time.Sleep(1 * time.Second)
	fmt.Println(x)
}
