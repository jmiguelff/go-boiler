package main

import (
	"fmt"
	"time"

	"github.com/stianeikeland/go-rpio"
)

func main() {
	// Openg gpio, panic in case of error
	fmt.Println("Opening gpio")
	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpip", err.Error()))
	}

	// Close gpio after usage
	defer rpio.Close()

	// Configure pin as output
	pin := rpio.Pin(18)
	pin.Output()

	// Togle ping 20 times
	for x := 0; x < 20; x++ {
		pin.Toggle()
		time.Sleep(time.Second / 5)
	}
}
