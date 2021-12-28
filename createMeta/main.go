package main

import (
	"encoding/json"
	"os"
	"strconv"
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
	d := Metadata{
		Name:        "IOC#",
		Description: "Fucking your APE since this morning",
		Image:       "ipfs://QmYGva6afU3obmwYSQd7pqfJN5nAAxExBbyU4Ur2s8Wyyn",
		Attributes: []Attribute{
			{
				Trait_type: "background",
				Value:      "notYet",
			},
			{
				Trait_type: "feet",
				Value:      "notYet",
			},
			{
				Trait_type: "body",
				Value:      "notYet",
			},
			{
				Trait_type: "eyes",
				Value:      "notYet",
			},
			{
				Trait_type: "arms",
				Value:      "notYet",
			},
			{
				Trait_type: "clothes",
				Value:      "notYet",
			},
			{
				Trait_type: "hat",
				Value:      "notYet",
			},
			{
				Trait_type: "extra",
				Value:      "notYet",
			},
			{
				Trait_type: "glasses",
				Value:      "notYet",
			},
			{
				Trait_type: "mouth",
				Value:      "notYet",
			},
		},
	}

	for i := 0; i < 100; i++ {
		idx := strconv.Itoa(i)
		d.Name = "IOC#" + idx
		out, _ := json.Marshal(d)

		fileName := "data/" + idx + ".json"
		err := os.WriteFile(fileName, out, 0644)
		if err != nil {
			panic(err)
		}
	}
}
