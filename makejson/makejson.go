package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	m := make(map[string]string)

	fmt.Println("Insert name pls")
	scanner.Scan()
	usrInput := scanner.Text()
	m["name"] = usrInput

	fmt.Println("Insert address pls")
	scanner.Scan()
	usrInput = scanner.Text()
	m["address"] = usrInput

	barr, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprint("Fail json marshalling", err.Error()))
	}
	fmt.Printf("%s", barr)
}
