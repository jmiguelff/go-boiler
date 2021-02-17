package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Insert string:")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	usrInput := scanner.Text()

	lw := strings.ToLower(usrInput)

	if strings.HasPrefix(lw, "i") && strings.HasSuffix(lw, "n") && strings.ContainsAny(lw, "a") {
		fmt.Println("Found!")
	} else {
		fmt.Println("Not Found!")
	}
}
