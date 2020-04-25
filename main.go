// TODO: package docs
package main

import (
	"fmt"

	"github.com/bradford-hamilton/dora/pkg/dora"
)

// Notes:
// Since you can iterpolate at the call site, going to start with _only_ dot notation for object and bracket notation for arrays
// For now I've whipped up a much more juvenile parser for the queries because it doesn't need to do a whole lot

// Currently using main as an example to follow
func main() {
	var exampleJSON = `{ "string": "a neat string", "bool": true, "PI": 3.14159 }`

	c, err := dora.NewFromString(exampleJSON)
	if err != nil {
		fmt.Printf("\nError creating client: %v\n", err)
	}

	str, err := c.GetString("$.string")
	if err != nil {
		fmt.Println(err)
	}

	boolean, err := c.GetBool("$.bool")
	if err != nil {
		fmt.Println(err)
	}

	float, err := c.GetFloat64("$.PI")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(str)     // a neat string
	fmt.Println(boolean) // true
	fmt.Println(float)   // 3.14159
}
