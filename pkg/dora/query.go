package dora

import (
	"errors"
	"fmt"

	"github.com/bradford-hamilton/dora/pkg/ast"
)

// prepareQuery validates the query root, sets the query on the client struct, and parses the query into parsedQuery of some sort.
// TODO: later decide if this is useful
func (c *Client) prepareQuery(query string, rootNodeType ast.RootNodeType) error {
	if err := validateQueryRoot(query, c.program.Type); err != nil {
		return err
	}
	c.setQuery([]rune(query))
	if err := c.parseQuery(); err != nil {
		return err
	}
	return nil
}

func (c *Client) setQuery(query []rune) {
	c.query = query
}

func (c *Client) parseQuery() error {
	tokens, err := scanQueryTokens(c.query)
	if err != nil {
		return err
	}
	c.parsedQuery = tokens
	return nil
}

func (c *Client) executeQuery() error {
	rootVal := *c.program.RootValue
	obj, _ := rootVal.(ast.Object)
	arr, ok := rootVal.(ast.Array)
	currentType := "object"
	if ok {
		currentType = "array"
	}
	parsedQueryLen := len(c.parsedQuery)

	for i := 0; i < parsedQueryLen; i++ {
		// Final iteration
		if i == parsedQueryLen-1 {
			if currentType == "object" {
				r := c.parsedQuery[i].keyReq
				for _, v := range obj.Children {
					if r == v.Key.Value {
						switch val := v.Value.(type) {
						case ast.Literal:
							switch lit := val.Value.(type) {
							case string:
								c.result = lit
							case int:
								c.result = string(lit)
							case bool:
								c.result = fmt.Sprintf("%v", lit)
							case nil:
								c.result = "null"
							}

						case ast.Object:
							c.result = string(c.input[val.Start:val.End])
						case ast.Array:
							c.result = string(c.input[val.Start:val.End])
						}
					}
				}
			} else {
				// c.result = arr
			}
		}

		if c.parsedQuery[i].accessType == ObjectAccess {
			if currentType != "object" {
				return fmt.Errorf("TODO: error")
			}
			var found bool

			for _, v := range obj.Children {
				if v.Key.Value == c.parsedQuery[i].keyReq {
					found = true
					o, astObj := v.Value.(ast.Object)
					a, astArr := v.Value.(ast.Array)
					if astObj {
						obj = o
						currentType = "object"
						break
					}
					if astArr {
						arr = a
						currentType = "array"
						break
					}
				}
			}
			if !found {
				return fmt.Errorf("Sorry, could not find a key with that value. Key: %s", c.parsedQuery[i].keyReq)
			}
		} else {
			if currentType != "array" {
				return fmt.Errorf("TODO: error")
			}
			qt := c.parsedQuery[i]
			val := arr.Children[qt.indexReq]

			// TODO: Left off here. Think about:
			// case object: set obj = t
			// case array: set arr = t
			// case literal: check if it's the end of a query? If not isn't that an error?
			// Also think about how to return the final object/array/literal when on the last query in c.parsedQuery
			// 		- the above will mean I need to go back and add token locations of some sort
			//		- the purpose being so that if the query needs to return anything besides a single string (string or keyword (true,false,null)), it can ask for the bytes from the JSON
			//				- OR -> "please use iterator methods to act on anything besides a single return string/keyword"
			// 						- Is it useful to get a sub chunk of some JSON as a string/bytes? Like would a user ever need to retrieve (GET specifically) something like this?
			// 								- If not, should this resolve to one string/keyword and anything else iterator/explorer methods are used to act on parts of the JSON?
			switch v := val.(type) {
			case ast.Object:
				obj = v
				currentType = "object"
				break
			case ast.Array:
				arr = v
				currentType = "array"
				break
			case ast.Literal:
				fmt.Print(v)
			}
		}
	}

	return nil
}

func validateQueryRoot(query string, rootNodeType ast.RootNodeType) error {
	if query[0] != '$' {
		return ErrNoDollarSignRoot
	}

	// The query root after the `$` must be a `.` if the rootNodeType is an object
	validObjQueryRoot := query[1] == '.'
	if rootNodeType == ast.ObjectRoot && !validObjQueryRoot {
		return ErrWrongObjectRootSelector
	}

	// The query root after the `$` must be a `[` if the rootNodeType is an array
	validArrayQueryRoot := query[1] == '['
	if rootNodeType == ast.ArrayRoot && !validArrayQueryRoot {
		return ErrWrongArrayRootSelector
	}

	return nil
}

// ErrNoDollarSignRoot is used for telling the user the very first character must be a `$`
var ErrNoDollarSignRoot = errors.New(
	"Incorrect syntax, query must start with `$` representing the root object or array",
)

// ErrWrongObjectRootSelector is used for telling the user their JSON root is an object and the selector found was not a `.`
var ErrWrongObjectRootSelector = errors.New(
	"Incorrect syntax. Your root JSON type is an object. Therefore, path queries must" +
		"begin by selecting a `key` from your root object. Ex: `$.keyOnRootObject` or `$[\"keyOnRootObject\"]`",
)

// ErrWrongArrayRootSelector is used for telling the user their JSON root is an array and the selector found was not a `[`
var ErrWrongArrayRootSelector = errors.New(
	"Incorrect syntax. Your root JSON type is an array. Therefore, path queries must" +
		"begin by selecting an item by index on the root array. Ex: `$[0]` or `$[1]`",
)

func errSelectorSytax(operator string) error {
	return fmt.Errorf(
		"Error parsing query, expected either a `.` for selections on an object or a `[` for selections on an array. Got: %s",
		operator,
	)
}
