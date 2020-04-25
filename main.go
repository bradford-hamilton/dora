// TODO: package docs
package main

import (
	"fmt"

	"github.com/bradford-hamilton/dora/pkg/dora"
)

// Notes:
// Since you can iterpolate at the call site, going to start with _only_ dot notation for object and bracket notation for arrays
// For now I've whipped up a much more juvenile parser for the queries because it doesn't need to do much right now

// Currently using main as my own testing ground as if dora was 3rd party
func main() {
	c, err := dora.NewFromString(testJSONObject)
	if err != nil {
		fmt.Printf("\nError creating client: %v\n", err)
	}

	result, err := c.Get("$.item1[2].some")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}

const testJSONArray = `[
	"item1",
	"item2",
	{"item3": "item3value", "item4": {"innerkey": "innervalue"}},
	["item1", ["array"]]
]`

const testJSONObject = `{
	"item1": ["aryitem1", "aryitem2", {"some": {"thing": "coolObj"}}],
	"item2": "simplestringvalue",
	"item3": {
		"item4": {
			"item5": {
				"item6": ["thing1", 2],
				"item7": {"reallyinnerobjkey": {"is": "anobject"}}
			}
		}
	}
}`
