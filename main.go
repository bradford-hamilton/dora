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

const testJSON = `{ "key0": [], "key1": "simplestringvalue" }`

// const testJSON = `{
// 	"thing": [{
// 		"insidekey": "value"
// 	}]
// }`

// const testJSON = `{
// 	"thing": [{
// 		"insidekey": "value",
// 		"insidekey2": "value"
// 	}]
// }`

// Some test json
// const testJSON = `{
// 	"id": "0001",
// 	"names": ["catstack"]
// }`

// const testJSON = `{
// 	"id": "0001",
// 	"batters": { "ling": "lang" }
// }`
