package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Attribute struct {
	Trait_type string `json:"trait_type"`
	Value      string `json:"value"`
}

type Metadata struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Image       string      `json:"image"`
	Attributes  []Attribute `json:"attributes"`
}

func main() {

	dat, err := os.ReadFile("data/metadata.json")
	if err != nil {
		panic(err)
	}

	var res map[string]Metadata
	json.Unmarshal(dat, &res)

	for key, val := range res {
		fmt.Print(key)
		fmt.Print(" -> ")
		fmt.Println(val)

		fileName := "data/" + key + ".json"
		out, _ := json.Marshal(val)
		err = os.WriteFile(fileName, out, 0644)
		if err != nil {
			panic(err)
		}
	}
}
