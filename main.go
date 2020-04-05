// TODO: package docs
package main

import (
	"fmt"

	"github.com/bradford-hamilton/parsejson/pkg/lexer"
)

func main() {
	l := lexer.New(testJSON)
	fmt.Println(l)
}

// Some test json
const testJSON = `{
	"items": {
		"item": [{
			"id": "0001",
			"type": "donut",
			"name": "Cake",
			"ppu": 0.55,
			"batters": {
				"batter": [{
					"id": "1001",
					"type": "Regular",
					"fun": "true"
				}]
			},
			"topping": [{
				"id": "5001",
				"type": "null",
				"fun": "false"
			}]
		}]
	}
}`
