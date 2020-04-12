// TODO: package docs
package main

import (
	"fmt"

	"github.com/bradford-hamilton/dora/pkg/dora"
)

func main() {
	c, err := dora.NewFromString(testJSONObject)
	if err != nil {
		fmt.Printf("\nError creating client: %v\n", err)
	}

	result, err := c.GetByPath("$.item1[2].innerKey.anotherValue")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}

// TODO
// First query to add to tests tomorrow that's working (parsing into queryTokens): "$.item1[2].innerKey"
// $.item1[2].innerKey.anotherValue however does not work
// Since you can iterpolate at the call site, going to start with _only_ dot notation for object and bracket notation for arrays
// For now I've whipped up a much more juvenile parser for the queries because I don't think it's going to need to do much

const testJSONArray = `[
	"item1",
	"item2",
	{"item3": "item3value", "item4": {"innerkey": "innervalue"}},
	["item1", ["array"]]
]`

const testJSONObject = `{
	"item1": ["aryitem1", "aryitem2", {"some": "object"}],
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
