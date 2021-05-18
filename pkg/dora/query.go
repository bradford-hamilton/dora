package dora

import (
	"errors"
	"fmt"

	"github.com/bradford-hamilton/dora/pkg/ast"
	"github.com/bradford-hamilton/dora/pkg/danger"
)

var (
	// ErrNoDollarSignRoot is used for telling the user the very first character must be a `$`
	ErrNoDollarSignRoot = errors.New(
		"Incorrect syntax, query must start with `$` representing the root object or array",
	)
	// ErrWrongObjectRootSelector is used for telling the user their JSON root is an object and the selector found was not a `.`
	ErrWrongObjectRootSelector = errors.New(
		"Incorrect syntax. Your root JSON type is an object. Therefore, path queries must" +
			"begin by selecting a `key` from your root object. Ex: `$.keyOnRootObject` or `$[\"keyOnRootObject\"]`",
	)
	// ErrWrongArrayRootSelector is used for telling the user their JSON root is an array and the selector found was not a `[`
	ErrWrongArrayRootSelector = errors.New(
		"Incorrect syntax. Your root JSON type is an array. Therefore, path queries must" +
			"begin by selecting an item by index on the root array. Ex: `$[0]` or `$[1]`",
	)
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
	c.setQuery(danger.StringToBytes(query))
	if err := c.parseQuery(); err != nil {
		return err
	}
	return nil
}

func (c *Client) setQuery(query []byte) {
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

// executeQuery is called after the JSON and the query are parsed into their respective
// tokens. We then iterate over the query tokens, and traverse our tree attempting to
// find the result the user is looking for.
func (c *Client) executeQuery() error {
	rootVal := *c.tree.RootValue
	obj, _ := rootVal.Content.(ast.Object)
	arr, ok := rootVal.Content.(ast.Array)
	currentType := ast.ObjectType
	if ok {
		currentType = ast.ArrayType
	}
	parsedQueryLen := len(c.parsedQuery)

	for i := 0; i < parsedQueryLen; i++ {
		// If the query token we're on is asking for an object
		if c.parsedQuery[i].accessType == ObjectAccess {
			if currentType != ast.ObjectType {
				return errors.New("incorrect syntax, your query asked for an object but found array")
			}
			var found bool

			for _, v := range obj.Children {
				if v.Key.Value == c.parsedQuery[i].key {
					found = true
					val := v.Value

					// If i == parsedQueryLen-1, we are on the final iteration
					if i == parsedQueryLen-1 {
						c.setResultFromValue(val)
						return nil
					}

					o, astObj := val.Content.(ast.Object)
					if astObj {
						obj = o
						currentType = ast.ObjectType
						break
					}
					a, astArr := val.Content.(ast.Array)
					if astArr {
						arr = a
						currentType = ast.ArrayType
						break
					}
				}
			}
			if !found {
				return fmt.Errorf("sorry, could not find a key with that value. Key: %s", c.parsedQuery[i].key)
			}
		} else { // If the query token we're on is asking for an array
			if currentType != ast.ArrayType {
				return errors.New("incorrect syntax, your query asked for an array but found object")
			}
			qt := c.parsedQuery[i]
			val := arr.Children[qt.index].Value

			// If i == parsedQueryLen-1, we are on the final iteration
			if i == parsedQueryLen-1 {
				c.setResultFromValue(val)
				return nil
			}

			switch v := val.(type) {
			case ast.Object:
				obj = v
				currentType = ast.ObjectType
			case ast.Array:
				arr = v
				currentType = ast.ArrayType
			case ast.Literal:
				// If we're on the final value, return it
				if i == parsedQueryLen-1 {
					c.setResultFromLiteral(v.Value)
				} else {
					return errors.New("sorry, it looks like your query isn't quite right")
				}
			}
		}
	}

	return nil
}

// setResultFromValue switches on an ast.Value type and assigns the appropriate result to the client
func (c *Client) setResultFromValue(value ast.ValueContent) {
	if v2, ok := value.(ast.ArrayItem); ok {
		// unwrap the ArrayItem
		value = v2.Value
	}
	if v2, ok := value.(ast.Value); ok {
		// unwrap the Value
		value = v2.Content
	}
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
func (c *Client) setResultFromLiteral(value ast.ValueContent) {
	switch lit := value.(type) {
	case string:
		c.result = lit
	case float64:
		c.result = fmt.Sprintf("%f", lit)
	case int, int64:
		c.result = fmt.Sprintf("%d", lit)
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

func errSelectorSytax(operator string) error {
	return fmt.Errorf(
		"error parsing query, expected either a `.` for selections on an object or a `[` for selections on an array. Got: %s",
		operator,
	)
}
