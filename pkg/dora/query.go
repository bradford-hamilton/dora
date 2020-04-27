package dora

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/bradford-hamilton/dora/pkg/ast"
)

// prepAndExecQuery prepares and executes a passed in query
func (c *Client) prepAndExecQuery(query string) error {
	if err := c.prepareQuery(query, c.tree.Type); err != nil {
		return err
	}
	if err := c.executeQuery(); err != nil {
		return err
	}
	return nil
}

// prepareQuery validates the query root, sets the query on the client struct, and parses the query.
func (c *Client) prepareQuery(query string, rootNodeType ast.RootNodeType) error {
	if err := validateQueryRoot(query, c.tree.Type); err != nil {
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

// parseQuery is pretty straight forward. Scan the query into tokens, set the the tokns
// to the `parsedQuery` field on the client.
func (c *Client) parseQuery() error {
	tokens, err := scanQueryTokens(c.query)
	if err != nil {
		return err
	}
	c.parsedQuery = tokens
	return nil
}

// Get takes a dora query, prepares and validates it, executes the query, and returns the result or an error.
func (c *Client) get(query string) (string, error) {
	if err := c.prepAndExecQuery(query); err != nil {
		return "", err
	}
	return c.result, nil
}

// setstring sets a string value at the location specified in the query
func (c *Client) setstring(query string, value string) error {
	return errors.New("implement me")
}

// setbool sets a boolean value at the location specified in the query
func (c *Client) setbool(query string, value bool) error {
	return errors.New("implement me")
}

// setbool sets a float64 value at the location specified in the query
func (c *Client) setfloat64(query string, value float64) error {
	return errors.New("implement me")
}

// Object is a simple constant used throughout executeQuery for checking conditions
const Object = "object"

// Array is a simple constant used throughout executeQuery for checking conditions
const Array = "array"

// executeQuery is called after the JSON and the query are parsed into their respective
// tokens. We then iterate over the query tokens, and traverse our tree attempting to
// find the result the user is looking for.
func (c *Client) executeQuery() error {
	rootVal := *c.tree.RootValue
	obj, _ := rootVal.(ast.Object)
	arr, ok := rootVal.(ast.Array)
	currentType := ast.ObjectType
	if ok {
		currentType = ast.ArrayType
	}
	parsedQueryLen := len(c.parsedQuery)

	for i := 0; i < parsedQueryLen; i++ {
		// If i == parsedQueryLen-1, we are on the final iteration
		if i == parsedQueryLen-1 {
			c.setFinalValue(currentType, i, obj, arr)
		}

		// If the query token we're on is asking for an object
		if c.parsedQuery[i].accessType == ObjectAccess {
			if currentType != ast.ObjectType {
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
						currentType = ast.ObjectType
						break
					}
					if astArr {
						arr = a
						currentType = ast.ArrayType
						break
					}
				}
			}
			if !found {
				return fmt.Errorf("Sorry, could not find a key with that value. Key: %s", c.parsedQuery[i].keyReq)
			}
		} else { // If the query token we're on is asking for an array
			if currentType != ast.ArrayType {
				return fmt.Errorf("TODO: error")
			}
			qt := c.parsedQuery[i]
			val := arr.Children[qt.indexReq]

			switch v := val.(type) {
			case ast.Object:
				obj = v
				currentType = ast.ObjectType
				break
			case ast.Array:
				arr = v
				currentType = ast.ArrayType
				break
			case ast.Literal:
				// If we're on the final value, return it
				if i == parsedQueryLen-1 {
					c.setResultFromLiteral(v.Value)
				} else {
					fmt.Println("TODO: error?")
				}
			}
		}
	}

	return nil
}

// setFinalValue is called when we are on the final queryToken. It handles narrowing down what
// needs to be returned and sets the result to the Client
func (c *Client) setFinalValue(currentType ast.Type, index int, obj ast.Object, arr ast.Array) {
	if currentType == ast.ObjectType {
		r := c.parsedQuery[index].keyReq
		for _, v := range obj.Children {
			if r == v.Key.Value {
				c.setResultFromValue(v.Value)
				break
			}
		}
		return
	}
	ind := c.parsedQuery[index].indexReq
	c.setResultFromValue(arr.Children[ind])
}

// setResultFromValue switches on an ast.Value type and assigns the appropriate result to the client
func (c *Client) setResultFromValue(value ast.Value) {
	switch val := value.(type) {
	case ast.Literal:
		c.setResultFromLiteral(val.Value)
	case ast.Object:
		c.result = string(c.input[val.Start:val.End])
	case ast.Array:
		c.result = string(c.input[val.Start:val.End])
	}
}

// setResultFromLiteral is very similar to setResultFromValue, except it we know the value we're switching over
// must be a Literal, meaning the assigned result will either be a string, number, boolean, or null
func (c *Client) setResultFromLiteral(value ast.Value) {
	switch lit := value.(type) {
	case string:
		c.result = lit
	case float64:
		c.result = fmt.Sprintf("%f", lit)
	case int:
		c.result = strconv.Itoa(lit)
	case bool:
		c.result = fmt.Sprintf("%v", lit)
	case nil:
		c.result = "null"
	}
}

// validateQueryRoot handles some very simple validation around the root of the query
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
