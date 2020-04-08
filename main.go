// TODO: package docs
package main

import (
	"fmt"

	"github.com/bradford-hamilton/parsejson/pkg/parsejson"
)

func main() {
	pc, err := parsejson.NewFromString(testJSON)
	if err != nil {
		fmt.Printf("\nError creating client: %v\n", err)
	}

	fmt.Println(pc)
}

// Some test json
const testJSON = `{
	"items": {
		"item": [{
			"id": "0001",
			"type": "donut",
			"name": "Cake",
			"cpu": 55,
			"batters": {
				"batter": [{
					"id": false,
					"name": null,
					"fun": true
				}]
			},
			"names": ["catstack", "lampcat", "langlang"]
		}]
	},
	"version": 0.1,
	"number": 11.4,
	"negativeNum": -5
}`
