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

	res, err := c.GetByFullPath("$.item1[0]")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)
}

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
